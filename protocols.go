package sleego

import (
	"context"
	"time"
)

// AppConfig is the struct that will be used to store the configuration of each app
type AppConfig struct {
	Name        string `json:"name"`
	AllowedFrom string `json:"allowed_from"` // AllowedFrom is the initial hour that the app is allowed to be used
	AllowedTo   string `json:"allowed_to"`   // AllowedTo is the final hour that the app is allowed to be used
}

// ConfigLoader defines the behavior for loading application configurations.
type ConfigLoader interface {
	Load(path string) ([]AppConfig, error)
}

// ProcessInfo contains the information of a process
type ProcessInfo struct {
	Name string
	Pid  int
}

// Process defines the behavior of a process
type Process interface {
	GetInfo() (ProcessInfo, error)
	Kill() error
}

// ProcessorMonitor will be used to interact with the system processes
type ProcessorMonitor interface {
	GetRunningProcesses() ([]Process, error)
}

// ProcessPolicy controls when the application process will be terminated
type ProcessPolicy interface {
	Apply(ctx context.Context, appsConfig []AppConfig) error
}

// ShutdownPolicy defines the behavior for shutting down the system
type ShutdownPolicy interface {
	Apply(endTime time.Time) error
}
