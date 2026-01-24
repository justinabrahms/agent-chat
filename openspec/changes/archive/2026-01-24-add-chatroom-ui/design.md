# Design: Chatroom UI for Agent Messages

## Context
Gas Town workers and multiclaude agents communicate via filesystem-based "mail" systems. Users need visibility into these conversations without digging through directories. This is a local-only developer tool with no auth requirements.

**Constraints:**
- Must work offline/locally
- Must support existing Gas Town mail format (`gt mail`)
- Must support multiclaude filesystem message format
- Should feel familiar (Slack/Discord-like)

## Goals / Non-Goals

**Goals:**
- Real-time display of agent messages grouped by workspace/channel
- Support both Gas Town and multiclaude message formats
- Simple single-binary deployment (Go)
- Minimal UI dependencies (HTMX + vanilla CSS)

**Non-Goals:**
- Sending messages (read-only view initially)
- Multi-user access / auth
- Message persistence beyond filesystem
- Mobile-responsive design (desktop developer tool)

## Decisions

### Architecture: Single Go binary + embedded templates

**Decision:** Go HTTP server with embedded HTML templates, HTMX for interactivity, SSE for real-time updates.

**Alternatives considered:**
- SPA (React/Vue) - Overkill for this use case, adds build complexity
- Static HTML + File System Access API - Browser compatibility issues, permission friction
- Python/Flask - Works but Go gives single-binary deployment

### Message Source Abstraction

**Decision:** Define a `MessageSource` interface that both Gas Town and multiclaude adapters implement.

```go
type Message struct {
    ID        string
    Workspace string    // maps to "channel" in UI
    From      string    // agent identifier
    To        string    // recipient (may be empty for broadcasts)
    Body      string
    Timestamp time.Time
    Source    string    // "gastown" or "multiclaude"
}

type MessageSource interface {
    Name() string
    Watch(ctx context.Context) (<-chan Message, error)
    List(workspace string) ([]Message, error)
}
```

This allows adding more message sources later without changing the core.

### Real-time Updates: Server-Sent Events (SSE)

**Decision:** Use SSE rather than WebSockets.

**Rationale:**
- HTMX has first-class SSE support (`hx-ext="sse"`)
- Simpler than WebSockets for one-way server→client push
- Automatic reconnection built into browser SSE API
- Sufficient for our read-only use case

### UI Structure

```
┌────────────────────────────────────────────────────────────┐
│  Agent Chat                                          [⚙️]  │
├──────────────┬─────────────────────────────────────────────┤
│ WORKSPACES   │  #project-alpha                             │
│              │─────────────────────────────────────────────│
│ #project-a   │  [worker-1] 10:32                           │
│ #project-b   │  Starting task on molecule XYZ              │
│ #infra       │                                             │
│              │  [worker-2] 10:33                           │
│              │  Acknowledged. Taking over from worker-1    │
│              │                                             │
│              │  [orchestrator] 10:35                       │
│              │  All workers synced. Proceeding to build.   │
│              │                                             │
├──────────────┴─────────────────────────────────────────────┤
│ SSE connected ● 3 sources active                           │
└────────────────────────────────────────────────────────────┘
```

### Configuration

**Decision:** Config via flags and/or environment variables. No config file initially.

```bash
agent-chat \
  --gastown-mail-dir ~/.gastown/mail \
  --multiclaude-dir ~/.multiclaude \
  --port 8080
```

Environment equivalents: `GASTOWN_MAIL_DIR`, `MULTICLAUDE_DIR`, `PORT`

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| fsnotify may miss events on some filesystems | Periodic rescan fallback (every 30s) |
| Large message volumes could overwhelm UI | Pagination + limit displayed messages per channel |
| Gas Town / multiclaude format changes | Abstract behind MessageSource interface |

## Open Questions

1. **What is the exact Gas Town mail directory structure?** Need to inspect `~/.gastown/mail` or similar to understand format.
2. **What is the multiclaude message format?** Need to inspect existing multiclaude filesystem messages.
3. **Should we support filtering messages by agent?** Deferred to future enhancement.
