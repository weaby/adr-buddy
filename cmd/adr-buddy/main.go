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
	Short: "Manage Architecture Decision Records for AI-assisted development",
	Long: `adr-buddy manages Architecture Decision Records (ADRs) stored as markdown files
with YAML frontmatter in .claude/rules/decisions/. These context files help AI
assistants understand your project's architectural decisions.`,
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

var newCmd = &cobra.Command{
	Use:   "new <name>",
	Short: "Create a new ADR",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		category, _ := cmd.Flags().GetString("category")
		return cli.NewCommand(".", args[0], category)
	},
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Regenerate the decisions index",
	RunE: func(cmd *cobra.Command, args []string) error {
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		format, _ := cmd.Flags().GetString("format")
		return cli.SyncWithFormat(".", dryRun, format, os.Stdout)
	},
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Validate all ADR files",
	RunE: func(cmd *cobra.Command, args []string) error {
		strict, _ := cmd.Flags().GetBool("strict")
		format, _ := cmd.Flags().GetString("format")
		return cli.CheckWithFormat(".", strict, format, os.Stdout)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all ADRs",
	RunE: func(cmd *cobra.Command, args []string) error {
		category, _ := cmd.Flags().GetString("category")
		return cli.ListCommand(".", category, nil)
	},
}

func init() {
	initCmd.Flags().String("claude-skill", "", "Install Claude Code skill: project, user, or skip")

	newCmd.Flags().String("category", "", "ADR category")

	syncCmd.Flags().Bool("dry-run", false, "Show what would change without writing files")
	syncCmd.Flags().String("format", "text", "Output format: text or json")

	checkCmd.Flags().Bool("strict", false, "Treat warnings as errors")
	checkCmd.Flags().String("format", "text", "Output format: text or json")

	listCmd.Flags().String("category", "", "Filter by category")

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(listCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
