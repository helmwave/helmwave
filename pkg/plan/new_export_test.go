package plan

import (
	"path/filepath"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/helmwave/helmwave/pkg/repo"
	"github.com/helmwave/helmwave/pkg/template"
	"github.com/stretchr/testify/mock"
	helm "helm.sh/helm/v3/pkg/cli"
	helmRelease "helm.sh/helm/v3/pkg/release"
	helmRepo "helm.sh/helm/v3/pkg/repo"
)

type mockReleaseConfig struct {
	mock.Mock
}

func (r *mockReleaseConfig) Uniq() uniqname.UniqName {
	return uniqname.UniqName(r.Name() + uniqname.Separator + r.Namespace())
}

func (r *mockReleaseConfig) HandleDependencies(_ []release.Config) {
	panic("not implemented") // TODO: Implement
}

func (r *mockReleaseConfig) Sync() (*helmRelease.Release, error) {
	panic("not implemented") // TODO: Implement
}

func (r *mockReleaseConfig) NotifySuccess() {
	panic("not implemented") // TODO: Implement
}

func (r *mockReleaseConfig) NotifyFailed() {
	panic("not implemented") // TODO: Implement
}

func (r *mockReleaseConfig) DryRun(_ bool) {
	panic("not implemented") // TODO: Implement
}

func (r *mockReleaseConfig) ChartDepsUpd() error {
	panic("not implemented") // TODO: Implement
}

func (r *mockReleaseConfig) In(a []release.Config) bool {
	panic("not implemented") // TODO: Implement
}

func (r *mockReleaseConfig) BuildValues(dir, templater string) error {
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

func (r *mockReleaseConfig) Uninstall() (*helmRelease.UninstallReleaseResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (r *mockReleaseConfig) Get() (*helmRelease.Release, error) {
	panic("not implemented") // TODO: Implement
}

func (r *mockReleaseConfig) List() (*helmRelease.Release, error) {
	panic("not implemented") // TODO: Implement
}

func (r *mockReleaseConfig) Rollback() error {
	panic("not implemented") // TODO: Implement
}

func (r *mockReleaseConfig) Status() (*helmRelease.Release, error) {
	panic("not implemented") // TODO: Implement
}

func (r *mockReleaseConfig) Name() string {
	return r.Called().String(0)
}

func (r *mockReleaseConfig) Namespace() string {
	return r.Called().String(0)
}

func (r *mockReleaseConfig) Chart() release.Chart {
	panic("not implemented") // TODO: Implement
}

func (r *mockReleaseConfig) DependsOn() []string {
	panic("not implemented") // TODO: Implement
}

func (r *mockReleaseConfig) Tags() []string {
	panic("not implemented") // TODO: Implement
}

func (r *mockReleaseConfig) Repo() string {
	return r.Called().String(0)
}

func (r *mockReleaseConfig) Values() []release.ValuesReference {
	return r.Called().Get(0).([]release.ValuesReference)
}

type mockRepoConfig struct {
	mock.Mock
}

func (r *mockRepoConfig) In(_ []repo.Config) bool {
	panic("not implemented") // TODO: Implement
}

func (r *mockRepoConfig) Install(_ *helm.EnvSettings, _ *helmRepo.File) error {
	panic("not implemented") // TODO: Implement
}

func (r *mockRepoConfig) Name() string {
	return r.Called().String(0)
}

func (r *mockRepoConfig) URL() string {
	return r.Called().String(0)
}
