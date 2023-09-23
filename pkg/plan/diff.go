package plan

import (
	"context"
	"io/fs"
	"reflect"
	"strings"
	"sync"

	"github.com/databus23/helm-diff/v3/diff"
	"github.com/databus23/helm-diff/v3/manifest"
	"github.com/helmwave/helmwave/pkg/helper"
	logSetup "github.com/helmwave/helmwave/pkg/log"
	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	structDiff "github.com/r3labs/diff/v3"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/chart"
	live "helm.sh/helm/v3/pkg/release"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/cli-runtime/pkg/resource"
)

// SkippedAnnotations is a map with all annotations to be skipped by differ.
//
//nolint:gochecknoglobals // can't make this const
var SkippedAnnotations = map[string][]string{
	live.HookAnnotation:               {string(live.HookTest), "test-success", "test-failure"},
	helper.RootAnnoName + "skip-diff": {"true"},
}

// DiffPlan show diff between 2 plans.
func (p *Plan) DiffPlan(b *Plan, showSecret bool, diffWide int) {
	visited := make(map[uniqname.UniqName]bool)
	k := 0
	opts := &diff.Options{
		ShowSecrets:   showSecret,
		OutputContext: diffWide,
		OutputFormat:  logSetup.Default.Format(),
	}

	for _, rel := range append(p.body.Releases, b.body.Releases...) {
		if visited[rel.Uniq()] {
			continue
		}
		visited[rel.Uniq()] = true

		oldSpecs := parseManifests(b.manifests[rel.Uniq()], rel.Namespace())
		newSpecs := parseManifests(p.manifests[rel.Uniq()], rel.Namespace())

		change := diff.Manifests(oldSpecs, newSpecs, opts, log.StandardLogger().Out)
		if !change {
			k++
			log.Info("🆚 ❎ ", rel.Uniq(), " no changes")
			p.unchanged = append(p.unchanged, rel)
		}
	}

	visitedNames := make([]uniqname.UniqName, 0, len(visited))
	for n := range visited {
		visitedNames = append(visitedNames, n)
	}

	showChangesReport(p.body.Releases, visitedNames, k)
}

// DiffLive show diff with production releases in k8s-cluster.
func (p *Plan) DiffLive(ctx context.Context, baseFS fs.FS, showSecret bool, diffWide int, threeWayMerge bool) {
	alive, _, err := p.GetLive(ctx)
	if err != nil {
		log.Fatalf("Something went wrong with getting releases in the kubernetes cluster: %v", err)
	}

	visited := make([]uniqname.UniqName, 0, len(p.body.Releases))
	k := 0
	opts := &diff.Options{
		ShowSecrets:   showSecret,
		OutputContext: diffWide,
	}

	for _, rel := range p.body.Releases {
		visited = append(visited, rel.Uniq())
		if active, ok := alive[rel.Uniq()]; ok {
			newManifest := p.manifests[rel.Uniq()]
			oldManifest := active.Manifest
			if threeWayMerge {
				oldManifest = get3WayMergeManifests(rel, active.Manifest)
			}
			// I dont use manifest.ParseRelease
			// Because Structs are different.
			oldSpecs := parseManifests(oldManifest, rel.Namespace())
			newSpecs := parseManifests(newManifest, rel.Namespace())

			change := diff.Manifests(oldSpecs, newSpecs, opts, rel.Logger().Logger.Out)
			chartChange := diffCharts(ctx, baseFS, active.Chart, rel, rel.Logger())

			if !change && !chartChange {
				k++
				rel.Logger().Info("🆚 ❎ no changes")
				p.unchanged = append(p.unchanged, rel)
			}
		}
	}

	showChangesReport(p.body.Releases, visited, k)
}

