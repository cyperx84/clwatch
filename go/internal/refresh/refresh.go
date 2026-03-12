package refresh

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cyperx/clwatch/internal/manifest"
	"github.com/cyperx/clwatch/internal/state"
)

// Payload is the full tool payload from changelogs.info/api/refs/<tool>.json
type Payload struct {
	Schema     string          `json:"schema"`
	Tool       string          `json:"tool"`
	Version    string          `json:"version"`
	Generated  string          `json:"generated_at"`
	StaleAfter string          `json:"stale_after"`
	Verify     Verification    `json:"verification"`
	Delta      Delta           `json:"delta"`
	RawPayload json.RawMessage `json:"payload,omitempty"`
}

type Verification struct {
	Status           string   `json:"status"`
	Confidence       float64  `json:"confidence"`
	UnverifiedFields []string `json:"unverified_fields"`
}

type Delta struct {
	FromVersion        string   `json:"from_version"`
	ToVersion          string   `json:"to_version"`
	NewFeatures        []string `json:"new_features"`
	NewCommands        []string `json:"new_commands"`
	NewFlags           []string `json:"new_flags"`
	NewEnvVars         []string `json:"new_env_vars"`
	DeprecatedCommands []string `json:"deprecated_commands"`
	DeprecatedFlags    []string `json:"deprecated_flags"`
	BreakingChanges    []string `json:"breaking_changes"`
}

// FetchPayload fetches the full payload JSON for a tool.
func FetchPayload(m *manifest.Manifest, toolID string) (*Payload, []byte, error) {
	url := payloadURL(m, toolID)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, nil, fmt.Errorf("fetching payload for %s: %w", toolID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("payload for %s returned HTTP %d", toolID, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("reading payload body for %s: %w", toolID, err)
	}

	var p Payload
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, nil, fmt.Errorf("parsing payload JSON for %s: %w", toolID, err)
	}

	return &p, body, nil
}

// payloadURL constructs the URL for a tool's payload.
func payloadURL(m *manifest.Manifest, toolID string) string {
	// Use payload_url from manifest if available
	if tool, ok := m.Tools[toolID]; ok && tool.PayloadURL != "" {
		return tool.PayloadURL
	}
	// Fallback: replace manifest.json with <tool>.json in manifest URL
	base := manifest.ManifestURL()
	return strings.TrimSuffix(base, "manifest.json") + toolID + ".json"
}

// UpdateStateForTool updates local state after a successful refresh.
func UpdateStateForTool(s *state.State, toolID string, version string) {
	s.Tools[toolID] = state.ToolState{
		KnownVersion: version,
		LastSeenAt:   time.Now().UTC().Format(time.RFC3339),
	}
}

// PrintSummary prints a human-readable summary of the payload.
func PrintSummary(p *Payload) {
	fmt.Printf("%s %s (%s)\n", p.Tool, p.Version, p.Verify.Status)
	fmt.Printf("Generated: %s\n", formatDate(p.Generated))
	fmt.Printf("Stale after: %s\n", formatDate(p.StaleAfter))

	d := &p.Delta
	newFeatureCount := len(d.NewFeatures)
	newFlagCount := len(d.NewFlags)
	newCmdCount := len(d.NewCommands)
	newEnvCount := len(d.NewEnvVars)
	deprecatedCount := len(d.DeprecatedCommands) + len(d.DeprecatedFlags)
	breakingCount := len(d.BreakingChanges)

	fmt.Printf("\nRecent delta (%s → %s):\n", d.FromVersion, d.ToVersion)
	fmt.Printf("  + %d new features\n", newFeatureCount)
	fmt.Printf("  + %d new commands\n", newCmdCount)
	fmt.Printf("  + %d new flags\n", newFlagCount)
	fmt.Printf("  + %d new env vars\n", newEnvCount)
	fmt.Printf("  ! %d deprecations\n", deprecatedCount)
	fmt.Printf("  ! %d breaking changes\n", breakingCount)

	if newFeatureCount > 0 {
		fmt.Println("\nNew features:")
		for _, f := range d.NewFeatures {
			fmt.Printf("  • %s\n", f)
		}
	}

	if newCmdCount > 0 {
		fmt.Println("\nNew commands:")
		for _, c := range d.NewCommands {
			fmt.Printf("  • %s\n", c)
		}
	}

	if newFlagCount > 0 {
		fmt.Println("\nNew flags:")
		for _, f := range d.NewFlags {
			fmt.Printf("  • %s\n", f)
		}
	}

	if newEnvCount > 0 {
		fmt.Println("\nNew env vars:")
		for _, e := range d.NewEnvVars {
			fmt.Printf("  • %s\n", e)
		}
	}

	if deprecatedCount > 0 {
		fmt.Println("\nDeprecations:")
		for _, c := range d.DeprecatedCommands {
			fmt.Printf("  • %s (command)\n", c)
		}
		for _, f := range d.DeprecatedFlags {
			fmt.Printf("  • %s (flag)\n", f)
		}
	}

	if breakingCount > 0 {
		fmt.Println("\nBreaking changes:")
		for _, b := range d.BreakingChanges {
			fmt.Printf("  • %s\n", b)
		}
	} else {
		fmt.Println("\nBreaking changes: none")
	}
}

// PrintDiffOnly outputs only the delta block as JSON.
func PrintDiffOnly(p *Payload) error {
	data, err := json.MarshalIndent(p.Delta, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func formatDate(iso string) string {
	t, err := time.Parse(time.RFC3339, iso)
	if err != nil {
		return iso
	}
	return t.Format("2006-01-02")
}
