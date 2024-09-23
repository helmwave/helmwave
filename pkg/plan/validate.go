package plan

import (
	"os"

	"github.com/helmwave/helmwave/pkg/fileref"

	"github.com/helmwave/helmwave/pkg/monitor"
	"github.com/helmwave/helmwave/pkg/registry"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/helmwave/helmwave/pkg/repo"
	log "github.com/sirupsen/logrus"
)

// ValidateValuesImport checks whether all values files exist.
func (p *Plan) ValidateValuesImport() error {
	f := false
	for _, rel := range p.body.Releases {
		for i := range rel.Values() {
			y := rel.Values()[i].Dst
			_, err := os.Stat(y)
			if os.IsNotExist(err) {
				f = true
				rel.Logger().Errorf("❌ values %q", rel.Values()[i].Src)
			} else if err != nil {
				f = true
				rel.Logger().WithError(err).Errorf("failed to open values %s", y)
			}
		}
	}
	if !f {
		return nil
	}

	return ErrValidateFailed
}

// ValidateValuesBuild Dst now is a public method.
// Dst needs to marshal for export.
// Also, dst needs to unmarshal for import from plan.
func (p *Plan) ValidateValuesBuild() error {
	for _, rel := range p.body.Releases {
		err := fileref.ProhibitDst(rel.Values())
		if err != nil {
			return err
		}
	}

	return nil
}

// Validate validates releases and repositories in plan.
func (p *planBody) Validate() error {
	if len(p.Releases) == 0 && len(p.Repositories) == 0 && len(p.Registries) == 0 {
		log.Warn("releases, repositories and registries are empty")

		return nil
	}

	if err := p.ValidateRegistries(); err != nil {
		return err
	}

	if err := p.ValidateRepositories(); err != nil {
		return err
	}

	if err := p.ValidateReleases(); err != nil {
		return err
	}

	if err := p.ValidateMonitors(); err != nil {
		return err
	}

	return nil
}

// ValidateRepositories validates all repositories.
func (p *planBody) ValidateRepositories() error {
	a := make(map[string]bool)
	for _, r := range p.Repositories {
		err := r.Validate()
		if err != nil {
			return err
		}

		if a[r.Name()] {
			return repo.NewDuplicateError(r.Name())
		}

		a[r.Name()] = true
	}

	return nil
}

func (p *planBody) ValidateRegistries() error {
	a := make(map[string]bool)
	for _, r := range p.Registries {
		err := r.Validate()
		if err != nil {
			return err
		}

		if a[r.Host()] {
			return registry.NewDuplicateError(r.Host())
		}

		a[r.Host()] = true
	}

	return nil
}

// ValidateReleases validates all releases.
func (p *planBody) ValidateReleases() error {
	a := make(map[uniqname.UniqName]bool)
	for _, r := range p.Releases {
		err := r.Validate()
		if err != nil {
			return err
		}

		if a[r.Uniq()] {
			return release.NewDuplicateError(r.Uniq())
		}

		a[r.Uniq()] = true
	}

	_, err := p.generateDependencyGraph()
	if err != nil {
		return err
	}

	return nil
}

func (p *planBody) ValidateMonitors() error {
	a := make(map[string]bool)
	for _, r := range p.Monitors {
		err := r.Validate()
		if err != nil {
			return err
		}

		if a[r.Name()] {
			return monitor.NewDuplicateError(r.Name())
		}

		a[r.Name()] = true
	}

	for _, r := range p.Releases {
		mons := r.Monitors()
		for i := range mons {
			mon := mons[i]
			if !a[mon.Name] {
				return monitor.NewNotExistsError(mon.Name)
			}
		}
	}

	return nil
}
