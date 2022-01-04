package plan

import (
	"errors"
	"net/url"
	"os"
	"regexp"

	"github.com/helmwave/helmwave/pkg/release/uniqname"
	log "github.com/sirupsen/logrus"
)

// ErrValidateFailed is returned for failed values validation.
var ErrValidateFailed = errors.New("validate failed")

// ValidateValues checkes whether all values files exist.
func (p *Plan) ValidateValues() error {
	f := false
	for _, rel := range p.body.Releases {
		for i := range rel.Values() {
			p := rel.Values()[i].Get()
			_, err := os.Stat(p)
			if os.IsNotExist(err) {
				f = true
				log.WithError(err).Errorf("❌ %s values (%s)", rel.Uniq(), rel.Values()[i].Src)
			} else if err != nil {
				f = true
				log.WithError(err).Errorf("failed to open values %s", p)
			}
		}
	}
	if !f {
		return nil
	}

	return ErrValidateFailed
}

// Validate validates releases and repositories in plan.
func (p *planBody) Validate() error {
	if len(p.Releases) == 0 && len(p.Repositories) == 0 {
		return errors.New("releases and repositories are empty")
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
			return errors.New("repository name duplicate: " + r.Name())
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
			log.Warnf("namespace for %q is empty. I will use the namespace of your k8s context.", r.Uniq())
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
