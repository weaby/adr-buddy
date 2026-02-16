package cli

import (
	"fmt"
	"io"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/weaby/adr-buddy/internal/context"
)

func ListCommand(rootDir, category string, output io.Writer) error {
	if output == nil {
		output = os.Stdout
	}

	reader := context.NewReader(rootDir)
	adrs, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read ADRs: %w", err)
	}

	if len(adrs) == 0 {
		fmt.Fprintln(output, "No ADRs found.")
		return nil
	}

	if category != "" {
		var filtered []*context.ADR
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

	sort.Slice(adrs, func(i, j int) bool {
		return adrs[i].ID < adrs[j].ID
	})

	w := tabwriter.NewWriter(output, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tSTATUS\tCATEGORY")
	for _, adr := range adrs {
		cat := adr.Category
		if cat == "" {
			cat = "-"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			adr.ID, adr.Name, adr.Status, cat)
	}
	w.Flush()

	return nil
}
