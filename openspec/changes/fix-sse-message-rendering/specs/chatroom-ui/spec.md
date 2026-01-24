## MODIFIED Requirements

### Requirement: Real-time Message Updates
The system SHALL display new messages in real-time as they arrive via Server-Sent Events without requiring page refresh. SSE messages SHALL be rendered as HTML elements on the server side before transmission.

#### Scenario: New message arrives
- **WHEN** a new message is written to a watched directory
- **THEN** the message appears in the UI within 2 seconds as a properly rendered HTML element
- **AND** no manual refresh is required

#### Scenario: SSE message format
- **WHEN** the server sends a message via SSE
- **THEN** the message data SHALL be pre-rendered HTML (not raw JSON)
- **AND** the HTML SHALL match the message format used in the initial page load

#### Scenario: SSE connection status
- **WHEN** the SSE connection is active
- **THEN** a status indicator shows "connected"
- **WHEN** the SSE connection is lost
- **THEN** the status indicator shows "disconnected"
- **AND** the browser automatically attempts to reconnect
