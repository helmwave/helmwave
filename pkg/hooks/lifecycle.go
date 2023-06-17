package hooks

import log "github.com/sirupsen/logrus"

type Lifecycle struct {
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

// BUILD

func (l *Lifecycle) PreBuilding() {
	if len(l.PreBuild) != 0 {
		log.Info("🩼 Running pre-build hooks...")
		Run(l.PreBuild)
	}
}

func (l *Lifecycle) PostBuilding() {
	if len(l.PostBuild) != 0 {
		log.Info("🩼 Running post-build hooks...")
		Run(l.PostBuild)
	}
}

// UP

func (l *Lifecycle) PreUping() {
	if len(l.PreUp) != 0 {
		log.Info("🩼 Running pre-up hooks...")
		Run(l.PreUp)
	}
}

func (l *Lifecycle) PostUping() {
	if len(l.PostUp) != 0 {
		log.Info("🩼 Running post-up hooks...")
		Run(l.PostUp)
	}
}

// DOWN

func (l *Lifecycle) PreDowning() {
	if len(l.PreDown) != 0 {
		log.Info("🩼 Running pre-down hooks...")
		Run(l.PreDown)
	}
}

func (l *Lifecycle) PostDowning() {
	if len(l.PostDown) != 0 {
		log.Info("🩼 Running post-down hooks...")
		Run(l.PostDown)
	}
}

// ROLLBACK

func (l *Lifecycle) PreRolling() {
	if len(l.PreRollback) != 0 {
		log.Info("🩼 Running pre-rollback hooks...")
		Run(l.PreRollback)
	}
}

func (l *Lifecycle) PostRolling() {
	if len(l.PostRollback) != 0 {
		log.Info("🩼 Running post-rollback hooks...")
		Run(l.PostRollback)
	}
}
