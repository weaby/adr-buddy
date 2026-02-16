package cli

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheck_ValidADRs(t *testing.T) {
	tmpDir := t.TempDir()

	createTestADR(t, tmpDir, "adr-001-test.md", `---
adr_id: adr-001
name: Test Decision
status: accepted
date: 2024-01-15
---

## Decision

Test decision.
`)

	var buf bytes.Buffer
	err := CheckWithFormat(tmpDir, false, "text", &buf)
	require.NoError(t, err)

	assert.Contains(t, buf.String(), "Validated 1 ADR(s) successfully")
}

func TestCheck_MissingRequiredFields(t *testing.T) {
	tmpDir := t.TempDir()

	createTestADR(t, tmpDir, "adr-001-bad.md", `---
adr_id: adr-001
status: accepted
date: 2024-01-15
---

## Decision

Missing name field.
`)

	var buf bytes.Buffer
	err := CheckWithFormat(tmpDir, false, "text", &buf)
	assert.Error(t, err)
	assert.Contains(t, buf.String(), "ERROR")
}

func TestCheck_NoADRs(t *testing.T) {
	tmpDir := t.TempDir()

	var buf bytes.Buffer
	err := CheckWithFormat(tmpDir, false, "text", &buf)
	require.NoError(t, err)

	assert.Contains(t, buf.String(), "No ADRs found")
}

func TestCheck_JSONFormat(t *testing.T) {
	tmpDir := t.TempDir()

	createTestADR(t, tmpDir, "adr-001-test.md", `---
adr_id: adr-001
name: Test Decision
status: accepted
date: 2024-01-15
---

## Decision

Test decision.
`)

	var buf bytes.Buffer
	err := CheckWithFormat(tmpDir, false, "json", &buf)
	require.NoError(t, err)

	var result CheckResult
	err = json.Unmarshal(buf.Bytes(), &result)
	require.NoError(t, err)

	assert.Equal(t, "pass", result.Status)
	assert.Equal(t, 1, result.Summary.TotalADRs)
	assert.Equal(t, 1, result.Summary.ValidADRs)
}

func TestCheck_JSONFormat_WithErrors(t *testing.T) {
	tmpDir := t.TempDir()

	createTestADR(t, tmpDir, "adr-001-bad.md", `---
adr_id: adr-001
status: accepted
date: 2024-01-15
---

## Decision

Missing name field.
`)

	var buf bytes.Buffer
	err := CheckWithFormat(tmpDir, false, "json", &buf)
	assert.Error(t, err)

	var result CheckResult
	err = json.Unmarshal(buf.Bytes(), &result)
	require.NoError(t, err)

	assert.Equal(t, "fail", result.Status)
	assert.Equal(t, 1, result.Summary.ErrorCount)
}
