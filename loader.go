package sleego

import (
	"encoding/json"
	"os"
)

type Loader struct {
}

func (l *Loader) Load(path string) ([]AppConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return []AppConfig{}, err
	}
	defer file.Close()
	var appConfigs []AppConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&appConfigs)
	if err != nil {
		return []AppConfig{}, err
	}
	return appConfigs, err
}

func (l *Loader) Save(path string, configs []AppConfig) error {
	data, err := json.MarshalIndent(configs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

var _ ConfigLoader = &Loader{}
