---
name: adr-review
description: Discover and document undocumented architectural decisions in your codebase
---

# ADR Review

Discover and document undocumented architectural decisions in a codebase.

## When to Use This Skill

Use this skill when the user:
- Wants to discover undocumented architectural decisions
- Asks to review the codebase for decisions
- Wants to document existing technology choices
- Is onboarding to a new project and wants to understand its architecture
- Mentions "ADR review", "decision review", or "document existing decisions"

## Prerequisites

Before starting:

1. Verify adr-buddy is initialized:
   ```bash
   test -f .adr-buddy/config.yml && echo "Ready" || echo "Run: adr-buddy init"
   ```

2. Read `.adr-buddy/config.yml` to find:
   - `output_dir` - where ADR files are stored
   - `scan_paths` - directories to analyze

3. Find the next available ADR ID:
   ```bash
   ls -1 decisions/ 2>/dev/null | grep -E '^adr-[0-9]+' | sort -V | tail -1
   ```

## Step 1: Choose Analysis Depth

Ask the user: "How deep should I analyze your codebase?"

| Level | Name | What it discovers |
|-------|------|-------------------|
| **1** | **Dependencies** | Package manifests (go.mod, package.json) + config files (docker-compose, Dockerfile, CI) |
| **2** | **Structure** | Level 1 + architectural patterns from folder/file organization |
| **3** | **Patterns** | Level 2 + design patterns visible in code |

Present as multiple choice:
1. Level 1 - Dependencies only (fast)
2. Level 2 - Dependencies + structural patterns
3. Level 3 - Full analysis including code patterns

Default to Level 1 if user wants a quick scan.

## Step 2: Scan for Decisions

### Level 1: Dependencies

Scan these files to discover technology choices:

**Dependency Manifests:**

| File | Command | What to extract |
|------|---------|-----------------|
| `go.mod` | `grep -E "^\t[a-z]" go.mod` | Direct dependencies (not indented with //) |
| `package.json` | Read `dependencies` and `devDependencies` keys | Package names |
| `requirements.txt` | `cat requirements.txt` | Package names (ignore version specs) |
| `pyproject.toml` | Read `[project.dependencies]` section | Package names |
| `Cargo.toml` | Read `[dependencies]` section | Crate names |
| `Gemfile` | `grep "^gem " Gemfile` | Gem names |

**Configuration Files:**

| File | What to extract |
|------|-----------------|
| `docker-compose.yml` | Service images (postgres:*, redis:*, etc.) |
| `Dockerfile` | Base image (FROM line) |
| `.github/workflows/*.yml` | Actions used, test commands |
| `.env.example` | External service references |
| `Makefile` | Build tools referenced |

**For each discovered technology, record:**
- Name (e.g., "postgresql")
- Category (database, cache, framework, library, etc.)
- Source file where found
- Usage context (e.g., "docker-compose service" or "direct dependency")
