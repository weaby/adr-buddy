---
name: adr
description: Automatically document architectural decisions when coding
---

# ADR Buddy

**You MUST use this skill automatically whenever you make architectural decisions while coding.** Do not wait for the user to ask — document decisions as you make them.

## When to Add Decision Annotations

Add a `@decision` annotation whenever you:

| Trigger | Examples |
|---------|----------|
| **Choose a technology** | Database, queue, framework, library |
| **Implement a pattern** | Retry logic, caching, error handling, auth |
| **Make a trade-off** | Performance vs simplicity, consistency vs availability |
| **Configure non-obvious values** | Timeouts, limits, retry counts, thresholds |
| **Structure code significantly** | Module boundaries, API design, data models |
| **Pick between alternatives** | Any time you consider multiple options |

**Rule of thumb:** If a future developer would ask "why is this like this?" — add an annotation.

## Automatic Behavior

When you write code that involves a decision:

1. **Pause and acknowledge:** "This is an architectural decision. Let me document it."
2. **State alternatives you considered**
3. **Add the annotation to the code**
4. **List all files affected by this decision**

Do NOT skip this step. Do NOT wait to be asked.

## Before Creating Annotations

1. Check `.adr-buddy/config.yml` for:
   - `output_dir` - where ADR files are stored
   - `template` - path to custom template (defaults to `internal/template/default.go`)

2. The ADR template format is defined in `internal/template/default.go`. This is the authoritative source for how ADRs are rendered.

3. Scan the output directory to find the highest existing ADR ID:
   ```bash
   ls -1 decisions/ | grep -E '^adr-[0-9]+' | sort -V | tail -1
   ```
   Then use the next number.

4. Check for existing decisions on the same topic to avoid duplicates or conflicts.

## Annotation Syntax

Place annotations in code comments. Use `//` for Go/JS/TS/Java/C/Rust or `#` for Python/Ruby/YAML/Shell.

### Required Fields

| Field | Description |
|-------|-------------|
| `@decision.id` | Unique identifier (e.g., "adr-001") |
| `@decision.name` | Short title (under 60 chars) |
| `@decision.context` | WHY this decision was needed |
| `@decision.decision` | WHAT was decided |

### Recommended Fields

| Field | Description |
|-------|-------------|
| `@decision.alternatives` | Other options considered and why rejected |
| `@decision.consequences` | Trade-offs, what becomes easier/harder |
| `@decision.refs` | Other files affected by this decision |

### Optional Fields

| Field | Description |
|-------|-------------|
| `@decision.status` | proposed, accepted, rejected, deprecated, superseded (default: accepted) |
| `@decision.category` | infrastructure, data, security, architecture, etc. |
| `@decision.supersedes` | ID of decision this replaces (e.g., "adr-002") |

### Multi-line Values

Continue on the next line with indentation:

```go
// @decision.context: We needed a message queue that supports
//   replay capability for debugging payment failures and allows
//   multiple consumers to receive the same events.
```

### Listing Alternatives

```go
// @decision.alternatives:
//   - SQS: Simpler but no replay capability, harder multi-consumer
//   - RabbitMQ: Good but less ecosystem support for our Go stack
//   - Redis Streams: Considered but persistence guarantees unclear
```

### Listing Affected Files

```go
// @decision.refs:
//   - internal/payments/publisher.go
//   - internal/payments/consumer.go
//   - config/kafka.yaml
//   - docker-compose.yml
```

## Complete Example

```go
// @decision.id: adr-003
// @decision.name: Use Kafka for payment events
// @decision.status: accepted
// @decision.category: infrastructure
// @decision.context: We need to publish payment events to multiple
//   consumers (analytics, fraud detection, notifications). Events
//   must be replayable for debugging failed payments.
// @decision.decision: Use Apache Kafka over SQS or RabbitMQ.
//   Configure with 3 partitions and 7-day retention.
// @decision.alternatives:
//   - SQS: Simpler, but no replay and complex multi-consumer setup
//   - RabbitMQ: Good features but team lacks expertise
//   - Redis Streams: Insufficient durability guarantees for payments
// @decision.consequences: Requires Kafka expertise on team. Adds
//   operational complexity. Enables event replay and multi-consumer
//   patterns. 7-day retention allows debugging historical issues.
// @decision.refs:
//   - internal/payments/event_publisher.go
//   - internal/payments/event_consumer.go
//   - internal/analytics/payment_listener.go
//   - deployments/kafka.yaml
func NewPaymentEventPublisher(kafka *kafka.Producer) *PaymentEventPublisher {
```

## Good vs Bad Annotations

### ✅ Good: Explains WHY with alternatives

```go
// @decision.id: adr-007
// @decision.name: 5 retries with exponential backoff for payments
// @decision.context: Payment gateway has occasional timeouts.
//   Need retry logic that handles transient failures without
//   overwhelming the gateway during outages.
// @decision.decision: 5 retries with exponential backoff (1s, 2s,
//   4s, 8s, 16s). Payments are idempotent via payment_id.
// @decision.alternatives:
//   - Fixed delay: Simpler but can cause thundering herd
//   - 3 retries: Insufficient for p99 recovery based on metrics
//   - Circuit breaker only: Doesn't handle single-request failures
// @decision.consequences: Max 31s latency for retried requests.
//   Safe because payment_id ensures idempotency.
```

### ❌ Bad: No context, no alternatives

```go
// @decision.id: adr-007
// @decision.name: Use retries
// @decision.decision: Retry 5 times
```

### ❌ Bad: States WHAT but not WHY

```go
// @decision.id: adr-007
// @decision.name: Use Kafka
// @decision.decision: We are using Kafka for events
```

## When NOT to Add Annotations

Skip annotations for:
- Trivial implementation details (loop styles, variable names)
- Standard library usage with no alternatives
- Obvious choices anyone would make
- Temporary or experimental code

## Checking Existing Decisions

Before making a decision in an area with existing ADRs:

```bash
# Search for related decisions
adr-buddy list | grep -i "kafka\|queue\|messaging"

# Or search annotations in code
grep -r "@decision" --include="*.go" | grep -i "kafka"
```

If a relevant decision exists:
- **Follow it** if still applicable
- **Supersede it** if requirements changed (add `@decision.supersedes: adr-XXX`)

## After Adding Annotations

Remind the user to sync:

```bash
# Validate annotations
adr-buddy check

# Generate/update ADR files
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
> Let me add the annotation and implement it."

Then add the annotation with full context, alternatives, and refs.