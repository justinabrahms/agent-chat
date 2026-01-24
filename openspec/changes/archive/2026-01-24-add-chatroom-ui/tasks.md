# Tasks: Add Chatroom UI

## 1. Project Setup
- [x] 1.1 Initialize Go module (`go mod init`)
- [x] 1.2 Create directory structure (`cmd/`, `internal/`, `web/`)
- [x] 1.3 Add fsnotify dependency

## 2. Core Types and Interfaces
- [x] 2.1 Define `Message` struct
- [x] 2.2 Define `MessageSource` interface
- [x] 2.3 Define `Workspace` type

## 3. Message Sources
- [x] 3.1 Implement Gas Town mail source adapter
- [x] 3.2 Implement multiclaude message source adapter
- [x] 3.3 Implement filesystem watcher with fsnotify
- [x] 3.4 Add periodic rescan fallback (30s interval)

## 4. HTTP Server
- [x] 4.1 Create main HTTP server with routing
- [x] 4.2 Implement SSE endpoint for real-time updates
- [x] 4.3 Implement workspace list endpoint
- [x] 4.4 Implement messages endpoint (filtered by workspace)

## 5. UI Templates
- [x] 5.1 Create base HTML layout with HTMX
- [x] 5.2 Create sidebar workspace list component
- [x] 5.3 Create message list component
- [x] 5.4 Create individual message component
- [x] 5.5 Add SSE connection status indicator
- [x] 5.6 Style with CSS (Slack/Discord-like appearance)

## 6. Configuration
- [x] 6.1 Add CLI flags for directory paths and port
- [x] 6.2 Add environment variable support
- [x] 6.3 Validate configured directories exist

## 7. Integration
- [x] 7.1 Wire up message sources to SSE broadcaster
- [x] 7.2 Test with sample Gas Town messages
- [x] 7.3 Test with sample multiclaude messages
- [x] 7.4 Verify real-time updates work end-to-end

**Status:** ✅ Completed (initial implementation in early commits)
