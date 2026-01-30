package template

// DefaultTemplate returns the embedded default ADR template
func DefaultTemplate() string {
	return `# {{.ID}}: {{.Name}}

**Status:** {{.Status}}
**Date:** {{.Date}}
{{if .Category}}**Category:** {{.Category}}{{end}}

## Context
{{if .Context}}
{{range .Context}}
{{.}}

{{end}}
{{else}}
<!-- TODO: Add context - what is the issue we're facing? -->
{{end}}

## Decision
{{if .Decision}}
{{range .Decision}}
{{.}}

{{end}}
{{else}}
<!-- TODO: Document the decision and rationale -->
{{end}}

## Alternatives Considered
{{if .Alternatives}}
{{range .Alternatives}}
{{.}}

{{end}}
{{else}}
<!-- TODO: What alternatives were considered and why were they rejected? -->
{{end}}

## Consequences
{{if .Consequences}}
{{range .Consequences}}
{{.}}

{{end}}
{{else}}
<!-- TODO: What are the positive/negative outcomes? -->
{{end}}

## Code Locations
{{range .Locations}}
- {{.File}}:{{.Line}}
{{end}}
`
}
