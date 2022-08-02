package plan

import (
	"errors"
	"fmt"
	"github.com/helmwave/helmwave/pkg/release"
	"net/url"
	"os"
	"regexp"

	"github.com/helmwave/helmwave/pkg/release/uniqname"
)

// ErrValidateFailed is returned for failed values validation.
var ErrValidateFailed = errors.New("validate failed")

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

// ValidateValuesBuild Dst now is public method.
// Dst needs to marshal for export.
// Also, dst needs to unmarshal for import from plan.
func (p *Plan) ValidateValuesBuild() error {
	for _, rel := range p.body.Releases {
		err := release.ProhibitDst(rel.Values())
		if err != nil {
			return err
		}
	}

	return nil
}

// Validate validates releases and repositories in plan.
func (p *planBody) Validate() error {
	if len(p.Releases) == 0 && len(p.Repositories) == 0 {
		return errors.New("releases and repositories are empty")
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

	return nil
}

// ValidateRepositories validates all repositories.
func (p *planBody) ValidateRepositories() error {
	a := make(map[string]int8)
	for _, r := range p.Repositories {
		if r.Name() == "" {
			return errors.New("repository name is empty")
		}

		if r.URL() == "" {
			return errors.New("repository url is empty")
		}

		if _, err := url.Parse(r.URL()); err != nil {
			return errors.New("cant parse url: " + r.URL())
		}

		a[r.Name()]++
		if a[r.Name()] > 1 {
			return fmt.Errorf("repository %s duplicate", r.Name())
		}
	}

	return nil
}

func (p *planBody) ValidateRegistries() error {
	a := make(map[string]int8)
	for _, r := range p.Registries {
		if r.Host() == "" {
			return errors.New("registry name is empty")
		}

		a[r.Host()]++
		if a[r.Host()] > 1 {
			return fmt.Errorf("registry %s duplicate", r.Host())
		}
	}

	return nil
}

// ValidateReleases validates all releases.
func (p *planBody) ValidateReleases() error {
	a := make(map[uniqname.UniqName]int8)
	for _, r := range p.Releases {
		if r.Name() == "" {
			return errors.New("release name is empty")
		}

		if r.Namespace() == "" {
			r.Logger().Warnf("namespace is empty. I will use the namespace of your k8s context.")
		}

		if !validateNS(r.Namespace()) {
			return errors.New("bad namespace: " + r.Namespace())
		}

		if err := r.Uniq().Validate(); err != nil {
			return errors.New("bad uniqname: " + string(r.Uniq()))
		}

		a[r.Uniq()]++
		if a[r.Uniq()] > 1 {
			return errors.New("release duplicate: " + string(r.Uniq()))
		}
	}

	return nil
}

func validateNS(ns string) bool {
	r := regexp.MustCompile("[a-z0-9]([-a-z0-9]*[a-z0-9])?")

	return r.MatchString(ns)
}
