# Change: Add Unread Count in Browser Tab Title

## Why
Users cannot see at a glance if they have unread messages when the Agent Chat tab is not in focus. Displaying the unread count in the browser tab title (e.g., "(3) Agent Chat") provides immediate visibility without switching tabs.

## What Changes
- Add JavaScript logic to track unread messages across all workspaces
- Update `document.title` dynamically when unread count changes
- Integrate with existing SSE message handling to update count in real-time
- Reset count when all workspaces are viewed

## Impact
- Affected specs: browser-notifications (new capability)
- Affected code: `internal/server/templates/index.html` (JavaScript)
