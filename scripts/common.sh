#!/usr/bin/env bash
# Common utilities for GitHub Action scripts

set -euo pipefail

# Post a comment on a pull request
# Args: $1 = comment body (markdown)
post_pr_comment() {
  local comment_body="$1"

  # Check if we're in a pull request context
  if [ -z "${GITHUB_EVENT_NAME:-}" ] || [ "$GITHUB_EVENT_NAME" != "pull_request" ]; then
    echo "Not in a pull request context, skipping comment"
    return 0
  fi

  # Get PR number from event
  local pr_number
  pr_number=$(jq -r '.pull_request.number' "$GITHUB_EVENT_PATH")

  if [ -z "$pr_number" ] || [ "$pr_number" = "null" ]; then
    echo "Could not determine PR number, skipping comment"
    return 0
  fi

  # Post comment using GitHub API
  local repo="${GITHUB_REPOSITORY}"
  local api_url="https://api.github.com/repos/${repo}/issues/${pr_number}/comments"

  curl -X POST \
    -H "Authorization: token ${GITHUB_TOKEN}" \
    -H "Accept: application/vnd.github.v3+json" \
    -d "{\"body\": $(jq -Rs . <<< "$comment_body")}" \
    "$api_url" > /dev/null

  echo "Posted comment to PR #${pr_number}"
}

# Set action output
# Args: $1 = key, $2 = value
set_output() {
  local key="$1"
  local value="$2"
  echo "${key}=${value}" >> /tmp/adr-action-outputs.env
}
