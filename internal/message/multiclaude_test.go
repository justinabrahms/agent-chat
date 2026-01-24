package message

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWatch_DetectsNewDirectories(t *testing.T) {
	// Create a temporary messages directory structure
	tmpDir := t.TempDir()
	messagesDir := filepath.Join(tmpDir, "messages")
	if err := os.MkdirAll(messagesDir, 0755); err != nil {
		t.Fatalf("failed to create messages dir: %v", err)
	}

	// Create source
	source, err := NewMulticlaudeSource(tmpDir)
	if err != nil {
		t.Fatalf("failed to create source: %v", err)
	}

	// Start watching
	msgChan, err := source.Watch(t.Context())
	if err != nil {
		t.Fatalf("failed to start watch: %v", err)
	}

	// Give watcher time to initialize
	time.Sleep(100 * time.Millisecond)

	// Create a new repo/agent directory structure AFTER watch has started
	newRepoDir := filepath.Join(messagesDir, "new-repo")
	newAgentDir := filepath.Join(newRepoDir, "test-agent")
	if err := os.MkdirAll(newAgentDir, 0755); err != nil {
		t.Fatalf("failed to create new agent dir: %v", err)
	}

	// Give watcher time to detect the new directory
	time.Sleep(200 * time.Millisecond)

	// Create a message file in the new directory
	msg := MulticlaudeMessage{
		ID:        "test-msg-1",
		From:      "test-sender",
		To:        "test-receiver",
		Timestamp: time.Now(),
		Body:      "test message body",
		Status:    "pending",
	}
	msgData, _ := json.Marshal(msg)
	msgPath := filepath.Join(newAgentDir, "test-msg-1.json")
	if err := os.WriteFile(msgPath, msgData, 0644); err != nil {
		t.Fatalf("failed to write message file: %v", err)
	}

	// The message should be detected quickly (not waiting for 30s rescan)
	select {
	case received := <-msgChan:
		if received.ID != "test-msg-1" {
			t.Errorf("expected message ID 'test-msg-1', got '%s'", received.ID)
		}
		if received.Workspace != "mc-new-repo" {
			t.Errorf("expected workspace 'mc-new-repo', got '%s'", received.Workspace)
		}
	case <-time.After(2 * time.Second):
		t.Error("timeout waiting for message - new directory may not have been watched")
	}
}

func TestWatch_DetectsNestedDirectoriesCreatedSimultaneously(t *testing.T) {
	// Create a temporary messages directory structure
	tmpDir := t.TempDir()
	messagesDir := filepath.Join(tmpDir, "messages")
	if err := os.MkdirAll(messagesDir, 0755); err != nil {
		t.Fatalf("failed to create messages dir: %v", err)
	}

	// Create source
	source, err := NewMulticlaudeSource(tmpDir)
	if err != nil {
		t.Fatalf("failed to create source: %v", err)
	}

	// Start watching
	msgChan, err := source.Watch(t.Context())
	if err != nil {
		t.Fatalf("failed to start watch: %v", err)
	}

	// Give watcher time to initialize
	time.Sleep(100 * time.Millisecond)

	// Create nested directories simultaneously (like mkdir -p)
	nestedAgentDir := filepath.Join(messagesDir, "another-repo", "nested-agent")
	if err := os.MkdirAll(nestedAgentDir, 0755); err != nil {
		t.Fatalf("failed to create nested agent dir: %v", err)
	}

	// Give watcher time to detect
	time.Sleep(200 * time.Millisecond)

	// Create a message file
	msg := MulticlaudeMessage{
		ID:        "nested-msg-1",
		From:      "sender",
		To:        "receiver",
		Timestamp: time.Now(),
		Body:      "nested test",
		Status:    "pending",
	}
	msgData, _ := json.Marshal(msg)
	msgPath := filepath.Join(nestedAgentDir, "nested-msg-1.json")
	if err := os.WriteFile(msgPath, msgData, 0644); err != nil {
		t.Fatalf("failed to write message file: %v", err)
	}

	// Should detect without 30s wait
	select {
	case received := <-msgChan:
		if received.ID != "nested-msg-1" {
			t.Errorf("expected message ID 'nested-msg-1', got '%s'", received.ID)
		}
	case <-time.After(2 * time.Second):
		t.Error("timeout waiting for message in nested directory")
	}
}
