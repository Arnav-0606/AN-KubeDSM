package sim

import (
	"testing"

	"github.com/amsen20/ecmus/internal/model/testing_tool"
)

func TestSimulation(t *testing.T) {
	// 1. Initialize the testing builder
	builder := testing_tool.New()

	// 2. Register all required deployments including "service_a"
	deploymentsDesc := []*testing_tool.DeploymentDesc{
		{
			Name:      "service_a",
			Cpu:       2.0,
			Memory:    2.0,
			EdgeShare: 0.8,
		},
		// Add any other deployments your simulation scenarios expect here
	}
	builder.ImportDeployments(deploymentsDesc)

	// 3. Define your mock edge node descriptors and their deployed pods
	edgeNodes := map[*testing_tool.NodeDesc][]string{
		{Cpu: 10.0, Memory: 10.0}: {"service_a"},
	}

	// 4. Define cloud-offloaded pods if applicable
	cloudPods := []string{}

	// 5. Build and initialize the cluster state
	clusterState := builder.GetCluster(edgeNodes, cloudPods)

	if clusterState == nil {
		t.Fatal("Failed to initialize cluster state in simulation test")
	}
}