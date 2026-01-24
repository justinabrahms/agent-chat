# Message Display

## ADDED Requirements

### Requirement: Display relative timestamps for messages

The system SHALL display relative timestamps ("2 minutes ago" style) instead of static HH:MM format.

#### Scenario: Message sent just now
Given a message was sent less than 1 minute ago
When the message is displayed
Then the timestamp shows "just now"

#### Scenario: Message sent minutes ago
Given a message was sent 5 minutes ago
When the message is displayed
Then the timestamp shows "5 minutes ago"

#### Scenario: Message sent hours ago
Given a message was sent 3 hours ago
When the message is displayed
Then the timestamp shows "3 hours ago"

#### Scenario: Message sent yesterday
Given a message was sent yesterday
When the message is displayed
Then the timestamp shows "yesterday"

#### Scenario: Message sent within past week
Given a message was sent 3 days ago on Wednesday
When the message is displayed
Then the timestamp shows the day name "Wed"

#### Scenario: Message sent more than a week ago
Given a message was sent on January 10th
When the message is displayed
Then the timestamp shows the date "Jan 10"

### Requirement: Show absolute time on hover

The system SHALL show the full absolute date and time when hovering over a timestamp.

#### Scenario: Hover shows full timestamp
Given a message with a relative timestamp displayed
When the user hovers over the timestamp
Then a tooltip shows the full date and time (e.g., "Jan 24, 2026, 2:30 PM")

### Requirement: Timestamps update automatically

The system SHALL update timestamps periodically without requiring a page reload.

#### Scenario: Timestamps refresh every 30 seconds
Given a message displayed with timestamp "just now"
When 2 minutes pass
Then the timestamp updates to "2 minutes ago" without page reload

#### Scenario: New SSE messages get relative timestamps
Given a user is viewing the message panel
When a new message arrives via SSE
Then the message displays with a relative timestamp
