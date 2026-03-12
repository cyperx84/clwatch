# clwatch (Go CLI)

Track coding tool updates from [changelogs.info](https://changelogs.info).

## Install

### Homebrew

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

### Go install

```bash
go install github.com/cyperx/clwatch/cmd/clwatch@latest
```

### Build from source

```bash
cd go/
make build        # builds ./clwatch
make install      # installs to $GOPATH/bin/clwatch
```

Requires Go 1.22+.

## Usage

```bash
# Check for updates
clwatch diff

# Check with JSON output
clwatch diff --json

# Show all tools (including current)
clwatch diff --verbose

# Check without updating local state
clwatch diff --no-update

# List all tracked tools
clwatch list

# List as JSON
clwatch list --json

# Print version
clwatch version
```

## Exit codes

- `clwatch diff` exits **0** if all tools are current, **1** if any changes detected (useful for scripting).

## Environment

| Variable | Description |
|---|---|
| `CLWATCH_MANIFEST_URL` | Override the manifest URL (default: `https://changelogs.info/api/refs/manifest.json`) |

## Local state

State is stored at `~/.clwatch/state.json` and tracks known versions and last-checked timestamps.
