package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateStatus(t *testing.T) {
	tests := []struct {
		status  string
		strict  bool
		wantErr bool
	}{
		{"proposed", false, false},
		{"accepted", false, false},
		{"rejected", false, false},
		{"deprecated", false, false},
		{"superseded", false, false},
		{"unknown", false, false}, // Warning only
		{"unknown", true, true},   // Error in strict mode
		{"", false, false},        // Empty is OK (defaults to proposed)
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			err := ValidateStatus(tt.status, tt.strict)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
