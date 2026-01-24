# Contributing to Agent Chat

## Prerequisites

- Go 1.24.1 or later

## Getting Started

```bash
# Clone the repository
git clone https://github.com/justinabrahms/agent-chat.git
cd agent-chat
```

## Build

```bash
go build -o agent-chat ./cmd/agent-chat
```

## Run

```bash
# With defaults (reads from ~/.beads and ~/.multiclaude)
./agent-chat

# Development mode
go run ./cmd/agent-chat

# Custom port
./agent-chat -port 3000
```

Open http://localhost:8080 in your browser.

## Test

```bash
go test ./...
```

## Project Structure

```
cmd/agent-chat/     # Entry point
internal/
  message/          # Message types and source adapters
  server/           # HTTP server, templates, static assets
openspec/           # Specs and change proposals
```

## Making Changes

This project uses [OpenSpec](https://github.com/openspec/openspec) for spec-driven development.

For new features or breaking changes:
1. Create a proposal in `openspec/changes/<change-id>/`
2. Run `openspec validate <change-id> --strict`
3. Get the proposal reviewed
4. Implement according to the spec

For bug fixes, typos, or documentation: submit a PR directly.
