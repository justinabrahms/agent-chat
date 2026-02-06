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

// Avatar colors palette - works well in both dark and light themes.
// Each color is an HSL hue value (0-360).
var avatarHues = []int{
	0,   // red
	25,  // orange
	45,  // gold
	120, // green
	180, // cyan
	210, // blue
	260, // purple
	300, // magenta
	330, // pink
}

// stringHash returns a consistent hash value for a string.
func stringHash(s string) uint32 {
	var h uint32 = 0
	for _, c := range s {
		h = h*31 + uint32(c)
	}
	return h
}

// avatarColorIndex returns a consistent index into the color palette for a name.
func avatarColorIndex(name string) int {
	return int(stringHash(name) % uint32(len(avatarHues)))
}

// linkifyURLs replaces URLs with clickable links.
func linkifyURLs(s string) string {
	return urlRegex.ReplaceAllStringFunc(s, func(url string) string {
		return fmt.Sprintf(`<a href="%s" target="_blank" rel="noopener noreferrer">%s</a>`, url, template.HTMLEscapeString(url))
	})
}

// linkifyIssueRefs replaces GitHub-style issue references with links.
// If repoURL is provided (e.g., "https://github.com/owner/repo"), links go to that repo's PR/issue.
// Otherwise, links go to a generic GitHub search.
func linkifyIssueRefs(s string, repoURL string) string {
	return issueRefRegex.ReplaceAllStringFunc(s, func(ref string) string {
		num := ref[1:] // Remove the # prefix
		var href string
		if repoURL != "" {
			// Convert git URL to web URL and link to the pull request
			// Handle both https://github.com/owner/repo.git and https://github.com/owner/repo
			webURL := strings.TrimSuffix(repoURL, ".git")
			href = fmt.Sprintf("%s/pull/%s", webURL, num)
		} else {
			href = fmt.Sprintf("https://github.com/search?q=%s&type=issues", num)
		}
		return fmt.Sprintf(`<a href="%s" target="_blank" rel="noopener noreferrer" class="issue-ref" data-issue="%s">%s</a>`, href, num, ref)
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

// renderMarkdown converts markdown-ish text to HTML with code block support.
// repoURL is used to create proper GitHub links for issue/PR references.
func renderMarkdown(s string, repoURL string) string {
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
	s = linkifyIssueRefs(s, repoURL)

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

	// RepoURLs maps workspace names (e.g., "mc-agent-chat") to GitHub URLs.
	RepoURLs map[string]string

	// SSE subscribers
	mu          sync.RWMutex
	subscribers map[chan message.Message]bool
}

// New creates a new server with the given message aggregator.
// repoURLs maps workspace names (e.g., "mc-agent-chat") to GitHub URLs.
func New(agg *message.Aggregator, repoURLs map[string]string) (*Server, error) {
	srv := &Server{
		aggregator:  agg,
		RepoURLs:    repoURLs,
		subscribers: make(map[chan message.Message]bool),
	}

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
			case "claude-teams":
				return "👥"
			default:
				return "📨"
			}
		},
		"avatar": func(name, workspace, source string) template.HTML {
			if name == "" {
				name = "?"
			}
			// For AI agent sources (robots), use Robohash
			if source == "multiclaude" || source == "claude-teams" {
				// Use workspace (channel) and name (agent) for unique robot avatars
				robohashText := fmt.Sprintf("%s-%s", workspace, name)
				// URL encode the text for safety
				img := fmt.Sprintf(`<img class="avatar" width="24" height="24" src="https://robohash.org/%s?size=48x48" alt="%s avatar">`,
					template.URLQueryEscaper(robohashText), template.HTMLEscapeString(name))
				return template.HTML(img)
			}
			// For non-robots (humans), use the SVG avatar with initials
			// Get first character for initial
			initial := strings.ToUpper(string([]rune(name)[0]))
			// Get color from palette
			idx := avatarColorIndex(name)
			hue := avatarHues[idx]
			// Generate inline SVG avatar
			svg := fmt.Sprintf(`<svg class="avatar" width="24" height="24" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
				<circle cx="12" cy="12" r="12" fill="hsl(%d, 65%%, 50%%)"/>
				<text x="12" y="16" text-anchor="middle" fill="white" font-size="12" font-weight="600" font-family="-apple-system, BlinkMacSystemFont, sans-serif">%s</text>
			</svg>`, hue, initial)
			return template.HTML(svg)
		},
		"senderColorClass": func(name string) string {
			idx := avatarColorIndex(name)
			return fmt.Sprintf("sender-color-%d", idx)
		},
		"markdown": func(body, workspace string) template.HTML {
			repoURL := srv.RepoURLs[workspace]
			return template.HTML(renderMarkdown(body, repoURL))
		},
		"isStatusMessage": func(body string) bool {
			return strings.HasPrefix(body, "[status] ")
		},
		"stripStatusPrefix": func(body string) string {
			return strings.TrimPrefix(body, "[status] ")
		},
	}

	tmpl, err := template.New("").Funcs(funcMap).ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("parsing templates: %w", err)
	}

	srv.templates = tmpl
	return srv, nil
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
