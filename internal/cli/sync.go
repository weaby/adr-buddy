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
	"github.com/weaby/adr-buddy/internal/template"
)

// SyncCommand scans code and generates/updates ADR files (text output)
func SyncCommand(rootDir string, dryRun bool, strict bool) error {
	return SyncWithFormat(rootDir, dryRun, "text", os.Stdout)
}

// SyncWithFormat scans and syncs with specified output format
func SyncWithFormat(rootDir string, dryRun bool, format string, output io.Writer) error {
	cfg, err := config.Load(rootDir)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if format == "text" {
		fmt.Fprintln(output, "Scanning for annotations...")
	}

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

	if format == "text" {
		fmt.Fprintf(output, "Found %d annotation(s)\n", len(allAnnotations))
	}

	if len(allAnnotations) == 0 {
		if format == "json" {
			result := &model.SyncResult{
				ChangesDetected: false,
				Files: model.FileChanges{
					Created:  []string{},
					Modified: []string{},
					Deleted:  []string{},
				},
				ADRs: []model.ADRChange{},
			}
			encoder := json.NewEncoder(output)
			encoder.SetIndent("", "  ")
			return encoder.Encode(result)
		}
		fmt.Fprintln(output, "No annotations found.")
		return nil
	}

	// Aggregate into ADRs
	adrs, err := model.Aggregate(allAnnotations)
	if err != nil {
		return fmt.Errorf("aggregation failed: %w", err)
	}

	if format == "text" {
		fmt.Fprintf(output, "Generated %d ADR(s)\n\n", len(adrs))
	}

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

	result := &model.SyncResult{
		ChangesDetected: false,
		Files: model.FileChanges{
			Created:  []string{},
			Modified: []string{},
			Deleted:  []string{},
		},
		ADRs: []model.ADRChange{},
	}

	for _, adr := range adrs {
		outputPath := filepath.Join(outputDir, adr.OutputPath(""))
		relPath, _ := filepath.Rel(rootDir, outputPath)

		var action string
		if _, err := os.Stat(outputPath); err == nil {
			action = "update"
			result.Files.Modified = append(result.Files.Modified, relPath)
		} else {
			action = "create"
			result.Files.Created = append(result.Files.Created, relPath)
		}

		result.ADRs = append(result.ADRs, model.ADRChange{
			ID:       adr.ID,
			Name:     adr.Name,
			Action:   action,
			FilePath: relPath,
		})

		if dryRun {
			if format == "text" {
				fmt.Fprintf(output, "[DRY RUN] Would write: %s\n", relPath)
			}
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
			if format == "text" {
				fmt.Fprintf(output, "Updated: %s\n", relPath)
			}
		} else {
			// Render new
			content, err = template.Render(adr, tmplStr)
			if err != nil {
				return fmt.Errorf("failed to render %s: %w", outputPath, err)
			}
			if format == "text" {
				fmt.Fprintf(output, "Created: %s\n", relPath)
			}
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

	result.ChangesDetected = len(result.Files.Created) > 0 || len(result.Files.Modified) > 0

	// Output based on format
	if format == "json" {
		encoder := json.NewEncoder(output)
		encoder.SetIndent("", "  ")
		return encoder.Encode(result)
	}

	if !dryRun {
		fmt.Fprintln(output, "\n✓ Sync complete")
	}

	return nil
}
