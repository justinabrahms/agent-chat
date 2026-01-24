## ADDED Requirements

### Requirement: Code Block Syntax Highlighting
The system SHALL render code blocks with syntax highlighting when messages contain fenced code blocks.

#### Scenario: Code block with language hint
- **WHEN** a message contains a fenced code block with a language hint (e.g., ```go)
- **THEN** the code block SHALL be rendered with syntax highlighting appropriate for that language

#### Scenario: Code block without language hint
- **WHEN** a message contains a fenced code block without a language hint (e.g., ```)
- **THEN** the system SHALL attempt to auto-detect the language and apply appropriate highlighting

#### Scenario: Theme-aware highlighting
- **WHEN** the user switches between light and dark themes
- **THEN** the syntax highlighting colors SHALL adapt to remain readable in both themes

#### Scenario: New messages via SSE
- **WHEN** a new message with code blocks arrives via SSE
- **THEN** syntax highlighting SHALL be applied to the code blocks without requiring a page refresh
