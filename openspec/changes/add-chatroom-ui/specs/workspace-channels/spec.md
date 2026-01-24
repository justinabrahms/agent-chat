## ADDED Requirements

### Requirement: Workspace Discovery
The system SHALL automatically discover workspaces from message directories and display them as channels in the sidebar.

#### Scenario: Discover workspaces from messages
- **WHEN** messages are loaded from a source
- **THEN** unique workspace identifiers are extracted
- **AND** each workspace appears as a channel in the sidebar

#### Scenario: New workspace appears
- **WHEN** a message arrives for a previously unseen workspace
- **THEN** a new channel is added to the sidebar

### Requirement: Workspace Selection State
The system SHALL maintain the currently selected workspace and update the message view accordingly.

#### Scenario: Initial workspace selection
- **WHEN** the application loads with available workspaces
- **THEN** the first workspace alphabetically is selected by default

#### Scenario: Persist selection during updates
- **WHEN** new messages arrive
- **THEN** the currently selected workspace remains selected
- **AND** the view does not jump to a different workspace

### Requirement: Workspace Message Filtering
The system SHALL filter displayed messages to only show those belonging to the currently selected workspace.

#### Scenario: Filter messages by workspace
- **WHEN** a workspace is selected
- **THEN** only messages with a matching workspace identifier are displayed
- **AND** messages from other workspaces are hidden
