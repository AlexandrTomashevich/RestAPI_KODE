package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Database struct {
		ConnectionString string `yaml:"connection_string"`
	} `yaml:"database"`
}

func LoadConfig(filename string) (Config, error) {
	var config Config

	file, err := os.ReadFile(filename) // Изменил на ReadFile для удобства
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
