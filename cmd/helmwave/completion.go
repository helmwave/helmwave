package main

import (
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"
)

const (
	bash = `#!/bin/bash

: ${PROG:=helmwave}

_helm() {
  if [[ "${COMP_WORDS[0]}" != "source" ]]; then
    local cur opts base
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    if [[ "$cur" == "-"* ]]; then
      opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} ${cur} --generate-bash-completion )
    else
      opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} --generate-bash-completion )
    fi
    COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
    return 0
  fi
}

complete -o bashdefault -o default -o nospace -F _helm $PROG
unset PROG
`

	zsh = `#compdef helmwave

_helmwave() {

  local -a opts
  local cur
  cur=${words[-1]}
  if [[ "$cur" == "-"* ]]; then
    opts=("${(@f)$(_CLI_ZSH_AUTOCOMPLETE_HACK=1 ${words[@]:0:#words[@]-1} ${cur} --generate-bash-completion)}")
  else
    opts=("${(@f)$(_CLI_ZSH_AUTOCOMPLETE_HACK=1 ${words[@]:0:#words[@]-1} --generate-bash-completion)}")
  fi

  if [[ "${opts[1]}" != "" ]]; then
    _describe 'values' opts
  else
    _files
  fi

  return
}

# don't run the completion function when being source-ed or eval-ed
if [ "$funcstack[1]" = "_helmwave" ]; then
    _helmwave
fi

compdef _helmwave helmwave
`
)

var (
	// ErrWrongShell is an error for unsupported shell.
	ErrWrongShell = errors.New("wrong shell")

	// ErrNotChose is an error for not provided shell name.
	ErrNotChose = errors.New("you did not specify a shell")
)

func completion() *cli.Command {
	return &cli.Command{
		Name:  "completion",
		Usage: "generate completion script",
		Description: `
			 echo "source <(helmwave completion bash)" >> ~/.bashrc
			 echo "source <(helmwave completion zsh)" >> ~/.zshrc"
		`,
		Action: func(c *cli.Context) error {
			if c.Args().Len() == 0 {
				return ErrNotChose
			}

			switch c.Args().First() {
			case "bash":
				fmt.Print(bash) //nolint:forbidigo

				return nil
			case "zsh":
				fmt.Print(zsh) //nolint:forbidigo

				return nil
			default:
				return ErrWrongShell
			}
		},
	}
}
