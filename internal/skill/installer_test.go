package skill

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstallProjectLevel(t *testing.T) {
	tmpDir := t.TempDir()

	err := InstallProjectLevel(tmpDir)
	require.NoError(t, err)

	skillPath := filepath.Join(tmpDir, ".claude", "skills", "adr.md")
	content, err := os.ReadFile(skillPath)
	require.NoError(t, err)

	assert.Contains(t, string(content), "name: adr")
	assert.Contains(t, string(content), "@decision.id")
}

func TestInstallUserLevel(t *testing.T) {
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	err := InstallUserLevel()
	require.NoError(t, err)

	skillPath := filepath.Join(tmpDir, ".claude", "skills", "adr.md")
	content, err := os.ReadFile(skillPath)
	require.NoError(t, err)

	assert.Contains(t, string(content), "name: adr")
}

func TestExists_ProjectLevel(t *testing.T) {
	tmpDir := t.TempDir()

	assert.False(t, ExistsProjectLevel(tmpDir))

	err := InstallProjectLevel(tmpDir)
	require.NoError(t, err)

	assert.True(t, ExistsProjectLevel(tmpDir))
}

func TestExists_UserLevel(t *testing.T) {
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	assert.False(t, ExistsUserLevel())

	err := InstallUserLevel()
	require.NoError(t, err)

	assert.True(t, ExistsUserLevel())
}

func TestInstallProjectLevel_Idempotent(t *testing.T) {
	tmpDir := t.TempDir()

	err := InstallProjectLevel(tmpDir)
	require.NoError(t, err)

	skillPath := filepath.Join(tmpDir, ".claude", "skills", "adr.md")
	err = os.WriteFile(skillPath, []byte("custom content"), 0644)
	require.NoError(t, err)

	err = InstallProjectLevel(tmpDir)
	require.NoError(t, err)

	content, err := os.ReadFile(skillPath)
	require.NoError(t, err)
	assert.Equal(t, "custom content", string(content))
}
