package config

import (
	"fmt"
	validator "github.com/asaskevich/govalidator"
	"gopkg.in/yaml.v3"
	"os"
)

const (
	defaultLogLevel = "info"
)

type Config struct {
	Log struct {
		Level string `yaml:"level"`
	} `yaml:"log"`

	Temporal struct {
		Host          string `yaml:"host" valid:"host,required"`
		Port          uint32 `yaml:"port" valid:"port"`
		Namespace     string `yaml:"namespace" valid:"stringlength(1|100)"`
		Queue         string `yaml:"queue"`
		MaxConcurrent uint32 `yaml:"max_concurrent"`
	} `yaml:"temporal"`
}

func (c *Config) Validate() error {
	if c.Log.Level == "" {
		c.Log.Level = defaultLogLevel
	}

	if valid, err := validator.ValidateStruct(c); !valid || err != nil {
		return fmt.Errorf("invalid configuration: %v", err)
	}

	return nil
}

func (c *Config) Load(path string) error {
	if path != "" {
		configBytes, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to load configuration file at '%s': %v", path, err)
		}

		if err = yaml.Unmarshal(configBytes, c); err != nil {
			return fmt.Errorf("failed to parse configuration: %v", err)
		}
	}

	return nil
}
