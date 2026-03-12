package workspace

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config is the .clwatch.json file format.
type Config struct {
	Schema           string   `json:"schema"`
	Tools            []string `json:"tools"`
	ManifestURL      string   `json:"manifestUrl"`
	ReferenceDir     string   `json:"referenceDir"`
	Tier2Threshold   string   `json:"tier2Threshold"`
	NotifyOnBreaking bool     `json:"notifyOnBreaking"`
	StateFile        string   `json:"stateFile"`
}

// DefaultConfig returns the default clwatch config.
func DefaultConfig(tools []string) Config {
	if len(tools) == 0 {
		tools = []string{"claude-code", "codex-cli", "gemini-cli", "opencode", "openclaw"}
	}
	return Config{
		Schema:           "clwatch.config.v1",
		Tools:            tools,
		ManifestURL:      "https://changelogs.info/api/refs/manifest.json",
		ReferenceDir:     "references/",
		Tier2Threshold:   "medium",
		NotifyOnBreaking: true,
		StateFile:        "~/.clwatch/state.json",
	}
}

// Init scaffolds a clwatch workspace in the given directory.
// Returns a list of action messages for display.
func Init(dir string, tools []string, force bool) ([]string, error) {
	var actions []string

	configPath := filepath.Join(dir, ".clwatch.json")

	// Check for existing config
	if _, err := os.Stat(configPath); err == nil && !force {
		return nil, fmt.Errorf(".clwatch.json already exists (use --force to overwrite)")
	}

	// Write config
	cfg := DefaultConfig(tools)
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshaling config: %w", err)
	}
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return nil, fmt.Errorf("writing .clwatch.json: %w", err)
	}
	actions = append(actions, "Created .clwatch.json")

	// Create references directory
	refDir := filepath.Join(dir, cfg.ReferenceDir)
	if err := os.MkdirAll(refDir, 0755); err != nil {
		return nil, fmt.Errorf("creating references dir: %w", err)
	}
	actions = append(actions, "Created "+cfg.ReferenceDir)

	// Create placeholder reference files
	for _, tool := range cfg.Tools {
		slug := tool
		filename := slug + "-features.md"
		refPath := filepath.Join(refDir, filename)

		if _, err := os.Stat(refPath); err == nil {
			// File already exists, skip
			continue
		}

		displayName := toolDisplayName(tool)
		content := fmt.Sprintf("# %s Features Reference\n\n<!-- clwatch managed: last updated NEVER, version unknown -->\n<!-- Run `clwatch refresh %s` to populate this file -->\n", displayName, tool)

		if err := os.WriteFile(refPath, []byte(content), 0644); err != nil {
			return nil, fmt.Errorf("writing %s: %w", filename, err)
		}
		actions = append(actions, "Created "+cfg.ReferenceDir+filename+" (placeholder)")
	}

	return actions, nil
}

// toolDisplayName converts a tool slug to a display name.
func toolDisplayName(slug string) string {
	parts := strings.Split(slug, "-")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, " ")
}
