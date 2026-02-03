package context

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"gopkg.in/yaml.v3"
)

type Writer struct {
	rootDir string
}

func NewWriter(rootDir string) *Writer {
	return &Writer{rootDir: rootDir}
}

func (w *Writer) DecisionsDir() string {
	return filepath.Join(w.rootDir, ".claude", "rules", "decisions")
}

func (w *Writer) EnsureDecisionsDir() error {
	return os.MkdirAll(w.DecisionsDir(), 0755)
}

func (w *Writer) WriteADR(adr *ADR) error {
	if err := w.EnsureDecisionsDir(); err != nil {
		return fmt.Errorf("failed to create decisions directory: %w", err)
	}

	content, err := RenderADR(adr)
	if err != nil {
		return fmt.Errorf("failed to render ADR: %w", err)
	}

	filename := generateFilename(adr.ID, adr.Name)
	path := filepath.Join(w.DecisionsDir(), filename)

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write ADR file: %w", err)
	}

	return nil
}

func RenderADR(adr *ADR) (string, error) {
	var buf bytes.Buffer

	buf.WriteString("---\n")

	frontmatter := struct {
		ID       string   `yaml:"adr_id"`
		Name     string   `yaml:"name"`
		Status   string   `yaml:"status"`
		Category string   `yaml:"category,omitempty"`
		Date     string   `yaml:"date"`
		Globs    []string `yaml:"globs,omitempty"`
	}{
		ID:       adr.ID,
		Name:     adr.Name,
		Status:   adr.Status,
		Category: adr.Category,
		Date:     adr.Date,
		Globs:    adr.Globs,
	}

	if frontmatter.Status == "" {
		frontmatter.Status = "proposed"
	}

	if frontmatter.Date == "" {
		frontmatter.Date = time.Now().Format("2006-01-02")
	}

	yamlData, err := yaml.Marshal(frontmatter)
	if err != nil {
		return "", fmt.Errorf("failed to marshal frontmatter: %w", err)
	}
	buf.Write(yamlData)
	buf.WriteString("---\n\n")

	tmpl := `# {{ .ID | upper }}: {{ .Name }}

## Context

{{ .Context | default "[Why this decision was needed]" }}

## Decision

{{ .Decision | default "[What was decided]" }}

## Alternatives Considered

{{ .Alternatives | default "[Other options that were considered]" }}

## Consequences

{{ .Consequences | default "**Positive:** [Benefits]\n**Negative:** [Drawbacks]" }}
`

	funcMap := template.FuncMap{
		"upper": strings.ToUpper,
		"default": func(def, val string) string {
			if val == "" {
				return def
			}
			return val
		},
	}

	t, err := template.New("adr").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	if err := t.Execute(&buf, adr); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

func generateFilename(id, name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")

	var cleaned strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			cleaned.WriteRune(r)
		}
	}

	cleanName := cleaned.String()
	if len(cleanName) > 50 {
		cleanName = cleanName[:50]
	}

	cleanName = strings.TrimRight(cleanName, "-")

	return fmt.Sprintf("%s-%s.md", id, cleanName)
}

func NewADR(id, name string) *ADR {
	return &ADR{
		ID:     id,
		Name:   name,
		Status: "proposed",
		Date:   time.Now().Format("2006-01-02"),
	}
}

func NextADRID(adrs []*ADR) string {
	maxNum := 0
	for _, adr := range adrs {
		var num int
		_, err := fmt.Sscanf(adr.ID, "adr-%d", &num)
		if err == nil && num > maxNum {
			maxNum = num
		}
	}
	return fmt.Sprintf("adr-%03d", maxNum+1)
}
