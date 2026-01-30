package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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
		skillFlag, _ := cmd.Flags().GetString("claude-skill")

		var skillLocation cli.SkillLocation
		switch skillFlag {
		case "project":
			skillLocation = cli.SkillLocationProject
		case "user":
			skillLocation = cli.SkillLocationUser
		case "skip":
			skillLocation = cli.SkillLocationSkip
		case "":
			skillLocation = promptSkillLocation()
		default:
			return fmt.Errorf("invalid --claude-skill value: %s (must be project, user, or skip)", skillFlag)
		}

		return cli.InitWithSkill(".", skillLocation)
	},
}

func promptSkillLocation() cli.SkillLocation {
	fmt.Println()
	fmt.Println("Would you like to install the Claude Code skill?")
	fmt.Println("  [1] Project-level (.claude/skills/adr.md) - for this project only")
	fmt.Println("  [2] User-level (~/.claude/skills/adr.md) - available in all projects")
	fmt.Println("  [3] Skip - don't install the skill")
	fmt.Println()
	fmt.Print("Choice [1/2/3]: ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	switch input {
	case "1":
		return cli.SkillLocationProject
	case "2":
		return cli.SkillLocationUser
	default:
		return cli.SkillLocationSkip
	}
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
	initCmd.Flags().String("claude-skill", "", "Install Claude Code skill: project, user, or skip")

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
