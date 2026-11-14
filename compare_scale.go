package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/amsen20/ecmus/internal/ankubedsm"
)

type PodWorkload struct {
	Id          string
	Nu          float64
	RequiredCPU float64
}

type TierStats struct {
	Total      int
	EdgeCount  int
	CloudCount int
}

func main() {
	fmt.Println("=======================================================================================")
	fmt.Println("             AN-KubeDSM vs. Baseline: Scale Summary Report (N=100)                     ")
	fmt.Println("=======================================================================================")

	rand.Seed(time.Now().UnixNano())

	var workloadsBaseline []PodWorkload
	var workloadsAnKube []PodWorkload

	for i := 1; i <= 100; i++ {
		var nu float64
		r := rand.Float64()
		if r < 0.3 {
			nu = 1.0
		} else if r < 0.6 {
			nu = 0.5
		} else {
			nu = 0.0
		}

		pod := PodWorkload{Id: fmt.Sprintf("pod-%d", i), Nu: nu, RequiredCPU: 2.0}
		workloadsBaseline = append(workloadsBaseline, pod)
		workloadsAnKube = append(workloadsAnKube, pod)
	}

	sort.Slice(workloadsAnKube, func(i, j int) bool {
		return workloadsAnKube[i].Nu > workloadsAnKube[j].Nu
	})

	getFreshNodes := func() []ankubedsm.NodeInfo {
		return []ankubedsm.NodeInfo{
			{Id: "core-edge (8ms)", IsCloud: false, FreeCPU: 35.0, Latency: ankubedsm.CoreEdgeLatency},
			{Id: "metro-edge (12ms)", IsCloud: false, FreeCPU: 35.0, Latency: ankubedsm.MetroEdgeLatency},
			{Id: "cloud-dc (100ms)", IsCloud: true, FreeCPU: 1000.0, Latency: ankubedsm.CloudLatency},
		}
	}

	nodesBase := getFreshNodes()
	baseLat, baseCrit, baseNorm, baseBatch := runExclusiveSimulation(workloadsBaseline, &nodesBase, false)

	nodesAnKube := getFreshNodes()
	anLat, anCrit, anNorm, anBatch := runExclusiveSimulation(workloadsAnKube, &nodesAnKube, true)

	// Print Clean Single-Line Columnar Table
	fmt.Printf("%-32s | %-24s | %-24s\n", "Evaluation Metric", "Baseline KubeDSM", "AN-KubeDSM (Ours)")
	fmt.Println("-----------------------------------------------------------------------------------------------")
	fmt.Printf("%-32s | %.1f ms                  | %.1f ms\n", "Average Latency", baseLat, anLat)
	printRow("Critical Tier (nu=1.0) Edge", baseCrit, anCrit)
	printRow("Normal Tier (nu=0.5) Edge", baseNorm, anNorm)
	printRow("Batch Tier (nu=0.0) Cloud", baseBatch, anBatch)
	fmt.Println("=======================================================================================")
	fmt.Println("[Success] Clean scale evaluation successfully completed.")
}

func runExclusiveSimulation(workloads []PodWorkload, nodes *[]ankubedsm.NodeInfo, isSmart bool) (float64, TierStats, TierStats, TierStats) {
	var totalLatency float64
	crit := TierStats{}
	norm := TierStats{}
	batch := TierStats{}

	for _, wl := range workloads {
		if wl.Nu == 1.0 { crit.Total++ }
		if wl.Nu == 0.5 { norm.Total++ }
		if wl.Nu == 0.0 { batch.Total++ }

		var nodeId string
		var lat float64
		var success bool

		if isSmart {
			nodeId, lat, success = ankubedsm.SelectOptimalNodeExternal(wl.Id, wl.Nu, wl.RequiredCPU, *nodes)
		} else {
			nodeId, lat = simulateDynamicAllocation(nodes, wl.RequiredCPU)
			success = (nodeId != "none")
		}

		if success {
			for i := range *nodes {
				if (*nodes)[i].Id == nodeId {
					(*nodes)[i].FreeCPU -= wl.RequiredCPU
				}
			}
			totalLatency += lat

			isEdge := (nodeId == "core-edge (8ms)" || nodeId == "metro-edge (12ms)")

			if wl.Nu == 1.0 {
				if isEdge { crit.EdgeCount++ } else { crit.CloudCount++ }
			} else if wl.Nu == 0.5 {
				if isEdge { norm.EdgeCount++ } else { norm.CloudCount++ }
			} else if wl.Nu == 0.0 {
				if isEdge { batch.EdgeCount++ } else { batch.CloudCount++ }
			}
		}
	}

	return totalLatency / float64(len(workloads)), crit, norm, batch
}

func printRow(label string, base, an TierStats) {
	baseStr := fmt.Sprintf("%d/%d (%.1f%%)", base.EdgeCount, base.Total, safePct(base.EdgeCount, base.Total))
	anStr := fmt.Sprintf("%d/%d (%.1f%%)", an.EdgeCount, an.Total, safePct(an.EdgeCount, an.Total))
	if label == "Batch Tier (nu=0.0) Cloud" {
		baseStr = fmt.Sprintf("%d/%d (%.1f%%)", base.CloudCount, base.Total, safePct(base.CloudCount, base.Total))
		anStr = fmt.Sprintf("%d/%d (%.1f%%)", an.CloudCount, an.Total, safePct(an.CloudCount, an.Total))
	}
	fmt.Printf("%-32s | %-24s | %-24s\n", label, baseStr, anStr)
}

func safePct(count, total int) float64 {
	if total == 0 { return 0.0 }
	return float64(count) * 100.0 / float64(total)
}

func simulateDynamicAllocation(nodes *[]ankubedsm.NodeInfo, requiredCPU float64) (string, float64) {
	for i := range *nodes {
		if (*nodes)[i].FreeCPU >= requiredCPU {
			(*nodes)[i].FreeCPU -= requiredCPU
			return (*nodes)[i].Id, (*nodes)[i].Latency
		}
	}
	return "none", 0.0
}