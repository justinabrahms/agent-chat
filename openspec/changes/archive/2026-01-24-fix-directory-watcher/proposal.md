# Change: Fix directory watcher to detect new directories immediately

## Why
The file watcher in `Watch()` only monitors directories that exist at startup. When new agent directories are created after the watcher starts, they aren't watched until the 30-second periodic rescan. This causes messages in new directories to be delayed by up to 30 seconds.

## What Changes
- When a new directory creation event is detected, recursively walk and watch all subdirectories
- Update the periodic rescan to also add newly discovered directories to the watcher
- Add helper function `watchDirRecursive` to encapsulate recursive directory watching

## Impact
- Affected code: `internal/message/multiclaude.go:160-274` (Watch function)
- No breaking changes - existing behavior is preserved, latency is improved
