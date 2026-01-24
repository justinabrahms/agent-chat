## ADDED Requirements

### Requirement: Message Grouping by Sender
The system SHALL group consecutive messages from the same sender, displaying the sender name and timestamp only on the first message in each group.

#### Scenario: First message in a group
- **WHEN** a message is displayed
- **AND** the previous message was from a different sender (or there is no previous message)
- **THEN** the message displays the sender name and full timestamp

#### Scenario: Continuation message in a group
- **WHEN** a message is displayed
- **AND** the previous message was from the same sender
- **THEN** the message displays only the body content
- **AND** the sender name and timestamp are hidden

#### Scenario: Hover on continuation message
- **WHEN** the user hovers over a continuation message (one without visible timestamp)
- **THEN** the timestamp is shown in a tooltip or inline indicator

### Requirement: Visual Separation Between Sender Groups
The system SHALL display a subtle visual separator between message groups from different senders to distinguish conversation flow.

#### Scenario: Separator between sender groups
- **WHEN** a message from a new sender follows messages from a different sender
- **THEN** a subtle visual separator (spacing or divider line) appears between the groups

#### Scenario: No separator within same sender group
- **WHEN** multiple consecutive messages are from the same sender
- **THEN** no separator appears between those messages
- **AND** they appear as a visually cohesive group
