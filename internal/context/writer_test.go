package context

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderADR(t *testing.T) {
	adr := &ADR{
		ID:           "adr-001",
		Name:         "Use PostgreSQL",
		Status:       "accepted",
		Category:     "infrastructure",
		Date:         "2024-01-15",
		Globs:        []string{"internal/database/**"},
		Context:      "We need a database.",
		Decision:     "Use PostgreSQL.",
		Alternatives: "MySQL was considered.",
		Consequences: "Good performance.",
	}

	content, err := RenderADR(adr, DefaultTemplate)
	require.NoError(t, err)

	assert.Contains(t, content, "adr_id: adr-001")
	assert.Contains(t, content, "name: Use PostgreSQL")
	assert.Contains(t, content, "status: accepted")
	assert.Contains(t, content, "category: infrastructure")
	assert.Contains(t, content, "date: \"2024-01-15\"")
	assert.Contains(t, content, "# ADR-001: Use PostgreSQL")
	assert.Contains(t, content, "## Context")
	assert.Contains(t, content, "We need a database.")
	assert.Contains(t, content, "## Decision")
	assert.Contains(t, content, "Use PostgreSQL.")
}

func TestRenderADR_DefaultValues(t *testing.T) {
	adr := &ADR{
		ID:   "adr-001",
		Name: "Test Decision",
	}

	content, err := RenderADR(adr, DefaultTemplate)
	require.NoError(t, err)

	assert.Contains(t, content, "status: proposed")
	assert.Regexp(t, `date: "\d{4}-\d{2}-\d{2}"`, content)
	assert.Contains(t, content, "[Why this decision was needed]")
	assert.Contains(t, content, "[What was decided]")
}

func TestRenderADR_CustomTemplate(t *testing.T) {
	adr := &ADR{
		ID:       "adr-001",
		Name:     "Test Decision",
		Status:   "accepted",
		Date:     "2024-01-15",
		Decision: "We decided X.",
	}

	customTmpl := `# {{ .ID | upper }}: {{ .Name }}

## Decision

{{ .Decision | default "[TBD]" }}
`

	content, err := RenderADR(adr, customTmpl)
	require.NoError(t, err)

	assert.Contains(t, content, "# ADR-001: Test Decision")
	assert.Contains(t, content, "We decided X.")
	assert.NotContains(t, content, "## Context")
	assert.NotContains(t, content, "## Alternatives")
}

func TestGenerateFilename(t *testing.T) {
	tests := []struct {
		id       string
		name     string
		expected string
	}{
		{"adr-001", "Use PostgreSQL", "adr-001-use-postgresql.md"},
		{"adr-002", "API Versioning Strategy", "adr-002-api-versioning-strategy.md"},
		{"adr-003", "Use Go 1.21+", "adr-003-use-go-121.md"},
		{"adr-004", "Test!@#Special$%^Characters", "adr-004-testspecialcharacters.md"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateFilename(tt.id, tt.name)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWriter_WriteADR(t *testing.T) {
	tmpDir := t.TempDir()
	writer := NewWriter(tmpDir)

	adr := &ADR{
		ID:       "adr-001",
		Name:     "Test Decision",
		Status:   "accepted",
		Date:     "2024-01-15",
		Context:  "Test context",
		Decision: "Test decision",
	}

	err := writer.WriteADR(adr)
	require.NoError(t, err)

	// Verify file was created
	expectedPath := filepath.Join(tmpDir, ".claude", "rules", "decisions", "adr-001-test-decision.md")
	_, err = os.Stat(expectedPath)
	require.NoError(t, err)

	// Verify content
	content, err := os.ReadFile(expectedPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "adr_id: adr-001")
	assert.Contains(t, string(content), "Test context")
}

func TestWriter_WriteADR_CustomTemplate(t *testing.T) {
	tmpDir := t.TempDir()

	tmplDir := filepath.Join(tmpDir, ".adr-buddy")
	require.NoError(t, os.MkdirAll(tmplDir, 0755))

	customTmpl := `# {{ .ID | upper }}: {{ .Name }}

## Decision

{{ .Decision | default "[TBD]" }}
`
	require.NoError(t, os.WriteFile(filepath.Join(tmplDir, "template.md"), []byte(customTmpl), 0644))

	writer := NewWriter(tmpDir)
	adr := &ADR{
		ID:       "adr-001",
		Name:     "Custom Template Test",
		Status:   "accepted",
		Date:     "2024-01-15",
		Decision: "We use custom templates.",
	}

	err := writer.WriteADR(adr)
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(tmpDir, ".claude", "rules", "decisions", "adr-001-custom-template-test.md"))
	require.NoError(t, err)

	assert.Contains(t, string(content), "We use custom templates.")
	assert.NotContains(t, string(content), "## Context")
}

func TestNewADR(t *testing.T) {
	adr := NewADR("adr-001", "Test Decision")

	assert.Equal(t, "adr-001", adr.ID)
	assert.Equal(t, "Test Decision", adr.Name)
	assert.Equal(t, "proposed", adr.Status)
	assert.NotEmpty(t, adr.Date)
}

func TestNextADRID(t *testing.T) {
	tests := []struct {
		name     string
		adrs     []*ADR
		expected string
	}{
		{
			name:     "no existing ADRs",
			adrs:     nil,
			expected: "adr-001",
		},
		{
			name: "sequential ADRs",
			adrs: []*ADR{
				{ID: "adr-001"},
				{ID: "adr-002"},
				{ID: "adr-003"},
			},
			expected: "adr-004",
		},
		{
			name: "gap in sequence",
			adrs: []*ADR{
				{ID: "adr-001"},
				{ID: "adr-005"},
			},
			expected: "adr-006",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NextADRID(tt.adrs)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWriter_EnsureDecisionsDir(t *testing.T) {
	tmpDir := t.TempDir()
	writer := NewWriter(tmpDir)

	err := writer.EnsureDecisionsDir()
	require.NoError(t, err)

	// Verify directory was created
	decisionsDir := writer.DecisionsDir()
	info, err := os.Stat(decisionsDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())

	// Calling again should not error
	err = writer.EnsureDecisionsDir()
	require.NoError(t, err)
}

func TestFilenameLength(t *testing.T) {
	longName := strings.Repeat("a", 100)
	filename := generateFilename("adr-001", longName)

	// Should be truncated but still valid
	assert.True(t, len(filename) < 70)
	assert.True(t, strings.HasPrefix(filename, "adr-001-"))
	assert.True(t, strings.HasSuffix(filename, ".md"))
}
