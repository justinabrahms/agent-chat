## ADDED Requirements

### Requirement: Immediate Directory Detection
The file watcher SHALL detect and begin monitoring new directories immediately when they are created, without waiting for periodic rescans.

#### Scenario: New agent directory created during watch
- **WHEN** a new agent directory is created under an existing repo directory
- **THEN** the watcher SHALL add the new directory to monitoring within 1 second
- **AND** subsequent message files in that directory SHALL be detected immediately

#### Scenario: Nested directory structure created simultaneously
- **WHEN** a new repo directory and agent subdirectory are created together (e.g., `mkdir -p repo/agent`)
- **THEN** the watcher SHALL recursively discover and monitor all new directories
- **AND** message files SHALL be detected without waiting for the 30-second rescan

### Requirement: Rescan Discovers New Directories
The periodic rescan SHALL add newly discovered directories to the active watcher, not just process new files.

#### Scenario: Directory created between rescans not caught by event
- **WHEN** a directory creation event is missed (e.g., due to timing)
- **THEN** the periodic rescan SHALL discover the directory and add it to the watcher
- **AND** subsequent file events in that directory SHALL be detected normally
