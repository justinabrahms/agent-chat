# Change: Add URL/Link Detection to Messages

## Why
Message bodies often contain URLs, GitHub PR numbers (#123), and issue references that should be clickable. Currently these are displayed as plain text, forcing users to copy/paste links manually.

## What Changes
- Enhance the `markdown` template function to detect and linkify URLs (http/https)
- Detect GitHub-style PR/issue references (#123) and link them
- All generated links open in new tabs (target="_blank" with rel="noopener noreferrer")

## Impact
- Affected specs: message-rendering (new capability)
- Affected code: `internal/server/server.go` (markdown function)
