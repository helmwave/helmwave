package plan

import (
	"fmt"
	"strings"

	"github.com/helmwave/asciigraph"
	"github.com/helmwave/asciigraph/ascii"
	"github.com/helmwave/asciigraph/core"
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
)

func buildGraphMD(releases release.Configs) string {
	md := "# Depends On\n\n" +
		"```mermaid\ngraph RL\n"

	for _, r := range releases {
		for _, dep := range r.DependsOn() {
			md += fmt.Sprintf(
				"\t%s[%q] --> %s[%q]\n",
				strings.ReplaceAll(r.Uniq().String(), "@", "_"), r.Uniq(),
				strings.ReplaceAll(dep.Uniq().String(), "@", "_"), dep.Uniq().String(),
			)
		}
	}

	md += "```"

	return md
}

func buildGraphASCII(releases release.Configs, width int) string {
	list := make([]core.NodeInput, 0, len(releases))

	maxLength := 0
	minLength := 9999

	for _, rel := range releases {
		deps := make([]string, len(rel.DependsOn()))
		for i, d := range rel.DependsOn() {
			deps[i] = d.Uniq().String()
		}

		if len(rel.Uniq().String()) > maxLength {
			maxLength = len(rel.Uniq().String())
		}

		if len(rel.Uniq().String()) < minLength {
			minLength = len(rel.Uniq().String())
		}

		l := core.NodeInput{
			Id:   rel.Uniq().String(),
			Next: deps,
		}

		list = append(list, l)
	}

	o := ascii.DrawOptions{
		CellHeight:   3,
		MinCellWidth: 3,
		MaxWidth:     18,
		Padding:      1,
	}

	switch {
	case 0 == width:
		if maxLength > o.MinCellWidth {
			o.MinCellWidth = maxLength
		}
		if maxLength > o.MaxWidth {
			o.MaxWidth = maxLength
		}
	case 1 == width:
		return ""
	case 1 < width:
		o.MinCellWidth = width
		o.MaxWidth = width
	case 0 > width:
		o.MinCellWidth = minLength + width
		o.MaxWidth = maxLength + width
	}

	log.WithFields(log.Fields{
		"graph-width":     width,
		"max-word-length": maxLength,
		"min-word-length": minLength,
		"min-cell-width":  o.MinCellWidth,
		"max-cell-width":  o.MaxWidth,
	}).Debug("graph draw options")

	canvas, err := dgraph.DrawGraph(list, o)
	if err != nil {
		log.Fatal(err)
	}

	return canvas.String()
}

func (p *Plan) BuildGraphASCII(width int) string {
	return buildGraphASCII(p.body.Releases, width)
}
