package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/weaby/adr-buddy/internal/context"
)

type CheckResult struct {
	Status   string            `json:"status"`
	Errors   []ValidationError `json:"errors"`
	Warnings []ValidationError `json:"warnings"`
	Summary  CheckSummary      `json:"summary"`
}

type ValidationError struct {
	File     string `json:"file"`
	Type     string `json:"type"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

type CheckSummary struct {
	TotalADRs  int `json:"total_adrs"`
	ValidADRs  int `json:"valid_adrs"`
	ErrorCount int `json:"error_count"`
}

func Check(rootDir string, strict bool) error {
	return CheckWithFormat(rootDir, strict, "text", os.Stdout)
}

func CheckWithFormat(rootDir string, strict bool, format string, output io.Writer) error {
	reader := context.NewReader(rootDir)
	adrs, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read ADRs: %w", err)
	}

	result := &CheckResult{
		Status:   "pass",
		Errors:   []ValidationError{},
		Warnings: []ValidationError{},
	}

	for _, adr := range adrs {
		if err := adr.Validate(); err != nil {
			result.Errors = append(result.Errors, ValidationError{
				File:     adr.FilePath,
				Type:     "validation_error",
				Message:  err.Error(),
				Severity: "error",
			})
			result.Summary.ErrorCount++
		} else {
			result.Summary.ValidADRs++
		}
	}

	result.Summary.TotalADRs = len(adrs)

	if result.Summary.ErrorCount > 0 {
		result.Status = "fail"
	}

	if format == "json" {
		encoder := json.NewEncoder(output)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(result); err != nil {
			return fmt.Errorf("failed to encode JSON: %w", err)
		}
	} else {
		fmt.Fprintf(output, "Found %d ADR(s)\n", len(adrs))

		if len(adrs) == 0 {
			fmt.Fprintln(output, "No ADRs found - nothing to validate")
			return nil
		}

		for _, e := range result.Errors {
			fmt.Fprintf(output, "ERROR: %s - %s\n", e.File, e.Message)
		}

		if result.Summary.ErrorCount > 0 {
			return fmt.Errorf("validation failed with %d error(s)", result.Summary.ErrorCount)
		}

		fmt.Fprintf(output, "Validated %d ADR(s) successfully\n", result.Summary.ValidADRs)
	}

	if result.Status == "fail" {
		return fmt.Errorf("validation failed")
	}

	return nil
}
