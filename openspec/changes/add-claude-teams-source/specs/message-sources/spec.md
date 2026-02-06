## ADDED Requirements

### Requirement: Claude Teams Message Source
The system SHALL provide a `claude-teams` message source that reads messages from Claude Teams inbox files stored at `{baseDir}/{team-name}/inboxes/{agent}.json`.

#### Scenario: List messages from team inboxes
- **WHEN** the Claude Teams directory contains team directories with inbox JSON files
- **THEN** all messages from all inboxes are returned with workspace `teams-{teamName}`

#### Scenario: Idle notification rendering
- **WHEN** an inbox message text starts with `{"type":"idle_notification"`
- **THEN** the message body SHALL be set to a human-readable status prefixed with `[status]`
- **AND** the UI SHALL render it with dimmed styling

#### Scenario: Watch for new messages
- **WHEN** new messages are written to inbox files
- **THEN** the watcher SHALL emit them via the Watch channel
- **AND** previously seen messages SHALL not be re-emitted

### Requirement: Claude Teams Configuration
The system SHALL accept a `claude-teams-dir` configuration via CLI flag, environment variable (`CLAUDE_TEAMS_DIR`), or config file (`claude-teams-dir` YAML key), defaulting to `~/.claude/teams`.

#### Scenario: Default directory
- **WHEN** no explicit configuration is provided
- **THEN** the system SHALL look for teams at `~/.claude/teams`
