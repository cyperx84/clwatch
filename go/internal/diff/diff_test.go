package diff

import (
	"testing"

	"github.com/cyperx/clwatch/internal/manifest"
	"github.com/cyperx/clwatch/internal/state"
)

func TestCompare_NoChanges(t *testing.T) {
	m := &manifest.Manifest{
		Schema:    "clwatch.manifest.v1",
		Generated: "2026-03-12T00:00:00Z",
		Tools: map[string]manifest.Tool{
			"claude-code": {CurrentVersion: "2.1.74"},
		},
	}
	s := &state.State{
		Schema: "clwatch.state.v1",
		Tools: map[string]state.ToolState{
			"claude-code": {KnownVersion: "2.1.74"},
		},
	}

	results := Compare(m, s)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != StatusCurrent {
		t.Errorf("expected status=current, got %s", results[0].Status)
	}
	if results[0].Tool != "claude-code" {
		t.Errorf("expected tool=claude-code, got %s", results[0].Tool)
	}
}

func TestCompare_VersionChanged(t *testing.T) {
	m := &manifest.Manifest{
		Schema:    "clwatch.manifest.v1",
		Generated: "2026-03-12T00:00:00Z",
		Tools: map[string]manifest.Tool{
			"claude-code": {CurrentVersion: "2.1.75"},
		},
	}
	s := &state.State{
		Schema: "clwatch.state.v1",
		Tools: map[string]state.ToolState{
			"claude-code": {KnownVersion: "2.1.74"},
		},
	}

	results := Compare(m, s)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != StatusUpdated {
		t.Errorf("expected status=updated, got %s", results[0].Status)
	}
	if results[0].PreviousVersion != "2.1.74" {
		t.Errorf("expected previous=2.1.74, got %s", results[0].PreviousVersion)
	}
	if results[0].CurrentVersion != "2.1.75" {
		t.Errorf("expected current=2.1.75, got %s", results[0].CurrentVersion)
	}
}

func TestCompare_NewTool(t *testing.T) {
	m := &manifest.Manifest{
		Schema:    "clwatch.manifest.v1",
		Generated: "2026-03-12T00:00:00Z",
		Tools: map[string]manifest.Tool{
			"new-tool": {CurrentVersion: "1.0.0"},
		},
	}
	s := &state.State{
		Schema: "clwatch.state.v1",
		Tools:  map[string]state.ToolState{},
	}

	results := Compare(m, s)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != StatusNew {
		t.Errorf("expected status=new, got %s", results[0].Status)
	}
}

func TestCompare_MultipleTools(t *testing.T) {
	m := &manifest.Manifest{
		Schema:    "clwatch.manifest.v1",
		Generated: "2026-03-12T00:00:00Z",
		Tools: map[string]manifest.Tool{
			"claude-code": {CurrentVersion: "2.1.74"},
			"codex-cli":   {CurrentVersion: "0.115.0"},
			"gemini-cli":  {CurrentVersion: "0.33.0"},
		},
	}
	s := &state.State{
		Schema: "clwatch.state.v1",
		Tools: map[string]state.ToolState{
			"claude-code": {KnownVersion: "2.1.74"}, // current
			"codex-cli":   {KnownVersion: "0.114.0"}, // updated
			// gemini-cli not in state → new
		},
	}

	results := Compare(m, s)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	statuses := map[string]Status{}
	for _, r := range results {
		statuses[r.Tool] = r.Status
	}

	if statuses["claude-code"] != StatusCurrent {
		t.Errorf("claude-code: expected current, got %s", statuses["claude-code"])
	}
	if statuses["codex-cli"] != StatusUpdated {
		t.Errorf("codex-cli: expected updated, got %s", statuses["codex-cli"])
	}
	if statuses["gemini-cli"] != StatusNew {
		t.Errorf("gemini-cli: expected new, got %s", statuses["gemini-cli"])
	}
}
