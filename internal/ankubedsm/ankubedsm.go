package ankubedsm

import (
	"math"
	"sort"
)

// Sensitivity Tiers & Priority Weights
func GetServiceWeight(nu float64) float64 {
	return 1.0 + nu
}
const (
	CoreEdgeLatency  = 8.0
	MetroEdgeLatency = 12.0
	OuterEdgeLatency = 20.0
	CloudLatency     = 100.0
)
func WeightedQoS(edgePodsCount, totalPodsCount int, targetQoS float64, nu float64) float64 {
	if totalPodsCount == 0 {
		return 0.0
	}
	actualRatio := float64(edgePodsCount) / float64(totalPodsCount)
	delta := actualRatio - targetQoS
	
	baseF := math.Max(0.0, 1.0-math.Abs(delta))
	return GetServiceWeight(nu) * baseF
}

type PodMetadata struct {
	Id       string
	NuWeight float64
}

// ReorderEdgeTieBreak sorts edge pods by nu ascending.
func ReorderEdgeTieBreak(edgePods []PodMetadata) []PodMetadata {
	sortedPods := make([]PodMetadata, len(edgePods))
	copy(sortedPods, edgePods)

	sort.Slice(sortedPods, func(i, j int) bool {
		if sortedPods[i].NuWeight == sortedPods[j].NuWeight {
			return sortedPods[i].Id < sortedPods[j].Id
		}
		return sortedPods[i].NuWeight < sortedPods[j].NuWeight
	})

	return sortedPods
}

// NodeInfo represents a node with its heterogeneous latency and capacity
type NodeInfo struct {
	Id      string
	IsCloud bool
	FreeCPU float64
	Latency float64
}

// SelectOptimalNodeExternal evaluates nodes using AN-KubeDSM rules
func SelectOptimalNodeExternal(podId string, nu float64, requiredCPU float64, nodes []NodeInfo) (string, float64, bool) {
	var bestNodeId string
	var bestLatency float64
	bestScore := -999999.0
	foundValidNode := false

	for _, node := range nodes {
		// 1. Filter by resource capacity
		if node.FreeCPU < requiredCPU {
			continue
		}

		var score float64

		if !node.IsCloud {
			proximityBonus := 500.0 / node.Latency
			score = node.FreeCPU + (proximityBonus * (1.0 + nu * 10.0))
		} else {
			cloudPenalty := node.Latency * (1.0 + nu * 50.0)
			score = node.FreeCPU - cloudPenalty
		}

		if score > bestScore {
			bestScore = score
			bestNodeId = node.Id
			bestLatency = node.Latency
			foundValidNode = true
		}
	}

	if !foundValidNode {
		return "", 0.0, false
	}

	return bestNodeId, bestLatency, true
}