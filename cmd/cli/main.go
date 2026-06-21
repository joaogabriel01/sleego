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
	categoryOp := sleego.GetCategoryOperator()
	config, err := loadConfig(*configPath, loader, categoryOp)
	if err != nil {
		loggerInstance.Error(err.Error())
		os.Exit(1)
	}

	monitor := &sleego.ProcessorMonitorImpl{}
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

func loadConfig(path string, loader sleego.ConfigLoader, categoryOp sleego.CategoryOperator) (sleego.FileConfig, error) {
	config, err := loader.Load(path)
	if err != nil {
		return sleego.FileConfig{}, fmt.Errorf("Error loading config file: %w", err)
	}

	if err := sleego.ValidateConfig(config); err != nil {
		return sleego.FileConfig{}, fmt.Errorf("Invalid config: %w", err)
	}

	categoryOp.SetProcessByCategories(config.Categories)
	return config, nil
}
