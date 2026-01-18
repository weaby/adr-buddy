package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	cfg := Default()

	assert.Equal(t, ".", cfg.ScanPaths[0])
	assert.Equal(t, "decisions", cfg.OutputDir)
	assert.False(t, cfg.StrictMode)
	assert.Contains(t, cfg.Exclude, "**/node_modules/**")
	assert.Contains(t, cfg.Exclude, "**/.git/**")
	assert.Contains(t, cfg.Exclude, "**/vendor/**")
}

func TestLoad(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".adr-buddy", "config.yml")

	// Create .adr-buddy directory
	err := os.MkdirAll(filepath.Dir(configPath), 0755)
	assert.NoError(t, err)

	// Write test config
	configContent := `scan_paths:
  - ./src
  - ./lib
output_dir: ./docs/decisions
strict_mode: true
`
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	assert.NoError(t, err)

	// Load config
	cfg, err := Load(tmpDir)
	assert.NoError(t, err)
	assert.Equal(t, []string{"./src", "./lib"}, cfg.ScanPaths)
	assert.Equal(t, "./docs/decisions", cfg.OutputDir)
	assert.True(t, cfg.StrictMode)
}

func TestLoad_NoConfigFile(t *testing.T) {
	tmpDir := t.TempDir()

	cfg, err := Load(tmpDir)
	assert.NoError(t, err)

	// Should return default config
	assert.Equal(t, Default().ScanPaths, cfg.ScanPaths)
	assert.Equal(t, Default().OutputDir, cfg.OutputDir)
}
