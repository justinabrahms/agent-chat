# Change: Add Claude Teams message source

## Why
Claude Code now supports "Agent Teams" with multiple Claude Code instances coordinating via shared task lists and mailbox-based messaging. Users need to follow these team conversations in Agent Chat alongside Gas Town and Multiclaude sources.

## What Changes
- Add new `ClaudeTeamsSource` implementing the `Source` interface
- Read messages from `~/.claude/teams/{team-name}/inboxes/{agent}.json`
- Detect idle notifications and render them as dimmed status messages
- Use Robohash avatars (same as multiclaude) with a team icon
- Add `claude-teams-dir` config/flag/env var

## Impact
- Affected specs: message-sources
- Affected code: `internal/message/claudeteams.go` (new), `internal/config/config.go`, `cmd/agent-chat/main.go`, `internal/server/server.go`, templates, CSS
