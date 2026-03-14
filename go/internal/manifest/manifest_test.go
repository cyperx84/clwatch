package manifest

import (
	"os"
	"testing"
)

func TestManifestURL_Default(t *testing.T) {
	os.Unsetenv("CLWATCH_MANIFEST_URL")
	url := ManifestURL()
	if url != DefaultManifestURL {
		t.Errorf("expected %s, got %s", DefaultManifestURL, url)
	}
}

func TestManifestURL_Override(t *testing.T) {
	os.Setenv("CLWATCH_MANIFEST_URL", "https://example.com/manifest.json")
	defer os.Unsetenv("CLWATCH_MANIFEST_URL")
	url := ManifestURL()
	if url != "https://example.com/manifest.json" {
		t.Errorf("expected override URL, got %s", url)
	}
}

func TestIsStale(t *testing.T) {
	tests := []struct {
		name     string
		staleAt  string
		expected bool
	}{
		{"future date", "2099-12-31T00:00:00Z", false},
		{"past date", "2020-01-01T00:00:00Z", true},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsStale(tt.staleAt)
			if result != tt.expected {
				t.Errorf("IsStale(%s) = %v, want %v", tt.staleAt, result, tt.expected)
			}
		})
	}
}

func TestToolStruct(t *testing.T) {
	tool := Tool{
		CurrentVersion: "2.1.74",
		StaleAfter:     "2026-03-19T00:00:00Z",
	}
	if tool.CurrentVersion != "2.1.74" {
		t.Errorf("expected version 2.1.74, got %s", tool.CurrentVersion)
	}
	if IsStale(tool.StaleAfter) {
		t.Error("future date should not be stale")
	}
}
