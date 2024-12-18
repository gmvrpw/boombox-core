package config

import (
	"os"

	"github.com/go-yaml/yaml"
	"github.com/google/uuid"
)

type Config struct {
	Request RequestsConfig `yaml:"requests"`
	Runners []RunnerConfig `yaml:"runners"`
}

func NewConfig(path string) (*Config, error) {
	cfg := Config{}

	if path != "" {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		err = yaml.NewDecoder(file).Decode(&cfg)
		if err != nil {
			return nil, err
		}
	}

	return &cfg, nil
}

type RequestsConfig struct {
	Url string `yaml:"url"`
}

type RunnerConfig struct {
	Name  string    `yaml:"name"`
	Owner uuid.UUID `yaml:"owner"`
	Url   string    `yaml:"url"`
	Test  string    `yaml:"test"`
}
