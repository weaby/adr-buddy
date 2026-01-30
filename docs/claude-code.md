# Claude Code Integration

ADR Buddy includes two Claude Code skills that make documentation effortless.

## Installation

Skills are installed during `adr-buddy init`:

```
Would you like to install the Claude Code skill?
  [1] Project-level (.claude/skills/adr.md) - for this project only
  [2] User-level (~/.claude/skills/adr.md) - available in all projects
  [3] Skip - don't install the skill
```

Or use the flag:

```bash
adr-buddy init --claude-skill=project  # This project only
adr-buddy init --claude-skill=user     # All projects
adr-buddy init --claude-skill=skip     # Don't install
```

**Installed skills:**

| Location | Files |
|----------|-------|
| Project-level | `.claude/skills/adr.md`, `.claude/skills/adr-review.md` |
| User-level | `~/.claude/skills/adr.md`, `~/.claude/skills/adr-review.md` |

---

## /adr — Document as You Code

This skill teaches Claude to automatically document architectural decisions while coding.

### When It Activates

Claude uses this skill when you:

- Choose a technology (database, framework, library)
- Implement a pattern (retry logic, caching, error handling)
- Make a trade-off (performance vs simplicity)
- Configure non-obvious values (timeouts, limits, thresholds)
- Structure code significantly (module boundaries, API design)

### How It Works

When Claude makes an architectural decision, it will:

1. Pause and acknowledge the decision
2. State alternatives considered
3. Add the annotation to the code
4. List affected files

### Example

**You:** "Add a caching layer for user profiles"

**Claude:**
> I'll add a caching layer for user profiles. This is an architectural decision, so let me document it.
>
> **Options I'm considering:**
> - Redis: Persistent, supports TTL, good for sessions too
> - Memcached: Simpler, but no persistence
> - In-memory: Fast but doesn't share across instances
>
> I'll go with **Redis** because we need cache sharing across API instances.

Claude then adds the annotation:

```go
// @decision.id: adr-004
// @decision.name: Redis for user profile caching
// @decision.status: accepted
// @decision.context: Need distributed caching for user profiles
//   to reduce database load across multiple API instances.
// @decision.decision: Use Redis with 15-minute TTL for user profiles.
// @decision.alternatives:
//   - Memcached: Simpler but no persistence if restart needed
//   - In-memory: Fast but doesn't share across instances
// @decision.consequences: Adds Redis dependency. Enables cache
//   sharing across instances. May use for sessions later.

func NewUserCache(redis *redis.Client) *UserCache {
```

### Trigger Keywords

The skill also activates when you say:

- "document this decision"
- "create an ADR"
- "architectural decision"
- "why did we choose"
- "record this choice"

---

## /adr-review — Discover Undocumented Decisions

This skill scans your codebase for decisions that were never documented.

### When to Use

- Starting a new project and want to establish ADR discipline
- Joining an existing project and want to understand its architecture
- Periodic "ADR health check" to catch undocumented decisions
- Before major refactoring to understand current state

### How It Works

1. **Choose depth** — How thoroughly to analyze

| Level | What It Discovers |
|-------|-------------------|
| 1 - Dependencies | Package manifests + config files |
| 2 - Structure | Level 1 + architectural patterns from folder organization |
| 3 - Patterns | Level 2 + design patterns visible in code |

2. **Scan** — Analyzes files based on selected depth

3. **Present findings** — Shows categorized summary:

```
Found 12 undocumented architectural decisions at depth level 2:

Dependencies (7)
 - pino (logging)
 - postgresql (database)
 - redis (caching)
 - express (web framework)
 - ...

Structural Patterns (3)
 - Layered architecture (/controllers, /services, /repositories)
 - Centralized error handling (/middleware/errors)
 - Feature-based modules (/features/auth, /features/billing)

Already documented: 3 decisions (ADR-001, ADR-002, ADR-003)
```

4. **Select** — You choose which to document ("all", specific numbers, or categories)

5. **Research** — Claude researches alternatives for each technology

6. **Guided conversation** — For each decision:
   - Why did you choose this?
   - What constraints influenced it?
   - Any trade-offs you're aware of?

7. **Annotate** — Claude adds annotations to your code

8. **Sync** — Runs `adr-buddy sync` to generate ADR files

### Example Session

**You:** `/adr-review`

**Claude:** "How deep should I analyze? (1) Dependencies only, (2) + Structure, (3) Full analysis"

**You:** "2"

**Claude:** *scans and presents findings*

**You:** "Document postgresql and the layered architecture"

**Claude:**
> Let's document PostgreSQL first.
>
> I've researched alternatives:
> - MySQL: Widely adopted, simpler replication
> - SQLite: Embedded, serverless
> - MongoDB: Document model, flexible schema
>
> Why did you choose PostgreSQL for this project?

**You:** "Team expertise and we needed JSONB support for flexible metadata"

**Claude:** *continues conversation, then adds annotation*

---

## Best Practices

### Let Claude Document Proactively

Don't wait to be asked. The `/adr` skill is designed to document decisions as they happen. If Claude isn't documenting automatically, remind it:

> "Remember to document architectural decisions as you make them"

### Review Generated ADRs

Claude's annotations are a starting point. Review and refine:

- Add context Claude might have missed
- Clarify trade-offs specific to your situation
- Link related decisions

### Run /adr-review Periodically

Set a reminder to run `/adr-review` monthly or before major releases. It catches decisions that slipped through.

### Commit Skills to Repository

If using project-level installation, commit `.claude/skills/` so team members get the skills automatically.
