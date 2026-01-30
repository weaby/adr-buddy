# Configuration

ADR Buddy works with zero configuration, but you can customize everything.

## Configuration File

After running `adr-buddy init`, you'll have `.adr-buddy/config.yml`:

```yaml
scan_paths:
  - .
output_dir: decisions
exclude:
  - "**/node_modules/**"
  - "**/.git/**"
  - "**/vendor/**"
  - "**/build/**"
  - "**/dist/**"
  - "**/.next/**"
  - "**/.adr-buddy/**"
  - "**/.claude/**"
  - "**/.github/**"
template: ""
strict_mode: false
```

## Options

### scan_paths

Directories to scan for annotations. Paths are relative to project root.

```yaml
scan_paths:
  - ./src
  - ./lib
  - ./internal
```

Default: `["."]` (entire project)

### output_dir

Where to generate ADR markdown files.

```yaml
output_dir: ./docs/decisions
```

Default: `decisions`

### exclude

Glob patterns for files/directories to skip.

```yaml
exclude:
  - "**/node_modules/**"
  - "**/vendor/**"
  - "**/*.test.go"
  - "**/testdata/**"
```

Default: Common build/dependency directories

### template

Path to a custom ADR template file.

```yaml
template: .adr-buddy/template.md
```

Default: `""` (uses built-in template)

### strict_mode

Treat warnings as errors during validation.

```yaml
strict_mode: true
```

Default: `false`

---

## Custom Templates

ADR Buddy uses Go's `text/template` syntax. You can customize how ADRs are rendered.

### Default Template

The built-in template (created at `.adr-buddy/template.md` during init):

```markdown
# {{.ID}}: {{.Name}}

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
```

### Available Variables

| Variable | Type | Description |
|----------|------|-------------|
| `{{.ID}}` | string | Decision ID (e.g., `adr-001`) |
| `{{.Name}}` | string | Decision title |
| `{{.Status}}` | string | Status (proposed, accepted, etc.) |
| `{{.Date}}` | string | Date ADR was created |
| `{{.Category}}` | string | Category (may be empty) |
| `{{.Context}}` | []string | Context paragraphs |
| `{{.Decision}}` | []string | Decision paragraphs |
| `{{.Alternatives}}` | []string | Alternatives considered |
| `{{.Consequences}}` | []string | Consequences/trade-offs |
| `{{.Locations}}` | []Location | Code locations |

Each location has:

- `{{.File}}` — File path
- `{{.Line}}` — Line number

### Template Examples

**Minimal template:**

```markdown
# {{.ID}}: {{.Name}}

{{range .Context}}{{.}}{{end}}

**Decision:** {{range .Decision}}{{.}}{{end}}

---
Found in: {{range .Locations}}{{.File}}:{{.Line}} {{end}}
```

**With conditional sections:**

```markdown
# {{.ID}}: {{.Name}}

**Status:** {{.Status}} | **Date:** {{.Date}}

## Context
{{range .Context}}{{.}}
{{end}}

## Decision
{{range .Decision}}{{.}}
{{end}}

{{if .Alternatives}}
## Alternatives Considered
{{range .Alternatives}}{{.}}
{{end}}
{{end}}

{{if .Consequences}}
## Consequences
{{range .Consequences}}{{.}}
{{end}}
{{end}}
```

### Using a Custom Template

1. Edit `.adr-buddy/template.md` with your format
2. Reference it in config (optional if using default path):

```yaml
template: .adr-buddy/template.md
```

3. Run `adr-buddy sync` — new and updated ADRs use your template
