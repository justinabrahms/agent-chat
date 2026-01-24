## ADDED Requirements

### Requirement: Repository-Specific PR Links
The system SHALL convert PR number references (e.g., `#123`) to clickable hyperlinks that navigate directly to the repository's pull request page when workspace repository context is available.

#### Scenario: PR link with repository context
- **WHEN** a message contains a PR reference like `#123`
- **AND** the workspace has an associated repository URL
- **THEN** the reference SHALL be rendered as a link to `https://github.com/{owner}/{repo}/pull/123`

#### Scenario: PR link without repository context
- **WHEN** a message contains a PR reference like `#123`
- **AND** the workspace does not have an associated repository URL
- **THEN** the reference SHALL be rendered as a link to GitHub search (`https://github.com/search?q=%23123&type=pullrequests`)

#### Scenario: Repository URL with .git suffix
- **WHEN** the workspace repository URL ends with `.git`
- **THEN** the system SHALL strip the `.git` suffix before constructing PR links
