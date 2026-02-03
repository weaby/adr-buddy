package context

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type ADR struct {
	ID       string   `yaml:"adr_id"`
	Name     string   `yaml:"name"`
	Status   string   `yaml:"status"`
	Category string   `yaml:"category"`
	Date     string   `yaml:"date"`
	Globs    []string `yaml:"globs"`

	Context      string
	Decision     string
	Alternatives string
	Consequences string

	FilePath string `yaml:"-"`
}

type Reader struct {
	rootDir string
}

func NewReader(rootDir string) *Reader {
	return &Reader{rootDir: rootDir}
}

func (r *Reader) DecisionsDir() string {
	return filepath.Join(r.rootDir, ".claude", "rules", "decisions")
}

func (r *Reader) ReadAll() ([]*ADR, error) {
	decisionsDir := r.DecisionsDir()

	if _, err := os.Stat(decisionsDir); os.IsNotExist(err) {
		return nil, nil
	}

	var adrs []*ADR

	err := filepath.Walk(decisionsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}

		if info.Name() == "decisions-index.md" {
			return nil
		}

		adr, err := r.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading %s: %w", path, err)
		}

		adrs = append(adrs, adr)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return adrs, nil
}

func (r *Reader) ReadFile(path string) (*ADR, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	adr, err := ParseADR(data)
	if err != nil {
		return nil, err
	}

	relPath, err := filepath.Rel(r.rootDir, path)
	if err != nil {
		relPath = path
	}
	adr.FilePath = relPath

	return adr, nil
}

func ParseADR(data []byte) (*ADR, error) {
	frontmatter, body, err := splitFrontmatter(data)
	if err != nil {
		return nil, err
	}

	adr := &ADR{}
	if err := yaml.Unmarshal(frontmatter, adr); err != nil {
		return nil, fmt.Errorf("error parsing frontmatter: %w", err)
	}

	parseMarkdownSections(adr, body)

	return adr, nil
}

func splitFrontmatter(data []byte) ([]byte, []byte, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))

	if !scanner.Scan() || strings.TrimSpace(scanner.Text()) != "---" {
		return nil, nil, fmt.Errorf("ADR file must start with YAML frontmatter (---)")
	}

	var frontmatter bytes.Buffer
	foundEnd := false

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "---" {
			foundEnd = true
			break
		}
		frontmatter.WriteString(line)
		frontmatter.WriteString("\n")
	}

	if !foundEnd {
		return nil, nil, fmt.Errorf("YAML frontmatter not properly closed (missing ---)")
	}

	var body bytes.Buffer
	for scanner.Scan() {
		body.WriteString(scanner.Text())
		body.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return frontmatter.Bytes(), body.Bytes(), nil
}

func parseMarkdownSections(adr *ADR, body []byte) {
	sectionRe := regexp.MustCompile(`(?m)^##\s+(.+)$`)

	content := string(body)
	matches := sectionRe.FindAllStringSubmatchIndex(content, -1)

	for i, match := range matches {
		headerStart, headerEnd := match[2], match[3]
		header := strings.TrimSpace(content[headerStart:headerEnd])

		contentStart := match[1]
		var contentEnd int
		if i+1 < len(matches) {
			contentEnd = matches[i+1][0]
		} else {
			contentEnd = len(content)
		}

		sectionContent := strings.TrimSpace(content[contentStart:contentEnd])

		headerLower := strings.ToLower(header)
		switch {
		case strings.Contains(headerLower, "context"):
			adr.Context = sectionContent
		case strings.Contains(headerLower, "decision"):
			adr.Decision = sectionContent
		case strings.Contains(headerLower, "alternative"):
			adr.Alternatives = sectionContent
		case strings.Contains(headerLower, "consequence"):
			adr.Consequences = sectionContent
		}
	}
}

func (a *ADR) Validate() error {
	if a.ID == "" {
		return fmt.Errorf("missing required field: adr_id")
	}
	if a.Name == "" {
		return fmt.Errorf("missing required field: name")
	}
	return nil
}

// MatchesPath returns true if no globs are set or the path matches any glob.
func (a *ADR) MatchesPath(path string) bool {
	if len(a.Globs) == 0 {
		return true
	}

	for _, glob := range a.Globs {
		matched, err := filepath.Match(glob, path)
		if err == nil && matched {
			return true
		}
	}
	return false
}
