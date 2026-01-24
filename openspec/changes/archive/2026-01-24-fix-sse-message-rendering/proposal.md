# Change: Fix SSE message rendering showing raw JSON

## Why
The UI currently displays raw JSON instead of rendered HTML message elements when new messages arrive via Server-Sent Events. This breaks the real-time message display functionality specified in the chatroom-ui capability.

## What Changes
- Modify the SSE handler to render messages as HTML using the message template before sending
- Create a single-message template partial for SSE rendering
- Messages sent via SSE will be properly formatted HTML that HTMX can insert directly into the DOM

## Impact
- Affected specs: chatroom-ui (clarifies SSE message format requirement)
- Affected code: `internal/server/server.go` (handleSSE function), `internal/server/templates/` (new message partial)
