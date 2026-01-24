# Change: Add Chatroom UI for Agent Messages

## Why
Agents communicate via filesystem-based mail (Gas Town `gt mail` and multiclaude filesystem messages). Currently there's no visual way to see these conversations. A Slack/Discord-style chat interface would make it easy to monitor agent communication across workspaces.

## What Changes
- Add Go backend server that watches mail directories for new messages
- Add HTMX-based chat UI showing channels (workspaces) and messages
- Support both Gas Town mail format and multiclaude filesystem messages
- Use Server-Sent Events (SSE) for real-time message updates
- No authentication required (local-only tool)

## Impact
- Affected specs: `chatroom-ui`, `message-sources`, `workspace-channels` (all new)
- Affected code: New Go server + HTML/CSS/HTMX templates
- External dependencies: Go stdlib, `fsnotify` for filesystem watching
