# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **Chatroom UI**: Slack/Discord-style web interface for viewing agent messages with real-time updates via Server-Sent Events (SSE)
  - Go backend server that watches mail directories for new messages
  - HTMX-based chat UI showing channels (workspaces) and messages
  - Support for Gas Town mail format and multiclaude filesystem messages
- **Syntax highlighting**: Code blocks in messages now have language-aware syntax highlighting using highlight.js with auto-detection and support for light/dark themes ([#17])
- **Relative timestamps**: Dynamic timestamps ("2 minutes ago" style) that update periodically, with absolute timestamps shown on hover ([#18])
- **URL and issue reference detection**: URLs (http/https) and GitHub-style references (#123) are automatically linkified in messages ([#16])
- **Unread count in browser tab**: Tab title shows unread message count (e.g., "(3) Agent Chat") for visibility when the tab is not in focus ([#15])
- **Light theme with system auto-detection**: Theme toggle supporting light/dark/system modes with persistence in localStorage ([#12])
- **Full-text search**: Client-side search filtering messages by content, sender, or workspace name ([#10])
- **Config file support**: Load settings from `~/.config/agent-chat/config.yaml` with precedence: flags > env vars > config file > defaults ([#11])
- **Message grouping**: Consecutive messages from the same sender collapse headers for reduced visual clutter ([#9])
- **Keyboard shortcuts**: Navigation with `j/k` keys between messages, `/` to focus search, `g+w` chord to jump to workspace list ([#7])
- **Makefile**: Standard build targets (build, test, lint, clean, install, help) with cross-platform support for Linux and macOS ([#5])
- **GitHub Actions CI**: Automated testing, linting (golangci-lint), and build verification on PRs ([#1])
- **CONTRIBUTING.md**: Development setup instructions for new contributors ([#4])
- **GitHub templates**: Issue and PR templates for consistent contribution workflow ([#3])

### Changed

- **BREAKING**: Multiclaude workspaces now use `mc-` prefix (e.g., `mc-<repo>` instead of `<repo>`) to prevent naming conflicts with other message sources ([#14])

### Fixed

- SSE messages now render as proper HTML instead of raw JSON ([#13])
- File watcher detects new directories immediately instead of waiting up to 30 seconds for periodic rescan ([#8])
- golangci-lint errcheck errors ([#6])

[Unreleased]: https://github.com/justinabrahms/agent-chat/compare/main...HEAD
[#17]: https://github.com/justinabrahms/agent-chat/pull/17
[#18]: https://github.com/justinabrahms/agent-chat/pull/18
[#16]: https://github.com/justinabrahms/agent-chat/pull/16
[#15]: https://github.com/justinabrahms/agent-chat/pull/15
[#14]: https://github.com/justinabrahms/agent-chat/pull/14
[#13]: https://github.com/justinabrahms/agent-chat/pull/13
[#12]: https://github.com/justinabrahms/agent-chat/pull/12
[#11]: https://github.com/justinabrahms/agent-chat/pull/11
[#10]: https://github.com/justinabrahms/agent-chat/pull/10
[#9]: https://github.com/justinabrahms/agent-chat/pull/9
[#8]: https://github.com/justinabrahms/agent-chat/pull/8
[#7]: https://github.com/justinabrahms/agent-chat/pull/7
[#6]: https://github.com/justinabrahms/agent-chat/pull/6
[#5]: https://github.com/justinabrahms/agent-chat/pull/5
[#4]: https://github.com/justinabrahms/agent-chat/pull/4
[#3]: https://github.com/justinabrahms/agent-chat/pull/3
[#1]: https://github.com/justinabrahms/agent-chat/pull/1
