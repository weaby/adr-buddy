package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	_ "github.com/weaby/adr-buddy/internal/model"
)

func TestDetectCommentStyle(t *testing.T) {
	tests := []struct {
		filename string
		want     string
	}{
		{"main.go", "//"},
		{"app.js", "//"},
		{"component.tsx", "//"},
		{"main.py", "#"},
		{"script.rb", "#"},
		{"deploy.sh", "#"},
		{"config.yml", "#"},
		{"unknown.txt", "//"},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			got := detectCommentStyle(tt.filename)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsAnnotationLine(t *testing.T) {
	tests := []struct {
		name         string
		line         string
		commentStyle string
		want         bool
	}{
		{
			name:         "valid annotation with //",
			line:         "// @decision.id: adr-1",
			commentStyle: "//",
			want:         true,
		},
		{
			name:         "valid annotation with #",
			line:         "# @decision.name: Using Pino",
			commentStyle: "#",
			want:         true,
		},
		{
			name:         "regular comment",
			line:         "// This is a regular comment",
			commentStyle: "//",
			want:         false,
		},
		{
			name:         "annotation wrong comment style",
			line:         "# @decision.id: adr-1",
			commentStyle: "//",
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isAnnotationLine(tt.line, tt.commentStyle)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseAnnotationField(t *testing.T) {
	tests := []struct {
		name         string
		line         string
		commentStyle string
		wantField    string
		wantValue    string
		wantOk       bool
	}{
		{
			name:         "simple field",
			line:         "// @decision.id: adr-1",
			commentStyle: "//",
			wantField:    "id",
			wantValue:    "adr-1",
			wantOk:       true,
		},
		{
			name:         "field with spaces in value",
			line:         "// @decision.name: Using Pino for logging",
			commentStyle: "//",
			wantField:    "name",
			wantValue:    "Using Pino for logging",
			wantOk:       true,
		},
		{
			name:         "hash comment style",
			line:         "# @decision.status: accepted",
			commentStyle: "#",
			wantField:    "status",
			wantValue:    "accepted",
			wantOk:       true,
		},
		{
			name:         "not an annotation",
			line:         "// regular comment",
			commentStyle: "//",
			wantOk:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field, value, ok := parseAnnotationField(tt.line, tt.commentStyle)
			assert.Equal(t, tt.wantOk, ok)
			if ok {
				assert.Equal(t, tt.wantField, field)
				assert.Equal(t, tt.wantValue, value)
			}
		})
	}
}

func TestIsContinuationLine(t *testing.T) {
	tests := []struct {
		name         string
		line         string
		commentStyle string
		want         bool
	}{
		{
			name:         "continuation with //",
			line:         "//   for our high-throughput API",
			commentStyle: "//",
			want:         true,
		},
		{
			name:         "continuation with #",
			line:         "#   additional context here",
			commentStyle: "#",
			want:         true,
		},
		{
			name:         "new field not continuation",
			line:         "// @decision.status: accepted",
			commentStyle: "//",
			want:         false,
		},
		{
			name:         "no space after marker not continuation",
			line:         "//no space",
			commentStyle: "//",
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isContinuationLine(tt.line, tt.commentStyle)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExtractContinuationValue(t *testing.T) {
	tests := []struct {
		name         string
		line         string
		commentStyle string
		want         string
	}{
		{
			name:         "extract from // comment with 2 spaces",
			line:         "//   for our high-throughput API",
			commentStyle: "//",
			want:         "for our high-throughput API",
		},
		{
			name:         "extract from # comment with 2 spaces",
			line:         "#   additional context here",
			commentStyle: "#",
			want:         "additional context here",
		},
		{
			name:         "extract from // comment with 1 space",
			line:         "//  single space indent",
			commentStyle: "//",
			want:         "single space indent",
		},
		{
			name:         "extract from // comment with multiple spaces",
			line:         "//     many spaces before text",
			commentStyle: "//",
			want:         "many spaces before text",
		},
		{
			name:         "extract from # comment with 1 space",
			line:         "#  Python continuation",
			commentStyle: "#",
			want:         "Python continuation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractContinuationValue(tt.line, tt.commentStyle)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseFile(t *testing.T) {
	annotations, err := ParseFile("../../testdata/fixtures/example.js")

	assert.NoError(t, err)
	assert.Len(t, annotations, 1)

	ann := annotations[0]
	assert.Equal(t, "adr-1", ann.ID)
	assert.Equal(t, "Using Pino for logging", ann.Name)
	assert.Equal(t, "accepted", ann.Status)
	assert.Equal(t, "infrastructure", ann.Category)
	assert.Contains(t, ann.Context, "We needed structured logging")
	assert.Contains(t, ann.Context, "Pino provided the best performance")
	assert.Contains(t, ann.Decision, "Adopt Pino")
	assert.Contains(t, ann.Consequences, "migrate from Winston")
	assert.Equal(t, "../../testdata/fixtures/example.js", ann.Location.File)
	assert.Equal(t, 1, ann.Location.Line)
}

func TestScanDirectory(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()

	// Create test files
	srcDir := filepath.Join(tmpDir, "src")
	err := os.MkdirAll(srcDir, 0755)
	assert.NoError(t, err)

	// Create a JS file with annotation
	jsContent := `// @decision.id: adr-1
// @decision.name: Test decision
console.log('test');
`
	err = os.WriteFile(filepath.Join(srcDir, "app.js"), []byte(jsContent), 0644)
	assert.NoError(t, err)

	// Create node_modules (should be excluded)
	nodeModules := filepath.Join(tmpDir, "node_modules")
	err = os.MkdirAll(nodeModules, 0755)
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(nodeModules, "lib.js"), []byte(jsContent), 0644)
	assert.NoError(t, err)

	// Scan directory
	excludePatterns := []string{"**/node_modules/**"}
	annotations, err := ScanDirectory(tmpDir, excludePatterns)

	assert.NoError(t, err)
	assert.Len(t, annotations, 1)
	assert.Equal(t, "adr-1", annotations[0].ID)
	assert.Contains(t, annotations[0].Location.File, "src/app.js")
}
