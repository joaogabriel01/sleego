package sleego

import (
	"encoding/json"
	"os"
)

// ConfigLoader defines the behavior for loading application configurations.
type ConfigLoader interface {
	Load(path string) (FileConfig, error)
}

// FileConfig is the struct that will be used to store the configuration of the apps
type FileConfig struct {
	Apps       []AppConfig         `json:"apps"`
	Shutdown   string              `json:"shutdown"`
	Categories map[string][]string `json:"categories"`
}

// AppConfig is the struct that will be used to store the configuration of each app
type AppConfig struct {
	Name        string `json:"name"`
	AllowedFrom string `json:"allowed_from"` // AllowedFrom is the initial hour that the app is allowed to be used
	AllowedTo   string `json:"allowed_to"`   // AllowedTo is the final hour that the app is allowed to be used
}

type Loader struct {
}

func (l *Loader) Load(path string) (FileConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return FileConfig{}, err
	}
	defer file.Close()
	var fileConfig FileConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&fileConfig)
	if err != nil {
		return FileConfig{}, err
	}
	return fileConfig, err
}

func (l *Loader) Save(path string, config FileConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

var _ ConfigLoader = &Loader{}
