# cliwatch

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
npm install -g cliwatch
```

## Usage

### Show recent changes

```bash
# All changes in last 7 days
cliwatch diff

# Changes for specific tool
cliwatch diff claude-code

# Changes in last 30 days
cliwatch diff --since 30d
```

### Check your config

```bash
# Scan config file for outdated options
cliwatch check .claude/CLAUDE.md
cliwatch check .cursorrules
cliwatch check .aider.conf.yml
```

### List models

```bash
# All tracked models
cliwatch models

# Only new models (last 90 days)
cliwatch models --new

# Filter by provider
cliwatch models --provider anthropic
```

### Check model compatibility

```bash
cliwatch compat claude-sonnet-4-6
```

### View deprecations

```bash
# All deprecations
cliwatch deprecations

# For specific harness
cliwatch deprecations --harness cursor
```

### Interactive TUI

```bash
cliwatch tui
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
