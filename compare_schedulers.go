package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/amsen20/ecmus/internal/ankubedsm"
)

type PodWorkload struct {
	Id          string
	Nu          float64 // Sensitivity tier (0.0, 0.5, 1.0)
	RequiredCPU float64
}

type SimulationResult struct {
	PodId        string
	Nu           float64
	BaselineNode string
	BaselineLat  float64
	AnKubeNode   string
	AnKubeLat    float64
}

func main() {
	fmt.Println("======================================================================")
	fmt.Println("   AN-KubeDSM (Priority Sorted) vs. Baseline KubeDSM                  ")
	fmt.Println("======================================================================")

	nodesBaseline := []ankubedsm.NodeInfo{
		{Id: "core-edge (8ms)", IsCloud: false, FreeCPU: 4.0, Latency: ankubedsm.CoreEdgeLatency},
		{Id: "metro-edge (12ms)", IsCloud: false, FreeCPU: 4.0, Latency: ankubedsm.MetroEdgeLatency},
		{Id: "cloud-dc (100ms)", IsCloud: true, FreeCPU: 100.0, Latency: ankubedsm.CloudLatency},
	}

	nodesAnKube := []ankubedsm.NodeInfo{
		{Id: "core-edge (8ms)", IsCloud: false, FreeCPU: 4.0, Latency: ankubedsm.CoreEdgeLatency},
		{Id: "metro-edge (12ms)", IsCloud: false, FreeCPU: 4.0, Latency: ankubedsm.MetroEdgeLatency},
		{Id: "cloud-dc (100ms)", IsCloud: true, FreeCPU: 100.0, Latency: ankubedsm.CloudLatency},
	}

	workloadsBaseline := []PodWorkload{
		{"backup-job-1", 0.0, 2.0},
		{"log-processor", 0.0, 2.0},
		{"auth-service", 1.0, 2.0},
		{"telemetry-collector", 1.0, 2.0},
	}

	workloadsAnKube := []PodWorkload{
		{"backup-job-1", 0.0, 2.0},
		{"log-processor", 0.0, 2.0},
		{"auth-service", 1.0, 2.0},
		{"telemetry-collector", 1.0, 2.0},
	}

	// AN-KubeDSM sorts incoming workloads by sensitivity (nu descending)
	sort.Slice(workloadsAnKube, func(i, j int) bool {
		return workloadsAnKube[i].Nu > workloadsAnKube[j].Nu
	})

	fmt.Printf("%-22s | %-12s | %-25s | %-25s\n", "Workload Pod", "Tier (nu)", "Baseline KubeDSM", "AN-KubeDSM (Ours)")
	fmt.Println("------------------------------------------------------------------------------------------------------")

	csvFile, _ := os.Create("contention_results.csv")
	defer csvFile.Close()
	csvFile.WriteString("PodId,SensitivityNu,BaselineNode,BaselineLatency,AnKubeNode,AnKubeLatency\n")

	var baselineTotalLatency, anKubeTotalLatency float64
	results := make([]SimulationResult, len(workloadsBaseline))

	// 1. Simulate Baseline KubeDSM (Arrival Order)
	for idx, wl := range workloadsBaseline {
		nodeId, lat := simulateDynamicAllocation(&nodesBaseline, wl.RequiredCPU)
		baselineTotalLatency += lat
		results[idx].PodId = wl.Id
		results[idx].Nu = wl.Nu
		results[idx].BaselineNode = nodeId
		results[idx].BaselineLat = lat
	}

	// 2. Simulate AN-KubeDSM (Priority Sorted)
	for _, wl := range workloadsAnKube {
		nodeId, lat, success := ankubedsm.SelectOptimalNodeExternal(wl.Id, wl.Nu, wl.RequiredCPU, nodesAnKube)
		if success {
			for idx := range nodesAnKube {
				if nodesAnKube[idx].Id == nodeId {
					nodesAnKube[idx].FreeCPU -= wl.RequiredCPU
				}
			}
		} else {
			nodeId = "Unscheduled"
			lat = 0.0
		}
		anKubeTotalLatency += lat

		// Map back to original workload index for clean side-by-side display
		for idx, origWl := range workloadsBaseline {
			if origWl.Id == wl.Id {
				results[idx].AnKubeNode = nodeId
				results[idx].AnKubeLat = lat
			}
		}
	}

	// Print comparison table & write CSV (Fixed unused index variable `i`)
	for _, res := range results {
		fmt.Printf("%-22s | nu=%.1f      | %-25s | %-25s\n",
			res.PodId, res.Nu,
			fmt.Sprintf("%s (%.0fms)", res.BaselineNode, res.BaselineLat),
			fmt.Sprintf("%s (%.0fms)", res.AnKubeNode, res.AnKubeLat),
		)
		csvFile.WriteString(fmt.Sprintf("%s,%.1f,%s,%.1f,%s,%.1f\n",
			res.PodId, res.Nu, res.BaselineNode, res.BaselineLat, res.AnKubeNode, res.AnKubeLat))
	}

	fmt.Println("------------------------------------------------------------------------------------------------------")
	fmt.Printf("Average Latency  |              | %-25s | %-25s\n",
		fmt.Sprintf("%.1f ms", baselineTotalLatency/float64(len(workloadsBaseline))),
		fmt.Sprintf("%.1f ms", anKubeTotalLatency/float64(len(workloadsBaseline))),
	)
	fmt.Println("======================================================================")
	fmt.Println("[Success] Accurate contention metrics exported to 'contention_results.csv'")
}

func simulateDynamicAllocation(nodes *[]ankubedsm.NodeInfo, requiredCPU float64) (string, float64) {
	for i := range *nodes {
		if (*nodes)[i].FreeCPU >= requiredCPU {
			(*nodes)[i].FreeCPU -= requiredCPU
			return (*nodes)[i].Id, (*nodes)[i].Latency
		}
	}
	return "None", 0.0
}