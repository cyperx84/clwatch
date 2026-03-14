package state

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_DefaultState(t *testing.T) {
	// Clear env to use default path
	state, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if state.Schema != "clwatch.state.v1" {
		t.Errorf("expected schema=clwatch.state.v1, got %s", state.Schema)
	}
	if state.Tools == nil {
		t.Error("Tools map should be initialized")
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, "state.json")
	os.Setenv("CLWATCH_STATE_FILE", statePath)
	defer os.Unsetenv("CLWATCH_STATE_FILE")

	s := &State{
		Schema: "clwatch.state.v1",
		Tools: map[string]ToolState{
			"claude-code": {
				KnownVersion: "2.1.74",
				LastSeenAt:   "2026-03-12T00:00:00Z",
			},
		},
	}

	err := Save(s)
	if err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	ts, ok := loaded.Tools["claude-code"]
	if !ok {
		t.Fatal("claude-code not found in loaded state")
	}
	if ts.KnownVersion != "2.1.74" {
		t.Errorf("expected version=2.1.74, got %s", ts.KnownVersion)
	}
}

func TestUpdateState(t *testing.T) {
	s := &State{
		Schema: "clwatch.state.v1",
		Tools:  map[string]ToolState{},
	}

	// Simulate what UpdateState does
	s.Tools["claude-code"] = ToolState{
		KnownVersion: "2.1.75",
		LastSeenAt:   "2026-03-12T00:00:00Z",
	}

	ts, ok := s.Tools["claude-code"]
	if !ok {
		t.Fatal("tool not found after update")
	}
	if ts.KnownVersion != "2.1.75" {
		t.Errorf("expected 2.1.75, got %s", ts.KnownVersion)
	}
}
