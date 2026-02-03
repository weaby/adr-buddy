package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit_CreatesDecisionsDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	err := Init(tmpDir)
	require.NoError(t, err)

	decisionsDir := filepath.Join(tmpDir, ".claude", "rules", "decisions")
	info, err := os.Stat(decisionsDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestInit_CreatesIndexFile(t *testing.T) {
	tmpDir := t.TempDir()
	err := Init(tmpDir)
	require.NoError(t, err)

	indexPath := filepath.Join(tmpDir, ".claude", "rules", "decisions", "decisions-index.md")
	content, err := os.ReadFile(indexPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "# Architecture Decision Records")
}

func TestInit_IdempotentExecution(t *testing.T) {
	tmpDir := t.TempDir()

	err := Init(tmpDir)
	require.NoError(t, err)

	// Running again should not error
	err = Init(tmpDir)
	require.NoError(t, err)
}

func TestInit_InvalidDirectory(t *testing.T) {
	err := Init("/nonexistent/path")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

func TestInitWithSkill_ProjectLevel(t *testing.T) {
	tmpDir := t.TempDir()
	err := InitWithSkill(tmpDir, SkillLocationProject)
	require.NoError(t, err)

	adrSkill := filepath.Join(tmpDir, ".claude", "skills", "adr", "SKILL.md")
	_, err = os.Stat(adrSkill)
	assert.NoError(t, err)

	reviewSkill := filepath.Join(tmpDir, ".claude", "skills", "adr-review", "SKILL.md")
	_, err = os.Stat(reviewSkill)
	assert.NoError(t, err)
}

func TestInitWithSkill_UserLevel(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	projectDir := t.TempDir()
	err := InitWithSkill(projectDir, SkillLocationUser)
	require.NoError(t, err)

	adrSkill := filepath.Join(tmpDir, ".claude", "skills", "adr", "SKILL.md")
	_, err = os.Stat(adrSkill)
	assert.NoError(t, err)

	reviewSkill := filepath.Join(tmpDir, ".claude", "skills", "adr-review", "SKILL.md")
	_, err = os.Stat(reviewSkill)
	assert.NoError(t, err)
}

func TestInitWithSkill_Skip(t *testing.T) {
	tmpDir := t.TempDir()
	err := InitWithSkill(tmpDir, SkillLocationSkip)
	require.NoError(t, err)

	adrSkill := filepath.Join(tmpDir, ".claude", "skills", "adr", "SKILL.md")
	_, err = os.Stat(adrSkill)
	assert.True(t, os.IsNotExist(err))
}
