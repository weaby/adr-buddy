package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSyncCommand(t *testing.T) {
	tmpDir := t.TempDir()

	createTestADR(t, tmpDir, "adr-001-test.md", `---
adr_id: adr-001
name: Test Decision
status: accepted
date: 2024-01-15
---

## Decision

Test decision.
`)

	var buf bytes.Buffer
	err := SyncWithFormat(tmpDir, false, "text", &buf)
	require.NoError(t, err)

	// Check output
	assert.Contains(t, buf.String(), "Found 1 ADR(s)")
	assert.Contains(t, buf.String(), "Regenerated decisions-index.md")

	// Check index file was created
	indexPath := filepath.Join(tmpDir, ".claude", "rules", "decisions", "decisions-index.md")
	content, err := os.ReadFile(indexPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "adr-001")
	assert.Contains(t, string(content), "Test Decision")
}

func TestSyncCommand_DryRun(t *testing.T) {
	tmpDir := t.TempDir()

	createTestADR(t, tmpDir, "adr-001-test.md", `---
adr_id: adr-001
name: Test Decision
status: accepted
date: 2024-01-15
---

## Decision

Test decision.
`)

	var buf bytes.Buffer
	err := SyncWithFormat(tmpDir, true, "text", &buf)
	require.NoError(t, err)

	assert.Contains(t, buf.String(), "[DRY RUN]")
}

func TestSync_JSONFormat(t *testing.T) {
	tmpDir := t.TempDir()

	createTestADR(t, tmpDir, "adr-001-test.md", `---
adr_id: adr-001
name: Test Decision
status: accepted
date: 2024-01-15
---

## Decision

Test decision.
`)

	var buf bytes.Buffer
	err := SyncWithFormat(tmpDir, false, "json", &buf)
	require.NoError(t, err)

	var result SyncResult
	err = json.Unmarshal(buf.Bytes(), &result)
	require.NoError(t, err)

	assert.True(t, result.IndexUpdated)
	assert.Equal(t, 1, result.ADRCount)
}

func TestSyncCommand_NoADRs(t *testing.T) {
	tmpDir := t.TempDir()

	var buf bytes.Buffer
	err := SyncWithFormat(tmpDir, false, "text", &buf)
	require.NoError(t, err)

	assert.Contains(t, buf.String(), "Found 0 ADR(s)")
}
