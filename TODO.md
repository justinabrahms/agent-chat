# Agent Chat TODO

A roadmap for making this project delightful to use and a healthy open source project.

> **Implementation Note:** All features below should have OpenSpec specifications created as part of their implementation. Create proposals in `openspec/changes/` before writing code.

---

## User Experience

### Search & Navigation
- [ ] Full-text search across all messages
- [ ] Search within current workspace
- [ ] Keyboard shortcuts (j/k navigation, / for search, g+w for workspace list)
- [ ] Jump to date picker
- [ ] "Jump to unread" button when scrolled up

### Message Display
- [ ] Collapsible long messages with "Show more"
- [ ] Syntax highlighting for code blocks in messages
- [ ] Linkify URLs, PR numbers, issue references
- [ ] Relative timestamps ("2 minutes ago") with hover for absolute
- [ ] Message grouping by sender (consecutive messages collapse headers)
- [ ] Thread view for related messages (reply chains)

### Filtering & Organization
- [ ] Filter by agent/sender
- [ ] Filter by source (Gas Town only, multiclaude only)
- [ ] Filter by date range
- [ ] Pin important workspaces to top
- [ ] Archive/hide inactive workspaces
- [ ] Custom workspace aliases/nicknames

### Notifications & Awareness
- [ ] Desktop notifications for new messages (opt-in)
- [ ] Sound alerts (configurable, off by default)
- [ ] Unread count in browser tab title
- [ ] "New messages below" indicator when scrolled up
- [ ] Activity indicator showing which workspaces are "hot"

### Visual Polish
- [ ] Light theme option
- [ ] System theme auto-detection
- [ ] Customizable accent colors
- [ ] Compact vs comfortable message density toggle
- [ ] Smooth scroll animations
- [ ] Loading skeletons instead of blank states

---

## Core Features

### Message Sources
- [ ] Generic filesystem message source (configurable JSON/YAML format)
- [ ] Slack export import
- [ ] Discord export import
- [ ] Plugin system for custom sources
- [ ] Source health indicators (last sync time, error states)

### Data Management
- [ ] Message caching for faster startup
- [ ] Pagination for workspaces with thousands of messages
- [ ] Memory-efficient handling of large message volumes
- [ ] Export conversations to markdown/JSON
- [ ] Clear read state / mark all as read

### Real-time
- [ ] Reconnection indicator with retry countdown
- [ ] Offline mode with cached messages
- [ ] Optimistic UI updates

---

## Developer Experience

### Build & Distribution
- [ ] Makefile with standard targets (build, test, install, clean)
- [ ] Cross-platform builds (Linux, macOS, Windows)
- [ ] Homebrew formula
- [ ] Docker image
- [ ] Pre-built binaries on GitHub Releases
- [ ] Version command (`agent-chat --version`)
- [ ] Single-binary with embedded assets (already done, but verify)

### Configuration
- [ ] Config file support (~/.config/agent-chat/config.yaml)
- [ ] Environment variable documentation
- [ ] Example config file with all options commented
- [ ] Config validation on startup with helpful errors
- [ ] `--config` flag to specify alternate config path

### Observability
- [ ] Structured logging with levels (--verbose, --debug)
- [ ] Metrics endpoint (/metrics) for Prometheus
- [ ] Health check improvements (source status, message counts)

---

## Documentation

### For Users
- [ ] README with screenshots/GIF demo
- [ ] Installation instructions (all methods)
- [ ] Quick start guide
- [ ] Configuration reference
- [ ] Troubleshooting guide
- [ ] FAQ

### For Contributors
- [ ] CONTRIBUTING.md with development setup
- [ ] Architecture overview document
- [ ] Code style guide (or adopt Go standard)
- [ ] How to add a new message source
- [ ] How to add a new theme
- [ ] Release process documentation

### Project Health
- [ ] LICENSE file (pick a license)
- [ ] CODE_OF_CONDUCT.md
- [ ] SECURITY.md for vulnerability reporting
- [ ] Issue templates (bug report, feature request)
- [ ] Pull request template
- [ ] CHANGELOG.md

---

## Testing & Quality

### Automated Testing
- [ ] Unit tests for message parsing
- [ ] Unit tests for source adapters
- [ ] Integration tests for HTTP endpoints
- [ ] End-to-end tests with browser automation
- [ ] Test coverage reporting

### CI/CD
- [ ] GitHub Actions workflow for tests
- [ ] Linting (golangci-lint)
- [ ] Build verification on PRs
- [ ] Automated releases on tag push
- [ ] Dependabot for dependency updates

### Code Quality
- [ ] Go module documentation (godoc comments)
- [ ] Error wrapping with context
- [ ] Graceful shutdown handling (already partial)
- [ ] Signal handling (SIGHUP for config reload)

---

## Accessibility

- [ ] Keyboard-only navigation
- [ ] Screen reader support (ARIA labels)
- [ ] High contrast mode
- [ ] Reduced motion option
- [ ] Focus indicators

---

## Performance

- [ ] Benchmark startup time
- [ ] Profile memory usage with large message sets
- [ ] Lazy load messages (virtual scrolling)
- [ ] Debounce filesystem watchers
- [ ] Connection pooling for SQLite

---

## Nice to Have (Someday)

- [ ] Mobile-responsive layout
- [ ] PWA support (installable, works offline)
- [ ] Message reactions/emoji
- [ ] Bookmarks for important messages
- [ ] Custom CSS injection for power users
- [ ] Multiple windows/tabs sync
- [ ] Message composer (for sources that support sending)
- [ ] AI-powered summary of unread messages
- [ ] Webhook support for integrations

---

## Priority Guide

**Start here (MVP polish):**
1. README with screenshots
2. LICENSE file
3. Makefile
4. Config file support
5. Keyboard shortcuts

**Next (usability):**
1. Full-text search
2. Message grouping
3. Light theme
4. Desktop notifications

**Then (distribution):**
1. GitHub Actions CI
2. Pre-built releases
3. Homebrew formula
