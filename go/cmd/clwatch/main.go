package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cyperx/clwatch/internal/diff"
	"github.com/cyperx/clwatch/internal/manifest"
	"github.com/cyperx/clwatch/internal/output"
	"github.com/cyperx/clwatch/internal/refresh"
	"github.com/cyperx/clwatch/internal/state"
	"github.com/cyperx/clwatch/internal/watcher"
	"github.com/cyperx/clwatch/internal/workspace"
)

var version = "dev"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "diff":
		os.Exit(runDiff(os.Args[2:]))
	case "list":
		os.Exit(runList(os.Args[2:]))
	case "refresh":
		os.Exit(runRefresh(os.Args[2:]))
	case "init":
		os.Exit(runInit(os.Args[2:]))
	case "ack":
		os.Exit(runAck(os.Args[2:]))
	case "watch":
		os.Exit(runWatch(os.Args[2:]))
	case "status":
		os.Exit(runStatus(os.Args[2:]))
	case "version":
		fmt.Printf("clwatch %s\n", version)
		os.Exit(0)
	case "service":
		os.Exit(runService(os.Args[2:]))
	case "completion":
		os.Exit(runCompletion(os.Args[2:]))
	case "diff-tool":
		os.Exit(runDiffTool(os.Args[2:]))
	case "help", "--help", "-h":
		printUsage()
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func runWatch(args []string) int {
	cfg := watcher.Config{
		ManifestURL: manifest.ManifestURL(),
		Interval:    6 * time.Hour,
	}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--interval":
			if i+1 < len(args) {
				d, err := watcher.ParseInterval(args[i+1])
				if err != nil {
					fmt.Fprintf(os.Stderr, "error: %v\n", err)
					return 1
				}
				cfg.Interval = d
				i++
			}
		case "--json":
			cfg.JSONOutput = true
		case "--webhook":
			if i+1 < len(args) {
				cfg.WebhookURL = args[i+1]
				i++
			}
		}
	}
	diffFn := func(ctx context.Context, url string, _ bool) ([]watcher.Update, error) {
		m, err := manifest.FetchFrom(url)
		if err != nil {
			return nil, err
		}
		s, _ := state.Load()
		results := diff.Compare(m, s)
		diff.UpdateState(s, results)
		state.Save(s)
		var updates []watcher.Update
		for _, r := range results {
			if r.Status != diff.StatusCurrent {
				updates = append(updates, watcher.Update{
					Tool:            r.Tool,
					Status:          string(r.Status),
					PreviousVersion: r.PreviousVersion,
					CurrentVersion:  r.CurrentVersion,
				})
			}
		}
		return updates, nil
	}
	if err := watcher.Run(cfg, diffFn); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	return 0
}

type statusResponse struct {
	Schema      string `json:"schema"`
	GeneratedAt string `json:"generated_at"`
	Pipeline    struct {
		LastRunAt    string `json:"last_run_at"`
		Status       string `json:"status"`
		ToolsChecked int    `json:"tools_checked"`
		ToolsUpdated int    `json:"tools_updated"`
		ToolsErrored int    `json:"tools_errored"`
	} `json:"pipeline"`
	Tools map[string]struct {
		Version            string `json:"version"`
		VerificationStatus string `json:"verification_status"`
		LastCheckedAt      string `json:"last_checked_at"`
		Stale              bool   `json:"stale"`
	} `json:"tools"`
}

func runStatus(args []string) int {
	jsonOutput := false
	statusURL := os.Getenv("CLWATCH_STATUS_URL")
	if statusURL == "" {
		statusURL = strings.Replace(manifest.ManifestURL(), "manifest.json", "status.json", 1)
	}
	for _, a := range args {
		if a == "--json" {
			jsonOutput = true
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, statusURL, nil)
	req.Header.Set("User-Agent", "clwatch/1.0")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error fetching status: %v\n", err)
		return 1
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "error: HTTP %d from %s\n", resp.StatusCode, statusURL)
		return 1
	}
	if jsonOutput {
		fmt.Println(string(body))
		return 0
	}
	var s statusResponse
	if err := json.Unmarshal(body, &s); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing status: %v\n", err)
		return 1
	}
	ago := func(ts string) string {
		if ts == "" {
			return "never"
		}
		t, err := time.Parse(time.RFC3339, ts)
		if err != nil {
			return ts
		}
		d := time.Since(t)
		if d < time.Minute {
			return "just now"
		} else if d < time.Hour {
			return fmt.Sprintf("%dm ago", int(d.Minutes()))
		} else if d < 24*time.Hour {
			return fmt.Sprintf("%dh ago", int(d.Hours()))
		}
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	}
	fmt.Println("changelogs.info status")
	fmt.Printf("Pipeline last ran: %s  (%s)\n\n", ago(s.Pipeline.LastRunAt), s.Pipeline.Status)
	fmt.Printf("%-14s %-10s %-10s %-6s %s\n", "TOOL", "VERSION", "VERIFIED", "STALE", "LAST CHECKED")
	for id, t := range s.Tools {
		verified := "✓"
		if t.VerificationStatus != "verified" {
			verified = "?"
		}
		stale := "no"
		if t.Stale {
			stale = "YES"
		}
		fmt.Printf("%-14s %-10s %-10s %-6s %s\n", id, t.Version, verified, stale, ago(t.LastCheckedAt))
	}
	return 0
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `clwatch — track coding tool updates

Usage:
  clwatch diff        [--json] [--verbose] [--no-update]
  clwatch list        [--json]
  clwatch refresh     <tool> [--json] [--diff-only] [--all]
  clwatch init        [--dir DIR] [--tools TOOLS] [--force]
  clwatch ack         <tool> <version>
  clwatch watch       [--interval 6h] [--json] [--webhook URL]
  clwatch status      [--json]
  clwatch diff-tool   <tool> <from> <to> [--json]
  clwatch service     <install|uninstall|start|stop|status|logs>
  clwatch completion  <bash|zsh|fish>
  clwatch version

Environment:
  CLWATCH_MANIFEST_URL   Override manifest URL (default: %s)
  CLWATCH_DIFF_URL      Override diff API URL (default: https://changelogs.info/api/refs/diff)
`, manifest.DefaultManifestURL)
}

