package plan

import (
	"context"
	"io/fs"
	"path/filepath"

	"github.com/helmwave/go-fsimpl"
	"github.com/helmwave/helmwave/pkg/monitor"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/helmwave/helmwave/pkg/template"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chartutil"
	helmRelease "helm.sh/helm/v3/pkg/release"
)

type MockReleaseConfig struct {
	mock.Mock
}

func (r *MockReleaseConfig) SetChartName(s string) {
	r.Called()
}

func (r *MockReleaseConfig) OfflineKubeVersion() *chartutil.KubeVersion {
	r.Called()

	v := &chartutil.KubeVersion{
		Major:   "1",
		Minor:   "22",
		Version: "1.22.0",
	}

	return v
}

func (r *MockReleaseConfig) Uniq() uniqname.UniqName {
	args := r.Called()

	if len(args) > 0 {
		return args.Get(0).(uniqname.UniqName)
	}

	u, _ := uniqname.Generate(r.Name(), r.Namespace())

	return u
}

func (r *MockReleaseConfig) Sync(_ context.Context, _ fs.FS) (*helmRelease.Release, error) {
	args := r.Called()

	return args.Get(0).(*helmRelease.Release), args.Error(1)
}

func (r *MockReleaseConfig) SyncDryRun(ctx context.Context, baseFS fs.FS) (*helmRelease.Release, error) {
	r.DryRun(true)
	defer r.DryRun(false)

	return r.Sync(ctx, baseFS)
}

func (r *MockReleaseConfig) DryRun(bool) {
	r.Called()
}

func (r *MockReleaseConfig) ChartDepsUpd(_ fs.StatFS) error {
	return r.Called().Error(0)
}

func (r *MockReleaseConfig) Equal(release.Config) bool {
	return r.Called().Bool(0)
}

func (r *MockReleaseConfig) ExportValues(statFS fs.StatFS, writeableFS fsimpl.WriteableFS, templater string) error {
	args := r.Called()
	if errReturn := args.Error(0); errReturn != nil {
		return errReturn
	}

	for i := len(r.Values()) - 1; i >= 0; i-- {
		v := r.Values()[i]
		dst := filepath.Join(Values, v.Src)
		err := template.Tpl2yml(statFS, writeableFS, v.Src, dst, nil, templater)
		if err != nil {
			return err
		}

		v.Src = dst
	}

	return nil
}

func (r *MockReleaseConfig) Uninstall(context.Context) (*helmRelease.UninstallReleaseResponse, error) {
	args := r.Called()

	return args.Get(0).(*helmRelease.UninstallReleaseResponse), args.Error(1)
}

func (r *MockReleaseConfig) Get(version int) (*helmRelease.Release, error) {
	args := r.Called(version)

	return args.Get(0).(*helmRelease.Release), args.Error(1)
}

func (r *MockReleaseConfig) List() (*helmRelease.Release, error) {
	args := r.Called()

	return args.Get(0).(*helmRelease.Release), args.Error(1)
}

func (r *MockReleaseConfig) Rollback(context.Context, int) error {
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

func (r *MockReleaseConfig) Chart() *release.Chart {
	return r.Called().Get(0).(*release.Chart)
}

func (r *MockReleaseConfig) DependsOn() []*release.DependsOnReference {
	return r.Called().Get(0).([]*release.DependsOnReference)
}

func (r *MockReleaseConfig) SetDependsOn(deps []*release.DependsOnReference) {
	r.Called(deps)
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

func (r *MockReleaseConfig) Logger() *logrus.Entry {
	return r.Called().Get(0).(*logrus.Entry)
}

func (r *MockReleaseConfig) AllowFailure() bool {
	return r.Called().Bool(0)
}

func (r *MockReleaseConfig) HelmWait() bool {
	return true
}

func (r *MockReleaseConfig) DownloadChart(_ fs.StatFS, _ fsimpl.WriteableFS, _ string) error {
	return r.Called().Error(0)
}

func (r *MockReleaseConfig) SetChart(string) {}

func (r *MockReleaseConfig) KubeContext() string {
	return r.Called().String(0)
}

func (r *MockReleaseConfig) Cfg() *action.Configuration {
	return r.Called().Get(0).(*action.Configuration)
}

func (r *MockReleaseConfig) HooksDisabled() bool {
	return r.Called().Bool(0)
}

func (r *MockReleaseConfig) Validate() error {
	return r.Called().Error(0)
}

func (r *MockReleaseConfig) Monitors() []release.MonitorReference {
	return r.Called().Get(0).([]release.MonitorReference)
}

func (r *MockReleaseConfig) NotifyMonitorsFailed(context.Context, ...monitor.Config) {
	r.Called()
}
