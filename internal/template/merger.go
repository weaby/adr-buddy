package template

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"github.com/weaby/adr-buddy/internal/model"
)

// ParsedADR represents a parsed existing ADR file
type ParsedADR struct {
	Frontmatter map[string]string
	Sections    map[string]string
}

// ParseExistingADR parses an existing ADR markdown file
func ParseExistingADR(content string) *ParsedADR {
	parsed := &ParsedADR{
		Frontmatter: make(map[string]string),
		Sections:    make(map[string]string),
	}

	// Parse frontmatter (Status, Date, Category)
	statusRe := regexp.MustCompile(`\*\*Status:\*\*\s*(.+)`)
	dateRe := regexp.MustCompile(`\*\*Date:\*\*\s*(.+)`)
	categoryRe := regexp.MustCompile(`\*\*Category:\*\*\s*(.+)`)

	if match := statusRe.FindStringSubmatch(content); match != nil {
		parsed.Frontmatter["Status"] = strings.TrimSpace(match[1])
	}
	if match := dateRe.FindStringSubmatch(content); match != nil {
		parsed.Frontmatter["Date"] = strings.TrimSpace(match[1])
	}
	if match := categoryRe.FindStringSubmatch(content); match != nil {
		parsed.Frontmatter["Category"] = strings.TrimSpace(match[1])
	}

	// Parse sections (Context, Decision, Consequences)
	// Match from ## SectionName until next ## or end of string
	parseSections := []string{"Context", "Decision", "Consequences"}
	for _, sectionName := range parseSections {
		// Pattern: ## SectionName followed by content until next ## or end
		pattern := fmt.Sprintf(`(?s)## %s\s*\n(.*?)(?:## |\z)`, sectionName)
		re := regexp.MustCompile(pattern)
		if match := re.FindStringSubmatch(content); match != nil {
			parsed.Sections[sectionName] = strings.TrimSpace(match[1])
		}
	}

	return parsed
}

// isPlaceholder checks if a section contains only a TODO placeholder
func isPlaceholder(content string) bool {
	trimmed := strings.TrimSpace(content)
	return strings.Contains(trimmed, "<!-- TODO:")
}

// Merge intelligently merges an ADR with existing content
// Rules:
// 1. Preserve Date from existing file
// 2. For each section (Context, Decision, Consequences):
//   - If annotation provides content → use annotation content (replace)
//   - If annotation empty AND section has manual content → preserve manual content
//   - If annotation empty AND section is placeholder → keep placeholder
//
// 3. Status: Always use status from annotation (updates allowed)
// 4. Locations: Always regenerate from current annotations
func Merge(adr *model.ADR, existingContent string, tmpl string) (string, error) {
	parsed := ParseExistingADR(existingContent)

	// Create merged ADR
	merged := &model.ADR{
		ID:           adr.ID,
		Name:         adr.Name,
		Status:       adr.Status, // Always use new status
		Category:     adr.Category,
		Date:         parsed.Frontmatter["Date"], // Preserve existing date
		Context:      adr.Context,
		Decision:     adr.Decision,
		Consequences: adr.Consequences,
		Locations:    adr.Locations, // Always use new locations
	}

	// If Date is empty in existing, use new date
	if merged.Date == "" {
		merged.Date = adr.Date
	}

	// Smart merge for Context
	if len(adr.Context) == 0 {
		// Annotation provides no context
		existingContext := parsed.Sections["Context"]
		if existingContext != "" && !isPlaceholder(existingContext) {
			// Preserve manual content
			merged.Context = []string{existingContext}
		}
		// Otherwise keep empty (will render placeholder)
	}

	// Smart merge for Decision
	if len(adr.Decision) == 0 {
		// Annotation provides no decision
		existingDecision := parsed.Sections["Decision"]
		if existingDecision != "" && !isPlaceholder(existingDecision) {
			// Preserve manual content
			merged.Decision = []string{existingDecision}
		}
		// Otherwise keep empty (will render placeholder)
	}

	// Smart merge for Consequences
	if len(adr.Consequences) == 0 {
		// Annotation provides no consequences
		existingConsequences := parsed.Sections["Consequences"]
		if existingConsequences != "" && !isPlaceholder(existingConsequences) {
			// Preserve manual content
			merged.Consequences = []string{existingConsequences}
		}
		// Otherwise keep empty (will render placeholder)
	}

	// Render the merged ADR using the template
	t, err := template.New("adr").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, merged); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
