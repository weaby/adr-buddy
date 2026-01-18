package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weaby/adr-buddy/internal/model"
)

func TestRender(t *testing.T) {
	adr := &model.ADR{
		ID:           "adr-1",
		Name:         "Using Pino for logging",
		Status:       "accepted",
		Date:         "2026-01-17",
		Category:     "infrastructure",
		Context:      []string{"Need fast logging"},
		Decision:     []string{"Use Pino"},
		Consequences: []string{"Better performance"},
		Locations: []model.SourceLocation{
			{File: "src/logger.js", Line: 10},
		},
	}

	result, err := Render(adr, DefaultTemplate())

	assert.NoError(t, err)
	assert.Contains(t, result, "# adr-1: Using Pino for logging")
	assert.Contains(t, result, "**Status:** accepted")
	assert.Contains(t, result, "**Date:** 2026-01-17")
	assert.Contains(t, result, "**Category:** infrastructure")
	assert.Contains(t, result, "Need fast logging")
	assert.Contains(t, result, "Use Pino")
	assert.Contains(t, result, "Better performance")
	assert.Contains(t, result, "- src/logger.js:10")
}

func TestRender_WithPlaceholders(t *testing.T) {
	adr := &model.ADR{
		ID:     "adr-2",
		Name:   "Minimal ADR",
		Status: "proposed",
		Date:   "2026-01-17",
		Locations: []model.SourceLocation{
			{File: "src/app.js", Line: 1},
		},
	}

	result, err := Render(adr, DefaultTemplate())

	assert.NoError(t, err)
	assert.Contains(t, result, "# adr-2: Minimal ADR")
	assert.Contains(t, result, "<!-- TODO: Add context")
	assert.Contains(t, result, "<!-- TODO: Document the decision")
	assert.Contains(t, result, "<!-- TODO: What are the positive/negative outcomes")
}
