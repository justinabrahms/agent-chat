# Tasks: Add Chatroom UI

## 1. Project Setup
- [ ] 1.1 Initialize Go module (`go mod init`)
- [ ] 1.2 Create directory structure (`cmd/`, `internal/`, `web/`)
- [ ] 1.3 Add fsnotify dependency

## 2. Core Types and Interfaces
- [ ] 2.1 Define `Message` struct
- [ ] 2.2 Define `MessageSource` interface
- [ ] 2.3 Define `Workspace` type

## 3. Message Sources
- [ ] 3.1 Implement Gas Town mail source adapter
- [ ] 3.2 Implement multiclaude message source adapter
- [ ] 3.3 Implement filesystem watcher with fsnotify
- [ ] 3.4 Add periodic rescan fallback (30s interval)

## 4. HTTP Server
- [ ] 4.1 Create main HTTP server with routing
- [ ] 4.2 Implement SSE endpoint for real-time updates
- [ ] 4.3 Implement workspace list endpoint
- [ ] 4.4 Implement messages endpoint (filtered by workspace)

## 5. UI Templates
- [ ] 5.1 Create base HTML layout with HTMX
- [ ] 5.2 Create sidebar workspace list component
- [ ] 5.3 Create message list component
- [ ] 5.4 Create individual message component
- [ ] 5.5 Add SSE connection status indicator
- [ ] 5.6 Style with CSS (Slack/Discord-like appearance)

## 6. Configuration
- [ ] 6.1 Add CLI flags for directory paths and port
- [ ] 6.2 Add environment variable support
- [ ] 6.3 Validate configured directories exist

## 7. Integration
- [ ] 7.1 Wire up message sources to SSE broadcaster
- [ ] 7.2 Test with sample Gas Town messages
- [ ] 7.3 Test with sample multiclaude messages
- [ ] 7.4 Verify real-time updates work end-to-end
