package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"text/tabwriter"
)

// DiffAPIResponse is the response from the diff API.
type DiffAPIResponse struct {
	Schema      string `json:"schema"`
	Tool        string `json:"tool"`
	FromVersion string `json:"from_version"`
	ToVersion   string `json:"to_version"`
	Message     string `json:"message,omitempty"`
	Changes     struct {
		Features struct {
			Added   []string `json:"added"`
			Removed []string `json:"removed"`
		} `json:"features"`
		Commands struct {
			Added   []string `json:"added"`
			Removed []string `json:"removed"`
		} `json:"commands"`
		Flags struct {
			Added   []string `json:"added"`
			Removed []string `json:"removed"`
		} `json:"flags"`
		EnvVars struct {
			Added   []string `json:"added"`
			Removed []string `json:"removed"`
		} `json:"env_vars"`
	} `json:"changes"`
	BreakingChanges struct {
		Added []string `json:"added"`
	} `json:"breaking_changes"`
}

var diffAPIURL = "https://changelogs.info/api/refs/diff"

func runDiffTool(args []string) int {
	fs := flag.NewFlagSet("diff-tool", flag.ExitOnError)
	jsonOut := fs.Bool("json", false, "output as JSON")
	fs.Parse(args)

	if fs.NArg() < 3 {
		fmt.Fprintf(os.Stderr, "Usage: clwatch diff-tool <tool> <from-version> <to-version> [--json]\n")
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  clwatch diff-tool claude-code 2.1.74 2.1.75\n")
		return 1
	}

	tool := fs.Arg(0)
	fromVer := fs.Arg(1)
	toVer := fs.Arg(2)

	// Build API URL
	apiURL := os.Getenv("CLWATCH_DIFF_URL")
	if apiURL == "" {
		apiURL = diffAPIURL
	}

	params := url.Values{}
	params.Set("tool", tool)
	params.Set("from", fromVer)
	params.Set("to", toVer)

	fullURL := apiURL + "?" + params.Encode()

	// Fetch diff
	resp, err := http.Get(fullURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error fetching diff: %v\n", err)
		return 1
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(os.Stderr, "API error %d: %s\n", resp.StatusCode, string(body))
		return 1
	}

	var result DiffAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing response: %v\n", err)
		return 1
	}

	// Output
	if *jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(result)
		return 0
	}

	// Human-readable output
	if result.Message != "" {
		fmt.Printf("%s: %s\n", tool, result.Message)
		return 0
	}

	fmt.Printf("%s %s → %s\n\n", tool, fromVer, toVer)

	// Summary counts
	totalAdded := len(result.Changes.Features.Added) +
		len(result.Changes.Commands.Added) +
		len(result.Changes.Flags.Added) +
		len(result.Changes.EnvVars.Added)
	totalRemoved := len(result.Changes.Features.Removed) +
		len(result.Changes.Commands.Removed) +
		len(result.Changes.Flags.Removed) +
		len(result.Changes.EnvVars.Removed)

	if totalAdded == 0 && totalRemoved == 0 {
		fmt.Println("No changes detected")
		return 0
	}

	// Print changes
	w := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)

	printSection(w, "Features", result.Changes.Features.Added, result.Changes.Features.Removed)
	printSection(w, "Commands", result.Changes.Commands.Added, result.Changes.Commands.Removed)
	printSection(w, "Flags", result.Changes.Flags.Added, result.Changes.Flags.Removed)
	printSection(w, "Env Vars", result.Changes.EnvVars.Added, result.Changes.EnvVars.Removed)

	w.Flush()

	// Breaking changes warning
	if len(result.BreakingChanges.Added) > 0 {
		fmt.Printf("\n⚠️  Breaking changes:\n")
		for _, bc := range result.BreakingChanges.Added {
			fmt.Printf("  - %s\n", bc)
		}
	}

	return 0
}

func printSection(w *tabwriter.Writer, name string, added, removed []string) {
	if len(added) == 0 && len(removed) == 0 {
		return
	}

	fmt.Fprintf(w, "%s:\n", name)

	if len(added) > 0 {
		sort.Strings(added)
		for _, a := range added {
			fmt.Fprintf(w, "  + %s\n", a)
		}
	}

	if len(removed) > 0 {
		sort.Strings(removed)
		for _, r := range removed {
			fmt.Fprintf(w, "  - %s\n", r)
		}
	}
}
