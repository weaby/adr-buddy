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

### Level 2: Structural Patterns

In addition to Level 1, analyze folder structure for architectural patterns:

**Pattern Detection:**

| Pattern | Indicator Folders | Decision to Document |
|---------|-------------------|---------------------|
| Layered Architecture | `/controllers`, `/services`, `/repositories`, `/models` | "Layered architecture pattern" |
| Hexagonal/Ports-Adapters | `/domain`, `/adapters`, `/ports`, `/application` | "Hexagonal architecture pattern" |
| Feature-based Modules | `/features/*`, `/modules/*` with self-contained subfolders | "Feature-based module organization" |
| Clean Architecture | `/entities`, `/usecases`, `/interfaces`, `/frameworks` | "Clean architecture pattern" |
| MVC | `/models`, `/views`, `/controllers` | "MVC pattern" |
| CQRS | `/commands`, `/queries`, `/handlers` | "CQRS pattern" |

**Detection Command:**

```bash
find . -type d -maxdepth 3 | grep -E "(controller|service|repositor|domain|adapter|port|feature|module|entity|usecase|handler|command|query)" | head -20
```

**For each detected pattern, record:**
- Pattern name
- Evidence (which folders exist)
- Root location

### Level 3: Code Patterns

In addition to Levels 1-2, analyze code for design patterns:

**Pattern Detection:**

| Pattern | Search Strategy | Indicators |
|---------|-----------------|------------|
| Dependency Injection | Search for constructor params that are interfaces | `func New.*\(.*Interface\)` |
| Repository Pattern | Search for Repository suffix | `type.*Repository interface` or `class.*Repository` |
| Factory Pattern | Search for Factory suffix or New* funcs returning interfaces | `func New.*\(\).*Interface` |
| Singleton | Search for GetInstance or sync.Once usage | `sync.Once` or `getInstance` |
| Observer/Pub-Sub | Search for Subscribe/Publish patterns | `Subscribe\(`, `Publish\(`, `EventEmitter` |
| Middleware Chain | Search for middleware patterns | `func.*Middleware`, `app.use\(` |

**Detection Commands:**

```bash
# Repository pattern (Go)
grep -r "Repository interface" --include="*.go" | head -5

# Dependency injection (constructor with interface params)
grep -rn "func New" --include="*.go" | grep -E "\(.*[A-Z][a-z]+er\)" | head -5

# Middleware (various languages)
grep -rn "middleware\|Middleware" --include="*.go" --include="*.ts" --include="*.js" | head -5
```

**For each detected pattern, record:**
- Pattern name
- Example file and line
- How it's implemented

## Step 3: Present Findings

After scanning, present findings in this format:

```
Found [N] undocumented architectural decisions at depth level [L]:

Dependencies ([count])
 - [name] ([category])
 - [name] ([category])
 - ...

Structural Patterns ([count])  [only if Level 2+]
 - [pattern name] ([evidence folders])
 - ...

Code Patterns ([count])  [only if Level 3]
 - [pattern name] ([example location])
 - ...

Already documented: [N] decisions ([list IDs])
```

**Filtering:**
- Check existing `@decision.id` annotations in codebase
- Check existing ADR files in output directory
- Mark already-documented items and exclude from count

**Then ask:**
"Which would you like to document? You can say 'all', list numbers, or pick a category like 'just the dependencies'."

## Step 4: Research Alternatives

For each selected decision, research alternatives and considerations before asking the user.

**Research Process:**

1. **Identify the category** - What problem does this technology solve?
2. **Web search for alternatives** - Search: "[technology] alternatives [year]"
3. **Gather selection criteria** - What factors matter when choosing in this category?
4. **Summarize findings** - Present 3-5 alternatives with key differentiators

**Research Prompt Template:**

```
I'll document your choice of [TECHNOLOGY] for [CATEGORY].

I've researched alternatives and considerations:

**Alternatives:**
- [Alt 1] - [Key differentiator, trade-off]
- [Alt 2] - [Key differentiator, trade-off]
- [Alt 3] - [Key differentiator, trade-off]

**Key considerations for [CATEGORY] selection:**
- [Factor 1]
- [Factor 2]
- [Factor 3]
```

**Research Scope:**
- For common technologies: Use knowledge + brief web search to confirm current status
- For obscure/newer tools: More thorough web search for alternatives and community status
