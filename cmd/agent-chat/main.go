package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/justinabrahms/agent-chat/internal/message"
	"github.com/justinabrahms/agent-chat/internal/server"
)

func main() {
	var (
		port           int
		gastownDir     string
		multiclaudeDir string
	)

	homeDir, _ := os.UserHomeDir()
	defaultBeadsDir := filepath.Join(homeDir, ".beads")
	defaultMulticlaudeDir := filepath.Join(homeDir, ".multiclaude")

	flag.IntVar(&port, "port", 8080, "HTTP server port")
	flag.StringVar(&gastownDir, "gastown-dir", defaultBeadsDir, "Path to Gas Town .beads directory")
	flag.StringVar(&multiclaudeDir, "multiclaude-dir", defaultMulticlaudeDir, "Path to multiclaude directory")
	flag.Parse()

	// Override with environment variables if set
	if v := os.Getenv("PORT"); v != "" {
		fmt.Sscanf(v, "%d", &port)
	}
	if v := os.Getenv("GASTOWN_DIR"); v != "" {
		gastownDir = v
	}
	if v := os.Getenv("MULTICLAUDE_DIR"); v != "" {
		multiclaudeDir = v
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

	if len(sources) == 0 {
		log.Fatal("No message sources available. Configure at least one source.")
	}

	// Create aggregator
	agg := message.NewAggregator(sources...)

	// Create server
	srv, err := server.New(agg)
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
		httpServer.Shutdown(context.Background())
	}()

	log.Printf("Agent Chat running at http://localhost%s", addr)
	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server error: %v", err)
	}
}
