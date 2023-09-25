package hooks

import (
	"github.com/invopop/jsonschema"
	log "github.com/sirupsen/logrus"
)

var _ Hook = (*hook)(nil)

type Lifecycle struct {
	PreBuild  Hooks `yaml:"pre_build,omitempty" json:"pre_build,omitempty" jsonschema:"title=pre_build,description=pre_build hooks"`
	PostBuild Hooks `yaml:"post_build,omitempty" json:"post_build,omitempty" jsonschema:"title=post_build,description=post_build hooks"`

	PreUp  Hooks `yaml:"pre_up,omitempty" json:"pre_up,omitempty" jsonschema:"title=pre_up,description=pre_up hooks"`
	PostUp Hooks `yaml:"post_up,omitempty" json:"post_up,omitempty" jsonschema:"title=post_up,description=post_up hooks"`

	PreRollback  Hooks `yaml:"pre_rollback,omitempty" json:"pre_rollback,omitempty" jsonschema:"title=pre_rollback,description=pre_rollback hooks"`
	PostRollback Hooks `yaml:"post_rollback,omitempty" json:"post_rollback,omitempty" jsonschema:"title=post_rollback,description=post_rollback hooks"`

	PreDown  Hooks `yaml:"pre_down,omitempty" json:"pre_down,omitempty" jsonschema:"title=pre_down,description=pre_down hooks"`
	PostDown Hooks `yaml:"post_down,omitempty" json:"post_down,omitempty" jsonschema:"title=post_down,description=post_down hooks"`
}

type Hooks []Hook

func (Hooks) JSONSchema() *jsonschema.Schema {
	r := &jsonschema.Reflector{
		DoNotReference:             true,
		RequiredFromJSONSchemaTags: true,
	}
	var l []*hook

	return r.Reflect(&l)
}

type hook struct {
	Cmd          string   `yaml:"cmd" json:"cmd" jsonschema:"required,title=cmd,description=executable to run"`
	Args         []string `yaml:"args" json:"args" jsonschema:"title=args,description=arguments to pass to executable"`
	Show         bool     `yaml:"show" json:"show" jsonschema:"title=show,description=whether to log command stdout,default=true"`
	AllowFailure bool     `yaml:"allow_failure" json:"allow_failure" jsonschema:"title=allow_failure,description=whether to fail the whole helmwave if command fail,default=false"`
}

func (hook) JSONSchemaExtend(schema *jsonschema.Schema) {
	schema.OneOf = []*jsonschema.Schema{
		{
			Type: "string",
		},
		{
			Type: "object",
		},
	}
	schema.Type = ""
}

func (h *hook) Log() *log.Entry {
	return log.WithFields(log.Fields{
		"cmd":  h.Cmd,
		"args": h.Args,
	})
}
