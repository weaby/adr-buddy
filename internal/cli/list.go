package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"text/tabwriter"

	"github.com/weaby/adr-buddy/internal/config"
	"github.com/weaby/adr-buddy/internal/model"
	"github.com/weaby/adr-buddy/internal/parser"
)

// ListCommand lists all discovered ADRs in tabular format.
// It scans the configured paths in rootDir, aggregates annotations into ADRs,
// optionally filters by category, and writes the table to output.
// If output is nil, writes to os.Stdout.
// Returns an error if scanning or aggregation fails.
func ListCommand(rootDir, category string, output io.Writer) error {
	if output == nil {
		output = os.Stdout
	}

	cfg, err := config.Load(rootDir)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
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

	// Check if any annotations found
	if len(allAnnotations) == 0 {
		fmt.Fprintln(output, "No annotations found.")
		return nil
	}

	// Aggregate into ADRs
	adrs, err := model.Aggregate(allAnnotations)
	if err != nil {
		return fmt.Errorf("aggregation failed: %w", err)
	}

	// Filter by category if specified
	if category != "" {
		var filtered []*model.ADR
		for _, adr := range adrs {
			if adr.Category == category {
				filtered = append(filtered, adr)
			}
		}
		adrs = filtered

		if len(adrs) == 0 {
			fmt.Fprintf(output, "No ADRs found in category %q.\n", category)
			return nil
		}
	}

	// Sort ADRs by ID for deterministic output
	sort.Slice(adrs, func(i, j int) bool {
		return adrs[i].ID < adrs[j].ID
	})

	// Print table
	w := tabwriter.NewWriter(output, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tSTATUS\tCATEGORY\tLOCATIONS")
	for _, adr := range adrs {
		cat := adr.Category
		if cat == "" {
			cat = "-"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\n",
			adr.ID, adr.Name, adr.Status, cat, len(adr.Locations))
	}
	w.Flush()

	return nil
}
