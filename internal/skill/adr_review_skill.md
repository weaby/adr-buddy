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

## Step 4: Draft ADRs

For each selected decision, **you** research and write a complete draft ADR. Do NOT ask open-ended questions. Instead, do the work and present a finished draft for approval.

### 4.1 Research (you do this silently)

For each decision:

1. **Web search** for "[technology] vs alternatives [year]" to find current alternatives
2. **Analyze the codebase** for how the technology is used, configured, and integrated
3. **Infer the context** from the project structure, config files, and usage patterns

Use what you find to write the draft. Do NOT ask the user "why did you choose X?" — infer it from evidence and let them correct you.

### 4.2 Write a Complete Draft

For each decision, write a full ADR draft with all sections filled in:

```markdown
---
adr_id: adr-[NEXT_ID]
name: [Short title]
status: accepted
category: [category]
date: "[YYYY-MM-DD]"
globs:
    - "[relevant file patterns based on where technology is used]"
---

# ADR-[NEXT_ID]: [Short title]

## Context

[Write 2-4 sentences based on what you found in the codebase.
What problem does this solve? What requirements drove this choice?
Infer from usage patterns, config, and project structure.]

## Decision

[Write 1-3 sentences. What was chosen and how is it configured?
Be specific — reference actual config values, versions, or patterns
you found in the codebase.]

## Alternatives Considered

- **[Alt 1]:** [Why it's a reasonable alternative, why not chosen — infer from project context]
- **[Alt 2]:** [Same]
- **[Alt 3]:** [Same]

## Consequences

**Positive:** [Infer from how it's used — what does it enable?]
**Negative:** [Infer from common trade-offs for this technology]
```

### 4.3 Present for Review

Present the draft and ask ONE structured question:

```
Here's the draft ADR for [TECHNOLOGY]:

[Show full draft]

What would you like to change?
1. Looks good — create it
2. Edit context/reasoning (I'll explain what to change)
3. Different alternatives were considered
4. Skip this one
```

**IMPORTANT:** Always present a complete, ready-to-save draft. Never present a skeleton with `[fill in]` placeholders. The user should only need to approve or make small corrections.

### 4.4 Apply Corrections

If the user picks option 2 or 3:
- Apply their feedback to the draft
- Show the updated version
- Ask again: "Updated. Ready to create, or more changes?"

Limit to 2 revision rounds. If still not right after 2 rounds, create the file and let the user edit it directly.

## Step 5: Create ADR Files

For each approved decision:

1. Create the file:
   ```bash
   adr-buddy new "[Decision title]" --category [category]
   ```

2. Replace the generated template content with the approved draft content by editing the file directly.

3. Move to the next decision immediately — don't wait for additional confirmation.

## Step 6: Check for Duplicates

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

## Step 7: Sync and Complete

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
