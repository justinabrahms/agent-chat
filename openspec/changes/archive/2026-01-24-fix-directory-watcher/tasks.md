## 1. Implementation

- [x] 1.1 Add `watchDirRecursive` helper function to recursively add directories to fsnotify watcher
- [x] 1.2 Update directory creation handler (lines 259-260) to use recursive watching
- [x] 1.3 Update `rescan()` function to also add new directories to watcher, not just check files
- [x] 1.4 Add unit tests for new directory detection scenarios
- [x] 1.5 Verify no race conditions with concurrent directory creation
