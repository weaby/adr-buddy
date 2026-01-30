package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/weaby/adr-buddy/internal/config"
)

func TestInit_CreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	err := Init(tmpDir)
	require.NoError(t, err)

	// Check .adr-buddy directory exists
	adrDir := filepath.Join(tmpDir, ".adr-buddy")
	info, err := os.Stat(adrDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestInit_CreatesConfigFile(t *testing.T) {
	tmpDir := t.TempDir()

	err := Init(tmpDir)
	require.NoError(t, err)

	// Check config.yml exists and is valid
	configPath := filepath.Join(tmpDir, ".adr-buddy", "config.yml")
	data, err := os.ReadFile(configPath)
	require.NoError(t, err)

	// Parse YAML to verify it's valid
	var cfg config.Config
	err = yaml.Unmarshal(data, &cfg)
	require.NoError(t, err)

	// Verify default values
	expected := config.Default()
	assert.Equal(t, expected.ScanPaths, cfg.ScanPaths)
	assert.Equal(t, expected.OutputDir, cfg.OutputDir)
	assert.Equal(t, expected.StrictMode, cfg.StrictMode)
	assert.NotEmpty(t, cfg.Exclude)
}

func TestInit_CreatesTemplateFile(t *testing.T) {
	tmpDir := t.TempDir()

	err := Init(tmpDir)
	require.NoError(t, err)

	// Check template.md exists
	templatePath := filepath.Join(tmpDir, ".adr-buddy", "template.md")
	data, err := os.ReadFile(templatePath)
	require.NoError(t, err)

	// Should contain template content
	content := string(data)
	assert.Contains(t, content, "{{.ID}}")
	assert.Contains(t, content, "{{.Name}}")
	assert.Contains(t, content, "{{.Status}}")
	assert.Contains(t, content, "## Context")
	assert.Contains(t, content, "## Decision")
	assert.Contains(t, content, "## Consequences")
}

func TestInit_IdempotentExecution(t *testing.T) {
	tmpDir := t.TempDir()

	// Run init first time
	err := Init(tmpDir)
	require.NoError(t, err)

	// Modify config file
	configPath := filepath.Join(tmpDir, ".adr-buddy", "config.yml")
	err = os.WriteFile(configPath, []byte("scan_paths: [\"custom\"]\n"), 0644)
	require.NoError(t, err)

	// Run init second time
	err = Init(tmpDir)
	require.NoError(t, err)

	// Config should NOT be overwritten
	data, err := os.ReadFile(configPath)
	require.NoError(t, err)
	assert.Contains(t, string(data), "custom")
}

func TestInit_PreservesExistingFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create .adr-buddy directory manually
	adrDir := filepath.Join(tmpDir, ".adr-buddy")
	err := os.MkdirAll(adrDir, 0755)
	require.NoError(t, err)

	// Create custom config
	configPath := filepath.Join(adrDir, "config.yml")
	customConfig := "output_dir: custom\n"
	err = os.WriteFile(configPath, []byte(customConfig), 0644)
	require.NoError(t, err)

	// Create custom template
	templatePath := filepath.Join(adrDir, "template.md")
	customTemplate := "# Custom Template\n"
	err = os.WriteFile(templatePath, []byte(customTemplate), 0644)
	require.NoError(t, err)

	// Run init
	err = Init(tmpDir)
	require.NoError(t, err)

	// Verify files were NOT overwritten
	configData, err := os.ReadFile(configPath)
	require.NoError(t, err)
	assert.Equal(t, customConfig, string(configData))

	templateData, err := os.ReadFile(templatePath)
	require.NoError(t, err)
	assert.Equal(t, customTemplate, string(templateData))
}

func TestInit_InvalidDirectory(t *testing.T) {
	// Try to init in a non-existent directory
	err := Init("/nonexistent/path/that/does/not/exist")
	assert.Error(t, err)
}

func TestInitWithSkill_ProjectLevel(t *testing.T) {
	tmpDir := t.TempDir()

	err := InitWithSkill(tmpDir, SkillLocationProject)
	require.NoError(t, err)

	// Check adr skill
	adrSkillPath := filepath.Join(tmpDir, ".claude", "skills", "adr", "SKILL.md")
	_, err = os.Stat(adrSkillPath)
	require.NoError(t, err)

	// Check adr-review skill
	adrReviewSkillPath := filepath.Join(tmpDir, ".claude", "skills", "adr-review", "SKILL.md")
	_, err = os.Stat(adrReviewSkillPath)
	require.NoError(t, err)
}

func TestInitWithSkill_UserLevel(t *testing.T) {
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	projectDir := filepath.Join(tmpDir, "project")
	err := os.MkdirAll(projectDir, 0755)
	require.NoError(t, err)

	err = InitWithSkill(projectDir, SkillLocationUser)
	require.NoError(t, err)

	// Check adr skill
	adrSkillPath := filepath.Join(tmpDir, ".claude", "skills", "adr", "SKILL.md")
	_, err = os.Stat(adrSkillPath)
	require.NoError(t, err)

	// Check adr-review skill
	adrReviewSkillPath := filepath.Join(tmpDir, ".claude", "skills", "adr-review", "SKILL.md")
	_, err = os.Stat(adrReviewSkillPath)
	require.NoError(t, err)
}

func TestInitWithSkill_Skip(t *testing.T) {
	tmpDir := t.TempDir()

	err := InitWithSkill(tmpDir, SkillLocationSkip)
	require.NoError(t, err)

	// Verify adr skill was NOT created
	adrSkillPath := filepath.Join(tmpDir, ".claude", "skills", "adr", "SKILL.md")
	_, err = os.Stat(adrSkillPath)
	assert.True(t, os.IsNotExist(err))

	// Verify adr-review skill was NOT created
	adrReviewSkillPath := filepath.Join(tmpDir, ".claude", "skills", "adr-review", "SKILL.md")
	_, err = os.Stat(adrReviewSkillPath)
	assert.True(t, os.IsNotExist(err))
}

func TestInit_StillWorksWithoutSkill(t *testing.T) {
	tmpDir := t.TempDir()

	err := Init(tmpDir)
	require.NoError(t, err)

	adrDir := filepath.Join(tmpDir, ".adr-buddy")
	_, err = os.Stat(adrDir)
	require.NoError(t, err)
}
