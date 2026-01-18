package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/weaby/adr-buddy/internal/config"
	"github.com/weaby/adr-buddy/internal/model"
	"github.com/weaby/adr-buddy/internal/parser"
	"github.com/weaby/adr-buddy/internal/template"
)

// SyncCommand scans code and generates/updates ADR files
func SyncCommand(rootDir string, dryRun bool, strict bool) error {
	cfg, err := config.Load(rootDir)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Println("Scanning for annotations...")

	// Scan all configured paths
	var allAnnotations []*model.Annotation
	for _, scanPath := range cfg.ScanPaths {
		absPath := scanPath
		if !filepath.IsAbs(scanPath) {
			absPath = filepath.Join(rootDir, scanPath)
		}

		annotations, err := parser.ScanDirectory(absPath, cfg.Exclude)
		if err != nil {
			return fmt.Errorf("failed to scan %s: %w", scanPath, err)
		}
		allAnnotations = append(allAnnotations, annotations...)
	}

	fmt.Printf("Found %d annotation(s)\n", len(allAnnotations))

	if len(allAnnotations) == 0 {
		fmt.Println("No annotations found.")
		return nil
	}

	// Aggregate into ADRs
	adrs, err := model.Aggregate(allAnnotations)
	if err != nil {
		return fmt.Errorf("aggregation failed: %w", err)
	}

	fmt.Printf("Generated %d ADR(s)\n\n", len(adrs))

	// Load template
	tmplStr := template.DefaultTemplate()
	if cfg.Template != "" {
		customTmpl, err := os.ReadFile(filepath.Join(rootDir, cfg.Template))
		if err == nil {
			tmplStr = string(customTmpl)
		}
	}

	// Generate files
	outputDir := cfg.OutputDir
	if !filepath.IsAbs(outputDir) {
		outputDir = filepath.Join(rootDir, outputDir)
	}

	for _, adr := range adrs {
		outputPath := filepath.Join(outputDir, adr.OutputPath(""))

		if dryRun {
			fmt.Printf("[DRY RUN] Would write: %s\n", outputPath)
			continue
		}

		// Check if file exists
		var content string
		if existingContent, err := os.ReadFile(outputPath); err == nil {
			// Merge with existing
			content, err = template.Merge(adr, string(existingContent), tmplStr)
			if err != nil {
				return fmt.Errorf("failed to merge %s: %w", outputPath, err)
			}
			fmt.Printf("Updated: %s\n", outputPath)
		} else {
			// Render new
			content, err = template.Render(adr, tmplStr)
			if err != nil {
				return fmt.Errorf("failed to render %s: %w", outputPath, err)
			}
			fmt.Printf("Created: %s\n", outputPath)
		}

		// Create directory if needed
		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		// Write file
		if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", outputPath, err)
		}
	}

	if !dryRun {
		fmt.Println("\n✓ Sync complete")
	}

	return nil
}
