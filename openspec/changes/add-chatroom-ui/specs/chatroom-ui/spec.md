## ADDED Requirements

### Requirement: Chat Interface Display
The system SHALL display agent messages in a chat-style interface with a sidebar listing workspaces and a main panel showing messages for the selected workspace.

#### Scenario: View workspace list
- **WHEN** the user opens the application
- **THEN** a sidebar displays all discovered workspaces from configured message sources

#### Scenario: Select workspace
- **WHEN** the user clicks a workspace in the sidebar
- **THEN** the main panel updates to show messages for that workspace

#### Scenario: Display message content
- **WHEN** messages are displayed in the main panel
- **THEN** each message shows the sender agent name, timestamp, and message body

### Requirement: Real-time Message Updates
The system SHALL display new messages in real-time as they arrive via Server-Sent Events without requiring page refresh.

#### Scenario: New message arrives
- **WHEN** a new message is written to a watched directory
- **THEN** the message appears in the UI within 2 seconds
- **AND** no manual refresh is required

#### Scenario: SSE connection status
- **WHEN** the SSE connection is active
- **THEN** a status indicator shows "connected"
- **WHEN** the SSE connection is lost
- **THEN** the status indicator shows "disconnected"
- **AND** the browser automatically attempts to reconnect

### Requirement: Message Chronological Order
The system SHALL display messages in chronological order with newest messages at the bottom.

#### Scenario: Message ordering
- **WHEN** messages are displayed
- **THEN** they appear in ascending timestamp order (oldest first, newest last)

#### Scenario: New message positioning
- **WHEN** a new message arrives
- **THEN** it is appended to the bottom of the message list

### Requirement: Unread Message Indicators
The system SHALL visually indicate workspaces that contain messages the user has not yet viewed, using browser localStorage to track read state.

#### Scenario: Workspace has unread messages
- **WHEN** a workspace has messages newer than the user's last view
- **THEN** the workspace is highlighted in the sidebar with an unread indicator

#### Scenario: Mark workspace as read
- **WHEN** the user selects a workspace
- **THEN** that workspace is marked as read up to the current time
- **AND** the unread indicator is removed

#### Scenario: New message in non-selected workspace
- **WHEN** a new message arrives for a workspace the user is not viewing
- **THEN** an unread indicator appears on that workspace in the sidebar
