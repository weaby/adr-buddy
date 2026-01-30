# ADR Buddy

A language-agnostic CLI tool that generates Architecture Decision Records (ADRs) from code annotations.

## Features

- 📝 Write ADRs directly in your code using simple comment annotations
- 🔄 Automatically generate and sync ADR documentation
- 🌍 Language-agnostic: works with Go, JavaScript, TypeScript, Python, Ruby, and more
- 🎯 Smart merge: preserves manual edits while updating from annotations
- 📁 Flexible organization: optional categorization of ADRs
- ⚡ Zero configuration required to get started

## Installation

```bash
go install github.com/weaby/adr-buddy/cmd/adr-buddy@latest
```

## Quick Start

1. Initialize in your project:
```bash
adr-buddy init
```

2. Add annotations to your code:
```javascript
// @decision.id: adr-1
// @decision.name: Using Pino for logging
// @decision.status: accepted
// @decision.category: infrastructure
// @decision.context: We needed structured logging with low overhead
//   for our high-throughput API.
// @decision.decision: Adopt Pino as the standard logging library.
// @decision.consequences: Better performance but team needs training.

const pino = require('pino');
```

3. Generate ADR documentation:
```bash
adr-buddy sync
```

Your ADR will be created at `decisions/infrastructure/adr-1.md`!

## Claude Code Integration

ADR Buddy includes a Claude Code skill that teaches Claude how to create and manage ADRs in your codebase.

### Installing the Skill

During `adr-buddy init`, you'll be prompted to install the Claude Code skill:

```
Would you like to install the Claude Code skill?
  [1] Project-level (.claude/skills/adr.md) - for this project only
  [2] User-level (~/.claude/skills/adr.md) - available in all projects
  [3] Skip - don't install the skill
```

Or use the `--claude-skill` flag for non-interactive installation:

```bash
# Install for this project only
adr-buddy init --claude-skill=project

# Install globally for all projects
adr-buddy init --claude-skill=user

# Skip skill installation
adr-buddy init --claude-skill=skip
```

### What the Skill Does

Once installed, Claude Code will automatically:
- Recognize when you want to document an architectural decision
- Read your project's ADR template and configuration
- Create properly formatted annotations in your code
- Suggest the next available ADR ID
- Run `adr-buddy sync` to generate the ADR files

## Commands

- `adr-buddy init` - Initialize configuration in current directory
- `adr-buddy sync` - Scan code and generate/update ADR files
- `adr-buddy check` - Validate annotations without generating files
- `adr-buddy list` - List all discovered ADRs

## Annotation Reference

**Required:**
- `@decision.id` - Unique identifier (e.g., "adr-1")
- `@decision.name` - Decision title

**Optional:**
- `@decision.status` - proposed | accepted | rejected | deprecated | superseded
- `@decision.category` - Organizational category
- `@decision.context` - Background and rationale (multi-line)
- `@decision.decision` - The choice made (multi-line)
- `@decision.consequences` - Outcomes and impacts (multi-line)

## Multi-line Values

Use indentation for continuation:
```javascript
// @decision.context: This is the first line
//   and this continues on the next line
//   and this is the third line.
```

## Configuration

Optional `.adr-buddy/config.yml`:
```yaml
scan_paths:
  - ./src
  - ./lib
output_dir: ./decisions
exclude:
  - "**/node_modules/**"
  - "**/vendor/**"
```

## GitHub Actions Integration

ADR Buddy integrates seamlessly with GitHub Actions for automated validation and synchronization.

### PR Validation Workflow

Validates ADRs on every pull request:

```yaml
# .github/workflows/adr-validate.yml
name: Validate ADRs

on:
  pull_request:
    branches: [main]

jobs:
  validate:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      pull-requests: write
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      - uses: weaby/adr-buddy@v1
        with:
          mode: validate
```

### Auto-Sync Workflow

Automatically creates PRs with ADR updates after merges:

```yaml
# .github/workflows/adr-sync.yml
name: Sync ADRs

on:
  push:
    branches: [main]

jobs:
  sync:
    runs-on: ubuntu-latest
    permissions:
      contents: write
      pull-requests: write
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      - uses: weaby/adr-buddy@v1
        with:
          mode: sync
          reviewers: 'team-lead,architect'
```

### Configuration Options

**Validate Mode:**
- `mode: validate` - Run validation checks
- `strict: true` - Treat warnings as errors (default: true)
- `config-path: .adr-buddy/config.yml` - Custom config location

**Sync Mode:**
- `mode: sync` - Generate and sync ADR files
- `create-pr: true` - Create PR with changes (default: true)
- `reviewers: 'user1,user2'` - Comma-separated reviewers
- `token: ${{ secrets.GITHUB_TOKEN }}` - GitHub token

### Action Outputs

The action provides outputs you can use in subsequent steps:

```yaml
- uses: weaby/adr-buddy@v1
  id: adr
  with:
    mode: validate

- name: Check validation status
  if: steps.adr.outputs.validation-status == 'fail'
  run: echo "Validation failed"
```

Available outputs:
- `validation-status` - pass, warning, or fail
- `changes-detected` - true if ADRs need updating
- `pr-number` - PR number if created (sync mode)
- `pr-url` - PR URL if created (sync mode)

## Examples

### One-to-Many: Same ADR Referenced in Multiple Files

`src/logger.js`:
```javascript
// @decision.id: adr-1
// @decision.name: Using Pino
// @decision.context: Need fast logging
const logger = require('pino')();
```

`src/monitor.js`:
```javascript
// @decision.id: adr-1
// @decision.name: Using Pino
// @decision.decision: Use Pino everywhere
const logger = require('./logger');
```

Both annotations merge into a single `decisions/adr-1.md` with two code locations.

### Many-to-One: Multiple ADRs in One File

```python
# @decision.id: adr-5
# @decision.name: Microservices architecture
# ...

# @decision.id: adr-12
# @decision.name: gRPC for inter-service communication
# ...
```

Each annotation generates its own ADR file, both referencing the same source file.

## License

MIT
