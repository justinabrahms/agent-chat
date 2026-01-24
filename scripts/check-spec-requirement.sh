#!/bin/bash
# check-spec-requirement.sh
# Ensures substantial Go code changes are accompanied by specification documents
#
# "Substantial" is defined as:
# - More than 50 lines of Go code changed (additions + deletions), OR
# - More than 3 non-test Go files changed
#
# This check passes if:
# - Changes are not substantial, OR
# - Changes include openspec files (proposal, spec, or tasks)

set -e

# Configuration thresholds
LINE_THRESHOLD=50
FILE_THRESHOLD=3

# Get the base branch (default to origin/main for PRs)
BASE_BRANCH="${GITHUB_BASE_REF:-origin/main}"

echo "Checking for spec requirement against ${BASE_BRANCH}..."

# Ensure we have the base branch
git fetch origin main --depth=1 2>/dev/null || true

# Get changed files
CHANGED_FILES=$(git diff --name-only "${BASE_BRANCH}"...HEAD 2>/dev/null || git diff --name-only "${BASE_BRANCH}" HEAD)

# Filter for Go files
GO_FILES=$(echo "$CHANGED_FILES" | grep '\.go$' || true)
GO_NON_TEST_FILES=$(echo "$GO_FILES" | grep -v '_test\.go$' || true)

# Count Go files changed
GO_FILE_COUNT=$(echo "$GO_NON_TEST_FILES" | grep -c '.' || echo 0)

# Get line changes for Go files only
LINE_CHANGES=0
if [ -n "$GO_FILES" ]; then
    # Get stats for each Go file
    for file in $GO_FILES; do
        if [ -f "$file" ]; then
            STATS=$(git diff --numstat "${BASE_BRANCH}"...HEAD -- "$file" 2>/dev/null || git diff --numstat "${BASE_BRANCH}" HEAD -- "$file")
            ADDED=$(echo "$STATS" | awk '{print $1}' | grep -v '-' || echo 0)
            DELETED=$(echo "$STATS" | awk '{print $2}' | grep -v '-' || echo 0)
            LINE_CHANGES=$((LINE_CHANGES + ADDED + DELETED))
        fi
    done
fi

# Check for openspec changes
SPEC_CHANGES=$(echo "$CHANGED_FILES" | grep -E '^openspec/(changes|specs)/' || true)
HAS_SPEC_CHANGES=false
if [ -n "$SPEC_CHANGES" ]; then
    HAS_SPEC_CHANGES=true
fi

echo "Go files changed (non-test): $GO_FILE_COUNT"
echo "Lines changed in Go files: $LINE_CHANGES"
echo "Has spec changes: $HAS_SPEC_CHANGES"

# Determine if changes are substantial
IS_SUBSTANTIAL=false
REASON=""

if [ "$LINE_CHANGES" -gt "$LINE_THRESHOLD" ]; then
    IS_SUBSTANTIAL=true
    REASON="More than ${LINE_THRESHOLD} lines of Go code changed (${LINE_CHANGES} lines)"
fi

if [ "$GO_FILE_COUNT" -gt "$FILE_THRESHOLD" ]; then
    IS_SUBSTANTIAL=true
    if [ -n "$REASON" ]; then
        REASON="${REASON}, and more than ${FILE_THRESHOLD} Go files changed (${GO_FILE_COUNT} files)"
    else
        REASON="More than ${FILE_THRESHOLD} Go files changed (${GO_FILE_COUNT} files)"
    fi
fi

# Make decision
if [ "$IS_SUBSTANTIAL" = true ]; then
    echo ""
    echo "Substantial Go changes detected: $REASON"

    if [ "$HAS_SPEC_CHANGES" = true ]; then
        echo "Specification documents found. Check passed."
        exit 0
    else
        echo ""
        echo "ERROR: Substantial Go changes require specification documents."
        echo ""
        echo "Please add or update specs under openspec/changes/ or openspec/specs/"
        echo "See openspec/AGENTS.md for guidance on creating change proposals."
        exit 1
    fi
else
    echo "Changes are not substantial. Spec documents not required."
    exit 0
fi
