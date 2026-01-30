# Commands

Complete reference for all ADR Buddy CLI commands.

## adr-buddy init

Initialize ADR Buddy in the current directory.

```bash
adr-buddy init [flags]
```

**What it creates:**

```
.adr-buddy/
├── config.yml    # Configuration options
└── template.md   # ADR template (customizable)
```

**Flags:**

| Flag | Values | Description |
|------|--------|-------------|
| `--claude-skill` | `project`, `user`, `skip` | Install Claude Code skills |

**Examples:**

```bash
# Interactive mode (prompts for skill installation)
adr-buddy init

# Install skills for this project only
adr-buddy init --claude-skill=project

# Install skills globally
adr-buddy init --claude-skill=user

# Skip skill installation
adr-buddy init --claude-skill=skip
```

**Notes:**

- Safe to run multiple times — won't overwrite existing files
- Creates `decisions/` directory on first sync, not during init

---

## adr-buddy sync

Scan code for annotations and generate/update ADR files.

```bash
adr-buddy sync [flags]
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--dry-run` | `false` | Show what would change without writing files |
| `--format` | `text` | Output format: `text` or `json` |
| `--watch` | `false` | Re-run on file changes (not yet implemented) |

**Examples:**

```bash
# Generate/update ADR files
adr-buddy sync

# Preview changes without writing
adr-buddy sync --dry-run

# JSON output (useful for CI)
adr-buddy sync --format=json
```

**Output (text):**

```
Scanned 127 files
Found 5 ADR annotations
Created decisions/adr-001.md
Created decisions/adr-002.md
Updated decisions/adr-003.md
Unchanged decisions/adr-004.md
Unchanged decisions/adr-005.md
```

**Output (json):**

```json
{
  "scanned": 127,
  "found": 5,
  "created": ["decisions/adr-001.md", "decisions/adr-002.md"],
  "updated": ["decisions/adr-003.md"],
  "unchanged": ["decisions/adr-004.md", "decisions/adr-005.md"]
}
```

---

## adr-buddy check

Validate annotations without generating files.

```bash
adr-buddy check [flags]
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--strict` | `false` | Treat warnings as errors |
| `--format` | `text` | Output format: `text` or `json` |

**Examples:**

```bash
# Validate annotations
adr-buddy check

# Strict mode (useful for CI)
adr-buddy check --strict

# JSON output
adr-buddy check --format=json
```

**Output (text):**

```
Scanned 127 files
Found 5 ADR annotations

✓ adr-001: PostgreSQL for primary datastore
✓ adr-002: Redis for caching
⚠ adr-003: Missing @decision.context (warning)
✓ adr-004: Kafka for events
✗ adr-005: Missing required field @decision.name (error)

4 valid, 1 warning, 1 error
```

**Exit codes:**

| Code | Meaning |
|------|---------|
| 0 | All valid (warnings allowed unless `--strict`) |
| 1 | Errors found |

---

## adr-buddy list

List all discovered ADRs.

```bash
adr-buddy list [flags]
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--category` | `""` | Filter by category |

**Examples:**

```bash
# List all ADRs
adr-buddy list

# Filter by category
adr-buddy list --category=infrastructure
```

**Output:**

```
ID        Status    Category        Name
adr-001   accepted  infrastructure  PostgreSQL for primary datastore
adr-002   accepted  infrastructure  Redis for caching
adr-003   accepted  architecture    Layered architecture pattern
adr-004   proposed  infrastructure  Kafka for events
adr-005   accepted  security        JWT for authentication
```

---

## Global Behavior

### Configuration Discovery

All commands automatically look for `.adr-buddy/config.yml` in the current directory. If not found, defaults are used.

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error (validation failure, missing files, etc.) |

### JSON Output

Commands with `--format=json` are designed for CI/CD integration and scripting. They output structured data to stdout, with errors going to stderr.
