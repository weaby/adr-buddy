package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestADR(t *testing.T, dir, filename, content string) {
	t.Helper()
	decisionsDir := filepath.Join(dir, ".claude", "rules", "decisions")
	err := os.MkdirAll(decisionsDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(decisionsDir, filename), []byte(content), 0644)
	require.NoError(t, err)
}

func TestListCommand(t *testing.T) {
	tmpDir := t.TempDir()

	createTestADR(t, tmpDir, "adr-001-postgres.md", `---
adr_id: adr-001
name: Use PostgreSQL
status: accepted
category: infrastructure
date: 2024-01-15
---

## Context

We need a database.

## Decision

Use PostgreSQL.
`)

	createTestADR(t, tmpDir, "adr-002-api.md", `---
adr_id: adr-002
name: REST API Design
status: proposed
category: api
date: 2024-01-16
---

## Decision

Use REST.
`)

	var buf bytes.Buffer
	err := ListCommand(tmpDir, "", &buf)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "adr-001")
	assert.Contains(t, output, "Use PostgreSQL")
	assert.Contains(t, output, "accepted")
	assert.Contains(t, output, "infrastructure")
	assert.Contains(t, output, "adr-002")
	assert.Contains(t, output, "REST API Design")
}

func TestListCommand_CategoryFilter(t *testing.T) {
	tmpDir := t.TempDir()

	createTestADR(t, tmpDir, "adr-001-postgres.md", `---
adr_id: adr-001
name: Use PostgreSQL
status: accepted
category: infrastructure
date: 2024-01-15
---

## Decision

Use PostgreSQL.
`)

	createTestADR(t, tmpDir, "adr-002-api.md", `---
adr_id: adr-002
name: REST API Design
status: proposed
category: api
date: 2024-01-16
---

## Decision

Use REST.
`)

	var buf bytes.Buffer
	err := ListCommand(tmpDir, "infrastructure", &buf)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "adr-001")
	assert.Contains(t, output, "Use PostgreSQL")
	assert.NotContains(t, output, "adr-002")
}

func TestListCommand_NoADRs(t *testing.T) {
	tmpDir := t.TempDir()

	var buf bytes.Buffer
	err := ListCommand(tmpDir, "", &buf)
	require.NoError(t, err)

	assert.Contains(t, buf.String(), "No ADRs found")
}
