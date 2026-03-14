package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig([]string{"claude-code", "codex-cli"})
	if cfg.Schema != "clwatch.config.v1" {
		t.Errorf("expected schema=clwatch.config.v1, got %s", cfg.Schema)
	}
	if len(cfg.Tools) != 2 {
		t.Errorf("expected 2 tools, got %d", len(cfg.Tools))
	}
	if cfg.ReferenceDir != "references/" {
		t.Errorf("expected references/, got %s", cfg.ReferenceDir)
	}
	if cfg.Tier2Threshold != "medium" {
		t.Errorf("expected medium, got %s", cfg.Tier2Threshold)
	}
	if !cfg.NotifyOnBreaking {
		t.Error("expected NotifyOnBreaking=true")
	}
}

func TestInitWorkspace(t *testing.T) {
	tmpDir := t.TempDir()

	actions, err := Init(tmpDir, []string{"claude-code", "codex-cli"}, false)
	if err != nil {
		t.Fatalf("Init error: %v", err)
	}

	if len(actions) == 0 {
		t.Error("expected actions to be non-empty")
	}

	configPath := filepath.Join(tmpDir, ".clwatch.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error(".clwatch.json not created")
	}

	refDir := filepath.Join(tmpDir, "references")
	if _, err := os.Stat(refDir); os.IsNotExist(err) {
		t.Error("references/ directory not created")
	}
}

func TestInitForce(t *testing.T) {
	tmpDir := t.TempDir()

	_, err := Init(tmpDir, []string{"claude-code"}, false)
	if err != nil {
		t.Fatalf("First Init error: %v", err)
	}

	// Force overwrite
	actions, err := Init(tmpDir, []string{"claude-code", "codex-cli"}, true)
	if err != nil {
		t.Fatalf("Force Init error: %v", err)
	}
	if len(actions) == 0 {
		t.Error("expected actions from force init")
	}
}
