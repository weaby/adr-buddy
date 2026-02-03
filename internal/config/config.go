package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DecisionsDir string `yaml:"decisions_dir"`
	StrictMode   bool   `yaml:"strict_mode"`
}

func Default() *Config {
	return &Config{
		DecisionsDir: filepath.Join(".claude", "rules", "decisions"),
		StrictMode:   false,
	}
}

// Load returns default config if no config file exists.
func Load(rootDir string) (*Config, error) {
	configPath := filepath.Join(rootDir, ".adr-buddy", "config.yml")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return Default(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	cfg := Default()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) DecisionsPath(rootDir string) string {
	if filepath.IsAbs(c.DecisionsDir) {
		return c.DecisionsDir
	}
	return filepath.Join(rootDir, c.DecisionsDir)
}
