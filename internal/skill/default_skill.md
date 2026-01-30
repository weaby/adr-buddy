---
name: adr
description: Create and manage Architecture Decision Records using adr-buddy
---

# ADR Buddy

Use this skill when the user:
- Wants to document an architectural decision
- Asks to create an ADR
- Mentions "architectural decision" or "decision record"
- Asks why a certain choice was made
- Wants to record a technical choice

## Before Creating Annotations

1. Read `.adr-buddy/config.yml` to find:
   - `output_dir` - where ADR files are stored
   - `template` - path to custom template

2. Read `.adr-buddy/template.md` to discover available fields by looking for `{{.FieldName}}` patterns

3. Scan the output directory to find the highest existing ADR ID and suggest the next one

## Annotation Syntax

Place annotations in code comments. Use `//` for Go/JS/TS/Java/C/Rust or `#` for Python/Ruby/YAML/Shell.

### Required Fields
- `@decision.id` - Unique identifier (e.g., "adr-001")
- `@decision.name` - Short title for the decision

### Optional Fields
- `@decision.status` - One of: proposed, accepted, rejected, deprecated, superseded (default: proposed)
- `@decision.category` - Organizes ADRs into subdirectories
- `@decision.context` - Why this decision was needed (multi-line)
- `@decision.decision` - What was decided (multi-line)
- `@decision.consequences` - Impact of the decision (multi-line)

### Multi-line Values

Continue on the next line with indentation after the comment marker:

```go
// @decision.context: We needed a logging solution that provides
//   structured output for our observability stack while maintaining
//   low overhead in hot paths.
```

### Full Example

```go
// @decision.id: adr-003
// @decision.name: Use PostgreSQL for primary datastore
// @decision.status: accepted
// @decision.category: infrastructure
// @decision.context: We need a relational database that supports
//   JSONB for flexible schemas and has strong ecosystem support.
// @decision.decision: Adopt PostgreSQL 15+ as our primary datastore.
// @decision.consequences: Team needs PostgreSQL expertise. We gain
//   excellent query performance and JSONB flexibility.
```

## Running Commands

After creating or modifying annotations:
```bash
adr-buddy sync
```

To validate without generating files:
```bash
adr-buddy check
```

To see all ADRs:
```bash
adr-buddy list
```

## Workflow

1. Add annotations to code
2. Run `adr-buddy check` to validate
3. Run `adr-buddy sync` to generate ADR files
4. Commit both the annotated code and generated ADR files
