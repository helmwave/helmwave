package plan

import (
	"github.com/databus23/helm-diff/diff"
	"github.com/databus23/helm-diff/manifest"
	"github.com/helmwave/helmwave/pkg/helper"
	log "github.com/sirupsen/logrus"
	"os"
)

func (p *Plan) Diff(b *Plan, diffWide int) {
	var visited []string

	for _, rel := range p.body.Releases {
		m := rel.UniqName() + ".yml"
		visited = append(visited, m)
		oldSpecs := manifest.Parse(b.manifests[m], rel.Namespace)
		newSpecs := manifest.Parse(p.manifests[m], rel.Namespace)

		// Я пока так и не понял кто такой ваш context=10
		change := diff.Manifests(oldSpecs, newSpecs, []string{}, true, diffWide, os.Stdout)
		if !change {
			log.Info(rel.UniqName(), " no changes")
		}
	}

	for _, rel := range b.body.Releases {
		if !helper.Contains(rel.UniqName()+".yml", visited) {
			log.Warn(rel.UniqName(), " was found in previous planfile but not affected in new")
		}
	}
}
