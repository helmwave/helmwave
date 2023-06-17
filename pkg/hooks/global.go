package hooks

type Global struct {
	PreBuild  []hook `yaml:"pre_build"`
	PostBuild []hook `yaml:"post_build"`

	PreUp  []hook `yaml:"pre_up"`
	PostUp []hook `yaml:"post_up"`

	PreRollback  []hook `yaml:"pre_rollback"`
	PostRollback []hook `yaml:"post_rollback"`

	PreDown  []hook `yaml:"pre_down"`
	PostDown []hook `yaml:"post_down"`
}

func Run(hooks []hook) {
	for _, h := range hooks {
		h.Run()
	}
}
