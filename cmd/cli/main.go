package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/joaogabriel01/sleego"
	"github.com/joaogabriel01/sleego/internal/logger"
)

func main() {
	ctx := context.Background()
	configPath := flag.String("config", "./config.json", "Path to config file")
	logLevel := flag.String("loglevel", "info", "Log level (debug, info, warn, error)")
	flag.Parse()
	fmt.Println("Log level set to:", *logLevel)

	if *logLevel != "debug" && *logLevel != "info" && *logLevel != "warn" && *logLevel != "error" {
		log.Fatalf("Invalid log level: %s", logLevel)
	}

	logger.Init(*logLevel)
	loggerInstance, err := logger.Get()
	if err != nil {
		log.Fatalf("Error getting logger instance: %v", err)
	}

	loader := &sleego.Loader{}
	config, err := loader.Load(*configPath)
	if err != nil {
		loggerInstance.Error("Error loading config file: " + err.Error())
	}

	monitor := &sleego.ProcessorMonitorImpl{}
	appPolicy := sleego.NewProcessPolicyImpl(monitor, nil, nil)

	shutdownChannel := make(chan string)
	shutdownPolicy := sleego.NewShutdownPolicyImpl(shutdownChannel, []int{})
	shutdownTime, err := time.Parse("15:04", config.Shutdown)
	if err != nil {
		loggerInstance.Error("Error parsing shutdown time: " + err.Error())
	}

	loggerInstance.Info("Starting process policy with config: " + *configPath)
	loggerInstance.Info("Starting shutdown policy with config: " + *configPath)

	go appPolicy.Apply(ctx, config.Apps)
	go shutdownPolicy.Apply(ctx, shutdownTime)

	select {}

}
