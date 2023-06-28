package hooks

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"

	"github.com/helmwave/helmwave/pkg/helper"
	log "github.com/sirupsen/logrus"
)

func runHooks(ctx context.Context, hooks []Hook) error {
	for _, h := range hooks {
		err := h.Run(ctx)
		if err != nil {
			h.Log().WithError(err).Error("failed to run hook")

			return err
		}
	}

	return nil
}

func (h *hook) Run(ctx context.Context) error {
	err := h.run(ctx)

	if h.AllowFailure {
		h.Log().WithError(err).Warn("caught lifecycle error, skipping...")

		return nil
	}

	return err
}

func (h *hook) run(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, h.Cmd, h.Args...)

	const t = "🩼 running hook..."

	switch h.Show {
	case true:
		h.Log().Info(t)
	case false:
		h.Log().Debug(t)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	cmd.Env = h.getCommandEnviron(ctx)

	// start the command after having set up the pipe
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to run command: %w", err)
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
		return fmt.Errorf("failed to read command stdout: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command returned error: %w", err)
	}

	return nil
}

// BUILD

func (l *Lifecycle) RunPreBuild(ctx context.Context) error {
	ctx = helper.ContextWithLifecycleType(ctx, "pre-build")

	if len(l.PreBuild) != 0 {
		log.Info("🩼 Running pre-build hooks...")

		return runHooks(ctx, l.PreBuild)
	}

	return nil
}

func (l *Lifecycle) RunPostBuild(ctx context.Context) error {
	ctx = helper.ContextWithLifecycleType(ctx, "post-build")

	if len(l.PostBuild) != 0 {
		log.Info("🩼 Running post-build hooks...")

		return runHooks(ctx, l.PostBuild)
	}

	return nil
}

// UP

func (l *Lifecycle) RunPreUp(ctx context.Context) error {
	ctx = helper.ContextWithLifecycleType(ctx, "pre-up")

	if len(l.PreUp) != 0 {
		log.Info("🩼 Running pre-up hooks...")

		return runHooks(ctx, l.PreUp)
	}

	return nil
}

func (l *Lifecycle) RunPostUp(ctx context.Context) error {
	ctx = helper.ContextWithLifecycleType(ctx, "post-up")

	if len(l.PostUp) != 0 {
		log.Info("🩼 Running post-up hooks...")

		return runHooks(ctx, l.PostUp)
	}

	return nil
}

// DOWN

func (l *Lifecycle) RunPreDown(ctx context.Context) error {
	ctx = helper.ContextWithLifecycleType(ctx, "pre-down")

	if len(l.PreDown) != 0 {
		log.Info("🩼 Running pre-down hooks...")

		return runHooks(ctx, l.PreDown)
	}

	return nil
}

func (l *Lifecycle) RunPostDown(ctx context.Context) error {
	ctx = helper.ContextWithLifecycleType(ctx, "post-down")

	if len(l.PostDown) != 0 {
		log.Info("🩼 Running post-down hooks...")

		return runHooks(ctx, l.PostDown)
	}

	return nil
}

// ROLLBACK

func (l *Lifecycle) RunPreRollback(ctx context.Context) error {
	ctx = helper.ContextWithLifecycleType(ctx, "pre-rollback")

	if len(l.PreRollback) != 0 {
		log.Info("🩼 Running pre-rollback hooks...")

		return runHooks(ctx, l.PreRollback)
	}

	return nil
}

func (l *Lifecycle) RunPostRollback(ctx context.Context) error {
	ctx = helper.ContextWithLifecycleType(ctx, "post-rollback")

	if len(l.PostRollback) != 0 {
		log.Info("🩼 Running post-rollback hooks...")

		return runHooks(ctx, l.PostRollback)
	}

	return nil
}
