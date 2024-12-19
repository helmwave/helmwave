package main

import (
	"fmt"

	"github.com/helmwave/helmwave/pkg/action"
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

	fish = `function __fish_helmwave_generate_completions
    set -l args (commandline -opc)
    set -l current_token (commandline -ct)
    if test (string match -r "^-" -- $current_token)
        eval $args $current_token --generate-bash-completion
    else
        eval $args --generate-bash-completion
    end
end

function __fish_helmwave_complete
    set -l completions (__fish_helmwave_generate_completions)
    for opt in $completions
        echo "$opt"
    end
end

complete -c helmwave -f -a '(__fish_helmwave_complete)'
`
)

func completion() *cli.Command {
	return &cli.Command{
		Name:     "completion",
		Category: action.Step_,
		Usage:    "generate completion script",
		Description: `
			echo "source <(helmwave completion bash)" >> ~/.bashrc
			echo "source <(helmwave completion zsh)" >> ~/.zshrc
			helmwave completion fish > ~/.config/fish/functions/helmwave.fish
		`,
		Subcommands: []*cli.Command{
			{
				Name:     "bash",
				Category: action.Step_,
				Usage:    "generate bash completion script",
				Action: func(c *cli.Context) error {
					fmt.Print(bash)

					return nil
				},
			},
			{
				Name:     "zsh",
				Category: action.Step_,
				Usage:    "generate zsh completion script",
				Action: func(c *cli.Context) error {
					fmt.Print(zsh)

					return nil
				},
			},
			{
				Name:     "fish",
				Category: action.Step_,
				Usage:    "generate fish completion script",
				Action: func(c *cli.Context) error {
					fmt.Print(fish)

					return nil
				},
			},
		},
	}
}
