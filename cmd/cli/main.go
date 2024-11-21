package main

import (
	"context"
	"flag"
	"log"

	"github.com/joaogabriel01/sleego"
)

func main() {
	ctx := context.Background()
	configPath := flag.String("config", "./config.json", "Path to config file")
	flag.Parse()

	loader := &sleego.Loader{}
	config, err := loader.Load(*configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	monitor := &sleego.ProcessorMonitorImpl{}
	policy := sleego.NewProcessPolicyImpl(monitor, nil, nil)
	log.Printf("Starting process policy with config: %+v of path: %s", config, *configPath)
	policy.Apply(ctx, config)
}
