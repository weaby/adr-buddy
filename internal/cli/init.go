package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/weaby/adr-buddy/internal/config"
	"github.com/weaby/adr-buddy/internal/template"
)

// Init initializes the .adr-buddy directory with default config and template
func Init(rootDir string) error {
	// Verify root directory exists
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", rootDir)
	}

	// Create .adr-buddy directory
	adrDir := filepath.Join(rootDir, ".adr-buddy")
	if err := os.MkdirAll(adrDir, 0755); err != nil {
		return fmt.Errorf("failed to create .adr-buddy directory: %w", err)
	}

	// Create config.yml if it doesn't exist
	configPath := filepath.Join(adrDir, "config.yml")
	if err := createFileIfNotExists(configPath, func() ([]byte, error) {
		cfg := config.Default()
		data, err := yaml.Marshal(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal config: %w", err)
		}
		fmt.Println("Created config file:", configPath)
		return data, nil
	}); err != nil {
		return err
	}

	// Create template.md if it doesn't exist
	templatePath := filepath.Join(adrDir, "template.md")
	if err := createFileIfNotExists(templatePath, func() ([]byte, error) {
		fmt.Println("Created template file:", templatePath)
		return []byte(template.DefaultTemplate()), nil
	}); err != nil {
		return err
	}

	fmt.Println("Initialized .adr-buddy directory successfully!")
	return nil
}

// createFileIfNotExists creates a file only if it doesn't exist
func createFileIfNotExists(path string, contentFunc func() ([]byte, error)) error {
	// Check if file already exists
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("File already exists, skipping: %s\n", path)
		return nil
	}

	// Get content
	content, err := contentFunc()
	if err != nil {
		return err
	}

	// Write file
	if err := os.WriteFile(path, content, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}

	return nil
}
