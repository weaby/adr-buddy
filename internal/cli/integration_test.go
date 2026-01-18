//go:build integration
// +build integration

package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaby/adr-buddy/internal/model"
)

func TestIntegration_ValidateAndSyncWorkflow(t *testing.T) {
	tmpDir := t.TempDir()

	// Initialize project
	err := Init(tmpDir)
	require.NoError(t, err)

	// Create source file with annotations
	srcDir := filepath.Join(tmpDir, "src")
	require.NoError(t, os.MkdirAll(srcDir, 0755))

	sourceFile := filepath.Join(srcDir, "logger.js")
	sourceContent := `// @decision.id: adr-1
// @decision.name: Using Pino for logging
// @decision.status: accepted
// @decision.context: We needed structured logging
// @decision.decision: Use Pino as logging library
// @decision.consequences: Better performance
const pino = require('pino');
`
	require.NoError(t, os.WriteFile(sourceFile, []byte(sourceContent), 0644))

	// Update config to scan src directory
	configFile := filepath.Join(tmpDir, ".adr-buddy", "config.yml")
	configContent := `scan_paths:
  - ./src
output_dir: ./decisions
exclude:
  - "**/node_modules/**"
`
	require.NoError(t, os.WriteFile(configFile, []byte(configContent), 0644))

	// Run check with JSON output
	var checkBuf bytes.Buffer
	err = CheckWithFormat(tmpDir, false, "json", &checkBuf)
	require.NoError(t, err)

	var checkResult model.CheckResult
	err = json.Unmarshal(checkBuf.Bytes(), &checkResult)
	require.NoError(t, err)
	assert.Equal(t, "pass", checkResult.Status)
	assert.Equal(t, 1, checkResult.Summary.TotalAnnotations)
	assert.Equal(t, 1, checkResult.Summary.ValidAnnotations)

	// Run sync dry-run with JSON output
	var syncBuf bytes.Buffer
	err = SyncWithFormat(tmpDir, true, "json", &syncBuf)
	require.NoError(t, err)

	var syncResult model.SyncResult
	err = json.Unmarshal(syncBuf.Bytes(), &syncResult)
	require.NoError(t, err)
	assert.True(t, syncResult.ChangesDetected)
	assert.Equal(t, 1, len(syncResult.Files.Created))
	assert.Contains(t, syncResult.Files.Created[0], "adr-1.md")

	// Run actual sync
	var actualSyncBuf bytes.Buffer
	err = SyncWithFormat(tmpDir, false, "text", &actualSyncBuf)
	require.NoError(t, err)

	// Verify ADR file was created
	adrFile := filepath.Join(tmpDir, "decisions", "adr-1.md")
	assert.FileExists(t, adrFile)

	content, err := os.ReadFile(adrFile)
	require.NoError(t, err)
	assert.Contains(t, string(content), "Using Pino for logging")
	assert.Contains(t, string(content), "**Status:** accepted")
}

func TestIntegration_CheckValidationErrors(t *testing.T) {
	tmpDir := t.TempDir()

	// Initialize
	err := Init(tmpDir)
	require.NoError(t, err)

	// Create source with invalid annotation (missing name)
	sourceFile := filepath.Join(tmpDir, "test.js")
	sourceContent := `// @decision.id: adr-1
const x = 1;`
	require.NoError(t, os.WriteFile(sourceFile, []byte(sourceContent), 0644))

	// Run check with JSON
	var buf bytes.Buffer
	err = CheckWithFormat(tmpDir, false, "json", &buf)
	require.Error(t, err) // Should fail validation

	var result model.CheckResult
	err = json.Unmarshal(buf.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, "fail", result.Status)
	assert.Greater(t, result.Summary.ErrorCount, 0)
	assert.Contains(t, result.Errors[0].Message, "@decision.name")
}
