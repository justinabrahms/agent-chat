# Change: Add 'mc-' prefix to multiclaude workspaces

## Why
Currently, repositories in `~/.multiclaude/messages/<repo>/` create workspaces named just `<repo>`. This can lead to naming conflicts with other message sources and makes it unclear which workspaces come from multiclaude versus other sources. Adding an `mc-` prefix makes multiclaude workspaces clearly identifiable and prevents potential collisions.

## What Changes
- Workspace names from multiclaude source change from `<repo>` to `mc-<repo>`
- Both `List()` and `Watch()` functions in `internal/message/multiclaude.go` are updated to use the new prefix
- **BREAKING**: Existing workspace references using the old naming scheme will no longer match

## Impact
- Affected specs: `multiclaude-source` (new capability spec)
- Affected code: `internal/message/multiclaude.go` lines ~117-121 and ~214-220
- Migration: Users may need to update any saved workspace preferences
