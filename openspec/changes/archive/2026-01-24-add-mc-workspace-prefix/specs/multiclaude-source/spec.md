## ADDED Requirements

### Requirement: Workspace Prefix Convention
The multiclaude message source SHALL prefix all workspace names with `mc-` to distinguish them from other message sources.

#### Scenario: Repository creates prefixed workspace
- **WHEN** messages are read from `~/.multiclaude/messages/agent-chat/`
- **THEN** the workspace name SHALL be `mc-agent-chat`

#### Scenario: Workspace filtering with prefix
- **WHEN** filtering messages by workspace `mc-agent-chat`
- **THEN** only messages from the `agent-chat` repository directory SHALL be returned

#### Scenario: Watch events use prefixed workspace
- **WHEN** a new message file is created in `~/.multiclaude/messages/my-repo/agent/msg.json`
- **THEN** the emitted Message SHALL have Workspace set to `mc-my-repo`
