package manifest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const DefaultManifestURL = "https://changelogs.info/api/refs/manifest.json"

// Tool represents a single tool entry in the manifest.
type Tool struct {
	CurrentVersion string          `json:"version"`
	StaleAfter     string          `json:"stale_after"`
	Delta          json.RawMessage `json:"delta,omitempty"`
}

// Manifest is the top-level payload from changelogs.info/api/refs/manifest.json
type Manifest struct {
	Schema    string          `json:"schema"`
	Generated string          `json:"generated"`
	Tools     map[string]Tool `json:"tools"`
}

func ManifestURL() string {
	if u := os.Getenv("CLWATCH_MANIFEST_URL"); u != "" {
		return u
	}
	return DefaultManifestURL
}

func Fetch() (*Manifest, error) {
	url := ManifestURL()

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetching manifest: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("manifest returned HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading manifest body: %w", err)
	}

	var m Manifest
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, fmt.Errorf("parsing manifest JSON: %w", err)
	}

	return &m, nil
}

// IsStale returns true if the stale_after absolute ISO8601 timestamp has passed.
func IsStale(staleAfter string) bool {
	if staleAfter == "" {
		return false
	}
	t, err := time.Parse(time.RFC3339, staleAfter)
	if err != nil {
		return false
	}
	return time.Now().UTC().After(t)
}

// ParseStaleAfter parses a duration string like "7d" or "24h" into a time.Duration.
func ParseStaleAfter(s string) (time.Duration, error) {
	if s == "" {
		return 7 * 24 * time.Hour, nil // default 7 days
	}
	last := s[len(s)-1]
	switch last {
	case 'd':
		var days int
		if _, err := fmt.Sscanf(s, "%dd", &days); err != nil {
			return 0, fmt.Errorf("invalid stale_after: %s", s)
		}
		return time.Duration(days) * 24 * time.Hour, nil
	case 'h':
		return time.ParseDuration(s)
	default:
		return time.ParseDuration(s)
	}
}
