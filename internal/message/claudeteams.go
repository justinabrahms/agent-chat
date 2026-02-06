package message

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// ClaudeTeamsInboxMessage represents a single message in a Claude Teams inbox file.
type ClaudeTeamsInboxMessage struct {
	From      string    `json:"from"`
	Text      string    `json:"text"`
	Summary   string    `json:"summary"`
	Timestamp time.Time `json:"timestamp"`
	Color     string    `json:"color"`
	Read      bool      `json:"read"`
}

// idleNotification is the embedded JSON structure for idle notifications.
type idleNotification struct {
	Type       string `json:"type"`
	From       string `json:"from"`
	IdleReason string `json:"idleReason"`
}

// ClaudeTeamsSource reads messages from Claude Teams inbox directories.
type ClaudeTeamsSource struct {
	baseDir     string
	mu          sync.RWMutex
	seenMsgKeys map[string]bool
	workspaces  map[string]bool
}

// NewClaudeTeamsSource creates a new Claude Teams message source.
// baseDir should be ~/.claude/teams or similar.
func NewClaudeTeamsSource(baseDir string) (*ClaudeTeamsSource, error) {
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("claude teams directory not found at %s", baseDir)
	}

	return &ClaudeTeamsSource{
		baseDir:     baseDir,
		seenMsgKeys: make(map[string]bool),
		workspaces:  make(map[string]bool),
	}, nil
}

func (c *ClaudeTeamsSource) Name() string {
	return "claude-teams"
}

// msgKey returns a composite key for deduplication in Watch.
func msgKey(teamName, from string, timestamp time.Time) string {
	return fmt.Sprintf("%s|%s|%s", teamName, from, timestamp.Format(time.RFC3339Nano))
}

// msgID returns a stable message ID from the composite key.
func msgID(teamName, from string, timestamp time.Time) string {
	h := sha256.Sum256([]byte(msgKey(teamName, from, timestamp)))
	return fmt.Sprintf("%x", h[:8])
}

// parseIdleNotification checks if the text is an idle notification and returns
// a human-readable body with [status] prefix, or empty string if not idle.
func parseIdleNotification(text string) string {
	if !strings.HasPrefix(text, `{"type":"idle_notification"`) {
		return ""
	}
	var notif idleNotification
	if err := json.Unmarshal([]byte(text), &notif); err != nil {
		return ""
	}
	reason := notif.IdleReason
	if reason == "" {
		reason = "idle"
	}
	return fmt.Sprintf("[status] %s is now %s", notif.From, reason)
}

func (c *ClaudeTeamsSource) List(workspace string) ([]Message, error) {
	teams, err := os.ReadDir(c.baseDir)
	if err != nil {
		return nil, err
	}

	var messages []Message

	for _, teamEntry := range teams {
		if !teamEntry.IsDir() {
			continue
		}
		teamName := teamEntry.Name()
		ws := "teams-" + teamName

		if workspace != "" && ws != workspace {
			continue
		}

		inboxesDir := filepath.Join(c.baseDir, teamName, "inboxes")
		inboxes, err := os.ReadDir(inboxesDir)
		if err != nil {
			continue
		}

		for _, inboxEntry := range inboxes {
			if inboxEntry.IsDir() || !strings.HasSuffix(inboxEntry.Name(), ".json") {
				continue
			}

			recipient := strings.TrimSuffix(inboxEntry.Name(), ".json")
			inboxPath := filepath.Join(inboxesDir, inboxEntry.Name())

			data, err := os.ReadFile(inboxPath)
			if err != nil {
				continue
			}

			var inboxMsgs []ClaudeTeamsInboxMessage
			if err := json.Unmarshal(data, &inboxMsgs); err != nil {
				continue
			}

			for _, m := range inboxMsgs {
				key := msgKey(teamName, m.From, m.Timestamp)

				c.mu.Lock()
				c.workspaces[ws] = true
				c.seenMsgKeys[key] = true
				c.mu.Unlock()

				body := m.Text
				if statusBody := parseIdleNotification(m.Text); statusBody != "" {
					body = statusBody
				}

				messages = append(messages, Message{
					ID:        msgID(teamName, m.From, m.Timestamp),
					Workspace: ws,
					From:      m.From,
					To:        recipient,
					Body:      body,
					Timestamp: m.Timestamp,
					Source:    "claude-teams",
				})
			}
		}
	}

	return messages, nil
}

func (c *ClaudeTeamsSource) Workspaces() ([]string, error) {
	if _, err := c.List(""); err != nil {
		return nil, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	var result []string
	for ws := range c.workspaces {
		result = append(result, ws)
	}
	return result, nil
}

func (c *ClaudeTeamsSource) Watch(ctx context.Context) (<-chan Message, error) {
	out := make(chan Message, 100)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	watchDirRecursive := func(root string) {
		_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				_ = watcher.Add(path)
			}
			return nil
		})
	}

	watchDirRecursive(c.baseDir)

	// Initialize seen message keys
	_, _ = c.List("")

	go func() {
		defer watcher.Close()
		defer close(out)

		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		emitNew := func(teamName, recipient string, m ClaudeTeamsInboxMessage) {
			key := msgKey(teamName, m.From, m.Timestamp)

			c.mu.RLock()
			seen := c.seenMsgKeys[key]
			c.mu.RUnlock()

			if seen {
				return
			}

			ws := "teams-" + teamName

			c.mu.Lock()
			c.seenMsgKeys[key] = true
			c.workspaces[ws] = true
			c.mu.Unlock()

			body := m.Text
			if statusBody := parseIdleNotification(m.Text); statusBody != "" {
				body = statusBody
			}

			select {
			case out <- Message{
				ID:        msgID(teamName, m.From, m.Timestamp),
				Workspace: ws,
				From:      m.From,
				To:        recipient,
				Body:      body,
				Timestamp: m.Timestamp,
				Source:    "claude-teams",
			}:
			case <-ctx.Done():
				return
			}
		}

		scanInbox := func(path string) {
			if !strings.HasSuffix(path, ".json") {
				return
			}

			// Extract team name and recipient from path: {baseDir}/{team}/inboxes/{agent}.json
			rel, err := filepath.Rel(c.baseDir, path)
			if err != nil {
				return
			}
			parts := strings.Split(rel, string(filepath.Separator))
			if len(parts) < 3 || parts[1] != "inboxes" {
				return
			}
			teamName := parts[0]
			recipient := strings.TrimSuffix(parts[2], ".json")

			data, err := os.ReadFile(path)
			if err != nil {
				return
			}

			var inboxMsgs []ClaudeTeamsInboxMessage
			if err := json.Unmarshal(data, &inboxMsgs); err != nil {
				return
			}

			for _, m := range inboxMsgs {
				emitNew(teamName, recipient, m)
			}
		}

		rescan := func() {
			_ = filepath.Walk(c.baseDir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return nil
				}
				if info.IsDir() {
					_ = watcher.Add(path)
					return nil
				}
				scanInbox(path)
				return nil
			})
		}

		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
					if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
						watchDirRecursive(event.Name)
					} else {
						scanInbox(event.Name)
					}
				}
			case <-watcher.Errors:
				// Ignore errors
			case <-ticker.C:
				rescan()
			}
		}
	}()

	return out, nil
}