func get3WayMergeManifests(rel release.Config, oldManifest string) string { //nolint:funlen,gocognit
	cfg := rel.Cfg()

	err := cfg.KubeClient.IsReachable()
	if err != nil {
		rel.Logger().WithError(err).Warn("failed to connect to k8s to run 3-way merge, skipping")

		return oldManifest
	}

	oldResources, err := cfg.KubeClient.Build(strings.NewReader(oldManifest), false)
	if err != nil {
		rel.Logger().WithError(err).Warn("failed to build old resources list for 3-way merge, skipping")

		return oldManifest
	}

	updatedManifest := ""

	err = oldResources.Visit(func(r *resource.Info, err error) error {
		if err != nil {
			return err
		}

		h := resource.NewHelper(r.Client, r.Mapping)
		currentObject, err := h.Get(r.Namespace, r.Name)
		if err != nil {
			if !apiErrors.IsNotFound(err) {
				return err //nolint:wrapcheck
			}

			return nil
		}

		out, err := yaml.Marshal(currentObject)
		if err != nil {
			return err //nolint:wrapcheck
		}
		// currentObject stores everything under 'object' key.
		// We need to get everything from this field and drop some generated parts.
		var ra map[string]any
		_ = yaml.Unmarshal(out, &ra)
		obj := ra["object"].(map[string]any) //nolint:forcetypeassert
		delete(obj, "status")

		metadata := obj["metadata"].(map[string]any) //nolint:forcetypeassert
		delete(metadata, "creationTimestamp")
		delete(metadata, "generation")
		delete(metadata, "managedFields")
		delete(metadata, "resourceVersion")
		delete(metadata, "uid")

		if a := metadata["annotations"]; a != nil {
			annotations := a.(map[string]any) //nolint:forcetypeassert
			delete(annotations, "meta.helm.sh/release-name")
			delete(annotations, "meta.helm.sh/release-namespace")
			delete(annotations, "deployment.kubernetes.io/revision")

			if len(annotations) == 0 {
				delete(metadata, "annotations")
			}
		}

		out, _ = yaml.Marshal(obj)
		updatedManifest += "\n---\n" + string(out)

		return nil
	})
	if err != nil {
		rel.Logger().WithError(err).Warn("failed to get latest objects for 3-way merge, skipping")

		return oldManifest
	}

	return updatedManifest
}

//nolint:gocritic // cannot change argument types as it is required by diff library
func diffChartsFilter(path []string, _ reflect.Type, _ reflect.StructField) bool {
	return len(path) >= 1 && path[0] == "Metadata"
}

func diffCharts(ctx context.Context, baseFS fs.FS, oldChart *chart.Chart, rel release.Config, l log.FieldLogger) bool {
	l.Info("getting charts diff")

	dryRunRelease, err := rel.SyncDryRun(ctx, baseFS)
	if err != nil {
		l.WithError(err).Error("failed to get dry-run release")

		return false
	}

	newChart := dryRunRelease.Chart

	changelog, err := structDiff.Diff(oldChart, newChart, structDiff.Filter(diffChartsFilter))
	if err != nil {
		l.WithError(err).Error("failed to get diff of charts")

		return false
	}

	if len(changelog) == 0 {
		return false
	}

	for i := range changelog {
		change := changelog[i]
		l.WithField("path", strings.Join(change.Path, ".")).Infof("🆚 %q -> %q", change.From, change.To)
	}

	return true
}

func parseManifests(m, ns string) map[string]*manifest.MappingResult {
	manifests := manifest.Parse(m, ns, true)

	type annotationManifest struct {
		Metadata struct {
			Annotations map[string]string
		}
	}

	for k := range manifests {
		parsed := annotationManifest{}

		if err := yaml.Unmarshal([]byte(manifests[k].Content), &parsed); err != nil {
			log.WithError(err).WithField("content", manifests[k].Content).Debug("failed to decode manifest")

			continue
		}

		for anno := range parsed.Metadata.Annotations {
			if helper.Contains(parsed.Metadata.Annotations[anno], SkippedAnnotations[anno]) {
				log.WithFields(log.Fields{
					"resource":   manifests[k].Name,
					"annotation": anno,
				}).Debug("resource diff is skipped due to annotation")
				delete(manifests, k)

				continue
			}
		}
	}

	return manifests
}

// showChangesReport help function for reporting helm-diff.
func showChangesReport(releases []release.Config, visited []uniqname.UniqName, k int) {
	previous := false
	for _, rel := range releases {
		if !helper.In(rel.Uniq(), visited) {
			previous = true
			rel.Logger().Warn("🆚 release was found in previous plan but not affected in new")
		}
	}

	if k == len(releases) && !previous {
		log.Info("🆚 🌝 Plan has no changes")
	}
}

// GetLive returns maps of releases in a k8s-cluster.
func (p *Plan) GetLive(
	ctx context.Context,
) (found map[uniqname.UniqName]*live.Release, notFound []uniqname.UniqName, err error) {
	wg := parallel.NewWaitGroup()
	wg.Add(len(p.body.Releases))

	found = make(map[uniqname.UniqName]*live.Release)
	mu := &sync.Mutex{}

	for i := range p.body.Releases {
		go func(wg *parallel.WaitGroup, mu *sync.Mutex, rel release.Config) {
			defer wg.Done()

			r, err := rel.Get(0)

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				log.Warnf("I can't get release from k8s: %v", err)
				//nolint:revive // we are under mutex here
				notFound = append(notFound, rel.Uniq())
			} else {
				//nolint:revive // we are under mutex here
				found[rel.Uniq()] = r
			}
		}(wg, mu, p.body.Releases[i])
	}

	if err := wg.WaitWithContext(ctx); err != nil {
		return nil, nil, err
	}

	return found, notFound, nil
}
