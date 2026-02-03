package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	cfg := Default()
	assert.Equal(t, filepath.Join(".claude", "rules", "decisions"), cfg.DecisionsDir)
	assert.False(t, cfg.StrictMode)
}

func TestLoad(t *testing.T) {
	tmpDir := t.TempDir()

	configDir := filepath.Join(tmpDir, ".adr-buddy")
	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	configContent := `decisions_dir: custom/decisions
strict_mode: true
`
	err = os.WriteFile(filepath.Join(configDir, "config.yml"), []byte(configContent), 0644)
	require.NoError(t, err)

	cfg, err := Load(tmpDir)
	require.NoError(t, err)

	assert.Equal(t, "custom/decisions", cfg.DecisionsDir)
	assert.True(t, cfg.StrictMode)
}

func TestLoad_NoConfigFile(t *testing.T) {
	tmpDir := t.TempDir()

	cfg, err := Load(tmpDir)
	require.NoError(t, err)

	assert.Equal(t, Default().DecisionsDir, cfg.DecisionsDir)
	assert.False(t, cfg.StrictMode)
}

func TestDecisionsPath(t *testing.T) {
	cfg := Default()
	path := cfg.DecisionsPath("/project")
	assert.Equal(t, filepath.Join("/project", ".claude", "rules", "decisions"), path)
}
