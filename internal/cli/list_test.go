package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListCommand(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	content1 := `// @decision.id: adr-1
// @decision.name: First Decision
// @decision.status: accepted
`
	content2 := `// @decision.id: adr-2
// @decision.name: Second Decision
// @decision.category: backend
`

	err := os.WriteFile(filepath.Join(tmpDir, "test1.js"), []byte(content1), 0644)
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(tmpDir, "test2.js"), []byte(content2), 0644)
	assert.NoError(t, err)

	var buf bytes.Buffer
	err = ListCommand(tmpDir, "", &buf)
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "adr-1")
	assert.Contains(t, output, "First Decision")
	assert.Contains(t, output, "accepted")
	assert.Contains(t, output, "adr-2")
	assert.Contains(t, output, "Second Decision")
	assert.Contains(t, output, "backend")
}

func TestListCommand_CategoryFilter(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files with different categories
	contentBackend := `// @decision.id: adr-backend-1
// @decision.name: Backend Decision
// @decision.status: accepted
// @decision.category: backend
`
	contentFrontend := `// @decision.id: adr-frontend-1
// @decision.name: Frontend Decision
// @decision.status: proposed
// @decision.category: frontend
`
	contentNoCategory := `// @decision.id: adr-root-1
// @decision.name: Root Decision
// @decision.status: deprecated
`

	err := os.WriteFile(filepath.Join(tmpDir, "backend.js"), []byte(contentBackend), 0644)
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(tmpDir, "frontend.js"), []byte(contentFrontend), 0644)
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(tmpDir, "root.js"), []byte(contentNoCategory), 0644)
	assert.NoError(t, err)

	// Test filtering for backend category
	var buf bytes.Buffer
	err = ListCommand(tmpDir, "backend", &buf)
	assert.NoError(t, err)

	output := buf.String()
	// Should contain backend ADR
	assert.Contains(t, output, "adr-backend-1")
	assert.Contains(t, output, "Backend Decision")
	assert.Contains(t, output, "backend")

	// Should NOT contain frontend or root ADRs
	assert.NotContains(t, output, "adr-frontend-1")
	assert.NotContains(t, output, "Frontend Decision")
	assert.NotContains(t, output, "adr-root-1")
	assert.NotContains(t, output, "Root Decision")
}
