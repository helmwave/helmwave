package plan

import (
	"fmt"
	"strings"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/lempiy/dgraph"
	"github.com/lempiy/dgraph/core"
	log "github.com/sirupsen/logrus"
)

func buildGraphMD(releases []release.Config) string {
	md :=
		"# Depends On\n\n" +
			"```mermaid\ngraph RL\n"

	for _, r := range releases {
		for _, dep := range r.DependsOn() {
			md += fmt.Sprintf(
				"\t%s[%q] --> %s[%q]\n",
				strings.Replace(string(r.Uniq()), "@", "_", -1), r.Uniq(), // nolint:gocritic
				strings.Replace(dep, "@", "_", -1), dep, // nolint:gocritic
			)
		}
	}

	md += "```"
	return md
}

func buildGraphASCII(releases []release.Config) string {
	list := make([]core.NodeInput, 0, len(releases))

	for _, rel := range releases {
		l := core.NodeInput{
			Id:   string(rel.Uniq()),
			Next: rel.DependsOn(),
		}

		list = append(list, l)
	}

	canvas, err := dgraph.DrawGraph(list)
	if err != nil {
		log.Fatal(err)
	}

	return canvas.String()
}
