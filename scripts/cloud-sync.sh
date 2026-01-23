#!/bin/bash
set -e

CLOUD_URL="${1:-https://cloud.adr-buddy.dev}"
CLOUD_TOKEN="${2}"
REPO_NAME="${3}"
REPO_URL="${4}"
BRANCH="${5}"
COMMIT_SHA="${6}"
SYNC_OUTPUT="${7}"

if [ -z "$CLOUD_TOKEN" ]; then
  echo "No cloud token provided, skipping cloud sync"
  exit 0
fi

echo "🌥️  Syncing ADRs to ADR Buddy Cloud..."

# Build JSON payload
PAYLOAD=$(cat <<EOF
{
  "workspace_token": "$CLOUD_TOKEN",
  "repo": {
    "name": "$REPO_NAME",
    "url": "$REPO_URL",
    "branch": "$BRANCH",
    "commit_sha": "$COMMIT_SHA"
  },
  "sync_timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "adrs": $(cat "$SYNC_OUTPUT")
}
EOF
)

# Send to Cloud API
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$CLOUD_URL/api/v1/sync" \
  -H "Content-Type: application/json" \
  -d "$PAYLOAD")

HTTP_CODE=$(echo "$RESPONSE" | tail -n 1)
BODY=$(echo "$RESPONSE" | head -n -1)

if [ "$HTTP_CODE" -eq 200 ]; then
  echo "✅ Successfully synced to Cloud"
  echo "$BODY" | jq '.' || echo "$BODY"
else
  echo "⚠️  Cloud sync failed with HTTP $HTTP_CODE"
  echo "$BODY"
  echo "Continuing workflow (cloud sync is optional)"
fi
