package server

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/justinabrahms/agent-chat/internal/message"
)

// Regular expressions for link detection
var (
	// Match URLs starting with http:// or https://
	urlRegex = regexp.MustCompile(`https?://[^\s<>"'` + "`" + `]+[^\s<>"'` + "`" + `.,;:!?)}\]]+`)
	// Match GitHub-style issue/PR references like #123
	issueRefRegex = regexp.MustCompile(`#(\d+)\b`)
)

// linkifyURLs replaces URLs with clickable links.
func linkifyURLs(s string) string {
	return urlRegex.ReplaceAllStringFunc(s, func(url string) string {
		return fmt.Sprintf(`<a href="%s" target="_blank" rel="noopener noreferrer">%s</a>`, url, template.HTMLEscapeString(url))
	})
}

// linkifyIssueRefs replaces GitHub-style issue references with links.
func linkifyIssueRefs(s string) string {
	return issueRefRegex.ReplaceAllStringFunc(s, func(ref string) string {
		num := ref[1:] // Remove the # prefix
		return fmt.Sprintf(`<a href="https://github.com/search?q=%s&type=issues" target="_blank" rel="noopener noreferrer" class="issue-ref" data-issue="%s">%s</a>`, num, num, ref)
	})
}

//go:embed templates/*.html
var templatesFS embed.FS

//go:embed static/*
var staticFS embed.FS

// WorkspaceInfo contains workspace metadata for the UI.
type WorkspaceInfo struct {
	Name          string
	LatestMsgTime int64 // Unix milliseconds for JavaScript compatibility
	MessageCount  int
}

// GroupedMessage wraps a message with grouping metadata for the UI.
type GroupedMessage struct {
	message.Message
	IsGroupStart bool // True if this is the first message in a group from this sender
	IsGroupEnd   bool // True if this is the last message in a group from this sender
}

// groupMessages adds grouping metadata to a slice of chronologically sorted messages.
func groupMessages(msgs []message.Message) []GroupedMessage {
	if len(msgs) == 0 {
		return nil
	}

	grouped := make([]GroupedMessage, len(msgs))
	for i, msg := range msgs {
		grouped[i] = GroupedMessage{Message: msg}

		// Check if this is the start of a new group
		if i == 0 || msgs[i-1].From != msg.From {
			grouped[i].IsGroupStart = true
		}

		// Check if this is the end of a group
		if i == len(msgs)-1 || msgs[i+1].From != msg.From {
			grouped[i].IsGroupEnd = true
		}
	}
	return grouped
}

// codeBlockRegex matches fenced code blocks with optional language hint
var codeBlockRegex = regexp.MustCompile("(?s)```(\\w*)\\n?(.*?)```")

// renderMarkdown converts markdown-ish text to HTML with code block support
func renderMarkdown(s string) string {
	// First, extract and replace code blocks with placeholders
	var codeBlocks []string
	placeholder := "\x00CODE_BLOCK_%d\x00"

	s = codeBlockRegex.ReplaceAllStringFunc(s, func(match string) string {
		parts := codeBlockRegex.FindStringSubmatch(match)
		lang := parts[1]
		code := parts[2]

		// Trim trailing newline from code
		code = strings.TrimSuffix(code, "\n")

		// Escape HTML in code
		code = html.EscapeString(code)

		var block string
		if lang != "" {
			block = fmt.Sprintf("<pre><code class=\"language-%s\">%s</code></pre>", lang, code)
		} else {
			block = fmt.Sprintf("<pre><code>%s</code></pre>", code)
		}

		idx := len(codeBlocks)
		codeBlocks = append(codeBlocks, block)
		return fmt.Sprintf(placeholder, idx)
	})

	// Linkify URLs and issue references (before other processing)
	s = linkifyURLs(s)
	s = linkifyIssueRefs(s)

	// Process inline markdown (simple: just remove **)
	s = strings.ReplaceAll(s, "**", "")

	// Convert newlines to <br>
	s = strings.ReplaceAll(s, "\n", "<br>")

	// Restore code blocks
	for i, block := range codeBlocks {
		s = strings.Replace(s, fmt.Sprintf(placeholder, i), block, 1)
	}

	return s
}

// Server handles HTTP requests for the chat UI.
type Server struct {
	aggregator *message.Aggregator
	templates  *template.Template

	// SSE subscribers
	mu          sync.RWMutex
	subscribers map[chan message.Message]bool
}

// New creates a new server with the given message aggregator.
func New(agg *message.Aggregator) (*Server, error) {
	funcMap := template.FuncMap{
		"formatTime": func(t time.Time) string {
			return t.Format("15:04")
		},
		"formatDate": func(t time.Time) string {
			return t.Format("Jan 2")
		},
		"sourceIcon": func(source string) string {
			switch source {
			case "gastown":
				return "⛽"
			case "multiclaude":
				return "🤖"
			default:
				return "📨"
			}
		},
		"markdown": func(s string) template.HTML {
			return template.HTML(renderMarkdown(s))
		},
	}

	tmpl, err := template.New("").Funcs(funcMap).ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("parsing templates: %w", err)
	}

	return &Server{
		aggregator:  agg,
		templates:   tmpl,
		subscribers: make(map[chan message.Message]bool),
	}, nil
}

