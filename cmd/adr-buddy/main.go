package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/weaby/adr-buddy/internal/cli"
)

var rootCmd = &cobra.Command{
	Use:   "adr-buddy",
	Short: "Generate Architecture Decision Records from code annotations",
	Long: `adr-buddy is a CLI tool that scans your codebase for @decision annotations
and automatically generates and maintains ADR documentation.`,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize adr-buddy in the current directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cli.Init(".")
	},
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Scan code and generate/update ADR files",
	RunE: func(cmd *cobra.Command, args []string) error {
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		watch, _ := cmd.Flags().GetBool("watch")
		format, _ := cmd.Flags().GetString("format")

		if watch {
			return fmt.Errorf("watch mode not yet implemented")
		}

		return cli.SyncWithFormat(".", dryRun, format, os.Stdout)
	},
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Validate annotations without generating files",
	RunE: func(cmd *cobra.Command, args []string) error {
		strict, _ := cmd.Flags().GetBool("strict")
		format, _ := cmd.Flags().GetString("format")
		return cli.CheckWithFormat(".", strict, format, os.Stdout)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all discovered ADRs",
	RunE: func(cmd *cobra.Command, args []string) error {
		category, _ := cmd.Flags().GetString("category")
		return cli.ListCommand(".", category, nil)
	},
}

func init() {
	syncCmd.Flags().Bool("dry-run", false, "Show what would change without writing files")
	syncCmd.Flags().Bool("watch", false, "Continuous mode (re-run on file changes)")
	syncCmd.Flags().String("format", "text", "Output format: text or json")

	checkCmd.Flags().Bool("strict", false, "Treat warnings as errors")
	checkCmd.Flags().String("format", "text", "Output format: text or json")

	listCmd.Flags().String("category", "", "Filter by category")

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(listCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
