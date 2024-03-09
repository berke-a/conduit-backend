package shared

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Host     string `yaml:"POSTGRES_HOST"`
	Port     int    `yaml:"POSTGRES_PORT"`
	User     string `yaml:"POSTGRES_USER"`
	Password string `yaml:"POSTGRES_PASSWORD"`
	DBName   string `yaml:"POSTGRES_DBNAME"`
}

// LoadConfig reads a YAML-config file and unmarshals it into Config struct
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
