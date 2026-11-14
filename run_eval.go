package main

import (
	"os"

	"github.com/amsen20/ecmus/internal/config"
	"github.com/amsen20/ecmus/sim"
	"gopkg.in/yaml.v2"
	"fmt"
	"github.com/amsen20/ecmus/internal/ankubedsm"
	"github.com/amsen20/ecmus/internal/model"
	"github.com/amsen20/ecmus/internal/model/testing_tool"
)

func main() {
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}

	if err := yaml.UnmarshalStrict(yamlFile, &config.SchedulerGeneralConfig); err != nil {
		panic(err)
	}

	println("Running QASRE simulation...")
	sim.Start("QASRE")

	sim.ResetReport()

	println("Running Kubernetes Default simulation...")
	sim.Start("KUBERNETES_DEFAULT")

	println("Simulations completed successfully!")
}
