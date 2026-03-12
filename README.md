# clwatch

```
     _ _                 _       _
 ___| (_)_      ____ _  | |_ ___| |__
/ __| | \ \ /\ / / _` | | __/ __| '_ \
\__ \ | |\ V  V / (_| | | || (__| | | |
|___/_|_| \_/\_/ \__,_|  \__\___|_| |_|

Track AI coding tool changes from changelogs.info
```

`clwatch` is a CLI that monitors changelogs for AI coding tools (Claude Code, Codex CLI, Gemini CLI, OpenCode, OpenClaw). It detects version changes, shows you what's new, and helps you keep your agent configs in sync.

## Install

### Homebrew (macOS/Linux)

```bash
brew install cyperx84/tap/clwatch
```

### npm

```bash
npm install -g clwatch
```

### curl

```bash
curl -fsSL https://raw.githubusercontent.com/cyperx84/clwatch/main/install.sh | bash
```

### Go

```bash
go install github.com/cyperx84/clwatch/cmd/clwatch@latest
```

### Build from source

```bash
git clone https://github.com/cyperx84/clwatch.git
cd clwatch/go && go build -o clwatch ./cmd/clwatch
```

## Commands

### `clwatch diff`

Check the manifest for new versions. Returns exit code 1 if any tool has an update — perfect for scripting.

```bash
clwatch diff                 # check for updates (human output)
clwatch diff --json          # raw JSON output
clwatch diff --verbose       # include debug info
clwatch diff --no-update     # don't update local state
```

### `clwatch list`

Show all tracked tools with current versions and status.

```bash
clwatch list                 # table output
clwatch list --json          # JSON output
```

Output:
```
TOOL         VERSION   KNOWN    STATUS  LAST CHECKED
claude-code  2.1.74    2.1.74   current 2h ago
codex-cli    0.114.0   0.114.0  current 2h ago
gemini-cli   0.33.0    0.33.0   current 2h ago
opencode     1.2.24    1.2.24   current 2h ago
openclaw     2026.3.9  2026.3.9 current 2h ago
```

### `clwatch refresh <tool>`

Pull the latest payload and show a summary of what changed.

```bash
clwatch refresh claude-code          # human summary of latest payload + delta
clwatch refresh claude-code --json   # raw payload JSON
clwatch refresh claude-code --diff-only  # only the diff block (for piping)
clwatch refresh --all                # refresh all 5 tools
```

Output:
```
claude-code 2.1.74 (verified)

Recent delta (2.1.71 → 2.1.74):
  + 3 new features
  + 1 new commands
  ! 1 deprecations
  ! 0 breaking changes

New features:
  • autoMemoryDirectory
  • modelOverrides
  • context-suggestions
```

### `clwatch watch`

Run as a daemon — polls the manifest on a schedule and reports changes.

```bash
clwatch watch                            # poll every 6h (default)
clwatch watch --interval 1h              # custom interval (min 15m)
clwatch watch --json                     # JSON output
clwatch watch --webhook https://...      # POST to webhook on changes
```

Webhook payload:
```json
{
  "event": "tools_updated",
  "detected_at": "2026-03-12T00:00:00Z",
  "updates": [
    {
      "tool": "claude-code",
      "status": "updated",
      "previous_version": "2.1.74",
      "current_version": "2.1.75",
      "breaking": false
    }
  ]
}
```

### `clwatch status`

Show pipeline health from changelogs.info.

```bash
clwatch status               # table output
clwatch status --json        # raw JSON
```

Output:
```
changelogs.info status
Pipeline last ran: 1h ago  (ok)

TOOL         VERSION   VERIFIED  STALE  LAST CHECKED
claude-code  2.1.74    ✓         no     1h ago
codex-cli    0.114.0   ✓         no     1h ago
gemini-cli   0.33.0    ✓         no     1h ago
opencode     1.2.24    ✓         no     1h ago
openclaw     2026.3.9  ✓         no     1h ago
```

### `clwatch init`

Scaffold a workspace with `.clwatch.json` config + reference files.

```bash
clwatch init                 # create in current directory
clwatch init --dir refs/     # custom reference directory
clwatch init --tools claude-code,codex-cli  # specific tools
clwatch init --force         # overwrite existing
```

Creates:
```
.clwatch.json          # workspace config
references/            # directory for reference files
  claude-code-features.md
  codex-cli-features.md
  ...
```

### `clwatch ack <tool> <version>`

Mark a version as known (used after merging updates).

```bash
clwatch ack claude-code 2.1.74
```

## Environment variables

| Variable | Default | Description |
|---|---|---|
| `CLWATCH_MANIFEST_URL` | `https://changelogs.info/api/refs/manifest.json` | Override manifest source |
| `CLWATCH_STATUS_URL` | Derived from manifest URL | Override status endpoint |

## Configuration

### Workspace config (`.clwatch.json`)

Created by `clwatch init`:

```json
{
  "schema": "clwatch.config.v1",
  "tools": ["claude-code", "codex-cli", "gemini-cli", "opencode", "openclaw"],
  "manifestUrl": "https://changelogs.info/api/refs/manifest.json",
  "referenceDir": "references/",
  "tier2Threshold": "medium",
  "notifyOnBreaking": true,
  "stateFile": "~/.clwatch/state.json"
}
```

### Local state (`~/.clwatch/state.json`)

Automatically managed by `clwatch diff` and `clwatch ack`. Tracks which versions you've seen.

## Use with agents

Pair with the [clwatch skill](https://github.com/cyperx84/clwatch-skill) for automatic agent integration:

```bash
# Session start check (silent if up-to-date)
bash scripts/check-updates.sh

# Trigger tier 2 merge flow
bash scripts/tier2-merge.sh claude-code
```

## API

clwatch consumes the static JSON API at [changelogs.info/api/refs/](https://changelogs.info/api/refs/):

| Endpoint | Description |
|---|---|
| `/api/refs/manifest.json` | Tracking manifest (all tools + versions) |
| `/api/refs/status.json` | Pipeline health |
| `/api/refs/{tool}.json` | Full structured payload per tool |

## License

MIT
