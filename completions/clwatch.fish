# Fish shell completion for clwatch

set -l tools claude-code codex-cli gemini-cli opencode openclaw
set -l commands diff list refresh init ack watch status service completion version help

# Disable file completions for all clwatch subcommands
complete -c clwatch -f

# Top-level commands
complete -c clwatch -n "not __fish_seen_subcommand_from $commands" -a diff      -d "Show version changes since last check"
complete -c clwatch -n "not __fish_seen_subcommand_from $commands" -a list      -d "List all tracked tools and their versions"
complete -c clwatch -n "not __fish_seen_subcommand_from $commands" -a refresh   -d "Force-refresh a specific tool"
complete -c clwatch -n "not __fish_seen_subcommand_from $commands" -a init      -d "Initialize a workspace"
complete -c clwatch -n "not __fish_seen_subcommand_from $commands" -a ack       -d "Acknowledge a version"
complete -c clwatch -n "not __fish_seen_subcommand_from $commands" -a watch     -d "Watch for updates on interval"
complete -c clwatch -n "not __fish_seen_subcommand_from $commands" -a status    -d "Show pipeline status"
complete -c clwatch -n "not __fish_seen_subcommand_from $commands" -a service   -d "Manage background service"
complete -c clwatch -n "not __fish_seen_subcommand_from $commands" -a completion -d "Generate shell completions"
complete -c clwatch -n "not __fish_seen_subcommand_from $commands" -a version   -d "Show version"
complete -c clwatch -n "not __fish_seen_subcommand_from $commands" -a help      -d "Show help"

# diff flags
complete -c clwatch -n "__fish_seen_subcommand_from diff" -l json    -d "Output as JSON"
complete -c clwatch -n "__fish_seen_subcommand_from diff" -l verbose -d "Show all tools including up-to-date"
complete -c clwatch -n "__fish_seen_subcommand_from diff" -l no-update -d "Do not update local state"

# list flags
complete -c clwatch -n "__fish_seen_subcommand_from list" -l json -d "Output as JSON"

# refresh — tool names + flags
complete -c clwatch -n "__fish_seen_subcommand_from refresh" -a "$tools"
complete -c clwatch -n "__fish_seen_subcommand_from refresh" -l json      -d "Output as JSON"
complete -c clwatch -n "__fish_seen_subcommand_from refresh" -l diff-only -d "Show diff only"
complete -c clwatch -n "__fish_seen_subcommand_from refresh" -l all       -d "Refresh all tools"

# ack — tool names
complete -c clwatch -n "__fish_seen_subcommand_from ack" -a "$tools"

# init flags
complete -c clwatch -n "__fish_seen_subcommand_from init" -l dir   -d "Directory to initialize" -r -F
complete -c clwatch -n "__fish_seen_subcommand_from init" -l tools -d "Comma-separated tool list" -r
complete -c clwatch -n "__fish_seen_subcommand_from init" -l force -d "Overwrite existing config"

# watch flags
complete -c clwatch -n "__fish_seen_subcommand_from watch" -l interval -d "Check interval (e.g. 6h, 30m)" -r
complete -c clwatch -n "__fish_seen_subcommand_from watch" -l json     -d "Output as JSON"
complete -c clwatch -n "__fish_seen_subcommand_from watch" -l webhook  -d "Webhook URL" -r

# status flags
complete -c clwatch -n "__fish_seen_subcommand_from status" -l json -d "Output as JSON"

# service subcommands
complete -c clwatch -n "__fish_seen_subcommand_from service" -a install   -d "Install background service"
complete -c clwatch -n "__fish_seen_subcommand_from service" -a uninstall -d "Remove background service"
complete -c clwatch -n "__fish_seen_subcommand_from service" -a start     -d "Start service"
complete -c clwatch -n "__fish_seen_subcommand_from service" -a stop      -d "Stop service"
complete -c clwatch -n "__fish_seen_subcommand_from service" -a status    -d "Show service status"
complete -c clwatch -n "__fish_seen_subcommand_from service" -a logs      -d "Show service logs"

# completion shells
complete -c clwatch -n "__fish_seen_subcommand_from completion" -a "bash zsh fish"
