//go:build integration

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

func TestIntegration_FullWorkflow(t *testing.T) {
	tmpDir := t.TempDir()

	// Step 1: Init
	err := Init(tmpDir)
	require.NoError(t, err)

	// Decisions directory should exist
	decisionsDir := filepath.Join(tmpDir, ".claude", "rules", "decisions")
	_, err = os.Stat(decisionsDir)
	require.NoError(t, err)

	// Step 2: Create a new ADR
	err = NewCommand(tmpDir, "Use PostgreSQL", "infrastructure")
	require.NoError(t, err)

	// Step 3: Check passes
	var checkBuf bytes.Buffer
	err = CheckWithFormat(tmpDir, false, "json", &checkBuf)
	require.NoError(t, err)

	var checkResult CheckResult
	err = json.Unmarshal(checkBuf.Bytes(), &checkResult)
	require.NoError(t, err)
	assert.Equal(t, "pass", checkResult.Status)
	assert.Equal(t, 1, checkResult.Summary.TotalADRs)

	// Step 4: List shows the ADR
	var listBuf bytes.Buffer
	err = ListCommand(tmpDir, "", &listBuf)
	require.NoError(t, err)
	assert.Contains(t, listBuf.String(), "Use PostgreSQL")

	// Step 5: Sync regenerates index
	var syncBuf bytes.Buffer
	err = SyncWithFormat(tmpDir, false, "json", &syncBuf)
	require.NoError(t, err)

	var syncResult SyncResult
	err = json.Unmarshal(syncBuf.Bytes(), &syncResult)
	require.NoError(t, err)
	assert.True(t, syncResult.IndexUpdated)
	assert.Equal(t, 1, syncResult.ADRCount)

	// Verify index contains the ADR
	indexContent, err := os.ReadFile(filepath.Join(decisionsDir, "decisions-index.md"))
	require.NoError(t, err)
	assert.Contains(t, string(indexContent), "Use PostgreSQL")
}
