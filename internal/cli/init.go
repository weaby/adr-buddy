package cli

import (
	"fmt"
	"os"

	"github.com/weaby/adr-buddy/internal/context"
	"github.com/weaby/adr-buddy/internal/skill"
)

type SkillLocation int

const (
	SkillLocationSkip SkillLocation = iota
	SkillLocationProject
	SkillLocationUser
)

func Init(rootDir string) error {
	return InitWithSkill(rootDir, SkillLocationSkip)
}

func InitWithSkill(rootDir string, skillLocation SkillLocation) error {
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", rootDir)
	}

	writer := context.NewWriter(rootDir)
	if err := writer.EnsureDecisionsDir(); err != nil {
		return fmt.Errorf("failed to create decisions directory: %w", err)
	}
	fmt.Println("Created decisions directory:", writer.DecisionsDir())

	generator := context.NewIndexGenerator(rootDir)
	if err := generator.Generate(); err != nil {
		return fmt.Errorf("failed to generate index: %w", err)
	}
	fmt.Println("Generated decisions index:", generator.IndexPath())

	switch skillLocation {
	case SkillLocationProject:
		if err := skill.InstallProjectLevel(rootDir); err != nil {
			return fmt.Errorf("failed to install adr skill: %w", err)
		}
		fmt.Println("Installed Claude Code skill:", skill.GetSkillPath(rootDir))
		if err := skill.InstallAdrReviewProjectLevel(rootDir); err != nil {
			return fmt.Errorf("failed to install adr-review skill: %w", err)
		}
		fmt.Println("Installed Claude Code skill:", skill.GetAdrReviewSkillPath(rootDir))
	case SkillLocationUser:
		if err := skill.InstallUserLevel(); err != nil {
			return fmt.Errorf("failed to install adr skill: %w", err)
		}
		homeDir, _ := os.UserHomeDir()
		fmt.Println("Installed Claude Code skill:", skill.GetSkillPath(homeDir))
		if err := skill.InstallAdrReviewUserLevel(); err != nil {
			return fmt.Errorf("failed to install adr-review skill: %w", err)
		}
		fmt.Println("Installed Claude Code skill:", skill.GetAdrReviewSkillPath(homeDir))
	}

	fmt.Println("Initialized adr-buddy successfully!")
	return nil
}
