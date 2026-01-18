package model

// Status constants for check results
const (
	StatusPass    = "pass"
	StatusWarning = "warning"
	StatusFail    = "fail"
)

// ValidationError represents a single validation error or warning
type ValidationError struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Type     string `json:"type"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

// ValidationSummary provides aggregate validation statistics
type ValidationSummary struct {
	TotalAnnotations int `json:"total_annotations"`
	ValidAnnotations int `json:"valid_annotations"`
	ErrorCount       int `json:"error_count"`
	WarningCount     int `json:"warning_count"`
}

// CheckResult represents the complete output of the check command
type CheckResult struct {
	Status   string              `json:"status"`
	Errors   []ValidationError   `json:"errors"`
	Warnings []ValidationError   `json:"warnings"`
	Summary  ValidationSummary   `json:"summary"`
}
