package action

import (
	"context"
	"fmt"

	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/urfave/cli/v2"
)

var _ Action = (*Manifests)(nil)

// Manifests is a struct for running 'Manifests' command.
type Manifests struct {
	build     *Build
	names     cli.StringSlice
	autoBuild bool
}

// Run is the main function for 'status' command.
//
//nolint:forbidigo
func (l *Manifests) Run(ctx context.Context) error {
	p, err := plan.NewAndImport(ctx, l.build.plandir)
	if err != nil {
		return err
	}

	names := l.names.Value()

	// Don't use maps.Keys here because in array you have to copy each element inside "for"
	if len(names) == 0 {
		for _, m := range p.Manifests() {
			fmt.Println(m)
		}

		return nil
	}

	for _, name := range names {
		n, _ := uniqname.NewFromString(name)
		fmt.Println(p.Manifests()[n])
	}

	return nil
}

// Cmd returns 'status' *cli.Command.
func (l *Manifests) Cmd() *cli.Command {
	return &cli.Command{
		Name:     "manifests",
		Aliases:  []string{"manifest"},
		Category: Step1,
		Usage:    "show only manifests",
		Flags:    l.flags(),
		Action:   toCtx(l.Run),
	}
}

// flags return flag set of CLI urfave.
func (l *Manifests) flags() []cli.Flag {
	// Init sub-structures
	l.build = &Build{}

	self := []cli.Flag{
		flagAutoBuild(&l.autoBuild),
		&cli.StringSliceFlag{
			Name:     "uniqnames",
			Aliases:  []string{"u"},
			Usage:    "show manifest only for specific release: -u nginx@namespace -u nginx@namespace,redis@namespace",
			Category: "SELECTION",
			EnvVars:  EnvVars("UNIQNAMES"),
		},
	}

	return append(self, l.build.flags()...)
}

// func (l *Manifests) uniqnames() (r []uniqname.UniqName) {
//	for _, name := range l.names.Value() {
//		n, _ := uniqname.NewFromString(name)
//		r = append(r, n)
//	}
//
//	return r
//}
