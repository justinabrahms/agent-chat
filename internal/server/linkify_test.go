package server

import (
	"strings"
	"testing"
)

func TestLinkifyURLs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
		excludes []string
	}{
		{
			name:     "HTTP URL",
			input:    "Check http://example.com for details",
			contains: []string{`<a href="http://example.com"`, `target="_blank"`, `rel="noopener noreferrer"`},
		},
		{
			name:     "HTTPS URL",
			input:    "See https://github.com/org/repo/pull/123",
			contains: []string{`<a href="https://github.com/org/repo/pull/123"`},
		},
		{
			name:     "URL with query params",
			input:    "Visit https://example.com/path?query=value&other=1",
			contains: []string{`href="https://example.com/path?query=value&other=1"`, `&amp;other=1</a>`},
		},
		{
			name:     "Multiple URLs",
			input:    "See https://a.com and https://b.com",
			contains: []string{`href="https://a.com"`, `href="https://b.com"`},
		},
		{
			name:     "No URL",
			input:    "Just plain text",
			excludes: []string{`<a href=`},
		},
		{
			name:     "URL at end with punctuation",
			input:    "Check https://example.com.",
			contains: []string{`href="https://example.com"`},
			excludes: []string{`href="https://example.com."`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := linkifyURLs(tt.input)
			for _, want := range tt.contains {
				if !strings.Contains(result, want) {
					t.Errorf("linkifyURLs(%q) = %q, want to contain %q", tt.input, result, want)
				}
			}
			for _, notWant := range tt.excludes {
				if strings.Contains(result, notWant) {
					t.Errorf("linkifyURLs(%q) = %q, should not contain %q", tt.input, result, notWant)
				}
			}
		})
	}
}

func TestLinkifyIssueRefs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		repoURL  string
		contains []string
		excludes []string
	}{
		{
			name:     "Single issue reference without repo",
			input:    "Fixed in #123",
			repoURL:  "",
			contains: []string{`<a href="https://github.com/search?q=123`, `#123`, `target="_blank"`, `data-issue="123"`},
		},
		{
			name:     "Single issue reference with repo",
			input:    "Fixed in #123",
			repoURL:  "https://github.com/owner/repo",
			contains: []string{`<a href="https://github.com/owner/repo/pull/123"`, `#123`, `target="_blank"`, `data-issue="123"`},
		},
		{
			name:     "Issue reference with .git suffix in repo URL",
			input:    "See #456",
			repoURL:  "https://github.com/owner/repo.git",
			contains: []string{`href="https://github.com/owner/repo/pull/456"`},
			excludes: []string{`.git/pull`},
		},
		{
			name:     "Multiple issue references",
			input:    "See #123 and #456",
			repoURL:  "",
			contains: []string{`data-issue="123"`, `data-issue="456"`},
		},
		{
			name:     "No issue reference",
			input:    "Just plain text",
			repoURL:  "",
			excludes: []string{`<a href=`},
		},
		{
			name:     "Hash with non-numeric (color code)",
			input:    "Color #ffffff",
			repoURL:  "",
			excludes: []string{`<a href=`},
		},
		{
			name:     "Hash channel name",
			input:    "Channel #general",
			repoURL:  "",
			excludes: []string{`<a href=`},
		},
		{
			name:     "Issue at start of line",
			input:    "#42 is the answer",
			repoURL:  "",
			contains: []string{`data-issue="42"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := linkifyIssueRefs(tt.input, tt.repoURL)
			for _, want := range tt.contains {
				if !strings.Contains(result, want) {
					t.Errorf("linkifyIssueRefs(%q, %q) = %q, want to contain %q", tt.input, tt.repoURL, result, want)
				}
			}
			for _, notWant := range tt.excludes {
				if strings.Contains(result, notWant) {
					t.Errorf("linkifyIssueRefs(%q, %q) = %q, should not contain %q", tt.input, tt.repoURL, result, notWant)
				}
			}
		})
	}
}
