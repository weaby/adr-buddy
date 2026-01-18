package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weaby/adr-buddy/internal/model"
)

func TestParseExistingADR(t *testing.T) {
	existing := `# adr-1: Test Decision

**Status:** accepted
**Date:** 2026-01-15

## Context
Existing context that was manually added.

## Decision
<!-- TODO: Document the decision and rationale -->

## Consequences
Some manual consequences.

## Code Locations
- old/path.js:10
`

	parsed := ParseExistingADR(existing)

	assert.Equal(t, "accepted", parsed.Frontmatter["Status"])
	assert.Equal(t, "2026-01-15", parsed.Frontmatter["Date"])
	assert.Contains(t, parsed.Sections["Context"], "Existing context")
	assert.Contains(t, parsed.Sections["Decision"], "TODO")
	assert.Contains(t, parsed.Sections["Consequences"], "manual consequences")
}

func TestMerge(t *testing.T) {
	adr := &model.ADR{
		ID:           "adr-1",
		Name:         "Test",
		Status:       "accepted",
		Date:         "2026-01-17",
		Context:      []string{"New context from annotation"},
		Decision:     []string{},
		Consequences: []string{},
		Locations: []model.SourceLocation{
			{File: "new/path.js", Line: 20},
		},
	}

	existingContent := `# adr-1: Test

**Status:** proposed
**Date:** 2026-01-15

## Context
Old manual context.

## Decision
Old manual decision.

## Consequences
<!-- TODO: What are the positive/negative outcomes? -->

## Code Locations
- old/path.js:10
`

	result, err := Merge(adr, existingContent, DefaultTemplate())

	assert.NoError(t, err)
	// Status should be updated
	assert.Contains(t, result, "**Status:** accepted")
	// Date should be preserved
	assert.Contains(t, result, "**Date:** 2026-01-15")
	// Context should be replaced (annotation provides new content)
	assert.Contains(t, result, "New context from annotation")
	assert.NotContains(t, result, "Old manual context")
	// Decision should be preserved (annotation is empty, section has manual content)
	assert.Contains(t, result, "Old manual decision")
	// Consequences should show placeholder (annotation empty, section is placeholder)
	assert.Contains(t, result, "TODO")
	// Locations should be updated
	assert.Contains(t, result, "new/path.js:20")
	assert.NotContains(t, result, "old/path.js:10")
}

func TestMerge_WithCategory(t *testing.T) {
	adr := &model.ADR{
		ID:           "adr-2",
		Name:         "With Category",
		Status:       "accepted",
		Category:     "backend",
		Date:         "2026-01-17",
		Context:      []string{"Context from annotation"},
		Decision:     []string{},
		Consequences: []string{},
		Locations: []model.SourceLocation{
			{File: "backend/service.go", Line: 10},
		},
	}

	existingContent := `# adr-2: With Category

**Status:** proposed
**Date:** 2026-01-16
**Category:** backend

## Context
<!-- TODO: Add context - what is the issue we're facing? -->

## Decision
Manual decision added.

## Consequences
<!-- TODO: What are the positive/negative outcomes? -->

## Code Locations
- old/file.go:5
`

	result, err := Merge(adr, existingContent, DefaultTemplate())

	assert.NoError(t, err)
	assert.Contains(t, result, "**Status:** accepted")
	assert.Contains(t, result, "**Date:** 2026-01-16")
	assert.Contains(t, result, "**Category:** backend")
	assert.Contains(t, result, "Context from annotation")
	assert.Contains(t, result, "Manual decision added")
	assert.Contains(t, result, "backend/service.go:10")
}

func TestMerge_AllPlaceholders(t *testing.T) {
	adr := &model.ADR{
		ID:           "adr-3",
		Name:         "New Decision",
		Status:       "proposed",
		Date:         "2026-01-17",
		Context:      []string{},
		Decision:     []string{},
		Consequences: []string{},
		Locations: []model.SourceLocation{
			{File: "main.go", Line: 1},
		},
	}

	existingContent := `# adr-3: New Decision

**Status:** proposed
**Date:** 2026-01-17

## Context
<!-- TODO: Add context - what is the issue we're facing? -->

## Decision
<!-- TODO: Document the decision and rationale -->

## Consequences
<!-- TODO: What are the positive/negative outcomes? -->

## Code Locations
- main.go:1
`

	result, err := Merge(adr, existingContent, DefaultTemplate())

	assert.NoError(t, err)
	// All sections should still have placeholders
	assert.Contains(t, result, "<!-- TODO: Add context")
	assert.Contains(t, result, "<!-- TODO: Document the decision")
	assert.Contains(t, result, "<!-- TODO: What are the positive")
}

func TestMerge_AllManualContent(t *testing.T) {
	adr := &model.ADR{
		ID:           "adr-4",
		Name:         "Manual Decision",
		Status:       "accepted",
		Date:         "2026-01-17",
		Context:      []string{},
		Decision:     []string{},
		Consequences: []string{},
		Locations: []model.SourceLocation{
			{File: "app.js", Line: 50},
		},
	}

	existingContent := `# adr-4: Manual Decision

**Status:** proposed
**Date:** 2026-01-10

## Context
This is manually written context.

## Decision
This is a manually written decision.

## Consequences
These are manually written consequences.

## Code Locations
- old.js:10
`

	result, err := Merge(adr, existingContent, DefaultTemplate())

	assert.NoError(t, err)
	// Status should update
	assert.Contains(t, result, "**Status:** accepted")
	// Date should be preserved
	assert.Contains(t, result, "**Date:** 2026-01-10")
	// All manual content should be preserved
	assert.Contains(t, result, "This is manually written context")
	assert.Contains(t, result, "This is a manually written decision")
	assert.Contains(t, result, "These are manually written consequences")
	// Locations should update
	assert.Contains(t, result, "app.js:50")
	assert.NotContains(t, result, "old.js:10")
}

func TestIsPlaceholder(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "TODO placeholder",
			content:  "<!-- TODO: Add content -->",
			expected: true,
		},
		{
			name:     "TODO with whitespace",
			content:  "\n\n  <!-- TODO: Something -->  \n",
			expected: true,
		},
		{
			name:     "Manual content",
			content:  "This is manual content",
			expected: false,
		},
		{
			name:     "Empty string",
			content:  "",
			expected: false,
		},
		{
			name:     "Comment but not TODO",
			content:  "<!-- This is a comment -->",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isPlaceholder(tt.content)
			assert.Equal(t, tt.expected, result)
		})
	}
}
