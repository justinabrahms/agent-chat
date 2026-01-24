package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_ValidConfig(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")

	content := `port: 9000
gastown-dir: /custom/beads
multiclaude-dir: /custom/multiclaude
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(configPath, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Port != 9000 {
		t.Errorf("port = %d, want 9000", cfg.Port)
	}
	if cfg.GastownDir != "/custom/beads" {
		t.Errorf("gastown-dir = %s, want /custom/beads", cfg.GastownDir)
	}
	if cfg.MulticlaudeDir != "/custom/multiclaude" {
		t.Errorf("multiclaude-dir = %s, want /custom/multiclaude", cfg.MulticlaudeDir)
	}
}

func TestLoad_PartialConfig(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")

	content := `port: 8888
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(configPath, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Port != 8888 {
		t.Errorf("port = %d, want 8888", cfg.Port)
	}
	if cfg.GastownDir != "" {
		t.Errorf("gastown-dir = %s, want empty", cfg.GastownDir)
	}
}

func TestLoad_MissingFile_Explicit(t *testing.T) {
	_, err := Load("/nonexistent/config.yaml", true)
	if err == nil {
		t.Fatal("expected error for missing explicit config")
	}
}

func TestLoad_MissingFile_Default(t *testing.T) {
	cfg, err := Load("/nonexistent/config.yaml", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")

	content := `port: not-a-number
this is: [invalid yaml
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(configPath, true)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestLoad_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")

	if err := os.WriteFile(configPath, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(configPath, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Port != 0 {
		t.Errorf("port = %d, want 0", cfg.Port)
	}
}

func TestLoad_UnknownKeys(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")

	content := `port: 9000
unknown-key: some-value
another-unknown: 123
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(configPath, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Unknown keys should be ignored, known keys should work
	if cfg.Port != 9000 {
		t.Errorf("port = %d, want 9000", cfg.Port)
	}
}

func TestDefaultPath(t *testing.T) {
	path := DefaultPath()
	if path == "" {
		t.Skip("could not determine default path")
	}

	// Should contain agent-chat and config.yaml
	if filepath.Base(path) != "config.yaml" {
		t.Errorf("default path should end with config.yaml, got %s", path)
	}
	if filepath.Base(filepath.Dir(path)) != "agent-chat" {
		t.Errorf("default path should be in agent-chat dir, got %s", path)
	}
}
