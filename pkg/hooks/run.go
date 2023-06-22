package hooks

import (
	"bufio"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func Run(hooks []hook) {
	for _, h := range hooks {
		h.Run()
	}
}

func (h *hook) Run() {
	cmd := exec.Command(h.Cmd, h.Args...)

	const t = "🩼 running hook..."

	switch h.Show {
	case true:
		h.Log().Info(t)
	case false:
		h.Log().Debug(t)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	// start the command after having set up the pipe
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	// read command's stdout line by line
	in := bufio.NewScanner(stdout)

	for in.Scan() {
		switch h.Show {
		case true:
			log.Info(in.Text())
		case false:
			log.Debug(in.Text())
		}
	}
	if err := in.Err(); err != nil {
		log.Fatal(err)
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
