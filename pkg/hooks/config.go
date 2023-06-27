package hooks

import (
	"context"

	log "github.com/sirupsen/logrus"
)

type Lifecycle struct {
	PreBuild  []hook `yaml:"pre_build" json:"pre_build" jsonschema:"title=pre_build,description=pre_build hooks"`
	PostBuild []hook `yaml:"post_build" json:"post_build" jsonschema:"title=post_build,description=post_build hooks"`

	PreUp  []hook `yaml:"pre_up" json:"pre_up" jsonschema:"title=pre_up,description=pre_up hooks"`
	PostUp []hook `yaml:"post_up" json:"post_up" jsonschema:"title=post_up,description=post_up hooks"`

	PreRollback  []hook `yaml:"pre_rollback" json:"pre_rollback" jsonschema:"title=pre_rollback,description=pre_rollback hooks"`
	PostRollback []hook `yaml:"post_rollback" json:"post_rollback" jsonschema:"title=post_rollback,description=post_rollback hooks"`

	PreDown  []hook `yaml:"pre_down" json:"pre_down" jsonschema:"title=pre_down,description=pre_down hooks"`
	PostDown []hook `yaml:"post_down" json:"post_down" jsonschema:"title=post_down,description=post_down hooks"`
}

type Hook interface {
	Run(context.Context) error
	Log() *log.Entry
}

type hook struct {
	Cmd          string
	Args         []string
	Show         bool
	AllowFailure bool
}

func (h *hook) Log() *log.Entry {
	return log.WithFields(log.Fields{
		"cmd":  h.Cmd,
		"args": h.Args,
	})
}
