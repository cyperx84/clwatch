_clwatch_completions() {
  local cur prev
  COMPREPLY=()
  cur="${COMP_WORDS[COMP_CWORD]}"

  local commands="diff list refresh init ack watch status service completion version help"
  local tools="claude-code codex-cli gemini-cli opencode openclaw"

  if [[ ${COMP_CWORD} -eq 1 ]]; then
    COMPREPLY=( $(compgen -W "${commands}" -- ${cur}) )
    return 0
  fi

  case "${COMP_WORDS[1]}" in
    diff)
      COMPREPLY=( $(compgen -W "--json --verbose --no-update" -- ${cur}) )
      ;;
    list)
      COMPREPLY=( $(compgen -W "--json" -- ${cur}) )
      ;;
    refresh)
      if [[ ${COMP_CWORD} -eq 2 ]]; then
        COMPREPLY=( $(compgen -W "${tools} --all" -- ${cur}) )
      else
        COMPREPLY=( $(compgen -W "--json --diff-only --all" -- ${cur}) )
      fi
      ;;
    ack)
      if [[ ${COMP_CWORD} -eq 2 ]]; then
        COMPREPLY=( $(compgen -W "${tools}" -- ${cur}) )
      fi
      ;;
    init)
      COMPREPLY=( $(compgen -W "--dir --tools --force" -- ${cur}) )
      ;;
    watch)
      COMPREPLY=( $(compgen -W "--interval --json --webhook" -- ${cur}) )
      ;;
    status)
      COMPREPLY=( $(compgen -W "--json" -- ${cur}) )
      ;;
    service)
      COMPREPLY=( $(compgen -W "install uninstall start stop status logs" -- ${cur}) )
      ;;
    completion)
      COMPREPLY=( $(compgen -W "bash zsh fish" -- ${cur}) )
      ;;
  esac
}

complete -F _clwatch_completions clwatch
