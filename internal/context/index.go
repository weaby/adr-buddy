package context

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

type IndexGenerator struct {
	rootDir string
}

func NewIndexGenerator(rootDir string) *IndexGenerator {
	return &IndexGenerator{rootDir: rootDir}
}

func (g *IndexGenerator) DecisionsDir() string {
	return filepath.Join(g.rootDir, ".claude", "rules", "decisions")
}

func (g *IndexGenerator) IndexPath() string {
	return filepath.Join(g.DecisionsDir(), "decisions-index.md")
}

func (g *IndexGenerator) Generate() error {
	reader := NewReader(g.rootDir)
	adrs, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read ADRs: %w", err)
	}

	content, err := RenderIndex(adrs)
	if err != nil {
		return fmt.Errorf("failed to render index: %w", err)
	}

	if err := os.MkdirAll(g.DecisionsDir(), 0755); err != nil {
		return fmt.Errorf("failed to create decisions directory: %w", err)
	}

	if err := os.WriteFile(g.IndexPath(), []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write index file: %w", err)
	}

	return nil
}

type ADRsByCategory struct {
	Category string
	ADRs     []*ADR
}

func RenderIndex(adrs []*ADR) (string, error) {
	categoryMap := make(map[string][]*ADR)
	for _, adr := range adrs {
		category := adr.Category
		if category == "" {
			category = "general"
		}
		categoryMap[category] = append(categoryMap[category], adr)
	}

	var categories []string
	for cat := range categoryMap {
		categories = append(categories, cat)
	}
	sort.Strings(categories)

	for _, adrsInCat := range categoryMap {
		sort.Slice(adrsInCat, func(i, j int) bool {
			return adrsInCat[i].ID < adrsInCat[j].ID
		})
	}

	var grouped []ADRsByCategory
	for _, cat := range categories {
		grouped = append(grouped, ADRsByCategory{
			Category: cat,
			ADRs:     categoryMap[cat],
		})
	}

	tmpl := `# Architecture Decision Records

This file provides an index of all architecture decisions for this project.
Claude Code will automatically load these decisions when working in relevant areas of the codebase.

## Quick Reference

| ID | Decision | Status | Category |
|----|----------|--------|----------|
{{- range .AllADRs }}
| {{ .ID }} | [{{ .Name }}]({{ .FilePath | basename }}) | {{ .Status }} | {{ .Category | default "general" }} |
{{- end }}

{{- range .Grouped }}

## {{ .Category | title }}

{{- range .ADRs }}

### {{ .ID }}: {{ .Name }}

**Status:** {{ .Status }}
{{- if .Globs }}
**Applies to:** {{ .Globs | join ", " }}
{{- end }}

{{ .Decision | truncate 200 }}

[Read full decision]({{ .FilePath | basename }})
{{- end }}
{{- end }}
`

	funcMap := template.FuncMap{
		"title": func(s string) string {
			if s == "" {
				return ""
			}
			return strings.ToUpper(s[:1]) + s[1:]
		},
		"default": func(def, val string) string {
			if val == "" {
				return def
			}
			return val
		},
		"basename": filepath.Base,
		"join": func(sep string, elems []string) string {
			return strings.Join(elems, sep)
		},
		"truncate": func(max int, s string) string {
			if len(s) <= max {
				return s
			}
			return s[:max] + "..."
		},
	}

	t, err := template.New("index").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	data := struct {
		AllADRs []*ADR
		Grouped []ADRsByCategory
	}{
		AllADRs: adrs,
		Grouped: grouped,
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
