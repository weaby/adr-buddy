package model

import (
	"errors"
	"fmt"
	"path/filepath"
)

// SourceLocation represents a location in source code
type SourceLocation struct {
	File string // Relative path from project root
	Line int    // Line number where annotation starts
}

// String returns a formatted location string
func (s SourceLocation) String() string {
	return fmt.Sprintf("%s:%d", s.File, s.Line)
}

// Annotation represents a single annotation block found in code
type Annotation struct {
	ID           string            // Required
	Name         string            // Required
	Status       string            // Optional (default: "proposed")
	Category     string            // Optional (empty = root)
	Context      string            // Optional multi-line
	Decision     string            // Optional multi-line
	Alternatives string            // Optional multi-line
	Consequences string            // Optional multi-line
	CustomFields map[string]string // Future extensibility
	Location     SourceLocation    // Where this annotation appears
}

// Validate checks if the annotation has all required fields
func (a *Annotation) Validate() error {
	if a.ID == "" {
		return errors.New("missing required field: @decision.id")
	}
	if a.Name == "" {
		return errors.New("missing required field: @decision.name")
	}
	return nil
}

// DefaultStatus returns the default status if none is set
func (a *Annotation) DefaultStatus() string {
	if a.Status == "" {
		return "proposed"
	}
	return a.Status
}

// ADR represents an Architecture Decision Record
type ADR struct {
	ID           string
	Name         string
	Status       string
	Category     string
	Date         string           // Auto-generated on first creation
	Context      []string         // Merged from all annotations
	Decision     []string         // Merged from all annotations
	Alternatives []string         // Merged from all annotations
	Consequences []string         // Merged from all annotations
	Locations    []SourceLocation // All code locations
}

// OutputPath returns the file path where this ADR should be written
func (a *ADR) OutputPath(outputDir string) string {
	filename := a.ID + ".md"
	if a.Category == "" {
		return filepath.Join(outputDir, filename)
	}
	return filepath.Join(outputDir, a.Category, filename)
}
