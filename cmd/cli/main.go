package main

import (
	"context"
	"flag"
	"log"
	"time"

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
	appPolicy := sleego.NewProcessPolicyImpl(monitor, nil, nil)

	shutdownChannel := make(chan string)
	shutdownPolicy := sleego.NewShutdownPolicyImpl(shutdownChannel, []int{})
	shutdownTime, err := time.Parse("15:04", config.Shutdown)
	if err != nil {
		log.Fatalf("Error parsing shutdown time: %v", err)
	}

	log.Printf("Starting process policy with config: %+v of path: %s", config, *configPath)
	go appPolicy.Apply(ctx, config.Apps)
	go shutdownPolicy.Apply(ctx, shutdownTime)

}
