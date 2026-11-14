package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/amsen20/ecmus/internal/ankubedsm" // Your custom package
	"github.com/amsen20/ecmus/internal/config"
	"github.com/amsen20/ecmus/internal/connector"
	"github.com/amsen20/ecmus/internal/gui"
	"github.com/amsen20/ecmus/internal/model"
	"github.com/amsen20/ecmus/internal/scheduler"
	"github.com/amsen20/ecmus/logging"
	"github.com/amsen20/ecmus/statistics"
	"gopkg.in/yaml.v2"
)

var log = logging.Get()

func main() {
	config_file_path := flag.String("config_file", "", "Path to config file")
	runAnKubeDSM := flag.Bool("ankubedsm", false, "Run with AN-KubeDSM extension logic")
	flag.Parse()

	if *runAnKubeDSM {
		fmt.Println("[AN-KubeDSM] Running custom sensitivity-aware & heterogeneous evaluation mode...")
		
		// Quick verification of our heterogeneous nodes and priority weights
		nodes := []ankubedsm.NodeInfo{
			{Id: "core-edge-8ms", IsCloud: false, FreeCPU: 12.0, Latency: ankubedsm.CoreEdgeLatency},
			{Id: "cloud-100ms", IsCloud: true, FreeCPU: 100.0, Latency: ankubedsm.CloudLatency},
		}

		assignedNode, latency, success := ankubedsm.SelectOptimalNodeExternal("critical-service", 1.0, 2.0, nodes)
		if success {
			fmt.Printf("[AN-KubeDSM] Test Placement -> Node: %s, Latency: %.0f ms, Weight w[i]: %.1f\n", 
				assignedNode, latency, ankubeGetWeight(1.0))
		}
		return
	}

	if *config_file_path == "" {
		log.Error().Msg("config_file flag is required")
		os.Exit(1)
	}

	yamlFile, err := os.ReadFile(*config_file_path)
	if err != nil {
		log.Err(err).Msgf("could not load config")
		os.Exit(1)
	}

	if err := yaml.UnmarshalStrict(yamlFile, &config.SchedulerGeneralConfig); err != nil {
		log.Err(err).Msgf("could not load config")
		os.Exit(1)
	}

	clusterState := model.NewClusterState()

	var c connector.Connector
	switch config.SchedulerGeneralConfig.ConnectorKind {
	case "const":
		c = connector.NewConstantConnector(clusterState)
	case "kubernetes":
		c, err = connector.NewKubeConnector(clusterState)
		if err != nil {
			log.Err(err).Msg("could not init the connector")
			os.Exit(1)
		}
	default:
		log.Error().Msg("connector kind is not recognized")
		os.Exit(1)
	}

	sched, err := scheduler.New(clusterState, c)
	if err != nil {
		log.Err(err).Msg("could not initiate scheduler")
		os.Exit(1)
	}

	if err := sched.Start(); err != nil {
		log.Err(err).Msg("could not start scheduler")
		os.Exit(1)
	}

	schedulerContext := context.Background()
	schedulerBridge, err := sched.Run(schedulerContext)
	if err != nil {
		log.Err(err).Msg("could not run scheduler")
		os.Exit(1)
	}

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	statistics.Init()
	statistics.Set("restarts", 0)

	gui.SetUp(schedulerBridge)
	go gui.Run()

	<-signalChannel
	log.Info().Msgf("exiting gracefully...")
	log.Info().Msgf("\n%s", statistics.Display())
}

func ankubeGetWeight(nu float64) float64 {
	return ankubedsm.GetServiceWeight(nu)
}