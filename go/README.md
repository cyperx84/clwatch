# clwatch — Go CLI

Developer reference for building and developing the clwatch CLI. For user-facing docs, see the [root README](../README.md).

## Build

```bash
# Build for current platform
make build

# Install to $GOPATH/bin
make install

# Cross-platform release builds
bash scripts/build-release.sh [version]
```

Requires Go 1.22+.

## Commands

| Command | Description | Exit codes |
|---|---|---|
| `diff` | Check manifest for updates | 0 = current, 1 = changes |
| `list` | Show all tracked tools | 0 |
| `refresh <tool>` | Pull payload + show delta | 0 |
| `init` | Scaffold workspace config | 0 |
| `ack <tool> <version>` | Mark version as known | 0 |
| `watch` | Daemon mode, poll on interval | 0 (Ctrl+C) |
| `status` | Pipeline health from site | 0 |
| `version` | Print version | 0 |

## Architecture

```
go/
├── cmd/clwatch/main.go         # CLI entry point, command routing
├── internal/
│   ├── diff/diff.go            # Manifest vs local state comparison
│   ├── manifest/manifest.go    # Manifest fetching + parsing
│   ├── output/output.go        # Table/JSON output formatting
│   ├── refresh/refresh.go      # Payload fetching + delta display
│   ├── state/state.go          # Local state (read/write ~/.clwatch/state.json)
│   ├── watcher/watcher.go      # Watch daemon (polling + webhooks)
│   └── workspace/workspace.go  # Workspace config (.clwatch.json)
├── scripts/build-release.sh    # Cross-platform build script
├── go.mod
└── go.sum
```

## Testing

```bash
# Against live site
clwatch diff
clwatch list
clwatch refresh claude-code

# Against local server
cd ~/github/changelogs-info/public && python3 -m http.server 8889 &
CLWATCH_MANIFEST_URL="http://localhost:8889/api/refs/manifest.json" clwatch diff

# Watch daemon (Ctrl+C to stop)
clwatch watch --interval 30s

# Status from live site
clwatch status
```

## Environment variables

| Variable | Default | Description |
|---|---|---|
| `CLWATCH_MANIFEST_URL` | `https://changelogs.info/api/refs/manifest.json` | Manifest source |
| `CLWATCH_STATUS_URL` | Derived from manifest URL | Status endpoint |

## State file

`~/.clwatch/state.json` — auto-managed by `diff` and `ack` commands:

```json
{
  "schema": "clwatch.state.v1",
  "last_checked": "2026-03-12T00:00:00Z",
  "tools": {
    "claude-code": {
      "known_version": "2.1.74",
      "last_seen_at": "2026-03-12T00:00:00Z"
    }
  }
}
```

## Adding a command

1. Add handler in `cmd/clwatch/main.go` (switch case + `runXxx` function)
2. Implement logic in a new `internal/` package
3. Add to `printUsage()` help text
4. Add test cases
5. Update root `README.md` with user-facing docs

## Release process

1. `git tag v1.x.x && git push --tags`
2. GitHub Actions builds 5 binaries via `scripts/build-release.sh`
3. Release created automatically via `softprops/action-gh-release`
4. Update Homebrew formula in `cyperx84/homebrew-tap` with new checksums
