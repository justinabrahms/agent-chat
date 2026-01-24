package message

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	_ "modernc.org/sqlite"
)

// GasTownSource reads messages from the Gas Town beads SQLite database.
type GasTownSource struct {
	beadsDir    string
	dbPath      string
	mu          sync.RWMutex
	lastSeenID  string
	lastSeenAt  time.Time
	workspaces  map[string]bool
}

// NewGasTownSource creates a new Gas Town message source.
// beadsDir should be the path to the .beads directory (e.g., ~/.beads or the redirect target).
func NewGasTownSource(beadsDir string) (*GasTownSource, error) {
	// Check for redirect
	redirectPath := filepath.Join(beadsDir, "redirect")
	if data, err := os.ReadFile(redirectPath); err == nil {
		beadsDir = strings.TrimSpace(string(data))
	}

	dbPath := filepath.Join(beadsDir, "beads.db")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("beads.db not found at %s", dbPath)
	}

	return &GasTownSource{
		beadsDir:   beadsDir,
		dbPath:     dbPath,
		workspaces: make(map[string]bool),
	}, nil
}

func (g *GasTownSource) Name() string {
	return "gastown"
}

func (g *GasTownSource) openDB() (*sql.DB, error) {
	return sql.Open("sqlite", g.dbPath+"?mode=ro")
}

// parseWorkspace extracts the workspace (rig) from sender or assignee.
// Format: "rig/agent" or "mayor/" etc.
func parseWorkspace(addr string) string {
	if addr == "" {
		return "general"
	}
	parts := strings.SplitN(addr, "/", 2)
	if len(parts) > 0 && parts[0] != "" {
		return parts[0]
	}
	return "general"
}

func (g *GasTownSource) List(workspace string) ([]Message, error) {
	db, err := g.openDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT id, title, sender, assignee, description, created_at
		FROM issues
		WHERE issue_type = 'message'
		AND deleted_at IS NULL
		ORDER BY created_at ASC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var (
			id, title, description string
			sender, assignee       sql.NullString
			createdAt              string
		)

		if err := rows.Scan(&id, &title, &sender, &assignee, &description, &createdAt); err != nil {
			return nil, err
		}

		ts, _ := time.Parse(time.RFC3339Nano, createdAt)

		from := sender.String
		to := assignee.String
		ws := parseWorkspace(from)
		if ws == "general" {
			ws = parseWorkspace(to)
		}

		if workspace != "" && ws != workspace {
			continue
		}

		g.mu.Lock()
		g.workspaces[ws] = true
		g.mu.Unlock()

		messages = append(messages, Message{
			ID:        id,
			Workspace: ws,
			From:      from,
			To:        to,
			Body:      fmt.Sprintf("**%s**\n\n%s", title, description),
			Timestamp: ts,
			Source:    "gastown",
		})
	}

	return messages, nil
}

func (g *GasTownSource) Workspaces() ([]string, error) {
	// First populate workspaces by listing messages
	if _, err := g.List(""); err != nil {
		return nil, err
	}

	g.mu.RLock()
	defer g.mu.RUnlock()

	var result []string
	for ws := range g.workspaces {
		result = append(result, ws)
	}
	return result, nil
}

func (g *GasTownSource) Watch(ctx context.Context) (<-chan Message, error) {
	out := make(chan Message, 100)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	// Watch the db file for changes
	if err := watcher.Add(g.dbPath); err != nil {
		watcher.Close()
		return nil, err
	}

	// Also watch WAL file if it exists
	walPath := g.dbPath + "-wal"
	if _, err := os.Stat(walPath); err == nil {
		watcher.Add(walPath)
	}

	// Initialize last seen
	msgs, err := g.List("")
	if err == nil && len(msgs) > 0 {
		g.mu.Lock()
		g.lastSeenAt = msgs[len(msgs)-1].Timestamp
		g.lastSeenID = msgs[len(msgs)-1].ID
		g.mu.Unlock()
	}

	go func() {
		defer watcher.Close()
		defer close(out)

		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		checkForNew := func() {
			msgs, err := g.List("")
			if err != nil {
				return
			}

			g.mu.RLock()
			lastAt := g.lastSeenAt
			lastID := g.lastSeenID
			g.mu.RUnlock()

			for _, msg := range msgs {
				if msg.Timestamp.After(lastAt) || (msg.Timestamp.Equal(lastAt) && msg.ID != lastID) {
					select {
					case out <- msg:
					case <-ctx.Done():
						return
					}
					g.mu.Lock()
					g.lastSeenAt = msg.Timestamp
					g.lastSeenID = msg.ID
					g.mu.Unlock()
				}
			}
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
					checkForNew()
				}
			case <-watcher.Errors:
				// Ignore errors, continue watching
			case <-ticker.C:
				checkForNew()
			}
		}
	}()

	return out, nil
}
