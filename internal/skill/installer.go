package skill

import (
	"fmt"
	"os"
	"path/filepath"
)

func InstallProjectLevel(projectRoot string) error {
	skillDir := filepath.Join(projectRoot, ".claude", "skills")
	return install(skillDir)
}

func InstallUserLevel() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	skillDir := filepath.Join(homeDir, ".claude", "skills")
	return install(skillDir)
}

func install(skillDir string) error {
	skillPath := filepath.Join(skillDir, "adr.md")

	if _, err := os.Stat(skillPath); err == nil {
		return nil
	}

	if err := os.MkdirAll(skillDir, 0755); err != nil {
		return fmt.Errorf("failed to create skill directory: %w", err)
	}

	if err := os.WriteFile(skillPath, []byte(DefaultSkill), 0644); err != nil {
		return fmt.Errorf("failed to write skill file: %w", err)
	}

	return nil
}

func ExistsProjectLevel(projectRoot string) bool {
	skillPath := filepath.Join(projectRoot, ".claude", "skills", "adr.md")
	_, err := os.Stat(skillPath)
	return err == nil
}

func ExistsUserLevel() bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	skillPath := filepath.Join(homeDir, ".claude", "skills", "adr.md")
	_, err = os.Stat(skillPath)
	return err == nil
}
