package hooks

import (
	"bufio"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

type Hook interface {
	Run()
	Log() *log.Entry
}

type hook struct {
	Cmd  string
	Args []string
	Show bool
}

func (h *hook) Log() *log.Entry {
	return log.WithFields(log.Fields{
		"cmd":  h.Cmd,
		"args": h.Args,
	})
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
