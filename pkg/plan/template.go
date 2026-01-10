package plan

import (
	"fmt"
	"sync"

	gotemplate "text/template"

	"github.com/helmwave/helmwave/pkg/template"
	"gopkg.in/yaml.v3"
)

//nolint:gocognit,cyclop,funlen
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

	// `getValues` template function
	var values map[string]map[string]any
	funcMap["getValues"] = func(release string, filename string) (any, error) {
		if values == nil {
			mu.Lock()
			defer mu.Unlock()
			values = make(map[string]map[string]any)
			for uniq, valuesYaml := range p.values {
				releaseValues := make(map[string]any)
				for filename, value := range valuesYaml {
					value, err := template.FromYaml(value)
					if err != nil {
						return nil, fmt.Errorf("failed to unmarshal value: %w", err)
					}
					releaseValues[filename] = value
				}
				values[uniq.String()] = releaseValues
			}
		}
		releaseValues, found := values[release]
		if !found {
			return nil, fmt.Errorf("values for release %q not found", release)
		}

		value, found := releaseValues[filename]
		if !found {
			return nil, fmt.Errorf("values file %q for release %q not found", filename, release)
		}

		return value, nil
	}

	return funcMap
}
