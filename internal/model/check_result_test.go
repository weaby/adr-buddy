package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckResult_MarshalJSON(t *testing.T) {
	result := &CheckResult{
		Status: StatusPass,
		Errors: []ValidationError{
			{
				File:     "src/logger.js",
				Line:     10,
				Type:     "missing_required_field",
				Message:  "Missing @decision.name",
				Severity: "error",
			},
		},
		Warnings: []ValidationError{},
		Summary: ValidationSummary{
			TotalAnnotations: 15,
			ValidAnnotations: 14,
			ErrorCount:       1,
			WarningCount:     0,
		},
	}

	data, err := json.Marshal(result)
	require.NoError(t, err)

	var decoded CheckResult
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Equal(t, result.Status, decoded.Status)
	assert.Equal(t, 1, len(decoded.Errors))
	assert.Equal(t, "src/logger.js", decoded.Errors[0].File)
}

func TestCheckResult_StatusValues(t *testing.T) {
	assert.Equal(t, "pass", StatusPass)
	assert.Equal(t, "warning", StatusWarning)
	assert.Equal(t, "fail", StatusFail)
}
