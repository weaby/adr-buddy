package cli

import (
	"fmt"
	"path/filepath"

	"github.com/weaby/adr-buddy/internal/config"
	"github.com/weaby/adr-buddy/internal/model"
	"github.com/weaby/adr-buddy/internal/parser"
)

// Check validates all annotations without generating files
func Check(rootDir string, strict bool) error {
	// Load configuration
	cfg, err := config.Load(rootDir)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Scan for annotations
	var allAnnotations []*model.Annotation
	for _, scanPath := range cfg.ScanPaths {
		// Resolve scan path relative to root directory
		fullPath := filepath.Join(rootDir, scanPath)
		annotations, err := parser.ScanDirectory(fullPath, cfg.Exclude)
		if err != nil {
			return fmt.Errorf("failed to scan %s: %w", scanPath, err)
		}
		allAnnotations = append(allAnnotations, annotations...)
	}

	// Print summary
	fmt.Printf("Found %d annotation(s)\n", len(allAnnotations))

	// If no annotations found, nothing to validate
	if len(allAnnotations) == 0 {
		fmt.Println("No annotations found - nothing to validate")
		return nil
	}

	// Validate annotations by aggregating them
	// This checks for:
	// - Missing required fields (id, name)
	// - Conflicting names for same ID
	// - Conflicting categories for same ID
	adrs, err := model.Aggregate(allAnnotations)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Validate status values
	for _, ann := range allAnnotations {
		if err := parser.ValidateStatus(ann.Status, strict); err != nil {
			return fmt.Errorf("%s: %w", ann.Location, err)
		}
	}

	// Print success summary
	fmt.Printf("Validated %d ADR(s) successfully\n", len(adrs))
	for _, adr := range adrs {
		fmt.Printf("  - %s: %s (%d location(s))\n", adr.ID, adr.Name, len(adr.Locations))
	}

	return nil
}
