package config

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/demoManito/pulse/pkg/database"
	"github.com/demoManito/pulse/pkg/gitops"
	"github.com/demoManito/pulse/pkg/logger"
	"github.com/demoManito/pulse/pkg/wecom"
)

type Config struct {
	HTTP     HTTPConfig      `yaml:"http"`
	Database database.Config `yaml:"database"`
	WeCom    wecom.Config    `yaml:"wecom"`
	GitOps   gitops.Config   `yaml:"gitops"`
}

type HTTPConfig struct {
	Port    int    `yaml:"port"`
	Address string `yaml:"address"`
}

// LoadConfig loads config from file.
func LoadConfig(filepath string) (*Config, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Errorf("Failed to close config file: %v", err)
			// PASS
		}
	}(file)

	var cfg *Config
	err = yaml.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
