#!/usr/bin/env bash
# Validate mode script for ADR Buddy GitHub Action
# Parses check and sync preview results and posts PR comment

set -euo pipefail

# Source common utilities
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/common.sh"

# Read check results
CHECK_RESULT=$(cat /tmp/adr-check-result.json)
SYNC_PREVIEW=$(cat /tmp/adr-sync-preview.json)

# Extract status
STATUS=$(echo "$CHECK_RESULT" | jq -r '.status')
ERROR_COUNT=$(echo "$CHECK_RESULT" | jq -r '.summary.error_count')
WARNING_COUNT=$(echo "$CHECK_RESULT" | jq -r '.summary.warning_count')
TOTAL_ANNOTATIONS=$(echo "$CHECK_RESULT" | jq -r '.summary.total_annotations')
CHANGES_DETECTED=$(echo "$SYNC_PREVIEW" | jq -r '.changes_detected')

# Set outputs
set_output "validation-status" "$STATUS"
set_output "changes-detected" "$CHANGES_DETECTED"

# Build comment
COMMENT="## ADR Buddy Validation Results\n\n"

# Add validation section
if [ "$STATUS" = "pass" ]; then
  COMMENT+="### ✅ Validation Passed\n\n"
  COMMENT+="${TOTAL_ANNOTATIONS} annotations found, all valid.\n\n"
elif [ "$STATUS" = "warning" ]; then
  COMMENT+="### ⚠️ Validation Passed with Warnings\n\n"
  COMMENT+="${TOTAL_ANNOTATIONS} annotations found with ${WARNING_COUNT} warning(s).\n\n"

  # List warnings
  WARNINGS=$(echo "$CHECK_RESULT" | jq -r '.warnings[] | "- `\(.file):\(.line)` - \(.message)"')
  if [ -n "$WARNINGS" ]; then
    COMMENT+="**Warnings:**\n${WARNINGS}\n\n"
  fi
else
  COMMENT+="### ❌ Validation Failed\n\n"
  COMMENT+="${ERROR_COUNT} error(s) found.\n\n"

  # List errors
  ERRORS=$(echo "$CHECK_RESULT" | jq -r '.errors[] | "- `\(.file):\(.line)` - \(.message)"')
  if [ -n "$ERRORS" ]; then
    COMMENT+="**Errors:**\n${ERRORS}\n\n"
  fi
fi

# Add sync preview section if changes detected
if [ "$CHANGES_DETECTED" = "true" ]; then
  COMMENT+="### ⚠️ ADRs Need Updating\n\n"
  COMMENT+="The following ADR files will be updated after merge:\n\n"

  # List created files
  CREATED=$(echo "$SYNC_PREVIEW" | jq -r '.files.created[]?' 2>/dev/null || true)
  if [ -n "$CREATED" ]; then
    COMMENT+="**Created:**\n"
    while IFS= read -r file; do
      # Get ADR name for this file
      ADR_NAME=$(echo "$SYNC_PREVIEW" | jq -r ".adrs[] | select(.file_path == \"$file\") | .name")
      COMMENT+="- \`${file}\` - ${ADR_NAME}\n"
    done <<< "$CREATED"
    COMMENT+="\n"
  fi

  # List modified files
  MODIFIED=$(echo "$SYNC_PREVIEW" | jq -r '.files.modified[]?' 2>/dev/null || true)
  if [ -n "$MODIFIED" ]; then
    COMMENT+="**Modified:**\n"
    while IFS= read -r file; do
      ADR_NAME=$(echo "$SYNC_PREVIEW" | jq -r ".adrs[] | select(.file_path == \"$file\") | .name")
      COMMENT+="- \`${file}\` - ${ADR_NAME}\n"
    done <<< "$MODIFIED"
    COMMENT+="\n"
  fi
fi

# Add footer
COMMENT+="---\n"
COMMENT+="🤖 [ADR Buddy](https://github.com/weaby/adr-buddy)"

# Post comment
echo -e "$COMMENT"
post_pr_comment "$(echo -e "$COMMENT")"

# Exit with appropriate code
if [ "$STATUS" = "fail" ]; then
  exit 1
fi

exit 0
