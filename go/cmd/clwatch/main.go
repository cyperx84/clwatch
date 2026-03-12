package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cyperx/clwatch/internal/diff"
	"github.com/cyperx/clwatch/internal/manifest"
	"github.com/cyperx/clwatch/internal/output"
	"github.com/cyperx/clwatch/internal/state"
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
  clwatch diff    [--json] [--verbose] [--no-update]
  clwatch list    [--json]
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
