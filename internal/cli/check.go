package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/weaby/adr-buddy/internal/config"
	"github.com/weaby/adr-buddy/internal/model"
	"github.com/weaby/adr-buddy/internal/parser"
)

// Check validates all annotations without generating files (text output)
func Check(rootDir string, strict bool) error {
	return CheckWithFormat(rootDir, strict, "text", os.Stdout)
}

// CheckWithFormat validates annotations with specified output format
func CheckWithFormat(rootDir string, strict bool, format string, output io.Writer) error {
	cfg, err := config.Load(rootDir)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Scan for annotations
	var allAnnotations []*model.Annotation
	for _, scanPath := range cfg.ScanPaths {
		fullPath := filepath.Join(rootDir, scanPath)
		annotations, err := parser.ScanDirectory(fullPath, cfg.Exclude)
		if err != nil {
			return fmt.Errorf("failed to scan %s: %w", scanPath, err)
		}
		allAnnotations = append(allAnnotations, annotations...)
	}

	// Build result
	result := &model.CheckResult{
		Status:   model.StatusPass,
		Errors:   []model.ValidationError{},
		Warnings: []model.ValidationError{},
		Summary: model.ValidationSummary{
			TotalAnnotations: len(allAnnotations),
			ValidAnnotations: 0,
			ErrorCount:       0,
			WarningCount:     0,
		},
	}

	// Validate each annotation
	for _, ann := range allAnnotations {
		if err := ann.Validate(); err != nil {
			result.Errors = append(result.Errors, model.ValidationError{
				File:     ann.Location.File,
				Line:     ann.Location.Line,
				Type:     "missing_required_field",
				Message:  err.Error(),
				Severity: "error",
			})
			result.Summary.ErrorCount++
		} else {
			result.Summary.ValidAnnotations++
		}

		// Validate status
		if err := parser.ValidateStatus(ann.Status, strict); err != nil {
			severity := "warning"
			if strict {
				severity = "error"
				result.Summary.ErrorCount++
			} else {
				result.Summary.WarningCount++
			}

			validationErr := model.ValidationError{
				File:     ann.Location.File,
				Line:     ann.Location.Line,
				Type:     "invalid_status",
				Message:  err.Error(),
				Severity: severity,
			}

			if strict {
				result.Errors = append(result.Errors, validationErr)
			} else {
				result.Warnings = append(result.Warnings, validationErr)
			}
		}
	}

	// Aggregate to check for conflicts
	if len(allAnnotations) > 0 {
		_, err := model.Aggregate(allAnnotations)
		if err != nil {
			result.Errors = append(result.Errors, model.ValidationError{
				File:     "",
				Line:     0,
				Type:     "aggregation_error",
				Message:  err.Error(),
				Severity: "error",
			})
			result.Summary.ErrorCount++
		}
	}

	// Set final status
	if result.Summary.ErrorCount > 0 {
		result.Status = model.StatusFail
	} else if result.Summary.WarningCount > 0 {
		result.Status = model.StatusWarning
	}

	// Output based on format
	if format == "json" {
		encoder := json.NewEncoder(output)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(result); err != nil {
			return fmt.Errorf("failed to encode JSON: %w", err)
		}
	} else {
		// Text format (original behavior)
		fmt.Fprintf(output, "Found %d annotation(s)\n", len(allAnnotations))

		if len(allAnnotations) == 0 {
			fmt.Fprintln(output, "No annotations found - nothing to validate")
			return nil
		}

		// Print errors
		for _, err := range result.Errors {
			fmt.Fprintf(output, "ERROR: %s:%d - %s\n", err.File, err.Line, err.Message)
		}

		// Print warnings
		for _, warn := range result.Warnings {
			fmt.Fprintf(output, "WARNING: %s:%d - %s\n", warn.File, warn.Line, warn.Message)
		}

		if result.Summary.ErrorCount > 0 {
			return fmt.Errorf("validation failed with %d error(s)", result.Summary.ErrorCount)
		}

		fmt.Fprintf(output, "Validated %d ADR(s) successfully\n", result.Summary.ValidAnnotations)
	}

	// Return error if validation failed
	if result.Status == model.StatusFail {
		return fmt.Errorf("validation failed")
	}

	return nil
}
