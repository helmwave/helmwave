package plan

import (
	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/helmwave/helmwave/pkg/repo"
	"github.com/helmwave/helmwave/pkg/template"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

// Build Принимает на вход ямл и опции по работе с ним.
func (p *Plan) Build(yml string, tags []string, matchAll bool) error {
	// Create Body
	body, err := NewBody(yml)
	if err != nil {
		return err
	}
	p.body = body

	// Build Releases
	p.body.Releases = buildReleases(tags, p.body.Releases, matchAll)

	// Build Values
	err = p.buildValues(os.TempDir())
	if err != nil {
		return err
	}

	// Build Repositories
	p.body.Repositories, err = buildRepo(p.body.Releases, p.body.Repositories)
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

// Todo: Split this function for several more KISS
func (p *Plan) buildValues(dir string) error {

	wg := parallel.NewWaitGroup()
	wg.Add(len(p.body.Releases))

	for _, rel := range p.body.Releases {
		go func(wg *parallel.WaitGroup, rel *release.Config) {
			defer wg.Done()

			// Todo: parallel
			for i, v := range rel.Values {
				dst := dir + Values + string(rel.Uniq()) + "/" + string(rune(i)) + ".yml"
				err := v.Set(dst)
				if err != nil {
					log.Warn(v.Src, " skipping: ", err)
					continue
				}

				err = template.Tpl2yml(dst, dst, struct{ Release *release.Config }{rel})
				if err != nil {
					log.Error(err)
					continue
				}

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
			if err != nil {
				log.Fatal(err)
			}

			if r != nil {
				log.Trace(r.Manifest)
			}

			m := rel.Uniq() + ".yml"
			p.manifests[m] = r.Manifest

			//log.Debug(rel.Uniq(), "`s manifest was successfully built ")

		}(wg, rel)

	}

	return wg.Wait()
}

func (p *Plan) PrettyPlan() {
	var a []string
	for _, r := range p.body.Releases {
		a = append(a, string(r.Uniq()))
	}

	var b []string
	for _, r := range p.body.Repositories {
		b = append(b, r.Name)
	}

	log.WithFields(log.Fields{
		"releases":     a,
		"repositories": b,
	}).Info("🏗 Plan")
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

func addToPlan(plan []*release.Config, rel *release.Config, releases map[uniqname.UniqName]*release.Config) []*release.Config {
	if rel.In(plan) {
		return plan
	}

	r := append(plan, rel)

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

func buildRepo(releases []*release.Config, repositories []*repo.Config) (plan []*repo.Config, err error) {
	all := getRepositories(releases)

	for _, a := range all {
		found := false
		for _, b := range repositories {
			if a == b.Name {
				found = true
				if !b.InByName(plan) {
					plan = append(plan, b)
					log.Debugf("🗄 %q has been added to the plan", a)
				}
			}
		}

		if !found {
			log.Errorf("🗄 %q not found ", a)
			return plan, repo.ErrNotFound
		}
	}

	return plan, nil
}

// getRepositories for releases
func getRepositories(releases []*release.Config) (repos []string) {
	for _, rel := range releases {
		rep := strings.Split(rel.Chart.Name, "/")[0]
		deps, _ := rel.Repositories()

		all := deps
		if repoIsLocal(rep) {
			log.Infof("🗄 %q is local repo", rep)
		} else {
			all = append(all, rep)
		}

		for _, r := range all {
			if !helper.Contains(r, repos) {
				repos = append(repos, r)
			}
		}
	}

	return repos
}

// repoIsLocal
func repoIsLocal(repo string) bool {
	if repo == "" {
		return true
	}

	stat, err := os.Stat(repo)
	if (err == nil || !os.IsNotExist(err)) && stat.IsDir() {
		return true
	}

	return false
}
