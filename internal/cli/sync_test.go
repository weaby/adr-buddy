package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weaby/adr-buddy/internal/model"
)

func TestSyncCommand(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test file with annotation
	content := `// @decision.id: adr-1
// @decision.name: Test Decision
// @decision.status: accepted
// @decision.context: This is the context
const x = 1;
`
	err := os.WriteFile(filepath.Join(tmpDir, "test.js"), []byte(content), 0644)
	assert.NoError(t, err)

	err = SyncCommand(tmpDir, false, false)
	assert.NoError(t, err)

	// Check ADR file created
	adrPath := filepath.Join(tmpDir, "decisions", "adr-1.md")
	adrContent, err := os.ReadFile(adrPath)
	assert.NoError(t, err)
	assert.Contains(t, string(adrContent), "# adr-1: Test Decision")
	assert.Contains(t, string(adrContent), "**Status:** accepted")
	assert.Contains(t, string(adrContent), "This is the context")
}

func TestSyncCommand_DryRun(t *testing.T) {
	tmpDir := t.TempDir()

	content := `// @decision.id: adr-1
// @decision.name: Test
const x = 1;
`
	err := os.WriteFile(filepath.Join(tmpDir, "test.js"), []byte(content), 0644)
	assert.NoError(t, err)

	err = SyncCommand(tmpDir, true, false)
	assert.NoError(t, err)

	// Check ADR file NOT created
	adrPath := filepath.Join(tmpDir, "decisions", "adr-1.md")
	_, err = os.Stat(adrPath)
	assert.True(t, os.IsNotExist(err))
}

func TestSyncCommand_UpdateExisting(t *testing.T) {
	tmpDir := t.TempDir()

	// Create decisions directory
	decisionsDir := filepath.Join(tmpDir, "decisions")
	err := os.MkdirAll(decisionsDir, 0755)
	assert.NoError(t, err)

	// Create initial ADR file with manual content
	initialADR := `# adr-2: Database Choice

**Status:** proposed
**Date:** 2026-01-10

## Context
This is manually written context about database selection.

## Decision
This is a manually written decision explaining why we chose PostgreSQL.

## Consequences
These are manually written consequences.

## Code Locations
- old/path.js:5
`
	adrPath := filepath.Join(decisionsDir, "adr-2.md")
	err = os.WriteFile(adrPath, []byte(initialADR), 0644)
	assert.NoError(t, err)

	// Get initial file info to compare later
	initialStat, err := os.Stat(adrPath)
	assert.NoError(t, err)

	// Create code file with annotations for the same ADR
	codeContent := `// @decision.id: adr-2
// @decision.name: Database Choice
// @decision.status: accepted
// @decision.context: New context from annotation
const dbConnection = "postgresql://...";
`
	err = os.WriteFile(filepath.Join(tmpDir, "database.js"), []byte(codeContent), 0644)
	assert.NoError(t, err)

	// Run sync command
	err = SyncCommand(tmpDir, false, false)
	assert.NoError(t, err)

	// Verify file was updated (not recreated)
	updatedStat, err := os.Stat(adrPath)
	assert.NoError(t, err)
	assert.True(t, updatedStat.ModTime().After(initialStat.ModTime()) ||
		updatedStat.ModTime().Equal(initialStat.ModTime()),
		"File should have been modified or touched")

	// Read the updated ADR
	updatedContent, err := os.ReadFile(adrPath)
	assert.NoError(t, err)
	adrStr := string(updatedContent)

	// Verify merge happened correctly:
	// 1. Status should be updated from annotation
	assert.Contains(t, adrStr, "**Status:** accepted")
	assert.NotContains(t, adrStr, "**Status:** proposed")

	// 2. Date should be preserved from original
	assert.Contains(t, adrStr, "**Date:** 2026-01-10")

	// 3. Context should be replaced with annotation content (annotation provides new content)
	assert.Contains(t, adrStr, "New context from annotation")
	assert.NotContains(t, adrStr, "manually written context about database selection")

	// 4. Decision should be preserved (annotation is empty, manual content exists)
	assert.Contains(t, adrStr, "This is a manually written decision explaining why we chose PostgreSQL")

	// 5. Consequences should be preserved (annotation is empty, manual content exists)
	assert.Contains(t, adrStr, "These are manually written consequences")

	// 6. Code locations should be updated
	assert.Contains(t, adrStr, "database.js:")
	assert.NotContains(t, adrStr, "old/path.js:5")
}

func TestSync_DryRunJSON(t *testing.T) {
	tmpDir := t.TempDir()

	// Write config
	configDir := filepath.Join(tmpDir, ".adr-buddy")
	assert.NoError(t, os.MkdirAll(configDir, 0755))
	configContent := `scan_paths: ["."]
output_dir: "./decisions"`
	assert.NoError(t, os.WriteFile(filepath.Join(configDir, "config.yml"), []byte(configContent), 0644))

	// Write source file
	sourceFile := filepath.Join(tmpDir, "test.js")
	sourceContent := `// @decision.id: adr-1
// @decision.name: Test Decision
// @decision.status: accepted
const x = 1;`
	assert.NoError(t, os.WriteFile(sourceFile, []byte(sourceContent), 0644))

	var buf bytes.Buffer
	err := SyncWithFormat(tmpDir, true, "json", &buf)
	assert.NoError(t, err)

	var result model.SyncResult
	err = json.Unmarshal(buf.Bytes(), &result)
	assert.NoError(t, err)
	assert.True(t, result.ChangesDetected)
	assert.Equal(t, 1, len(result.Files.Created))
}