func runDiff(args []string) int {
	fs := flag.NewFlagSet("diff", flag.ExitOnError)
	jsonOut := fs.Bool("json", false, "output as JSON")
	verbose := fs.Bool("verbose", false, "show all tools including current")
	noUpdate := fs.Bool("no-update", false, "do not update local state")
	fs.Parse(args)

	m, err := manifest.Fetch()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	s, err := state.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	results := diff.Compare(m, s)

	if *jsonOut {
		if err := output.PrintDiffJSON(results); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
	} else {
		output.PrintDiffText(results, *verbose)
	}

	if !*noUpdate {
		diff.UpdateState(s, results)
		if err := state.Save(s); err != nil {
			fmt.Fprintf(os.Stderr, "error saving state: %v\n", err)
			return 1
		}
	}

	if diff.HasChanges(results) {
		return 1
	}
	return 0
}

func runList(args []string) int {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	jsonOut := fs.Bool("json", false, "output as JSON")
	fs.Parse(args)

	m, err := manifest.Fetch()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	s, err := state.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	entries := output.BuildListEntries(m, s)

	if *jsonOut {
		if err := output.PrintListJSON(entries); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
	} else {
		output.PrintListTable(entries)
	}

	return 0
}

func runRefresh(args []string) int {
	// Separate positional args from flags so flags can appear after tool name
	var flagArgs []string
	var posArgs []string
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			flagArgs = append(flagArgs, a)
		} else {
			posArgs = append(posArgs, a)
		}
	}

	fs := flag.NewFlagSet("refresh", flag.ExitOnError)
	jsonOut := fs.Bool("json", false, "output raw payload JSON")
	diffOnly := fs.Bool("diff-only", false, "output only the delta block as JSON")
	all := fs.Bool("all", false, "refresh all tools in manifest")
	fs.Parse(flagArgs)

	m, err := manifest.Fetch()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	s, err := state.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	var tools []string
	if *all {
		for id := range m.Tools {
			tools = append(tools, id)
		}
	} else {
		if len(posArgs) < 1 {
			fmt.Fprintf(os.Stderr, "usage: clwatch refresh <tool> [--json] [--diff-only]\n")
			fmt.Fprintf(os.Stderr, "       clwatch refresh --all\n")
			return 1
		}
		tools = []string{posArgs[0]}
	}

	exitCode := 0
	for i, toolID := range tools {
		if _, ok := m.Tools[toolID]; !ok {
			fmt.Fprintf(os.Stderr, "error: tool %q not found in manifest\n", toolID)
			exitCode = 1
			continue
		}

		p, rawBody, err := refresh.FetchPayload(m, toolID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			exitCode = 1
			continue
		}

		if *jsonOut {
			// Pretty-print the raw JSON
			var buf json.RawMessage
			if err := json.Unmarshal(rawBody, &buf); err != nil {
				fmt.Println(string(rawBody))
			} else {
				pretty, _ := json.MarshalIndent(buf, "", "  ")
				fmt.Println(string(pretty))
			}
		} else if *diffOnly {
			if err := refresh.PrintDiffOnly(p); err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				exitCode = 1
				continue
			}
		} else {
			if *all && i > 0 {
				fmt.Println()
			}
			refresh.PrintSummary(p)
		}

		// Update local state
		refresh.UpdateStateForTool(s, toolID, p.Version)
	}

	if err := state.Save(s); err != nil {
		fmt.Fprintf(os.Stderr, "error saving state: %v\n", err)
		return 1
	}

	return exitCode
}

func runInit(args []string) int {
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	dir := fs.String("dir", ".", "directory to initialize")
	toolsFlag := fs.String("tools", "", "comma-separated list of tools")
	force := fs.Bool("force", false, "overwrite existing .clwatch.json")
	fs.Parse(args)

	var tools []string
	if *toolsFlag != "" {
		for _, t := range strings.Split(*toolsFlag, ",") {
			t = strings.TrimSpace(t)
			if t != "" {
				tools = append(tools, t)
			}
		}
	}

	actions, err := workspace.Init(*dir, tools, *force)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	for _, a := range actions {
		fmt.Printf("✓ %s\n", a)
	}

	fmt.Println("\nNext steps:")
	fmt.Println("  clwatch diff              check for updates")
	fmt.Println("  clwatch refresh --all     pull latest payloads")

	return 0
}

func runAck(args []string) int {
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: clwatch ack <tool> <version>\n")
		return 1
	}

	toolID := args[0]
	ver := args[1]

	s, err := state.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	refresh.UpdateStateForTool(s, toolID, ver)

	if err := state.Save(s); err != nil {
		fmt.Fprintf(os.Stderr, "error saving state: %v\n", err)
		return 1
	}

	fmt.Printf("✓ Acknowledged %s %s\n", toolID, ver)
	return 0
}

func runCompletion(args []string) int {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: clwatch completion <bash|zsh|fish>\n")
		return 1
	}
	switch args[0] {
	case "bash":
		fmt.Print(bashCompletion)
	case "zsh":
		fmt.Print(zshCompletion)
	case "fish":
		fmt.Print(fishCompletion)
	default:
		fmt.Fprintf(os.Stderr, "unknown shell: %s (supported: bash, zsh, fish)\n", args[0])
		return 1
	}
	return 0
}
