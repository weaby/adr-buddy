package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultTemplate(t *testing.T) {
	tmpl := DefaultTemplate()

	assert.NotEmpty(t, tmpl)
	assert.Contains(t, tmpl, "{{.ID}}")
	assert.Contains(t, tmpl, "{{.Name}}")
	assert.Contains(t, tmpl, "## Context")
	assert.Contains(t, tmpl, "## Decision")
	assert.Contains(t, tmpl, "## Consequences")
	assert.Contains(t, tmpl, "## Code Locations")
}
