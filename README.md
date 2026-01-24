# Agent Chat

A Slack-like chat interface for viewing agent communication. Displays messages from [Gas Town](https://github.com/yourusername/gastown) and [multiclaude](https://github.com/dlorenc/multiclaude) in a unified, real-time UI.

<img width="954" height="678" alt="Screenshot 2026-01-24 at 12 51 16 AM" src="https://github.com/user-attachments/assets/54b0ecde-204f-43c6-ade9-4ffe97e92d11" />

## Features

- **Unified View** — See messages from multiple agent systems in one place
- **Real-time Updates** — New messages appear instantly via Server-Sent Events
- **Workspace Channels** — Messages organized by workspace/rig, like Slack channels
- **Unread Indicators** — Know which workspaces have new messages at a glance
- **Dark Theme** — Easy on the eyes for long monitoring sessions
- **Zero Config** — Works out of the box with default paths

## Installation

### From Source

```bash
# Clone the repo
git clone https://github.com/justinabrahms/agent-chat.git
cd agent-chat

# Build
go build -o agent-chat ./cmd/agent-chat

# Run
./agent-chat
```

### Pre-built Binaries

Coming soon — see [Releases](https://github.com/justinabrahms/agent-chat/releases).

## Usage

```bash
# Start with defaults (looks for ~/.beads and ~/.multiclaude)
./agent-chat

# Custom port
./agent-chat -port 3000

# Custom source directories
./agent-chat -gastown-dir /path/to/.beads -multiclaude-dir /path/to/.multiclaude
```

Then open http://localhost:8080 in your browser.

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | HTTP server port | `8080` |
| `GASTOWN_DIR` | Path to Gas Town .beads directory | `~/.beads` |
| `MULTICLAUDE_DIR` | Path to multiclaude directory | `~/.multiclaude` |

## Supported Message Sources

- **Gas Town** — Reads from the beads SQLite database (`beads.db`)
- **Multiclaude** — Reads JSON message files from `~/.multiclaude/messages/`

## Development

```bash
# Run in development
go run ./cmd/agent-chat

# Build
go build -o agent-chat ./cmd/agent-chat

# Run tests
go test ./...
```

### Project Structure

```
.
├── cmd/agent-chat/     # Main application entry point
├── internal/
│   ├── message/        # Message types and source adapters
│   └── server/         # HTTP server, templates, static assets
├── openspec/           # Specifications and change proposals
└── TODO.md             # Project roadmap
```

## Contributing

This project uses [OpenSpec](https://github.com/openspec/openspec) for spec-driven development. Before implementing new features:

1. Create a proposal in `openspec/changes/<change-id>/`
2. Get the proposal reviewed
3. Implement according to the spec

See [TODO.md](TODO.md) for ideas on what to work on.

## License

MIT — see [LICENSE](LICENSE) for details.

## Acknowledgments

Built for monitoring AI agent orchestration systems like Gas Town and multiclaude.
