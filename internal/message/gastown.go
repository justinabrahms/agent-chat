package message

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v3"
)

// GasTownSource reads messages from the Gas Town Dolt database.
type GasTownSource struct {
	dsn        string
	mu         sync.RWMutex
	lastSeenID string
	lastSeenAt time.Time
	workspaces map[string]bool
}

type gastownConfig struct {
	IssuePrefix string `yaml:"issue-prefix"`
}

// NewGasTownSource creates a new Gas Town message source.
// beadsDir should be the path to the .beads directory (e.g., ~/gt/.beads).
// It reads the dolt-server.port and config.yaml from that directory.
func NewGasTownSource(beadsDir string) (*GasTownSource, error) {
	portFile := filepath.Join(beadsDir, "dolt-server.port")
	portData, err := os.ReadFile(portFile)
	if err != nil {
		return nil, fmt.Errorf("dolt-server.port not found in %s: %w", beadsDir, err)
	}
	port, err := strconv.Atoi(strings.TrimSpace(string(portData)))
	if err != nil {
		return nil, fmt.Errorf("invalid port in dolt-server.port: %w", err)
	}

	configFile := filepath.Join(beadsDir, "config.yaml")
	configData, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("config.yaml not found in %s: %w", beadsDir, err)
	}
	var cfg gastownConfig
	if err := yaml.Unmarshal(configData, &cfg); err != nil {
		return nil, fmt.Errorf("invalid config.yaml in %s: %w", beadsDir, err)
	}
	if cfg.IssuePrefix == "" {
		return nil, fmt.Errorf("issue-prefix not set in %s/config.yaml", beadsDir)
	}

	dsn := fmt.Sprintf("root@tcp(127.0.0.1:%d)/%s?parseTime=true", port, cfg.IssuePrefix)

	// Verify the connection works
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open Dolt connection: %w", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("cannot connect to Dolt server at port %d: %w", port, err)
	}

	return &GasTownSource{
		dsn:        dsn,
		workspaces: make(map[string]bool),
	}, nil
}

func (g *GasTownSource) Name() string {
	return "gastown"
}

func (g *GasTownSource) openDB() (*sql.DB, error) {
	return sql.Open("mysql", g.dsn)
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

	// Mail messages are stored in both wisps (ephemeral, default) and issues (permanent).
	// The sender is in created_by (not sender), and there is no special issue_type='message'.
	// We filter for items that have both a sender and a recipient, excluding self-assignments.
	query := `
		SELECT id, title, created_by, assignee, description, created_at
		FROM wisps
		WHERE created_by != '' AND assignee != '' AND closed_at IS NULL
		UNION ALL
		SELECT id, title, created_by, assignee, description, created_at
		FROM issues
		WHERE created_by != '' AND assignee != '' AND created_by != assignee AND closed_at IS NULL
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
			createdAt              time.Time
		)

		if err := rows.Scan(&id, &title, &sender, &assignee, &description, &createdAt); err != nil {
			return nil, err
		}

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
			Timestamp: createdAt,
			Source:    "gastown",
		})
	}

	return messages, nil
}

func (g *GasTownSource) Workspaces() ([]string, error) {
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

	// Initialize last seen
	msgs, err := g.List("")
	if err == nil && len(msgs) > 0 {
		g.mu.Lock()
		g.lastSeenAt = msgs[len(msgs)-1].Timestamp
		g.lastSeenID = msgs[len(msgs)-1].ID
		g.mu.Unlock()
	}

	go func() {
		defer close(out)

		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				msgs, err := g.List("")
				if err != nil {
					continue
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
		}
	}()

	return out, nil
}
