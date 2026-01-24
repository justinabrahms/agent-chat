## ADDED Requirements

### Requirement: Message Search
The system SHALL provide a search input that filters displayed messages in real-time based on user input, matching against message body, sender, and workspace.

#### Scenario: Empty search shows all messages
- **WHEN** the search input is empty
- **THEN** all messages for the selected workspace are displayed

#### Scenario: Search filters messages by body
- **WHEN** the user types a search query
- **THEN** only messages containing that text in the body are shown
- **AND** matching is case-insensitive

#### Scenario: Search filters messages by sender
- **WHEN** the user types a sender name in the search
- **THEN** messages from or to that sender are shown

#### Scenario: Search filters messages by workspace
- **WHEN** the user types a workspace name in the search
- **THEN** messages from that workspace are shown

#### Scenario: No results found
- **WHEN** the search query matches no messages
- **THEN** an empty state is displayed indicating no matches

#### Scenario: Clear search restores messages
- **WHEN** the user clears the search input
- **THEN** all messages are displayed again
