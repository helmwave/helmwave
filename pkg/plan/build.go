package plan

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"strings"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/helmwave/helmwave/pkg/repo"
	"github.com/lempiy/dgraph"
	"github.com/lempiy/dgraph/core"
	log "github.com/sirupsen/logrus"
)

// Build plan with yml and tags/matchALL options
func (p *Plan) Build(yml string, tags []string, matchAll bool) error {
	// Create Body
	body, err := NewBody(yml)
	if err != nil {
		return err
	}
	p.body = body

	// Build Releases
	p.body.Releases = buildReleases(tags, p.body.Releases, matchAll)
	if len(p.body.Releases) == 0 {
		return nil
	}

	// Build graph
	p.graphMD = buildGraphMD(p.body.Releases)
	log.Infof("Depends On:\n%s", buildGraphASCII(p.body.Releases))

	// Build Values
	err = p.buildValues(os.TempDir())
	if err != nil {
		return err
	}

	// Build Repositories
	repoMap, err := buildRepoMap(p.body.Releases)
	if err != nil {
		return err
	}

	log.Trace(repoMap)

	p.body.Repositories, err = buildRepo(repoMap, p.body.Repositories)
	if err != nil {
		return err
	}

	// Sync Repo
	err = p.syncRepositories()
	if err != nil {
		return err
	}

	// Build Manifest
	err = p.buildManifest()
	if err != nil {
		return err
	}

	return nil
}

func buildGraphMD(releases []*release.Config) string {
	md :=
		"# Depends On\n\n" +
			"```mermaid\ngraph RL\n"

	for _, r := range releases {
		for _, dep := range r.DependsOn {
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

func buildGraphASCII(releases []*release.Config) string {
	list := make([]core.NodeInput, 0, len(releases))

	for _, rel := range releases {
		l := core.NodeInput{
			Id:   string(rel.Uniq()),
			Next: rel.DependsOn,
		}

		list = append(list, l)
	}

	canvas, err := dgraph.DrawGraph(list)
	if err != nil {
		log.Fatal(err)
	}

	return canvas.String()
}

func (p *Plan) buildValues(dir string) error {
	wg := parallel.NewWaitGroup()
	wg.Add(len(p.body.Releases))

	for _, rel := range p.body.Releases {
		go func(wg *parallel.WaitGroup, rel *release.Config) {
			defer wg.Done()
			err := rel.BuildValues(dir)
			if err != nil {
				log.Fatal(err)
			}
		}(wg, rel)
	}

	return wg.Wait()
}

func (p *Plan) buildManifest() error {
	wg := parallel.NewWaitGroup()
	wg.Add(len(p.body.Releases))

	for _, rel := range p.body.Releases {
		go func(wg *parallel.WaitGroup, rel *release.Config) {
			defer wg.Done()

			rel.DryRun(true)
			r, err := rel.Sync()
			rel.DryRun(false)
			if err != nil || r == nil {
				log.Error("i cant generate manifest for ", rel.Uniq())
				log.Fatal(err)
			}

			var hooksManifests []string

			for _, h := range r.Hooks {
				hooksManifests = append(hooksManifests, h.Manifest)
			}

			hm, _ := yaml.Marshal(hooksManifests)

			document := r.Manifest
			if len(hooksManifests) > 0 {
				document += "# Hooks\n" + string(hm)
			}

			log.Trace(document)

			m := rel.Uniq() + ".yml"
			p.manifests[m] = document

			// log.Debug(rel.Uniq(), "`s manifest was successfully built ")
		}(wg, rel)
	}

	return wg.Wait()
}

func buildReleases(tags []string, releases []*release.Config, matchAll bool) (plan []*release.Config) {
	if len(tags) == 0 {
		return releases
	}

	releasesMap := make(map[uniqname.UniqName]*release.Config)

	for _, r := range releases {
		releasesMap[r.Uniq()] = r
	}

	for _, r := range releases {
		if checkTagInclusion(tags, r.Tags, matchAll) {
			plan = addToPlan(plan, r, releasesMap)
		}
	}

	return plan
}

func addToPlan(plan []*release.Config, rel *release.Config,
	releases map[uniqname.UniqName]*release.Config) []*release.Config {
	if rel.In(plan) {
		return plan
	}

	r := append(plan, rel) // nolint:gocritic

	for _, depName := range rel.DependsOn {
		depUN := uniqname.UniqName(depName)

		if dep, ok := releases[depUN]; ok {
			r = addToPlan(r, dep, releases)
		} else {
			log.Warnf("cannot find dependency %s in available releases, skipping it", depName)
		}
	}

	return r
}

// checkTagInclusion checks where any of release tags are included in target tags.
func checkTagInclusion(targetTags, releaseTags []string, matchAll bool) bool {
	for _, t := range targetTags {
		contains := helper.Contains(t, releaseTags)
		if matchAll && !contains {
			return false
		}
		if !matchAll && contains {
			return true
		}
	}

	return matchAll
}

func releaseNames(a []*release.Config) (n []string) {
	for _, r := range a {
		n = append(n, string(r.Uniq()))
	}

	return n
}

func buildRepo(m map[string][]*release.Config, in []*repo.Config) (out []*repo.Config, err error) {
	for rep, releases := range m {
		rm := releaseNames(releases)
		log.WithField(rep, rm).Debug("repo dependencies")

		if repoIsLocal(rep) {
			log.Infof("🗄 %q is local repo", rep)
		} else if index, found := repo.IndexOfName(in, rep); found {
			out = append(out, in[index])
			log.Infof("🗄 %q has been added to the plan", rep)
		} else {
			log.WithField("releases", rm).
				Warn("🗄 you will not be able to install this")
			return nil, errors.New("🗄 not found " + rep)
		}
	}

	return out, nil
}

func buildRepoMap(releases []*release.Config) (m map[string][]*release.Config, err error) {
	for _, rel := range releases {

		reps, err := rel.RepositoriesNames()
		if err != nil {
			log.Fatal("eto ", err)
			return nil, err
		}

		log.WithFields(log.Fields{
			"release":      rel.Uniq(),
			"repositories": reps,
		}).Trace("Repositories names")

		for _, rep := range reps {
			m[rep] = append(m[rep], rel)
		}
	}

	return m, err

}

// allRepos for releases
func allRepos(releases []*release.Config) ([]string, error) {
	var all []string
	for _, rel := range releases {
		r, err := rel.RepositoriesNames()
		if err != nil {
			return nil, err
		}

		all = append(all, r...)
	}

	return all, nil
}

// repoIsLocal return true if repo is dir
func repoIsLocal(repoString string) bool {
	if repoString == "" {
		return true
	}

	stat, err := os.Stat(repoString)
	if (err == nil || !os.IsNotExist(err)) && stat.IsDir() {
		return true
	}

	return false
}
