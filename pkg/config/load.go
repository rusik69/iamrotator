package config

import (
	"os"

	yaml "gopkg.in/yaml.v3"
)

// Load loads the configuration from the given path
func Load(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	var cfg Config
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err
	}
	if cfg.AWS.Region == "" {
		cfg.AWS.Region = "us-east-1"
	}
	return cfg, nil
}
