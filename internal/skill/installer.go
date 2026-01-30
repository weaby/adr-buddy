package skill

import (
	"fmt"
	"os"
	"path/filepath"
)

func InstallProjectLevel(projectRoot string) error {
	skillDir := GetSkillPath(projectRoot)
	return install(skillDir)
}

func InstallUserLevel() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	skillDir := GetSkillPath(homeDir)
	return install(skillDir)
}

func install(skillFilePath string) error {
	if _, err := os.Stat(skillFilePath); err == nil {
		return nil
	}

	skillDir := filepath.Dir(skillFilePath)
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		return fmt.Errorf("failed to create skill directory: %w", err)
	}

	if err := os.WriteFile(skillFilePath, []byte(DefaultSkill), 0644); err != nil {
		return fmt.Errorf("failed to write skill file: %w", err)
	}

	return nil
}

func ExistsProjectLevel(projectRoot string) bool {
	skillPath := GetSkillPath(projectRoot)
	_, err := os.Stat(skillPath)
	return err == nil
}

func ExistsUserLevel() bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	skillPath := GetSkillPath(homeDir)
	_, err = os.Stat(skillPath)
	return err == nil
}

func GetSkillPath(rootDir string) string {
	return filepath.Join(rootDir, ".claude", "skills", "adr", "SKILL.md")
}

// ADR Review skill installer functions

func InstallAdrReviewProjectLevel(projectRoot string) error {
	skillDir := GetAdrReviewSkillPath(projectRoot)
	return installAdrReview(skillDir)
}

func InstallAdrReviewUserLevel() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	skillDir := GetAdrReviewSkillPath(homeDir)
	return installAdrReview(skillDir)
}

func installAdrReview(skillFilePath string) error {
	if _, err := os.Stat(skillFilePath); err == nil {
		return nil
	}

	skillDir := filepath.Dir(skillFilePath)
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		return fmt.Errorf("failed to create skill directory: %w", err)
	}

	if err := os.WriteFile(skillFilePath, []byte(AdrReviewSkill), 0644); err != nil {
		return fmt.Errorf("failed to write skill file: %w", err)
	}

	return nil
}

func ExistsAdrReviewProjectLevel(projectRoot string) bool {
	skillPath := GetAdrReviewSkillPath(projectRoot)
	_, err := os.Stat(skillPath)
	return err == nil
}

func ExistsAdrReviewUserLevel() bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	skillPath := GetAdrReviewSkillPath(homeDir)
	_, err = os.Stat(skillPath)
	return err == nil
}

func GetAdrReviewSkillPath(rootDir string) string {
	return filepath.Join(rootDir, ".claude", "skills", "adr-review", "SKILL.md")
}
