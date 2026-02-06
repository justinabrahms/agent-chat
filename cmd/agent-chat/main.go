package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/justinabrahms/agent-chat/internal/config"
	"github.com/justinabrahms/agent-chat/internal/message"
	"github.com/justinabrahms/agent-chat/internal/server"
	"github.com/justinabrahms/agent-chat/internal/version"
)

// loadRepoURLs reads the multiclaude state.json and extracts repo URL mappings.
// Returns a map from workspace name (e.g., "mc-agent-chat") to GitHub URL.
func loadRepoURLs(multiclaudeDir string) map[string]string {
	result := make(map[string]string)

	statePath := filepath.Join(multiclaudeDir, "state.json")
	data, err := os.ReadFile(statePath)
	if err != nil {
		return result
	}

	var state struct {
		Repos map[string]struct {
			GitHubURL string `json:"github_url"`
		} `json:"repos"`
	}

	if err := json.Unmarshal(data, &state); err != nil {
		return result
	}

	for repoName, repo := range state.Repos {
		workspace := "mc-" + repoName
		result[workspace] = repo.GitHubURL
	}

	return result
}

func main() {
	var (
		port           int
		gastownDir     string
		multiclaudeDir string
		claudeTeamsDir string
		configPath     string
		showVersion    bool
	)

	homeDir, _ := os.UserHomeDir()
	defaultBeadsDir := filepath.Join(homeDir, ".beads")
	defaultMulticlaudeDir := filepath.Join(homeDir, ".multiclaude")
	defaultClaudeTeamsDir := filepath.Join(homeDir, ".claude", "teams")

	flag.IntVar(&port, "port", 8080, "HTTP server port")
	flag.StringVar(&gastownDir, "gastown-dir", defaultBeadsDir, "Path to Gas Town .beads directory")
	flag.StringVar(&multiclaudeDir, "multiclaude-dir", defaultMulticlaudeDir, "Path to multiclaude directory")
	flag.StringVar(&claudeTeamsDir, "claude-teams-dir", defaultClaudeTeamsDir, "Path to Claude Teams directory")
	flag.StringVar(&configPath, "config", "", "Path to config file (default: ~/.config/agent-chat/config.yaml)")
	flag.BoolVar(&showVersion, "version", false, "Print version information and exit")
	flag.Parse()

	if showVersion {
		fmt.Println(version.String())
		os.Exit(0)
	}

	// Track which flags were explicitly set
	explicitFlags := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) {
		explicitFlags[f.Name] = true
	})

	// Load config file
	// If --config was specified, it's an error if the file doesn't exist
	cfg, err := config.Load(configPath, explicitFlags["config"])
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Apply precedence: flags > env vars > config file > defaults
	// Config file values (only if flag wasn't explicitly set)
	if !explicitFlags["port"] && cfg.Port != 0 {
		port = cfg.Port
	}
	if !explicitFlags["gastown-dir"] && cfg.GastownDir != "" {
		gastownDir = cfg.GastownDir
	}
	if !explicitFlags["multiclaude-dir"] && cfg.MulticlaudeDir != "" {
		multiclaudeDir = cfg.MulticlaudeDir
	}
	if !explicitFlags["claude-teams-dir"] && cfg.ClaudeTeamsDir != "" {
		claudeTeamsDir = cfg.ClaudeTeamsDir
	}

	// Environment variables override config file (but not flags)
	if v := os.Getenv("PORT"); v != "" && !explicitFlags["port"] {
		_, _ = fmt.Sscanf(v, "%d", &port)
	}
	if v := os.Getenv("GASTOWN_DIR"); v != "" && !explicitFlags["gastown-dir"] {
		gastownDir = v
	}
	if v := os.Getenv("MULTICLAUDE_DIR"); v != "" && !explicitFlags["multiclaude-dir"] {
		multiclaudeDir = v
	}
	if v := os.Getenv("CLAUDE_TEAMS_DIR"); v != "" && !explicitFlags["claude-teams-dir"] {
		claudeTeamsDir = v
	}

	// Collect available sources
	var sources []message.Source

	// Try Gas Town source
	if gs, err := message.NewGasTownSource(gastownDir); err == nil {
		log.Printf("Loaded Gas Town source from %s", gastownDir)
		sources = append(sources, gs)
	} else {
		log.Printf("Gas Town source not available: %v", err)
	}

	// Try multiclaude source
	if ms, err := message.NewMulticlaudeSource(multiclaudeDir); err == nil {
		log.Printf("Loaded multiclaude source from %s", multiclaudeDir)
		sources = append(sources, ms)
	} else {
		log.Printf("Multiclaude source not available: %v", err)
	}

	// Try Claude Teams source
	if ts, err := message.NewClaudeTeamsSource(claudeTeamsDir); err == nil {
		log.Printf("Loaded Claude Teams source from %s", claudeTeamsDir)
		sources = append(sources, ts)
	} else {
		log.Printf("Claude Teams source not available: %v", err)
	}

	if len(sources) == 0 {
		log.Fatal("No message sources available. Configure at least one source.")
	}

	// Create aggregator
	agg := message.NewAggregator(sources...)

	// Load repo URL mappings for PR links
	repoURLs := loadRepoURLs(multiclaudeDir)

	// Create server
	srv, err := server.New(agg, repoURLs)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Start message watching
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := srv.Start(ctx); err != nil {
		log.Fatalf("Failed to start message watching: %v", err)
	}

	// Start HTTP server
	addr := fmt.Sprintf(":%d", port)
	httpServer := &http.Server{
		Addr:    addr,
		Handler: srv.Handler(),
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("Shutting down...")
		cancel()
		if err := httpServer.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
	}()

	log.Printf("Agent Chat running at http://localhost%s", addr)
	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server error: %v", err)
	}
}
