# GitHub Actions Integration Guide

## Overview

ADR Buddy provides two workflows for GitHub Actions:
1. **Validation** - Runs on PRs to validate annotations
2. **Sync** - Runs after merge to create ADR update PRs

## Quick Start

### 1. Add Validation Workflow

Create `.github/workflows/adr-validate.yml`:

```yaml
name: Validate ADRs
on:
  pull_request:
    branches: [main]

jobs:
  validate:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      pull-requests: write
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - uses: weaby/adr-buddy@v1
        with:
          mode: validate
```

### 2. Add Sync Workflow

Create `.github/workflows/adr-sync.yml`:

```yaml
name: Sync ADRs
on:
  push:
    branches: [main]

jobs:
  sync:
    runs-on: ubuntu-latest
    permissions:
      contents: write
      pull-requests: write
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - uses: weaby/adr-buddy@v1
        with:
          mode: sync
```

## Configuration Options

### Validate Mode

| Input | Default | Description |
|-------|---------|-------------|
| `mode` | (required) | Set to `validate` |
| `strict` | `true` | Treat warnings as errors |
| `config-path` | `.adr-buddy/config.yml` | Path to config |
| `version` | `latest` | ADR Buddy version |

### Sync Mode

| Input | Default | Description |
|-------|---------|-------------|
| `mode` | (required) | Set to `sync` |
| `create-pr` | `true` | Create PR with changes |
| `reviewers` | `""` | Comma-separated usernames |
| `token` | `${{ github.token }}` | GitHub token |
| `version` | `latest` | ADR Buddy version |

## Troubleshooting

### Error: "go: command not found"

**Cause:** Go is not installed in the runner.

**Solution:** Add setup-go step before adr-buddy action:

```yaml
- uses: actions/setup-go@v5
  with:
    go-version: '1.24'
```

### Error: "permission denied" when creating PR

**Cause:** Insufficient permissions for GITHUB_TOKEN.

**Solution:** Add required permissions to job:

```yaml
jobs:
  sync:
    permissions:
      contents: write
      pull-requests: write
```

### Sync PR not created

**Possible causes:**

1. No ADR changes detected (expected behavior)
2. Missing permissions
3. `create-pr: false` in configuration

**Debug steps:**

1. Check workflow logs for "No ADR changes detected"
2. Verify permissions in workflow
3. Check action inputs

### Multiple sync PRs created

**Cause:** Multiple rapid merges to main branch (expected behavior).

**Solution:** This is intentional - review and merge PRs in order.

### Validation passes but sync fails

**Cause:** File system or git configuration issues.

**Debug steps:**

1. Check sync workflow logs
2. Verify output directory exists and is writable
3. Ensure git is configured in runner

## Advanced Usage

### Custom Configuration Path

```yaml
- uses: weaby/adr-buddy@v1
  with:
    mode: validate
    config-path: config/adr-buddy.yml
```

### Specific Version

```yaml
- uses: weaby/adr-buddy@v1
  with:
    mode: sync
    version: v1.2.3
```

### Using Action Outputs

```yaml
- uses: weaby/adr-buddy@v1
  id: adr
  with:
    mode: validate

- name: Post Slack notification
  if: steps.adr.outputs.validation-status == 'fail'
  run: |
    curl -X POST $SLACK_WEBHOOK \
      -d '{"text": "ADR validation failed!"}'
```

## Examples

See [examples directory](../examples/github-actions/) for complete workflow examples.

## Support

For issues or questions:
- GitHub Issues: https://github.com/weaby/adr-buddy/issues
- Documentation: https://github.com/weaby/adr-buddy#readme
