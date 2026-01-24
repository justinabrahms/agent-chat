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
		contains []string
		excludes []string
	}{
		{
			name:     "Single issue reference",
			input:    "Fixed in #123",
			contains: []string{`<a href=`, `#123`, `target="_blank"`, `data-issue="123"`},
		},
		{
			name:     "Multiple issue references",
			input:    "See #123 and #456",
			contains: []string{`data-issue="123"`, `data-issue="456"`},
		},
		{
			name:     "No issue reference",
			input:    "Just plain text",
			excludes: []string{`<a href=`},
		},
		{
			name:     "Hash with non-numeric (color code)",
			input:    "Color #ffffff",
			excludes: []string{`<a href=`},
		},
		{
			name:     "Hash channel name",
			input:    "Channel #general",
			excludes: []string{`<a href=`},
		},
		{
			name:     "Issue at start of line",
			input:    "#42 is the answer",
			contains: []string{`data-issue="42"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := linkifyIssueRefs(tt.input)
			for _, want := range tt.contains {
				if !strings.Contains(result, want) {
					t.Errorf("linkifyIssueRefs(%q) = %q, want to contain %q", tt.input, result, want)
				}
			}
			for _, notWant := range tt.excludes {
				if strings.Contains(result, notWant) {
					t.Errorf("linkifyIssueRefs(%q) = %q, should not contain %q", tt.input, result, notWant)
				}
			}
		})
	}
}
