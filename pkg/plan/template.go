package plan

import (
	"fmt"
	"sync"

	gotemplate "text/template"

	"github.com/helmwave/helmwave/pkg/template"
	"gopkg.in/yaml.v3"
)

//nolint:gocognit
func (p *Plan) templateFuncs(mu *sync.Mutex) gotemplate.FuncMap {
	funcMap := gotemplate.FuncMap{}

	// `getPlan` template function
	var plan map[string]any
	funcMap["getPlan"] = func() (map[string]any, error) {
		if plan == nil {
			mu.Lock()
			defer mu.Unlock()
			planYaml, err := yaml.Marshal(p.body)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal plan: %w", err)
			}
			plan, err = template.FromYaml(string(planYaml))
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal plan: %w", err)
			}
		}

		return plan, nil
	}

	// `getManifests` template function
	var manifests map[string][]any
	funcMap["getManifests"] = func(release string) ([]any, error) {
		if manifests == nil {
			mu.Lock()
			defer mu.Unlock()
			manifests = make(map[string][]any)
			for uniq, manifestYaml := range p.manifests {
				manifest, err := template.FromYamlAll(manifestYaml)
				if err != nil {
					return nil, fmt.Errorf("failed to unmarshal manifest: %w", err)
				}
				manifests[uniq.String()] = manifest
			}
		}
		manifest, found := manifests[release]
		if !found {
			return nil, fmt.Errorf("manifests for release %q not found", release)
		}

		return manifest, nil
	}

	return funcMap
}
