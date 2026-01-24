# Change: Add Full-Text Search for Messages

## Why
As the number of agent messages grows, users need a quick way to find specific messages. A search feature lets users filter messages by content, sender, or workspace name without scrolling through the entire history.

## What Changes
- Add search input field to the chat UI header
- Implement client-side filtering as user types
- Filter matches message body, sender (from/to), and workspace name
- Show/hide messages based on search query with instant feedback

## Impact
- Affected specs: `chatroom-ui` (existing)
- Affected code: `internal/server/templates/index.html`, `internal/server/static/style.css`
- No backend changes required - pure client-side filtering
