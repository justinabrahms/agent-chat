package message

import (
	"context"
	"time"
)

// Message represents a single message from any source.
type Message struct {
	ID        string    `json:"id"`
	Workspace string    `json:"workspace"`
	From      string    `json:"from"`
	To        string    `json:"to,omitempty"`
	Body      string    `json:"body"`
	Timestamp time.Time `json:"timestamp"`
	Source    string    `json:"source"` // "gastown" or "multiclaude"
}

// Source defines the interface for reading messages from different backends.
type Source interface {
	// Name returns the source identifier (e.g., "gastown", "multiclaude").
	Name() string

	// Watch starts watching for new messages and sends them to the returned channel.
	// The channel is closed when the context is cancelled.
	Watch(ctx context.Context) (<-chan Message, error)

	// List returns existing messages, optionally filtered by workspace.
	// If workspace is empty, returns messages from all workspaces.
	List(workspace string) ([]Message, error)

	// Workspaces returns all known workspace identifiers.
	Workspaces() ([]string, error)
}

// Aggregator combines multiple sources into a single stream.
type Aggregator struct {
	sources []Source
}

// NewAggregator creates a new aggregator from the given sources.
func NewAggregator(sources ...Source) *Aggregator {
	return &Aggregator{sources: sources}
}

// Watch starts watching all sources and multiplexes messages into a single channel.
func (a *Aggregator) Watch(ctx context.Context) (<-chan Message, error) {
	out := make(chan Message, 100)

	for _, src := range a.sources {
		ch, err := src.Watch(ctx)
		if err != nil {
			return nil, err
		}

		go func(ch <-chan Message) {
			for msg := range ch {
				select {
				case out <- msg:
				case <-ctx.Done():
					return
				}
			}
		}(ch)
	}

	go func() {
		<-ctx.Done()
		close(out)
	}()

	return out, nil
}

// List returns messages from all sources.
func (a *Aggregator) List(workspace string) ([]Message, error) {
	var all []Message
	for _, src := range a.sources {
		msgs, err := src.List(workspace)
		if err != nil {
			return nil, err
		}
		all = append(all, msgs...)
	}
	return all, nil
}

// Workspaces returns all unique workspaces from all sources.
func (a *Aggregator) Workspaces() ([]string, error) {
	seen := make(map[string]bool)
	var result []string

	for _, src := range a.sources {
		ws, err := src.Workspaces()
		if err != nil {
			return nil, err
		}
		for _, w := range ws {
			if !seen[w] {
				seen[w] = true
				result = append(result, w)
			}
		}
	}
	return result, nil
}

// Sources returns the list of configured sources.
func (a *Aggregator) Sources() []Source {
	return a.sources
}
