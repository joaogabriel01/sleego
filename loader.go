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

var _ ConfigLoader = &Loader{}
