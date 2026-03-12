package output

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/cyperx/clwatch/internal/diff"
	"github.com/cyperx/clwatch/internal/manifest"
	"github.com/cyperx/clwatch/internal/state"
)

func RelativeTime(t time.Time) string {
	if t.IsZero() {
		return "never"
	}
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		m := int(d.Minutes())
		return fmt.Sprintf("%dm ago", m)
	case d < 24*time.Hour:
		h := int(d.Hours())
		return fmt.Sprintf("%dh ago", h)
	default:
		days := int(d.Hours() / 24)
		return fmt.Sprintf("%dd ago", days)
	}
}

func PrintDiffText(results []diff.Result, verbose bool) {
	for _, r := range results {
		switch r.Status {
		case diff.StatusNew:
			fmt.Printf("✦ %s: new → %s\n", r.Tool, r.CurrentVersion)
		case diff.StatusUpdated:
			fmt.Printf("✦ %s: %s → %s\n", r.Tool, r.PreviousVersion, r.CurrentVersion)
		case diff.StatusStale:
			fmt.Printf("⚠ %s: stale (version %s)\n", r.Tool, r.CurrentVersion)
		case diff.StatusCurrent:
			if verbose {
				fmt.Printf("  %s: %s (current)\n", r.Tool, r.CurrentVersion)
			}
		}
	}
}

func PrintDiffJSON(results []diff.Result) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

type ListEntry struct {
	Tool        string `json:"tool"`
	Version     string `json:"version"`
	Known       string `json:"known"`
	Status      string `json:"status"`
	LastChecked string `json:"last_checked"`
}

func BuildListEntries(m *manifest.Manifest, s *state.State) []ListEntry {
	var entries []ListEntry
	for id, tool := range m.Tools {
		entry := ListEntry{
			Tool:    id,
			Version: tool.CurrentVersion,
		}

		if local, ok := s.Tools[id]; ok {
			entry.Known = local.KnownVersion

			if local.KnownVersion != tool.CurrentVersion {
				entry.Status = "updated"
			} else {
				entry.Status = "current"
			}

			t, err := local.LastSeenTime()
			if err == nil && !t.IsZero() {
				// Check staleness
				// use manifest IsStale for absolute timestamp check
				if manifest.IsStale(tool.StaleAfter) {
					entry.Status = "stale"
				}
				entry.LastChecked = RelativeTime(t)
			} else {
				entry.LastChecked = "never"
			}
		} else {
			entry.Known = "-"
			entry.Status = "new"
			entry.LastChecked = "never"
		}

		entries = append(entries, entry)
	}
	return entries
}

func PrintListTable(entries []ListEntry) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TOOL\tVERSION\tKNOWN\tSTATUS\tLAST CHECKED")
	for _, e := range entries {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", e.Tool, e.Version, e.Known, e.Status, e.LastChecked)
	}
	w.Flush()
}

func PrintListJSON(entries []ListEntry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
