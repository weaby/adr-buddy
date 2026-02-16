---
name: adr
description: Automatically document architectural decisions when coding
---

# ADR Buddy

**You MUST use this skill automatically whenever you make architectural decisions while coding.** Do not wait for the user to ask — document decisions as you make them.

## When to Create an ADR

Create an ADR whenever you:

| Trigger | Examples |
|---------|----------|
| **Choose a technology** | Database, queue, framework, library |
| **Implement a pattern** | Retry logic, caching, error handling, auth |
| **Make a trade-off** | Performance vs simplicity, consistency vs availability |
| **Configure non-obvious values** | Timeouts, limits, retry counts, thresholds |
| **Structure code significantly** | Module boundaries, API design, data models |
| **Pick between alternatives** | Any time you consider multiple options |

**Rule of thumb:** If a future developer would ask "why is this like this?" — create an ADR.

## Automatic Behavior

When you write code that involves a decision:

1. **Pause and acknowledge:** "This is an architectural decision. Let me document it."
2. **State alternatives you considered**
3. **Create the ADR file**
4. **List all files affected by this decision**

Do NOT skip this step. Do NOT wait to be asked.

## Before Creating ADRs

1. **Read `.adr-buddy/template.md`** — this is the authoritative template for ADR structure. The user may have customized sections, added fields, or changed placeholders. You MUST follow it exactly when filling in ADR content. If the template has custom sections (e.g., "Locations", "Impact"), fill those in too.

2. Check `.adr-buddy/config.yml` for:
   - `decisions_dir` - where ADR files are stored (default: `.claude/rules/decisions`)

3. Check existing decisions to find the next ID and avoid duplicates:
   ```bash
   adr-buddy list
   ```

## ADR File Format

ADRs are markdown files with YAML frontmatter, stored in `.claude/rules/decisions/`.

### Frontmatter Fields

| Field | Required | Description |
|-------|----------|-------------|
| `adr_id` | Yes | Unique identifier (e.g., "adr-001") |
| `name` | Yes | Short title (under 60 chars) |
| `status` | No | proposed, accepted, rejected, deprecated, superseded (default: proposed) |
| `category` | No | infrastructure, data, security, architecture, etc. |
| `date` | No | Date in YYYY-MM-DD format (auto-set) |
| `globs` | No | File patterns this decision applies to (e.g., `["internal/database/**"]`) |

### Markdown Sections

| Section | Description |
|---------|-------------|
| **Context** | WHY this decision was needed |
| **Decision** | WHAT was decided |
| **Alternatives Considered** | Other options and why they were rejected |
| **Consequences** | Trade-offs, what becomes easier/harder |
| **References** | Affected files and relevant code locations |

**Note:** The template in `.adr-buddy/template.md` may define additional or different sections. Always follow the template.

## Creating an ADR

Use `adr-buddy new` to create an ADR with the next available ID:

```bash
adr-buddy new "Use Kafka for payment events" --category infrastructure
```

This creates a file like `.claude/rules/decisions/adr-003-use-kafka-for-payment-events.md` with the template structure. Then fill in the sections.

## Complete Example

```markdown
---
adr_id: adr-003
name: Use Kafka for payment events
status: accepted
category: infrastructure
date: "2024-06-15"
globs:
    - "internal/payments/**"
    - "deployments/kafka.yaml"
---

# ADR-003: Use Kafka for payment events

## Context

We need to publish payment events to multiple consumers (analytics, fraud
detection, notifications). Events must be replayable for debugging failed
payments.

## Decision

Use Apache Kafka over SQS or RabbitMQ. Configure with 3 partitions and
7-day retention.

## Alternatives Considered

- **SQS:** Simpler, but no replay and complex multi-consumer setup
- **RabbitMQ:** Good features but team lacks expertise
- **Redis Streams:** Insufficient durability guarantees for payments

## Consequences

**Positive:** Enables event replay and multi-consumer patterns. 7-day
retention allows debugging historical issues.
**Negative:** Requires Kafka expertise on team. Adds operational complexity.

## References

- `internal/payments/event_publisher.go` — Kafka producer setup
- `internal/payments/event_consumer.go` — Consumer group configuration
- `internal/analytics/payment_listener.go` — Analytics consumer
- `deployments/kafka.yaml` — Kafka cluster configuration
- `docker-compose.yml` — Local Kafka service definition
```

## Good vs Bad ADRs

### Good: Explains WHY with alternatives

```markdown
## Context

Payment gateway has occasional timeouts. Need retry logic that handles
transient failures without overwhelming the gateway during outages.

## Decision

5 retries with exponential backoff (1s, 2s, 4s, 8s, 16s). Payments are
idempotent via payment_id.

## Alternatives Considered

- **Fixed delay:** Simpler but can cause thundering herd
- **3 retries:** Insufficient for p99 recovery based on metrics
- **Circuit breaker only:** Doesn't handle single-request failures

## Consequences

**Positive:** Safe because payment_id ensures idempotency.
**Negative:** Max 31s latency for retried requests.
```

### Bad: No context, no alternatives

```markdown
## Decision

Retry 5 times
```

### Bad: States WHAT but not WHY

```markdown
## Decision

We are using Kafka for events
```

## When NOT to Create ADRs

Skip ADRs for:
- Trivial implementation details (loop styles, variable names)
- Standard library usage with no alternatives
- Obvious choices anyone would make
- Temporary or experimental code

## Using Globs

Use the `globs` field to scope an ADR to specific files. Claude Code will load these decisions automatically when working in matching paths.

```yaml
globs:
    - "internal/database/**"
    - "migrations/**"
```

## After Creating ADRs

Sync the index:

```bash
# Validate ADR files
adr-buddy check

# Regenerate the decisions index
adr-buddy sync

# View all decisions
adr-buddy list
```

## Example Workflow

**User:** "Add a caching layer for user profiles"

**You (thinking):** This requires choosing a cache technology — this is an architectural decision.

**You (response):**
> "I'll add a caching layer for user profiles. This is an architectural decision, so let me document it.
>
> **Options I'm considering:**
> - Redis: Persistent, supports TTL, good for sessions too
> - Memcached: Simpler, but no persistence
> - In-memory (local): Fast but doesn't share across instances
>
> I'll go with **Redis** because we need cache sharing across multiple API instances and may want to use it for sessions later.
>
> Let me create the ADR and implement it."

Then create the ADR file with full context, alternatives, and consequences.
