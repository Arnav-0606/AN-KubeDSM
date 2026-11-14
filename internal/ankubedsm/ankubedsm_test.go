package ankubedsm

import (
	"testing"
)

func TestGetServiceWeight(t *testing.T) {
	tests := []struct {
		nu       float64
		expected float64
	}{
		{0.0, 1.0}, // Batch
		{0.5, 1.5}, // Normal
		{1.0, 2.0}, // Latency-critical
	}

	for _, tt := range tests {
		got := GetServiceWeight(tt.nu)
		if got != tt.expected {
			t.Errorf("GetServiceWeight(%.1f) = %.1f; expected %.1f", tt.nu, got, tt.expected)
		}
	}
}

func TestWeightedQoS(t *testing.T) {
	if score := WeightedQoS(0, 0, 0.8, 1.0); score != 0.0 {
		t.Errorf("Expected 0.0 for zero total pods, got %f", score)
	}

	score := WeightedQoS(8, 10, 0.8, 1.0)
	expected := 2.0
	if score != expected {
		t.Errorf("Expected weighted QoS %.1f, got %.1f", expected, score)
	}
}

func TestReorderEdgeTieBreak(t *testing.T) {
	pods := []PodMetadata{
		{Id: "critical-pod", NuWeight: 1.0},
		{Id: "batch-pod", NuWeight: 0.0},
		{Id: "normal-pod", NuWeight: 0.5},
	}

	sorted := ReorderEdgeTieBreak(pods)

	// order: batch (0.0) -> normal (0.5) -> critical (1.0)
	if sorted[0].Id != "batch-pod" || sorted[1].Id != "normal-pod" || sorted[2].Id != "critical-pod" {
		t.Errorf("Pods were not sorted correctly by sensitivity tier: %v", sorted)
	}
}

func TestSelectOptimalNodeExternal(t *testing.T) {
	nodes := []NodeInfo{
		{Id: "cloud-node", IsCloud: true, FreeCPU: 100.0, Latency: CloudLatency},
		{Id: "edge-node", IsCloud: false, FreeCPU: 4.0, Latency: MetroEdgeLatency},
	}

	// Critical pod (nu = 1.0) should prefer edge over cloud due to weight multiplier
	bestNode, latency, success := SelectOptimalNodeExternal("test-pod", 1.0, 2.0, nodes)
	if !success {
		t.Fatalf("Expected successful node selection")
	}

	if bestNode != "edge-node" {
		t.Errorf("Expected edge-node for critical workload, got %s", bestNode)
	}
	if latency != MetroEdgeLatency {
		t.Errorf("Expected latency %.1f, got %.1f", MetroEdgeLatency, latency)
	}
}