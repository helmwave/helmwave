package plan

import (
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/helmwave/helmwave/pkg/template"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	helmRelease "helm.sh/helm/v3/pkg/release"
)

type mockReleaseConfig struct {
	mock.Mock

	NameF   string
	ValuesF []release.ValuesReference
}

func (r *mockReleaseConfig) Uniq() uniqname.UniqName {
	panic("not implemented") // TODO: Implement
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

func (r *mockReleaseConfig) BuildValues(dir string, gomplate *template.GomplateConfig) error {
	panic("not implemented") // TODO: Implement
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
	return r.NameF
}

func (r *mockReleaseConfig) Namespace() string {
	panic("not implemented") // TODO: Implement
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
	panic("not implemented") // TODO: Implement
}

func (r *mockReleaseConfig) Values() []release.ValuesReference {
	return r.ValuesF
}

type ExportTestSuite struct {
	suite.Suite
}

func (s *ExportTestSuite) TestValuesEmpty() {
	tmpDir := s.T().TempDir()
	p := New(filepath.Join(tmpDir, Dir))

	p.body = &planBody{
		Releases:     releaseConfigs{},
		Repositories: repoConfigs{},
	}

	err := p.exportValues()
	s.Require().NoError(err)
}

func (s *ExportTestSuite) TestValuesOneRelease() {
	tmpDir := s.T().TempDir()
	p := New(filepath.Join(tmpDir, Dir))

	p.body = &planBody{
		Releases: releaseConfigs{
			&mockReleaseConfig{
				NameF: "bitnami",
			},
			&mockReleaseConfig{
				NameF: "redis",
			},
		},
	}

	err := p.exportValues()
	s.Require().NoError(err)
}

func TestExportTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ExportTestSuite))
}
