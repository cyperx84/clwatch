package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const Schema = "clwatch.state.v1"

type ToolState struct {
	KnownVersion string `json:"known_version"`
	LastSeenAt   string `json:"last_seen_at"`
}

type State struct {
	Schema      string               `json:"schema"`
	LastChecked string               `json:"last_checked"`
	Tools       map[string]ToolState `json:"tools"`
}

func statePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home dir: %w", err)
	}
	return filepath.Join(home, ".clwatch", "state.json"), nil
}

func Load() (*State, error) {
	p, err := statePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return &State{
				Schema: Schema,
				Tools:  make(map[string]ToolState),
			}, nil
		}
		return nil, fmt.Errorf("reading state file: %w", err)
	}

	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parsing state file: %w", err)
	}

	if s.Tools == nil {
		s.Tools = make(map[string]ToolState)
	}

	return &s, nil
}

func Save(s *State) error {
	p, err := statePath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		return fmt.Errorf("creating state dir: %w", err)
	}

	s.Schema = Schema
	s.LastChecked = time.Now().UTC().Format(time.RFC3339)

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling state: %w", err)
	}

	return os.WriteFile(p, data, 0644)
}

func (ts ToolState) LastSeenTime() (time.Time, error) {
	if ts.LastSeenAt == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, ts.LastSeenAt)
}
