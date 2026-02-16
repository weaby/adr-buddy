package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCommand(t *testing.T) {
	tmpDir := t.TempDir()

	err := NewCommand(tmpDir, "Use PostgreSQL", "")
	require.NoError(t, err)

	// Check that an ADR file was created
	decisionsDir := filepath.Join(tmpDir, ".claude", "rules", "decisions")
	entries, err := os.ReadDir(decisionsDir)
	require.NoError(t, err)

	// Should have ADR file + index
	adrFiles := 0
	for _, e := range entries {
		if e.Name() != "decisions-index.md" {
			adrFiles++
		}
	}
	assert.Equal(t, 1, adrFiles)
}

func TestNewCommand_AutoIncrementID(t *testing.T) {
	tmpDir := t.TempDir()

	err := NewCommand(tmpDir, "First Decision", "")
	require.NoError(t, err)

	err = NewCommand(tmpDir, "Second Decision", "")
	require.NoError(t, err)

	// Should have two ADR files
	decisionsDir := filepath.Join(tmpDir, ".claude", "rules", "decisions")
	entries, err := os.ReadDir(decisionsDir)
	require.NoError(t, err)

	adrFiles := 0
	for _, e := range entries {
		if e.Name() != "decisions-index.md" {
			adrFiles++
		}
	}
	assert.Equal(t, 2, adrFiles)
}

func TestNewCommand_WithCategory(t *testing.T) {
	tmpDir := t.TempDir()

	err := NewCommand(tmpDir, "Use Kafka", "infrastructure")
	require.NoError(t, err)

	// Read back and verify category
	decisionsDir := filepath.Join(tmpDir, ".claude", "rules", "decisions")
	entries, err := os.ReadDir(decisionsDir)
	require.NoError(t, err)

	for _, e := range entries {
		if e.Name() == "decisions-index.md" {
			continue
		}
		content, err := os.ReadFile(filepath.Join(decisionsDir, e.Name()))
		require.NoError(t, err)
		assert.Contains(t, string(content), "category: infrastructure")
	}
}
