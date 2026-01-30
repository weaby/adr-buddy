# ADR Buddy

<img src="docs/assets/logo.png" alt="ADR Buddy Logo" width="150">

[![Go 1.21+](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License](https://img.shields.io/badge/License-Elastic%202.0-blue)](LICENSE)
[![Docs](https://img.shields.io/badge/docs-GitHub%20Pages-blue)](https://weaby.github.io/adr-buddy)

**AI builds fast. Documentation doesn't keep up.**

Code is being written faster than ever with AI assistants—but architectural decisions are getting lost in the velocity. Six months from now, no one remembers *why* that database was chosen or *why* that retry logic exists.

ADR Buddy fixes this by capturing decisions inline in your code. And with Claude Code integration, it happens automatically as you build—the same AI that makes the decisions also documents them.

## Features

- **Inline annotations** — Decisions live in your code, right where they're implemented
- **Auto-generated docs** — Run `adr-buddy sync` to generate markdown ADRs from annotations
- **Claude Code skills** — AI documents decisions as it codes (`/adr`) and discovers undocumented ones (`/adr-review`)
- **GitHub Actions** — Validate ADRs on PRs, auto-sync on merge
- **Language-agnostic** — Works with Go, TypeScript, Python, Rust, Ruby, and any language with comments
- **Zero config** — Start with `adr-buddy init`, customize later if needed

## Quick Start

### Install

```bash
go install github.com/weaby/adr-buddy/cmd/adr-buddy@latest
```

### Initialize

```bash
adr-buddy init
```

This creates `.adr-buddy/` with config and template, and optionally installs Claude Code skills.

### Add your first decision

```go
// @decision.id: adr-001
// @decision.name: PostgreSQL for primary datastore
// @decision.context: Need ACID compliance and strong ecosystem support.
// @decision.decision: Use PostgreSQL over MySQL or MongoDB.
// @decision.consequences: Team familiar with it, good tooling available.

func NewDatabase() *sql.DB {
```

### Generate ADR files

```bash
adr-buddy sync
```

Your ADR is now at `decisions/adr-001.md`.

## Claude Code Integration

ADR Buddy includes two skills that make documentation effortless:

### `/adr` — Document as you code

Claude automatically recognizes architectural decisions and documents them inline. When you choose a library, implement a pattern, or make a trade-off, Claude adds the annotation without being asked.

### `/adr-review` — Discover undocumented decisions

Already have a codebase? This skill scans your project for technology choices and patterns that were never documented—then guides you through capturing them as ADRs.

Install skills during `adr-buddy init` or manually:
- Project-level: `.claude/skills/adr.md`
- User-level: `~/.claude/skills/adr.md`

[Learn more in the docs →](https://weaby.github.io/adr-buddy/claude-code)

## Commands

| Command | Description |
|---------|-------------|
| `adr-buddy init` | Initialize config, template, and optionally Claude skills |
| `adr-buddy sync` | Generate/update ADR files from annotations |
| `adr-buddy check` | Validate annotations without generating files |
| `adr-buddy list` | List all discovered ADRs |

## GitHub Actions

ADR Buddy integrates with GitHub Actions for automated validation and sync:

- **Validate on PRs** — Catch missing or malformed ADRs before merge
- **Auto-sync on push** — Keep ADR files in sync automatically

[See workflow examples →](https://weaby.github.io/adr-buddy/github-actions)

## Documentation

Full documentation available at [weaby.github.io/adr-buddy](https://weaby.github.io/adr-buddy)

## License

ADR Buddy is licensed under the [Elastic License 2.0](LICENSE).
