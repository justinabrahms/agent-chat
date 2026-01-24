package message

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// MulticlaudeMessage represents the JSON format used by multiclaude.
type MulticlaudeMessage struct {
	ID        string    `json:"id"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Timestamp time.Time `json:"timestamp"`
	Body      string    `json:"body"`
	Status    string    `json:"status"`
}

// MulticlaudeSource reads messages from the multiclaude messages directory.
type MulticlaudeSource struct {
	baseDir    string
	mu         sync.RWMutex
	seenFiles  map[string]bool
	workspaces map[string]bool
}

// NewMulticlaudeSource creates a new multiclaude message source.
// baseDir should be ~/.multiclaude or similar.
func NewMulticlaudeSource(baseDir string) (*MulticlaudeSource, error) {
	messagesDir := filepath.Join(baseDir, "messages")
	if _, err := os.Stat(messagesDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("multiclaude messages directory not found at %s", messagesDir)
	}

	return &MulticlaudeSource{
		baseDir:    baseDir,
		seenFiles:  make(map[string]bool),
		workspaces: make(map[string]bool),
	}, nil
}

func (m *MulticlaudeSource) Name() string {
	return "multiclaude"
}

func (m *MulticlaudeSource) messagesDir() string {
	return filepath.Join(m.baseDir, "messages")
}

// readMessage reads and parses a single message file.
func readMessage(path string) (*MulticlaudeMessage, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var msg MulticlaudeMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

func (m *MulticlaudeSource) List(workspace string) ([]Message, error) {
	messagesDir := m.messagesDir()

	var messages []Message

	// Walk: messages/<repo>/<agent>/*.json
	repos, err := os.ReadDir(messagesDir)
	if err != nil {
		return nil, err
	}

	for _, repoEntry := range repos {
		if !repoEntry.IsDir() {
			continue
		}
		repoName := repoEntry.Name()
		repoPath := filepath.Join(messagesDir, repoName)

		agents, err := os.ReadDir(repoPath)
		if err != nil {
			continue
		}

		for _, agentEntry := range agents {
			if !agentEntry.IsDir() {
				continue
			}
			agentPath := filepath.Join(repoPath, agentEntry.Name())

			files, err := os.ReadDir(agentPath)
			if err != nil {
				continue
			}

			for _, f := range files {
				if f.IsDir() || !strings.HasSuffix(f.Name(), ".json") {
					continue
				}

				filePath := filepath.Join(agentPath, f.Name())
				msg, err := readMessage(filePath)
				if err != nil {
					continue
				}

				// Workspace is the repo name with mc- prefix
				ws := "mc-" + repoName
				if workspace != "" && ws != workspace {
					continue
				}

				m.mu.Lock()
				m.workspaces[ws] = true
				m.seenFiles[filePath] = true
				m.mu.Unlock()

				messages = append(messages, Message{
					ID:        msg.ID,
					Workspace: ws,
					From:      msg.From,
					To:        msg.To,
					Body:      msg.Body,
					Timestamp: msg.Timestamp,
					Source:    "multiclaude",
				})
			}
		}
	}

	return messages, nil
}

func (m *MulticlaudeSource) Workspaces() ([]string, error) {
	// Populate by listing
	if _, err := m.List(""); err != nil {
		return nil, err
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []string
	for ws := range m.workspaces {
		result = append(result, ws)
	}
	return result, nil
}

func (m *MulticlaudeSource) Watch(ctx context.Context) (<-chan Message, error) {
	out := make(chan Message, 100)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	// watchDirRecursive adds a directory and all its subdirectories to the watcher.
	// It's safe to call on directories already being watched (fsnotify ignores duplicates).
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

	// Watch the messages directory recursively
	messagesDir := m.messagesDir()
	watchDirRecursive(messagesDir)

	// Initialize seen files
	_, _ = m.List("")

	go func() {
		defer watcher.Close()
		defer close(out)

		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		checkFile := func(path string) {
			if !strings.HasSuffix(path, ".json") {
				return
			}

			m.mu.RLock()
			seen := m.seenFiles[path]
			m.mu.RUnlock()

			if seen {
				return
			}

			msg, err := readMessage(path)
			if err != nil {
				return
			}

			// Extract workspace from path: messages/<repo>/<agent>/*.json
			rel, _ := filepath.Rel(m.messagesDir(), path)
			parts := strings.Split(rel, string(filepath.Separator))
			ws := "mc-general"
			if len(parts) >= 1 {
				ws = "mc-" + parts[0]
			}

			m.mu.Lock()
			m.seenFiles[path] = true
			m.workspaces[ws] = true
			m.mu.Unlock()

			select {
			case out <- Message{
				ID:        msg.ID,
				Workspace: ws,
				From:      msg.From,
				To:        msg.To,
				Body:      msg.Body,
				Timestamp: msg.Timestamp,
				Source:    "multiclaude",
			}:
			case <-ctx.Done():
				return
			}
		}

		rescan := func() {
			_ = filepath.Walk(m.messagesDir(), func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return nil
				}
				if info.IsDir() {
					// Add any new directories to the watcher (fsnotify ignores duplicates)
					_ = watcher.Add(path)
					return nil
				}
				checkFile(path)
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
					// If it's a directory, recursively add it and any subdirectories to the watch.
					// This handles cases where nested directories are created simultaneously.
					if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
						watchDirRecursive(event.Name)
					} else {
						checkFile(event.Name)
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
