# Change: Add repository-specific PR links

## Why
When users mention PR numbers like `#123` in multiclaude workspaces, they expect these to link directly to the actual repository's pull request. Previously, these linked to a generic GitHub search which required an extra click and could return ambiguous results.

## What Changes
- `linkifyIssueRefs()` accepts an optional repo URL parameter
- `renderMarkdown()` passes repo URL for proper link generation
- Server loads workspace-to-repository URL mapping from `~/.multiclaude/state.json`
- PR number references now link directly to the repository's PR page (e.g., `https://github.com/owner/repo/pull/123`)
- Graceful fallback to GitHub search when repo URL is unavailable

## Impact
- Affected specs: message-rendering (new capability)
- Affected code: `cmd/agent-chat/main.go`, `internal/server/server.go`, `internal/server/templates/`
