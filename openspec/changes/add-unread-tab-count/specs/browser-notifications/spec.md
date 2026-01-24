## ADDED Requirements

### Requirement: Unread Count Tab Title
The system SHALL display the total unread message count in the browser tab title when there are unread messages.

#### Scenario: Tab title shows unread count
- **WHEN** there are unread messages across any workspace
- **THEN** the browser tab title displays in format "(N) Agent Chat" where N is the total unread count

#### Scenario: Tab title resets when no unread
- **WHEN** all messages have been read (no workspaces have unread messages)
- **THEN** the browser tab title displays "Agent Chat" without a count prefix

### Requirement: Real-time Tab Title Updates
The system SHALL update the browser tab title in real-time as new messages arrive via SSE.

#### Scenario: Count increments on new message
- **WHEN** a new message arrives for a workspace that is not currently active
- **THEN** the unread count in the tab title increments

#### Scenario: Count updates when workspace marked read
- **WHEN** user clicks on a workspace with unread messages
- **THEN** the unread count in the tab title decreases by the number of unread messages in that workspace
