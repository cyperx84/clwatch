package diff

import (
	"encoding/json"
	"time"

	"github.com/cyperx/clwatch/internal/manifest"
	"github.com/cyperx/clwatch/internal/state"
)

// Status represents the state of a tool relative to local knowledge.
type Status string

const (
	StatusUpdated Status = "updated"
	StatusStale   Status = "stale"
	StatusCurrent Status = "current"
	StatusNew     Status = "new"
)

// Result holds the comparison outcome for a single tool.
type Result struct {
	Tool            string          `json:"tool"`
	Status          Status          `json:"status"`
	PreviousVersion string          `json:"previous_version,omitempty"`
	CurrentVersion  string          `json:"current_version"`
	Delta           json.RawMessage `json:"delta,omitempty"`
}

// Compare checks each manifest tool against local state and returns results.
func Compare(m *manifest.Manifest, s *state.State) []Result {
	var results []Result

	for id, tool := range m.Tools {
		local, exists := s.Tools[id]

		if !exists {
			results = append(results, Result{
				Tool:           id,
				Status:         StatusNew,
				CurrentVersion: tool.CurrentVersion,
				Delta:          tool.Delta,
			})
			continue
		}

		if local.KnownVersion != tool.CurrentVersion {
			results = append(results, Result{
				Tool:            id,
				Status:          StatusUpdated,
				PreviousVersion: local.KnownVersion,
				CurrentVersion:  tool.CurrentVersion,
				Delta:           tool.Delta,
			})
			continue
		}

		// stale_after is an absolute ISO8601 timestamp in the manifest.
		// If it has passed, the payload needs a fresh Tier 2 refresh.
		if manifest.IsStale(tool.StaleAfter) {
			results = append(results, Result{
				Tool:            id,
				Status:          StatusStale,
				PreviousVersion: local.KnownVersion,
				CurrentVersion:  tool.CurrentVersion,
			})
			continue
		}

		results = append(results, Result{
			Tool:            id,
			Status:          StatusCurrent,
			PreviousVersion: local.KnownVersion,
			CurrentVersion:  tool.CurrentVersion,
		})
	}

	return results
}

// UpdateState writes new versions and timestamps into the state.
func UpdateState(s *state.State, results []Result) {
	now := time.Now().UTC().Format(time.RFC3339)
	for _, r := range results {
		s.Tools[r.Tool] = state.ToolState{
			KnownVersion: r.CurrentVersion,
			LastSeenAt:   now,
		}
	}
}

// HasChanges returns true if any tool is not current.
func HasChanges(results []Result) bool {
	for _, r := range results {
		if r.Status != StatusCurrent {
			return true
		}
	}
	return false
}
