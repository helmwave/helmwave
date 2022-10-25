package plan

import (
	"fmt"
	"strings"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/lempiy/dgraph"
	"github.com/lempiy/dgraph/core"
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

func buildGraphASCII(releases release.Configs) string {
	list := make([]core.NodeInput, 0, len(releases))

	for _, rel := range releases {
		deps := make([]string, len(rel.DependsOn()))
		for i, d := range rel.DependsOn() {
			deps[i] = d.Uniq().String()
		}

		l := core.NodeInput{
			Id:   rel.Uniq().String(),
			Next: deps,
		}

		list = append(list, l)
	}

	canvas, err := dgraph.DrawGraph(list)
	if err != nil {
		log.Fatal(err)
	}

	return canvas.String()
}
