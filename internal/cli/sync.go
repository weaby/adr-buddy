package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/weaby/adr-buddy/internal/context"
)

type SyncResult struct {
	IndexUpdated bool `json:"index_updated"`
	ADRCount     int  `json:"adr_count"`
}

func SyncCommand(rootDir string, dryRun bool, strict bool) error {
	return SyncWithFormat(rootDir, dryRun, "text", os.Stdout)
}

func SyncWithFormat(rootDir string, dryRun bool, format string, output io.Writer) error {
	reader := context.NewReader(rootDir)
	adrs, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read ADRs: %w", err)
	}

	result := &SyncResult{
		ADRCount: len(adrs),
	}

	if format == "text" {
		fmt.Fprintf(output, "Found %d ADR(s)\n", len(adrs))
	}

	if dryRun {
		result.IndexUpdated = true
		if format == "text" {
			fmt.Fprintln(output, "[DRY RUN] Would regenerate decisions-index.md")
		}
	} else {
		generator := context.NewIndexGenerator(rootDir)
		if err := generator.Generate(); err != nil {
			return fmt.Errorf("failed to generate index: %w", err)
		}
		result.IndexUpdated = true

		if format == "text" {
			fmt.Fprintf(output, "Regenerated decisions-index.md with %d ADR(s)\n", len(adrs))
		}
	}

	if format == "json" {
		encoder := json.NewEncoder(output)
		encoder.SetIndent("", "  ")
		return encoder.Encode(result)
	}

	return nil
}
