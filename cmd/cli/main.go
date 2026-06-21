package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
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
		*logLevel = "info"
	}

	loggerInstance, err := logger.Get(*logLevel)
	if err != nil {
		log.Fatalf("Error getting logger instance: %v", err)
	}

	loader := &sleego.Loader{}
	config, err := loader.Load(*configPath)
	if err != nil {
		loggerInstance.Error("Error loading config file: " + err.Error())
		os.Exit(1)
	}

	if err := sleego.ValidateConfig(config); err != nil {
		loggerInstance.Error("Invalid config: " + err.Error())
		os.Exit(1)
	}

	monitor := &sleego.ProcessorMonitorImpl{}
	categoryOp := sleego.GetCategoryOperator()
	appPolicy := sleego.NewProcessPolicyImpl(monitor, categoryOp, nil, nil)

	loggerInstance.Info("Starting process policy with config: " + *configPath)
	go appPolicy.Apply(ctx, config.Apps)

	if config.Shutdown != "" {
		shutdownChannel := make(chan string)
		shutdownPolicy := sleego.NewShutdownPolicyImpl(shutdownChannel, []int{})
		shutdownTime, err := time.Parse("15:04", config.Shutdown)
		if err != nil {
			loggerInstance.Error("Error parsing shutdown time: " + err.Error())
			os.Exit(1)
		}

		loggerInstance.Info("Starting shutdown policy with config: " + *configPath)
		go shutdownPolicy.Apply(ctx, shutdownTime)
	}

	select {}

}
