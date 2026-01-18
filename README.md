# ADR Buddy

A language-agnostic CLI tool that generates Architecture Decision Records (ADRs) from code annotations.

## Features

- 📝 Write ADRs directly in your code using simple comment annotations
- 🔄 Automatically generate and sync ADR documentation
- 🌍 Language-agnostic: works with Go, JavaScript, TypeScript, Python, Ruby, and more
- 🎯 Smart merge: preserves manual edits while updating from annotations
- 📁 Flexible organization: optional categorization of ADRs
- ⚡ Zero configuration required to get started

## Installation

```bash
go install github.com/weaby/adr-buddy/cmd/adr-buddy@latest
```

## Quick Start

1. Initialize in your project:
```bash
adr-buddy init
```

2. Add annotations to your code:
```javascript
// @decision.id: adr-1
// @decision.name: Using Pino for logging
// @decision.status: accepted
// @decision.category: infrastructure
// @decision.context: We needed structured logging with low overhead
//   for our high-throughput API.
// @decision.decision: Adopt Pino as the standard logging library.
// @decision.consequences: Better performance but team needs training.

const pino = require('pino');
```

3. Generate ADR documentation:
```bash
adr-buddy sync
```

Your ADR will be created at `decisions/infrastructure/adr-1.md`!

## Commands

- `adr-buddy init` - Initialize configuration in current directory
- `adr-buddy sync` - Scan code and generate/update ADR files
- `adr-buddy check` - Validate annotations without generating files
- `adr-buddy list` - List all discovered ADRs

## Annotation Reference

**Required:**
- `@decision.id` - Unique identifier (e.g., "adr-1")
- `@decision.name` - Decision title

**Optional:**
- `@decision.status` - proposed | accepted | rejected | deprecated | superseded
- `@decision.category` - Organizational category
- `@decision.context` - Background and rationale (multi-line)
- `@decision.decision` - The choice made (multi-line)
- `@decision.consequences` - Outcomes and impacts (multi-line)

## Multi-line Values

Use indentation for continuation:
```javascript
// @decision.context: This is the first line
//   and this continues on the next line
//   and this is the third line.
```

## Configuration

Optional `.adr-buddy/config.yml`:
```yaml
scan_paths:
  - ./src
  - ./lib
output_dir: ./decisions
exclude:
  - "**/node_modules/**"
  - "**/vendor/**"
```

## Examples

### One-to-Many: Same ADR Referenced in Multiple Files

`src/logger.js`:
```javascript
// @decision.id: adr-1
// @decision.name: Using Pino
// @decision.context: Need fast logging
const logger = require('pino')();
```

`src/monitor.js`:
```javascript
// @decision.id: adr-1
// @decision.name: Using Pino
// @decision.decision: Use Pino everywhere
const logger = require('./logger');
```

Both annotations merge into a single `decisions/adr-1.md` with two code locations.

### Many-to-One: Multiple ADRs in One File

```python
# @decision.id: adr-5
# @decision.name: Microservices architecture
# ...

# @decision.id: adr-12
# @decision.name: gRPC for inter-service communication
# ...
```

Each annotation generates its own ADR file, both referencing the same source file.

## License

MIT
