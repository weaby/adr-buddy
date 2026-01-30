<div style="text-align: center; margin-bottom: 2rem;">
  <img src="assets/logo.png" alt="ADR Buddy" width="180">
</div>

# ADR Buddy

**AI builds fast. Documentation doesn't keep up.**

Code is being written faster than ever with AI assistants—but architectural decisions are getting lost in the velocity. Six months from now, no one remembers *why* that database was chosen or *why* that retry logic exists.

ADR Buddy captures decisions inline in your code using simple annotations. And with Claude Code integration, it happens automatically as you build.

## How It Works

1. **Annotate** — Add `@decision` comments where decisions are implemented
2. **Sync** — Run `adr-buddy sync` to generate markdown ADRs
3. **Automate** — Use GitHub Actions to validate and keep docs in sync

## Features

- **Inline annotations** — Decisions live in your code, not forgotten files
- **Auto-generated docs** — Markdown ADRs generated from annotations
- **Claude Code skills** — AI documents decisions as it codes
- **GitHub Actions** — Validate on PRs, auto-sync on merge
- **Language-agnostic** — Go, TypeScript, Python, Rust, Ruby, and more
- **Zero config** — Start immediately, customize later

## Quick Links

- [Getting Started](getting-started.md) — Install and create your first ADR
- [Annotation Reference](annotations.md) — All available fields and syntax
- [Configuration](configuration.md) — Customize paths, templates, and behavior
- [Claude Code Integration](claude-code.md) — AI-powered documentation
- [Commands](commands.md) — CLI reference
- [GitHub Actions](github-actions.md) — CI/CD integration
