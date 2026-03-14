package main

// Shell completion scripts embedded at build time.
// These mirror the files in completions/ at the repo root.

const bashCompletion = `_clwatch_completions() {
  local cur
  COMPREPLY=()
  cur="${COMP_WORDS[COMP_CWORD]}"

  local commands="diff list refresh init ack watch status service completion version help"
  local tools="claude-code codex-cli gemini-cli opencode openclaw"

  if [[ ${COMP_CWORD} -eq 1 ]]; then
    COMPREPLY=( $(compgen -W "${commands}" -- ${cur}) )
    return 0
  fi

  case "${COMP_WORDS[1]}" in
    diff)       COMPREPLY=( $(compgen -W "--json --verbose --no-update" -- ${cur}) ) ;;
    list)       COMPREPLY=( $(compgen -W "--json" -- ${cur}) ) ;;
    refresh)
      if [[ ${COMP_CWORD} -eq 2 ]]; then
        COMPREPLY=( $(compgen -W "${tools} --all" -- ${cur}) )
      else
        COMPREPLY=( $(compgen -W "--json --diff-only --all" -- ${cur}) )
      fi ;;
    ack)
      if [[ ${COMP_CWORD} -eq 2 ]]; then
        COMPREPLY=( $(compgen -W "${tools}" -- ${cur}) )
      fi ;;
    init)       COMPREPLY=( $(compgen -W "--dir --tools --force" -- ${cur}) ) ;;
    watch)      COMPREPLY=( $(compgen -W "--interval --json --webhook" -- ${cur}) ) ;;
    status)     COMPREPLY=( $(compgen -W "--json" -- ${cur}) ) ;;
    service)    COMPREPLY=( $(compgen -W "install uninstall start stop status logs" -- ${cur}) ) ;;
    completion) COMPREPLY=( $(compgen -W "bash zsh fish" -- ${cur}) ) ;;
  esac
}

complete -F _clwatch_completions clwatch
`

const zshCompletion = `#compdef clwatch

_clwatch() {
  local -a commands
  commands=(
    'diff:Show version changes since last check'
    'list:List all tracked tools and their current versions'
    'refresh:Force-refresh a specific tool from the source'
    'init:Initialize a clwatch workspace in the current directory'
    'ack:Acknowledge a version (mark as seen)'
    'watch:Watch for updates and run on interval'
    'status:Show pipeline status from changelogs.info'
    'service:Install or manage the clwatch background service'
    'completion:Generate shell completion scripts'
    'version:Show version information'
    'help:Show help'
  )

  local -a tools
  tools=(claude-code codex-cli gemini-cli opencode openclaw)

  _arguments '1: :->command' '*: :->args'

  case $state in
    command) _describe 'command' commands ;;
    args)
      case $words[2] in
        diff)    _arguments '--json[JSON output]' '--verbose[Show all tools]' '--no-update[No state update]' ;;
        list)    _arguments '--json[JSON output]' ;;
        refresh) _arguments '1: :->tool' '--json[JSON output]' '--diff-only[Diff only]' '--all[All tools]'
                 [[ $state == tool ]] && _describe 'tool' tools ;;
        ack)     _arguments '1: :->tool' '2: :->version'
                 [[ $state == tool ]] && _describe 'tool' tools
                 [[ $state == version ]] && _message 'version (e.g. 2.1.74)' ;;
        init)    _arguments '--dir[Directory]:dir:_directories' '--tools[Tool list]:tools' '--force[Force overwrite]' ;;
        watch)   _arguments '--interval[Interval]:interval' '--json[JSON output]' '--webhook[Webhook URL]:url' ;;
        status)  _arguments '--json[JSON output]' ;;
        service)
          local -a subs
          subs=('install' 'uninstall' 'start' 'stop' 'status' 'logs')
          _describe 'service command' subs ;;
        completion)
          local -a shells; shells=('bash' 'zsh' 'fish')
          _describe 'shell' shells ;;
      esac
      ;;
  esac
}

_clwatch
`

const fishCompletion = `# Fish shell completion for clwatch

set -l tools claude-code codex-cli gemini-cli opencode openclaw
set -l commands diff list refresh init ack watch status service completion version help

complete -c clwatch -f

complete -c clwatch -n "not __fish_seen_subcommand_from $commands" -a diff       -d "Show version changes since last check"
complete -c clwatch -n "not __fish_seen_subcommand_from $commands" -a list       -d "List tracked tools and versions"
complete -c clwatch -n "not __fish_seen_subcommand_from $commands" -a refresh    -d "Force-refresh a specific tool"
complete -c clwatch -n "not __fish_seen_subcommand_from $commands" -a init       -d "Initialize a workspace"
complete -c clwatch -n "not __fish_seen_subcommand_from $commands" -a ack        -d "Acknowledge a version"
complete -c clwatch -n "not __fish_seen_subcommand_from $commands" -a watch      -d "Watch for updates on interval"
complete -c clwatch -n "not __fish_seen_subcommand_from $commands" -a status     -d "Show pipeline status"
complete -c clwatch -n "not __fish_seen_subcommand_from $commands" -a service    -d "Manage background service"
complete -c clwatch -n "not __fish_seen_subcommand_from $commands" -a completion -d "Generate shell completions"
complete -c clwatch -n "not __fish_seen_subcommand_from $commands" -a version    -d "Show version"

complete -c clwatch -n "__fish_seen_subcommand_from diff"    -l json      -d "JSON output"
complete -c clwatch -n "__fish_seen_subcommand_from diff"    -l verbose   -d "Show all tools"
complete -c clwatch -n "__fish_seen_subcommand_from diff"    -l no-update -d "No state update"
complete -c clwatch -n "__fish_seen_subcommand_from list"    -l json      -d "JSON output"
complete -c clwatch -n "__fish_seen_subcommand_from refresh" -a "$tools"
complete -c clwatch -n "__fish_seen_subcommand_from refresh" -l json      -d "JSON output"
complete -c clwatch -n "__fish_seen_subcommand_from refresh" -l diff-only -d "Diff only"
complete -c clwatch -n "__fish_seen_subcommand_from refresh" -l all       -d "All tools"
complete -c clwatch -n "__fish_seen_subcommand_from ack"     -a "$tools"
complete -c clwatch -n "__fish_seen_subcommand_from init"    -l dir       -d "Directory" -r -F
complete -c clwatch -n "__fish_seen_subcommand_from init"    -l tools     -d "Tool list" -r
complete -c clwatch -n "__fish_seen_subcommand_from init"    -l force     -d "Force overwrite"
complete -c clwatch -n "__fish_seen_subcommand_from watch"   -l interval  -d "Interval" -r
complete -c clwatch -n "__fish_seen_subcommand_from watch"   -l json      -d "JSON output"
complete -c clwatch -n "__fish_seen_subcommand_from watch"   -l webhook   -d "Webhook URL" -r
complete -c clwatch -n "__fish_seen_subcommand_from status"  -l json      -d "JSON output"
complete -c clwatch -n "__fish_seen_subcommand_from service" -a "install uninstall start stop status logs"
complete -c clwatch -n "__fish_seen_subcommand_from completion" -a "bash zsh fish"
`
