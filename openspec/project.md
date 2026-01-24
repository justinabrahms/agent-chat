# Project Context

## Purpose
Agent Chat is a Slack-like chat interface for viewing agent communication. It displays messages from Gas Town and multiclaude in a unified, real-time UI, making it easy to monitor AI agent orchestration systems.

## Tech Stack
- **Backend:** Go (1.21+)
- **Frontend:** HTML + HTMX + vanilla JavaScript (embedded in binary)
- **Database:** SQLite (for Gas Town beads), JSON files (for multiclaude)
- **Real-time:** Server-Sent Events (SSE)
- **Build:** Makefile, GitHub Actions CI

## Project Conventions

### Code Style
- Standard Go formatting (`go fmt`)
- Use `go vet` and `golangci-lint` for linting
- Avoid external dependencies when stdlib suffices

### Architecture Patterns
- **Message Sources:** Adapter pattern - each source implements a common interface
- **HTTP Server:** Single binary with embedded static assets
- **Templates:** Go html/template with HTMX for interactivity
- **File Watching:** fsnotify with periodic rescan fallback

### Testing Strategy
- Unit tests for message parsing and adapters
- Integration tests for HTTP endpoints
- Run with `go test ./...` or `make test`

### Git Workflow
- Feature branches with PRs
- Squash merge to main
- OpenSpec proposals for new features
- All PRs run CI (tests, lint, build)

## Domain Context
- **Workspace:** A channel-like grouping of messages (corresponds to a rig or agent session)
- **Message Source:** An adapter that reads messages from Gas Town or multiclaude
- **Gas Town:** An agent orchestration system using a beads SQLite database
- **multiclaude:** A multi-agent Claude system with JSON message files

## Important Constraints
- Must work with zero configuration (sensible defaults)
- Single binary distribution (no external dependencies at runtime)
- Real-time updates without polling (use SSE)
- Support both light and dark themes

## External Dependencies
- Gas Town beads database (`~/.beads/beads.db`)
- multiclaude message files (`~/.multiclaude/messages/`)
