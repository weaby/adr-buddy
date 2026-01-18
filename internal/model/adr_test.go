package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSourceLocation_String(t *testing.T) {
	loc := SourceLocation{
		File: "src/main.go",
		Line: 42,
	}

	assert.Equal(t, "src/main.go:42", loc.String())
}

func TestAnnotation_Validate(t *testing.T) {
	tests := []struct {
		name    string
		ann     Annotation
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid annotation",
			ann: Annotation{
				ID:   "adr-1",
				Name: "Using Pino",
				Location: SourceLocation{
					File: "test.js",
					Line: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			ann: Annotation{
				Name: "Using Pino",
				Location: SourceLocation{
					File: "test.js",
					Line: 1,
				},
			},
			wantErr: true,
			errMsg:  "missing required field: @decision.id",
		},
		{
			name: "missing name",
			ann: Annotation{
				ID: "adr-1",
				Location: SourceLocation{
					File: "test.js",
					Line: 1,
				},
			},
			wantErr: true,
			errMsg:  "missing required field: @decision.name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.ann.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestADR_OutputPath(t *testing.T) {
	tests := []struct {
		name      string
		adr       ADR
		outputDir string
		wantPath  string
	}{
		{
			name: "no category",
			adr: ADR{
				ID:       "adr-1",
				Category: "",
			},
			outputDir: "decisions",
			wantPath:  "decisions/adr-1.md",
		},
		{
			name: "with category",
			adr: ADR{
				ID:       "adr-2",
				Category: "infrastructure",
			},
			outputDir: "decisions",
			wantPath:  "decisions/infrastructure/adr-2.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.adr.OutputPath(tt.outputDir)
			assert.Equal(t, tt.wantPath, got)
		})
	}
}