// Start begins broadcasting messages from the aggregator to SSE subscribers.
func (s *Server) Start(ctx context.Context) error {
	msgCh, err := s.aggregator.Watch(ctx)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgCh {
			s.broadcast(msg)
		}
	}()

	return nil
}

func (s *Server) broadcast(msg message.Message) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for ch := range s.subscribers {
		select {
		case ch <- msg:
		default:
			// Skip slow subscribers
		}
	}
}

func (s *Server) subscribe() chan message.Message {
	ch := make(chan message.Message, 100)
	s.mu.Lock()
	s.subscribers[ch] = true
	s.mu.Unlock()
	return ch
}

func (s *Server) unsubscribe(ch chan message.Message) {
	s.mu.Lock()
	delete(s.subscribers, ch)
	s.mu.Unlock()
	close(ch)
}

// getWorkspaceInfos returns workspace metadata including latest message timestamps.
func (s *Server) getWorkspaceInfos() []WorkspaceInfo {
	workspaces, _ := s.aggregator.Workspaces()
	sort.Strings(workspaces)

	infos := make([]WorkspaceInfo, 0, len(workspaces))
	for _, ws := range workspaces {
		msgs, _ := s.aggregator.List(ws)
		var latestTime int64
		for _, msg := range msgs {
			ts := msg.Timestamp.UnixMilli()
			if ts > latestTime {
				latestTime = ts
			}
		}
		infos = append(infos, WorkspaceInfo{
			Name:          ws,
			LatestMsgTime: latestTime,
			MessageCount:  len(msgs),
		})
	}
	return infos
}

// Handler returns an http.Handler for the server.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	// Static files
	mux.Handle("/static/", http.FileServer(http.FS(staticFS)))

	// Main page
	mux.HandleFunc("/", s.handleIndex)

	// Workspaces sidebar (HTMX partial)
	mux.HandleFunc("/workspaces", s.handleWorkspaces)

	// Messages for a workspace (HTMX partial)
	mux.HandleFunc("/messages", s.handleMessages)

	// SSE endpoint
	mux.HandleFunc("/events", s.handleSSE)

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	return mux
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	workspaceInfos := s.getWorkspaceInfos()

	selectedWS := r.URL.Query().Get("workspace")
	if selectedWS == "" && len(workspaceInfos) > 0 {
		selectedWS = workspaceInfos[0].Name
	}

	messages, _ := s.aggregator.List(selectedWS)
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Timestamp.Before(messages[j].Timestamp)
	})
	groupedMsgs := groupMessages(messages)

	data := map[string]any{
		"Workspaces":        workspaceInfos,
		"SelectedWorkspace": selectedWS,
		"Messages":          groupedMsgs,
		"Sources":           s.aggregator.Sources(),
	}

	if err := s.templates.ExecuteTemplate(w, "index.html", data); err != nil {
		log.Printf("template error: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

func (s *Server) handleWorkspaces(w http.ResponseWriter, r *http.Request) {
	workspaceInfos := s.getWorkspaceInfos()
	selectedWS := r.URL.Query().Get("selected")

	data := map[string]any{
		"Workspaces":        workspaceInfos,
		"SelectedWorkspace": selectedWS,
	}

	if err := s.templates.ExecuteTemplate(w, "workspaces.html", data); err != nil {
		log.Printf("template error: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

func (s *Server) handleMessages(w http.ResponseWriter, r *http.Request) {
	workspace := r.URL.Query().Get("workspace")

	messages, _ := s.aggregator.List(workspace)
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Timestamp.Before(messages[j].Timestamp)
	})
	groupedMsgs := groupMessages(messages)

	data := map[string]any{
		"Messages":  groupedMsgs,
		"Workspace": workspace,
	}

	if err := s.templates.ExecuteTemplate(w, "messages.html", data); err != nil {
		log.Printf("template error: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

func (s *Server) handleSSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	ch := s.subscribe()
	defer s.unsubscribe(ch)

	// Send initial connection event
	fmt.Fprintf(w, "event: connected\ndata: {\"status\": \"connected\"}\n\n")
	flusher.Flush()

	// Keep-alive ticker
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			// Render the message as HTML using the template
			var buf bytes.Buffer
			if err := s.templates.ExecuteTemplate(&buf, "message.html", msg); err != nil {
				log.Printf("SSE template error: %v", err)
				continue
			}
			// Send HTML as SSE data (newlines in data need to be prefixed with "data: ")
			html := strings.ReplaceAll(buf.String(), "\n", "\ndata: ")
			fmt.Fprintf(w, "event: message\ndata: %s\n\n", html)
			flusher.Flush()
		case <-ticker.C:
			fmt.Fprintf(w, ": keepalive\n\n")
			flusher.Flush()
		}
	}
}
