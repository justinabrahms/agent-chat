# Change: Add keyboard shortcuts for navigation

## Why
Power users expect keyboard navigation for efficient browsing of chat interfaces. Currently the UI requires mouse interaction for all navigation.

## What Changes
- Add `j/k` keys for navigating between messages
- Add `/` key to focus search (when search exists)
- Add `g+w` chord to jump focus to workspace list
- Vanilla JS implementation with no dependencies

## Impact
- Affected specs: `keyboard-shortcuts` (new capability)
- Affected code: `internal/server/templates/index.html`
