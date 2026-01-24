## ADDED Requirements

### Requirement: Message Navigation
The system SHALL allow users to navigate between messages using `j` (next) and `k` (previous) keys.

#### Scenario: Navigate to next message
- **WHEN** user presses `j` key
- **THEN** the next message in the list is selected and scrolled into view

#### Scenario: Navigate to previous message
- **WHEN** user presses `k` key
- **THEN** the previous message in the list is selected and scrolled into view

#### Scenario: Boundary behavior
- **WHEN** user presses `j` on the last message or `k` on the first message
- **THEN** selection remains on the current message (no wrap-around)

### Requirement: Workspace Focus Shortcut
The system SHALL allow users to jump to the workspace list using the `g+w` key chord.

#### Scenario: Jump to workspace list
- **WHEN** user presses `g` followed by `w` within 500ms
- **THEN** focus moves to the first workspace item in the sidebar

### Requirement: Search Focus Shortcut
The system SHALL allow users to focus the search input using the `/` key.

#### Scenario: Focus search when present
- **WHEN** user presses `/` key and a search input exists
- **THEN** focus moves to the search input

#### Scenario: No search present
- **WHEN** user presses `/` key and no search input exists
- **THEN** nothing happens (graceful no-op)

### Requirement: Selected Message Visual Indicator
The system SHALL visually indicate which message is currently selected.

#### Scenario: Selected message styling
- **WHEN** a message is selected via keyboard navigation
- **THEN** the message displays a distinct background color or border

### Requirement: Keyboard Shortcut Non-Interference
Keyboard shortcuts SHALL NOT trigger when user is typing in an input field.

#### Scenario: Input focus prevents shortcuts
- **WHEN** user is focused on a text input or textarea
- **THEN** pressing shortcut keys types characters instead of triggering navigation
