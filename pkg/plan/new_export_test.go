package plan

import (
	"path/filepath"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/helmwave/helmwave/pkg/repo"
	"github.com/helmwave/helmwave/pkg/template"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	helm "helm.sh/helm/v3/pkg/cli"
	helmRelease "helm.sh/helm/v3/pkg/release"
	helmRepo "helm.sh/helm/v3/pkg/repo"
)

type MockReleaseConfig struct {
	mock.Mock
}

func (r *MockReleaseConfig) Uniq() uniqname.UniqName {
	r.Called()

	u, _ := uniqname.Generate(r.Name(), r.Namespace())

	return u
}

func (r *MockReleaseConfig) Sync() (*helmRelease.Release, error) {
	args := r.Called()

	return args.Get(0).(*helmRelease.Release), args.Error(1)
}

func (r *MockReleaseConfig) DryRun(_ bool) {
	r.Called()
}

func (r *MockReleaseConfig) ChartDepsUpd() error {
	return r.Called().Error(0)
}

func (r *MockReleaseConfig) In(a []release.Config) bool {
	return r.Called().Bool(0)
}

func (r *MockReleaseConfig) BuildValues(dir, templater string) error {
	args := r.Called()
	if errReturn := args.Error(0); errReturn != nil {
		return errReturn
	}

	for i := len(r.Values()) - 1; i >= 0; i-- {
		v := r.Values()[i]
		dst := filepath.Join(dir, Values, filepath.Base(v.Src))
		err := template.Tpl2yml(v.Src, dst, nil, templater)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *MockReleaseConfig) Uninstall() (*helmRelease.UninstallReleaseResponse, error) {
	args := r.Called()

	return args.Get(0).(*helmRelease.UninstallReleaseResponse), args.Error(1)
}

func (r *MockReleaseConfig) Get() (*helmRelease.Release, error) {
	args := r.Called()

	return args.Get(0).(*helmRelease.Release), args.Error(1)
}

func (r *MockReleaseConfig) List() (*helmRelease.Release, error) {
	args := r.Called()

	return args.Get(0).(*helmRelease.Release), args.Error(1)
}

func (r *MockReleaseConfig) Rollback(n int) error {
	return r.Called().Error(0)
}

func (r *MockReleaseConfig) Status() (*helmRelease.Release, error) {
	args := r.Called()

	return args.Get(0).(*helmRelease.Release), args.Error(1)
}

func (r *MockReleaseConfig) Name() string {
	return r.Called().String(0)
}

func (r *MockReleaseConfig) Namespace() string {
	return r.Called().String(0)
}

func (r *MockReleaseConfig) Chart() release.Chart {
	return r.Called().Get(0).(release.Chart)
}

func (r *MockReleaseConfig) DependsOn() []string {
	return r.Called().Get(0).([]string)
}

func (r *MockReleaseConfig) Tags() []string {
	return r.Called().Get(0).([]string)
}

func (r *MockReleaseConfig) Repo() string {
	return r.Called().String(0)
}

func (r *MockReleaseConfig) Values() []release.ValuesReference {
	return r.Called().Get(0).([]release.ValuesReference)
}

func (r *MockReleaseConfig) Logger() *log.Entry {
	return r.Called().Get(0).(*log.Entry)
}

func (r *MockReleaseConfig) AllowFailure() bool {
	return r.Called().Bool(0)
}

type MockRepoConfig struct {
	mock.Mock
}

func (r *MockRepoConfig) In(_ []repo.Config) bool {
	return r.Called().Bool(0)
}

func (r *MockRepoConfig) Install(_ *helm.EnvSettings, _ *helmRepo.File) error {
	return r.Called().Error(0)
}

func (r *MockRepoConfig) Name() string {
	return r.Called().String(0)
}

func (r *MockRepoConfig) URL() string {
	return r.Called().String(0)
}

func (r *MockRepoConfig) Logger() *log.Entry {
	return r.Called().Get(0).(*log.Entry)
}

func (p *Plan) NewBody() *planBody {
	p.body = &planBody{}

	return p.body
}

func (p *Plan) SetReleases(r ...*MockReleaseConfig) {
	if p.body == nil {
		p.NewBody()
	}
	c := make(release.Configs, len(r))
	for i := range r {
		c[i] = r[i]
	}
	p.body.Releases = c
}

func (p *Plan) SetRepositories(r ...*MockRepoConfig) {
	if p.body == nil {
		p.NewBody()
	}
	c := make(repo.Configs, len(r))
	for i := range r {
		c[i] = r[i]
	}
	p.body.Repositories = c
}
