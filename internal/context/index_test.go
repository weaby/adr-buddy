package context

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderIndex(t *testing.T) {
	adrs := []*ADR{
		{
			ID:       "adr-001",
			Name:     "Use PostgreSQL",
			Status:   "accepted",
			Category: "infrastructure",
			Decision: "We will use PostgreSQL as our primary database.",
			FilePath: ".claude/rules/decisions/adr-001-use-postgresql.md",
			Globs:    []string{"internal/database/**"},
		},
		{
			ID:       "adr-002",
			Name:     "API Versioning",
			Status:   "proposed",
			Category: "api",
			Decision: "We will use URL path versioning for our API.",
			FilePath: ".claude/rules/decisions/adr-002-api-versioning.md",
		},
	}

	content, err := RenderIndex(adrs)
	require.NoError(t, err)

	// Check title
	assert.Contains(t, content, "# Architecture Decision Records")

	// Check quick reference table
	assert.Contains(t, content, "| adr-001 | [Use PostgreSQL]")
	assert.Contains(t, content, "| adr-002 | [API Versioning]")

	// Check category sections
	assert.Contains(t, content, "## Infrastructure")
	assert.Contains(t, content, "## Api")

	// Check ADR details
	assert.Contains(t, content, "### adr-001: Use PostgreSQL")
	assert.Contains(t, content, "**Status:** accepted")
	assert.Contains(t, content, "**Applies to:** internal/database/**")
}

func TestRenderIndex_EmptyList(t *testing.T) {
	content, err := RenderIndex(nil)
	require.NoError(t, err)

	assert.Contains(t, content, "# Architecture Decision Records")
	// Should have empty table headers but no rows
	assert.Contains(t, content, "| ID | Decision | Status | Category |")
}

func TestRenderIndex_NoCategory(t *testing.T) {
	adrs := []*ADR{
		{
			ID:       "adr-001",
			Name:     "Test Decision",
			Status:   "accepted",
			Decision: "Test decision content.",
			FilePath: "adr-001.md",
		},
	}

	content, err := RenderIndex(adrs)
	require.NoError(t, err)

	// Should use "general" as default category
	assert.Contains(t, content, "| general |")
	assert.Contains(t, content, "## General")
}

func TestIndexGenerator_Generate(t *testing.T) {
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

Test decision content here.
`
	err = os.WriteFile(filepath.Join(decisionsDir, "adr-001-test.md"), []byte(adrContent), 0644)
	require.NoError(t, err)

	// Generate index
	generator := NewIndexGenerator(tmpDir)
	err = generator.Generate()
	require.NoError(t, err)

	// Verify index file was created
	indexPath := generator.IndexPath()
	content, err := os.ReadFile(indexPath)
	require.NoError(t, err)

	assert.Contains(t, string(content), "# Architecture Decision Records")
	assert.Contains(t, string(content), "adr-001")
	assert.Contains(t, string(content), "Test Decision")
}

func TestIndexGenerator_Generate_NoADRs(t *testing.T) {
	tmpDir := t.TempDir()

	generator := NewIndexGenerator(tmpDir)
	err := generator.Generate()
	require.NoError(t, err)

	// Index should be created even with no ADRs
	indexPath := generator.IndexPath()
	content, err := os.ReadFile(indexPath)
	require.NoError(t, err)

	assert.Contains(t, string(content), "# Architecture Decision Records")
}

func TestRenderIndex_LongDecision(t *testing.T) {
	longDecision := "This is a very long decision text that should be truncated in the index because we only want to show a preview. " +
		"It continues with more text that explains the decision in great detail. " +
		"And even more detail about why this decision was made and what the implications are."

	adrs := []*ADR{
		{
			ID:       "adr-001",
			Name:     "Test",
			Status:   "accepted",
			Decision: longDecision,
			FilePath: "adr-001.md",
		},
	}

	content, err := RenderIndex(adrs)
	require.NoError(t, err)

	// Decision should be truncated with ...
	assert.Contains(t, content, "...")
	// Full long text should not appear
	assert.NotContains(t, content, "what the implications are")
}

func TestRenderIndex_SortedByCategory(t *testing.T) {
	adrs := []*ADR{
		{ID: "adr-001", Name: "Z First", Category: "zebra", FilePath: "a.md"},
		{ID: "adr-002", Name: "A Second", Category: "alpha", FilePath: "b.md"},
		{ID: "adr-003", Name: "M Third", Category: "middle", FilePath: "c.md"},
	}

	content, err := RenderIndex(adrs)
	require.NoError(t, err)

	// Categories should appear in alphabetical order
	alphaIdx := indexOf(content, "## Alpha")
	middleIdx := indexOf(content, "## Middle")
	zebraIdx := indexOf(content, "## Zebra")

	assert.True(t, alphaIdx < middleIdx, "Alpha should come before Middle")
	assert.True(t, middleIdx < zebraIdx, "Middle should come before Zebra")
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
