package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaby/adr-buddy/internal/model"
)

func TestCheck_ValidAnnotations(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test file with valid annotations
	testFile := filepath.Join(tmpDir, "test.go")
	content := `package main

// @decision.id: ADR-001
// @decision.name: Use PostgreSQL
// @decision.status: accepted
// @decision.context: We need a reliable database
func main() {}
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	// Initialize config
	err = Init(tmpDir)
	require.NoError(t, err)

	// Run check
	err = Check(tmpDir, false)
	assert.NoError(t, err)
}

func TestCheck_MissingRequiredFields(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test file with missing required field (name)
	testFile := filepath.Join(tmpDir, "test.go")
	content := `package main

// @decision.id: ADR-001
// @decision.status: proposed
func main() {}
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	// Initialize config
	err = Init(tmpDir)
	require.NoError(t, err)

	// Run check - should fail
	err = Check(tmpDir, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required field: @decision.name")
}

func TestCheck_ConflictingNames(t *testing.T) {
	tmpDir := t.TempDir()

	// Create two files with same ID but different names
	file1 := filepath.Join(tmpDir, "file1.go")
	content1 := `package main

// @decision.id: ADR-001
// @decision.name: Use PostgreSQL
func main() {}
`
	err := os.WriteFile(file1, []byte(content1), 0644)
	require.NoError(t, err)

	file2 := filepath.Join(tmpDir, "file2.go")
	content2 := `package main

// @decision.id: ADR-001
// @decision.name: Use MySQL
func main() {}
`
	err = os.WriteFile(file2, []byte(content2), 0644)
	require.NoError(t, err)

	// Initialize config
	err = Init(tmpDir)
	require.NoError(t, err)

	// Run check - should fail
	err = Check(tmpDir, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "conflicting names")
}

func TestCheck_ConflictingCategories(t *testing.T) {
	tmpDir := t.TempDir()

	// Create two files with same ID but different categories
	file1 := filepath.Join(tmpDir, "file1.go")
	content1 := `package main

// @decision.id: ADR-001
// @decision.name: Use PostgreSQL
// @decision.category: database
func main() {}
`
	err := os.WriteFile(file1, []byte(content1), 0644)
	require.NoError(t, err)

	file2 := filepath.Join(tmpDir, "file2.go")
	content2 := `package main

// @decision.id: ADR-001
// @decision.name: Use PostgreSQL
// @decision.category: infrastructure
func main() {}
`
	err = os.WriteFile(file2, []byte(content2), 0644)
	require.NoError(t, err)

	// Initialize config
	err = Init(tmpDir)
	require.NoError(t, err)

	// Run check - should fail
	err = Check(tmpDir, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "conflicting categories")
}

func TestCheck_InvalidStatus_StrictMode(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test file with invalid status
	testFile := filepath.Join(tmpDir, "test.go")
	content := `package main

// @decision.id: ADR-001
// @decision.name: Use PostgreSQL
// @decision.status: invalid-status
func main() {}
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	// Initialize config
	err = Init(tmpDir)
	require.NoError(t, err)

	// Enable strict mode
	configPath := filepath.Join(tmpDir, ".adr-buddy", "config.yml")
	configContent := `scan_paths: ["."]
output_dir: decisions
exclude:
  - "**/node_modules/**"
  - "**/.git/**"
  - "**/vendor/**"
  - "**/build/**"
  - "**/dist/**"
  - "**/.next/**"
  - "**/.adr-buddy/**"
template: ""
strict_mode: true
`
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// Run check - should fail in strict mode
	err = Check(tmpDir, true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid status")
}

func TestCheck_InvalidStatus_NonStrictMode(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test file with invalid status
	testFile := filepath.Join(tmpDir, "test.go")
	content := `package main

// @decision.id: ADR-001
// @decision.name: Use PostgreSQL
// @decision.status: custom-status
func main() {}
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	// Initialize config (strict_mode is false by default)
	err = Init(tmpDir)
	require.NoError(t, err)

	// Run check - should succeed with warning in non-strict mode
	err = Check(tmpDir, false)
	assert.NoError(t, err)
}

func TestCheck_NoAnnotations(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test file without annotations
	testFile := filepath.Join(tmpDir, "test.go")
	content := `package main

func main() {
	println("Hello, world!")
}
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	// Initialize config
	err = Init(tmpDir)
	require.NoError(t, err)

	// Run check - should succeed
	err = Check(tmpDir, false)
	assert.NoError(t, err)
}

func TestCheck_NoConfigFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test file with valid annotations
	testFile := filepath.Join(tmpDir, "test.go")
	content := `package main

// @decision.id: ADR-001
// @decision.name: Use PostgreSQL
func main() {}
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	// Run check without initializing (uses default config)
	err = Check(tmpDir, false)
	assert.NoError(t, err)
}

func TestCheck_MultipleValidAnnotations(t *testing.T) {
	tmpDir := t.TempDir()

	// Create multiple files with different ADRs
	file1 := filepath.Join(tmpDir, "file1.go")
	content1 := `package main

// @decision.id: ADR-001
// @decision.name: Use PostgreSQL
// @decision.status: accepted
func main() {}
`
	err := os.WriteFile(file1, []byte(content1), 0644)
	require.NoError(t, err)

	file2 := filepath.Join(tmpDir, "file2.go")
	content2 := `package main

// @decision.id: ADR-002
// @decision.name: Use REST API
// @decision.status: proposed
func handler() {}
`
	err = os.WriteFile(file2, []byte(content2), 0644)
	require.NoError(t, err)

	// Initialize config
	err = Init(tmpDir)
	require.NoError(t, err)

	// Run check
	err = Check(tmpDir, false)
	assert.NoError(t, err)
}

func TestCheck_SameIDMultipleLocations(t *testing.T) {
	tmpDir := t.TempDir()

	// Create two files with same ID and same name (valid)
	file1 := filepath.Join(tmpDir, "file1.go")
	content1 := `package main

// @decision.id: ADR-001
// @decision.name: Use PostgreSQL
// @decision.context: First location context
func db1() {}
`
	err := os.WriteFile(file1, []byte(content1), 0644)
	require.NoError(t, err)

	file2 := filepath.Join(tmpDir, "file2.go")
	content2 := `package main

// @decision.id: ADR-001
// @decision.name: Use PostgreSQL
// @decision.context: Second location context
func db2() {}
`
	err = os.WriteFile(file2, []byte(content2), 0644)
	require.NoError(t, err)

	// Initialize config
	err = Init(tmpDir)
	require.NoError(t, err)

	// Run check - should succeed
	err = Check(tmpDir, false)
	assert.NoError(t, err)
}

func TestCheck_JSONFormat(t *testing.T) {
	// Setup temp directory with valid annotations
	tmpDir := t.TempDir()

	// Write config
	configDir := filepath.Join(tmpDir, ".adr-buddy")
	require.NoError(t, os.MkdirAll(configDir, 0755))
	configContent := `scan_paths: ["."]
output_dir: "./decisions"`
	require.NoError(t, os.WriteFile(filepath.Join(configDir, "config.yml"), []byte(configContent), 0644))

	// Write source file with annotation
	sourceFile := filepath.Join(tmpDir, "test.js")
	sourceContent := `// @decision.id: adr-1
// @decision.name: Test Decision
const x = 1;`
	require.NoError(t, os.WriteFile(sourceFile, []byte(sourceContent), 0644))

	// Capture stdout
	var buf bytes.Buffer
	err := CheckWithFormat(tmpDir, false, "json", &buf)
	require.NoError(t, err)

	// Parse JSON output
	var result model.CheckResult
	err = json.Unmarshal(buf.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, "pass", result.Status)
	assert.Equal(t, 1, result.Summary.TotalAnnotations)
}
