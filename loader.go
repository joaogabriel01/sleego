package sleego

import (
	"encoding/json"
	"os"
)

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
