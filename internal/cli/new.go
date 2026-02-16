package cli

import (
	"fmt"
	"os"

	"github.com/weaby/adr-buddy/internal/context"
)

func NewCommand(rootDir string, name string, category string) error {
	reader := context.NewReader(rootDir)
	existing, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read existing ADRs: %w", err)
	}

	id := context.NextADRID(existing)

	adr := context.NewADR(id, name)
	if category != "" {
		adr.Category = category
	}

	writer := context.NewWriter(rootDir)
	if err := writer.WriteADR(adr); err != nil {
		return fmt.Errorf("failed to write ADR: %w", err)
	}

	generator := context.NewIndexGenerator(rootDir)
	if err := generator.Generate(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to regenerate index: %v\n", err)
	}

	fmt.Printf("Created %s: %s\n", id, name)
	return nil
}
