## ADDED Requirements

### Requirement: Message Source Abstraction
The system SHALL define a common interface for reading messages from different sources, allowing multiple message formats to be supported.

#### Scenario: Multiple sources configured
- **WHEN** both Gas Town and multiclaude directories are configured
- **THEN** messages from both sources appear in the UI
- **AND** each message indicates its source

### Requirement: Gas Town Mail Support
The system SHALL read messages from Gas Town mail directories in the format used by `gt mail`.

#### Scenario: Read Gas Town messages
- **WHEN** the Gas Town mail directory is configured
- **THEN** existing messages are loaded and displayed
- **AND** new messages are detected via filesystem watching

#### Scenario: Parse Gas Town message format
- **WHEN** a Gas Town mail file is read
- **THEN** the sender, recipient, timestamp, and body are extracted

### Requirement: Multiclaude Message Support
The system SHALL read messages from multiclaude filesystem message directories.

#### Scenario: Read multiclaude messages
- **WHEN** the multiclaude directory is configured
- **THEN** existing messages are loaded and displayed
- **AND** new messages are detected via filesystem watching

#### Scenario: Parse multiclaude message format
- **WHEN** a multiclaude message file is read
- **THEN** the sender, recipient, timestamp, and body are extracted

### Requirement: Filesystem Watching
The system SHALL watch configured directories for new message files using filesystem events with a periodic rescan fallback.

#### Scenario: Filesystem event detection
- **WHEN** a new message file is created in a watched directory
- **THEN** the system detects it via fsnotify events

#### Scenario: Periodic rescan fallback
- **WHEN** filesystem events may be unreliable
- **THEN** the system performs a periodic rescan every 30 seconds to catch missed messages
