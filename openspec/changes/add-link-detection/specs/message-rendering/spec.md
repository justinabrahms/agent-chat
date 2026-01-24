## ADDED Requirements

### Requirement: URL Linkification
The system SHALL detect URLs in message bodies and render them as clickable hyperlinks.

#### Scenario: HTTP URL detection
- **WHEN** a message body contains "Check http://example.com for details"
- **THEN** the URL is rendered as `<a href="http://example.com" target="_blank" rel="noopener noreferrer">http://example.com</a>`

#### Scenario: HTTPS URL detection
- **WHEN** a message body contains "See https://github.com/org/repo/pull/123"
- **THEN** the URL is rendered as a clickable link opening in a new tab

#### Scenario: URL with path and query parameters
- **WHEN** a message body contains "https://example.com/path?query=value&other=1"
- **THEN** the entire URL including path and query string is linkified

### Requirement: GitHub Issue/PR Reference Linkification
The system SHALL detect GitHub-style issue and PR references (#123) and render them as links.

#### Scenario: Issue reference in message
- **WHEN** a message body contains "Fixed in #123"
- **THEN** "#123" is rendered as a link to the GitHub issue/PR

#### Scenario: Multiple issue references
- **WHEN** a message body contains "See #123 and #456"
- **THEN** both references are rendered as separate clickable links

#### Scenario: Hash in non-reference context
- **WHEN** a message body contains "Color #ffffff" or "Channel #general"
- **THEN** these are NOT linkified (only numeric references are linked)

### Requirement: Link Security
The system SHALL ensure all external links are safe to click.

#### Scenario: New tab with security attributes
- **WHEN** any link is rendered
- **THEN** it includes target="_blank" and rel="noopener noreferrer" attributes
