package plan

import (
	"fmt"
	"strings"

	"github.com/helmwave/asciigraph"
	"github.com/helmwave/asciigraph/ascii"
	"github.com/helmwave/asciigraph/core"
	"github.com/helmwave/helmwave/pkg/helper"
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
	list, maxLength, minLength := getGraphNodesForReleases(releases)

	o := ascii.DrawOptions{
		CellHeight:   3,
		MinCellWidth: 3,
		MaxWidth:     3,
		Padding:      1,
	}

	switch {
	case 0 == width:
		if minLength < o.MinCellWidth {
			o.MinCellWidth = minLength
		}
		if maxLength > o.MaxWidth {
			o.MaxWidth = maxLength + 4 // need to add a little bit so that it won't be shortened
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

	if o.MinCellWidth < 0 {
		log.Error("cannot output graph with no width available")

		return ""
	}

	canvas, err := dgraph.DrawGraph(list, o)
	if err != nil {
		log.Fatal(err)
	}

	return canvas.String()
}

func getGraphNodesForReleases(releases release.Configs) ([]core.NodeInput, int, int) {
	list := make([]core.NodeInput, 0, len(releases))

	maxLength := 0
	minLength := 9999

	for _, rel := range releases {
		deps := helper.SlicesMap(rel.DependsOn(), func(d *release.DependsOnReference) string {
			return d.Uniq().String()
		})

		uniq := rel.Uniq().String()
		if len(uniq) > maxLength {
			maxLength = len(uniq)
		}

		if len(uniq) < minLength {
			minLength = len(uniq)
		}

		l := core.NodeInput{
			Id:   uniq,
			Next: deps,
		}

		list = append(list, l)
	}

	return list, maxLength, minLength
}

func (p *Plan) BuildGraphASCII(width int) string {
	return buildGraphASCII(p.body.Releases, width)
}
