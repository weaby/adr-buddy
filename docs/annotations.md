# Annotation Reference

ADR Buddy scans your code for `@decision.*` annotations in comments and generates markdown ADRs.

## Comment Syntax

Use your language's comment style:

| Languages | Syntax |
|-----------|--------|
| Go, JavaScript, TypeScript, Java, C, Rust | `//` |
| Python, Ruby, YAML, Shell | `#` |

## Required Fields

Every annotation must have these fields:

| Field | Description |
|-------|-------------|
| `@decision.id` | Unique identifier (e.g., `adr-001`) |
| `@decision.name` | Short title (under 60 characters) |

```go
// @decision.id: adr-001
// @decision.name: PostgreSQL for primary datastore
```

## Recommended Fields

These fields make ADRs useful:

| Field | Description |
|-------|-------------|
| `@decision.context` | Why was this decision needed? What problem does it solve? |
| `@decision.decision` | What was decided? Include key details. |
| `@decision.alternatives` | What other options were considered and why rejected? |
| `@decision.consequences` | What are the trade-offs? What becomes easier/harder? |

## Optional Fields

| Field | Description | Default |
|-------|-------------|---------|
| `@decision.status` | `proposed`, `accepted`, `rejected`, `deprecated`, `superseded` | `proposed` |
| `@decision.category` | Organizational category (creates subdirectory) | none |
| `@decision.refs` | Related files affected by this decision | none |
| `@decision.supersedes` | ID of decision this replaces | none |

## Multi-line Values

Continue on the next line with indentation:

```go
// @decision.context: We needed a message queue that supports
//   replay capability for debugging payment failures and allows
//   multiple consumers to receive the same events.
```

The indented lines are joined to the first line.

## Listing Items

Use indented list syntax for alternatives and refs:

```go
// @decision.alternatives:
//   - SQS: Simpler but no replay capability
//   - RabbitMQ: Good but team lacks expertise
//   - Redis Streams: Insufficient durability for payments
```

```go
// @decision.refs:
//   - internal/payments/publisher.go
//   - internal/payments/consumer.go
//   - config/kafka.yaml
```

## Complete Example

```go
// @decision.id: adr-003
// @decision.name: Kafka for payment events
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
//   - deployments/kafka.yaml

func NewPaymentEventPublisher(kafka *kafka.Producer) *PaymentEventPublisher {
```

## One ADR, Multiple Locations

The same ADR ID can appear in multiple files. ADR Buddy merges them:

`src/cache/client.go`:
```go
// @decision.id: adr-005
// @decision.name: Redis for caching
// @decision.context: Need distributed cache for API responses.
```

`src/cache/config.go`:
```go
// @decision.id: adr-005
// @decision.name: Redis for caching
// @decision.decision: Use Redis with 1-hour TTL for API responses.
```

Both locations appear in the generated `adr-005.md` under "Code Locations".

## Multiple ADRs in One File

A single file can contain multiple decision annotations:

```python
# @decision.id: adr-010
# @decision.name: FastAPI for web framework
# ...

# @decision.id: adr-011
# @decision.name: Pydantic for validation
# ...
```

Each generates its own ADR file.

## Categories

Use `@decision.category` to organize ADRs into subdirectories:

```go
// @decision.id: adr-001
// @decision.category: infrastructure
```

Output: `decisions/infrastructure/adr-001.md`

Without a category, ADRs are placed in the root output directory.
