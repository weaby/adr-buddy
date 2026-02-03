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
   test -d .claude/rules/decisions && echo "Ready" || echo "Run: adr-buddy init"
   ```

2. Read `.adr-buddy/config.yml` to find:
   - `decisions_dir` - where ADR files are stored

3. Check existing ADRs:
   ```bash
   adr-buddy list
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
- Check existing ADR files in decisions directory
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

## Step 5: Guided Documentation Conversation

For each selected decision, follow this conversation flow:

### 5.1 Present Research (from Step 4)

Show the alternatives and considerations you researched.

### 5.2 Ask "Why"

```
Why did you choose [TECHNOLOGY] for this project?
```

Wait for response. This becomes the core of the Context and Decision sections.

### 5.3 Explore Constraints

```
What constraints or factors influenced this decision?
(e.g., team expertise, performance requirements, existing infrastructure, cost, timeline)
```

Wait for response. This enriches the Context section.

### 5.4 Capture Trade-offs

```
Any trade-offs or concerns you're aware of with this choice?
```

Wait for response. This becomes the Consequences section.

### 5.5 Confirm Alternatives

```
From the alternatives I listed, were any of these seriously considered?
[List the alternatives you researched]
```

Wait for response. This becomes the Alternatives Considered section.

### 5.6 Draft ADR

Synthesize all responses into a draft ADR file:

```markdown
---
adr_id: adr-[NEXT_ID]
name: [Short title from technology + category]
status: accepted
category: [category]
date: "[YYYY-MM-DD]"
globs:
    - "[relevant file patterns]"
---

# ADR-[NEXT_ID]: [Short title]

## Context

[Synthesized from user's "why" and constraints]

## Decision

[What was chosen and key configuration]

## Alternatives Considered

- **[Alt 1]:** [Why not chosen]
- **[Alt 2]:** [Why not chosen]

## Consequences

**Positive:** [Benefits]
**Negative:** [From trade-offs discussion]
```

Show the draft and ask: "Does this look right? Any changes?"

Wait for confirmation before writing the file.

## Step 6: Create ADR Files

For each confirmed decision:

```bash
adr-buddy new "[Decision title]" --category [category]
```

Then edit the generated file to fill in the sections with the content from Step 5.

## Step 7: Check for Duplicates

Before creating each ADR, check for existing ADRs on the same topic:

```bash
adr-buddy list | grep -i "[TECHNOLOGY]"
```

**If duplicate found:**
```
I found an existing ADR that might cover this:
- [ADR-XXX]: [Title]

Options:
1. Skip - this is already documented
2. Update - add more context to the existing ADR
3. Create new - this is a different aspect of the same technology
```

## Step 8: Sync and Complete

After all selected decisions are documented:

```bash
# Validate ADR files
adr-buddy check

# Regenerate the decisions index
adr-buddy sync
```

**Report completion:**
```
Created [N] new ADR files:

- .claude/rules/decisions/adr-[ID]-[name].md ([Title])
- .claude/rules/decisions/adr-[ID]-[name].md ([Title])
...

Index updated. All decisions documented and ready for commit.
```
