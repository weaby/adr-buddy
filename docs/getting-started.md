# Getting Started

Get ADR Buddy running in under 5 minutes.

## Installation

```bash
go install github.com/weaby/adr-buddy/cmd/adr-buddy@latest
```

Verify the installation:

```bash
adr-buddy --help
```

## Initialize Your Project

Run init in your project root:

```bash
adr-buddy init
```

This creates:

```
.adr-buddy/
├── config.yml    # Configuration options
└── template.md   # ADR template (customizable)
```

You'll be prompted to install Claude Code skills:

```
Would you like to install the Claude Code skill?
  [1] Project-level (.claude/skills/adr.md) - for this project only
  [2] User-level (~/.claude/skills/adr.md) - available in all projects
  [3] Skip - don't install the skill
```

Or skip the prompt with a flag:

```bash
adr-buddy init --claude-skill=project  # or user, skip
```

## Add Your First Decision

Add an annotation where an architectural decision is implemented:

```go
// @decision.id: adr-001
// @decision.name: PostgreSQL for primary datastore
// @decision.status: accepted
// @decision.context: Need ACID compliance, strong ecosystem,
//   and team familiarity for our transactional workload.
// @decision.decision: Use PostgreSQL over MySQL or MongoDB.
// @decision.consequences: Proven technology, good tooling.
//   Team already has expertise.

func NewDatabase(connStr string) (*sql.DB, error) {
```

## Generate ADR Files

```bash
adr-buddy sync
```

Output:

```
Scanned 42 files
Found 1 ADR annotation
Created decisions/adr-001.md
```

Your ADR is now at `decisions/adr-001.md`.

## Validate Annotations

Check for issues without generating files:

```bash
adr-buddy check
```

## List All ADRs

```bash
adr-buddy list
```

## Next Steps

- [Annotation Reference](annotations.md) — Learn all available fields
- [Configuration](configuration.md) — Customize output paths and templates
- [Claude Code Integration](claude-code.md) — Let AI document decisions for you
- [GitHub Actions](github-actions.md) — Automate validation and sync
