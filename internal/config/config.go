package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the adr-buddy configuration
type Config struct {
	ScanPaths  []string `yaml:"scan_paths"`
	OutputDir  string   `yaml:"output_dir"`
	Exclude    []string `yaml:"exclude"`
	Template   string   `yaml:"template"`
	StrictMode bool     `yaml:"strict_mode"`
}

// Default returns the default configuration
func Default() *Config {
	return &Config{
		ScanPaths: []string{"."},
		OutputDir: "decisions",
		Exclude: []string{
			"**/node_modules/**",
			"**/.git/**",
			"**/vendor/**",
			"**/build/**",
			"**/dist/**",
			"**/.next/**",
			"**/.adr-buddy/**",
			"**/.claude/**",
			"**/.github/**",
		},
		Template:   "",
		StrictMode: false,
	}
}

// Load loads configuration from the specified directory
// Returns default config if no config file exists
func Load(rootDir string) (*Config, error) {
	configPath := filepath.Join(rootDir, ".adr-buddy", "config.yml")

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return Default(), nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// Parse YAML
	cfg := Default()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
