# clwatch

```
     _ _                 _       _
 ___| (_)_      ____ _  | |_ ___| |__
/ __| | \ \ /\ / / _` | | __/ __| '_ \
\__ \ | |\ V  V / (_| | | || (__| | | |
|___/_|_| \_/\_/ \__,_|  \__\___|_| |_|

Track AI coding tool changes from changelogs.info
```

## Installation

```bash
npm install -g clwatch
```

## Usage

### Show recent changes

```bash
# All changes in last 7 days
clwatch diff

# Changes for specific tool
clwatch diff claude-code

# Changes in last 30 days
clwatch diff --since 30d
```

### Check your config

```bash
# Scan config file for outdated options
clwatch check .claude/CLAUDE.md
clwatch check .cursorrules
clwatch check .aider.conf.yml
```

### List models

```bash
# All tracked models
clwatch models

# Only new models (last 90 days)
clwatch models --new

# Filter by provider
clwatch models --provider anthropic
```

### Check model compatibility

```bash
clwatch compat claude-sonnet-4-6
```

### View deprecations

```bash
# All deprecations
clwatch deprecations

# For specific harness
clwatch deprecations --harness cursor
```

### Interactive TUI

```bash
clwatch tui
```

## Commands

| Command | Description |
|---------|-------------|
| `diff [harness] [--since <duration>]` | Show what changed recently |
| `check <config-path>` | Analyze config for issues |
| `models [--new] [--provider <name>]` | List AI models |
| `compat <model-id>` | Show harness compatibility |
| `deprecations [--harness <name>]` | Show deprecation timeline |
| `tui` | Interactive dashboard |

## Development

```bash
# Install dependencies
npm install

# Run in dev mode
npm run dev -- diff

# Build
npm run build

# Sync data from changelogs-info repo
npm run sync-data
```

## Data Source

Currently reads from bundled JSON files. Future versions will query changelogs.info API.

## License

MIT
