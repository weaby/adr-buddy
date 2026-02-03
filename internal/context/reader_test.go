package context

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseADR(t *testing.T) {
	content := []byte(`---
adr_id: adr-001
name: Use PostgreSQL for primary database
status: accepted
category: infrastructure
date: 2024-01-15
globs:
  - "internal/database/**"
  - "**/repository*.go"
---

# ADR-001: Use PostgreSQL for primary database

## Context

We need a reliable database for our application that supports complex queries
and transactions.

## Decision

We will use PostgreSQL as our primary database.

## Alternatives Considered

- **MySQL**: Good option but less feature-rich
- **SQLite**: Not suitable for production

## Consequences

**Positive:** Robust, well-supported, excellent query performance
**Negative:** More operational overhead than SQLite
`)

	adr, err := ParseADR(content)
	require.NoError(t, err)

	assert.Equal(t, "adr-001", adr.ID)
	assert.Equal(t, "Use PostgreSQL for primary database", adr.Name)
	assert.Equal(t, "accepted", adr.Status)
	assert.Equal(t, "infrastructure", adr.Category)
	assert.Equal(t, "2024-01-15", adr.Date)
	assert.Equal(t, []string{"internal/database/**", "**/repository*.go"}, adr.Globs)

	assert.Contains(t, adr.Context, "reliable database")
	assert.Contains(t, adr.Decision, "PostgreSQL")
	assert.Contains(t, adr.Alternatives, "MySQL")
	assert.Contains(t, adr.Consequences, "Positive")
}

func TestParseADR_MissingFrontmatter(t *testing.T) {
	content := []byte(`# Just a markdown file

No frontmatter here.
`)

	_, err := ParseADR(content)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must start with YAML frontmatter")
}

func TestParseADR_UnclosedFrontmatter(t *testing.T) {
	content := []byte(`---
adr_id: adr-001
name: Test

Missing closing marker
`)

	_, err := ParseADR(content)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not properly closed")
}

func TestADR_Validate(t *testing.T) {
	tests := []struct {
		name    string
		adr     ADR
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid ADR",
			adr: ADR{
				ID:   "adr-001",
				Name: "Test decision",
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			adr: ADR{
				Name: "Test decision",
			},
			wantErr: true,
			errMsg:  "adr_id",
		},
		{
			name: "missing name",
			adr: ADR{
				ID: "adr-001",
			},
			wantErr: true,
			errMsg:  "name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.adr.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestADR_MatchesPath(t *testing.T) {
	adr := &ADR{
		Globs: []string{"internal/database/*", "*.go"},
	}

	assert.True(t, adr.MatchesPath("internal/database/repo.go"))
	assert.True(t, adr.MatchesPath("main.go"))
	assert.False(t, adr.MatchesPath("internal/api/handler.go"))

	// ADR with no globs matches everything
	adrNoGlobs := &ADR{}
	assert.True(t, adrNoGlobs.MatchesPath("any/path/here.go"))
}

func TestReader_ReadAll(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()
	decisionsDir := filepath.Join(tmpDir, ".claude", "rules", "decisions")
	err := os.MkdirAll(decisionsDir, 0755)
	require.NoError(t, err)

	// Create test ADR file
	adrContent := `---
adr_id: adr-001
name: Test Decision
status: accepted
category: test
date: 2024-01-15
globs: []
---

# Test Decision

## Context

Test context.

## Decision

Test decision.
`
	err = os.WriteFile(filepath.Join(decisionsDir, "adr-001-test.md"), []byte(adrContent), 0644)
	require.NoError(t, err)

	// Create index file (should be skipped)
	err = os.WriteFile(filepath.Join(decisionsDir, "decisions-index.md"), []byte("# Index"), 0644)
	require.NoError(t, err)

	// Read all ADRs
	reader := NewReader(tmpDir)
	adrs, err := reader.ReadAll()
	require.NoError(t, err)

	assert.Len(t, adrs, 1)
	assert.Equal(t, "adr-001", adrs[0].ID)
	assert.Equal(t, "Test Decision", adrs[0].Name)
}

func TestReader_ReadAll_NoDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	reader := NewReader(tmpDir)
	adrs, err := reader.ReadAll()
	require.NoError(t, err)
	assert.Nil(t, adrs)
}
