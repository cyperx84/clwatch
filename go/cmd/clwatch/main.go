package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/cyperx/clwatch/internal/diff"
	"github.com/cyperx/clwatch/internal/manifest"
	"github.com/cyperx/clwatch/internal/output"
	"github.com/cyperx/clwatch/internal/refresh"
	"github.com/cyperx/clwatch/internal/state"
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
	case "version":
		fmt.Printf("clwatch %s\n", version)
		os.Exit(0)
	case "help", "--help", "-h":
		printUsage()
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `clwatch — track coding tool updates

Usage:
  clwatch diff      [--json] [--verbose] [--no-update]
  clwatch list      [--json]
  clwatch refresh   <tool> [--json] [--diff-only] [--all]
  clwatch init      [--dir DIR] [--tools TOOLS] [--force]
  clwatch ack       <tool> <version>
  clwatch version

Environment:
  CLWATCH_MANIFEST_URL   Override manifest URL (default: %s)
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
