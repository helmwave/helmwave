package plan

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/fileref"

	"github.com/helmwave/helmwave/pkg/hooks"
	"github.com/helmwave/helmwave/pkg/monitor"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/helmwave/helmwave/pkg/template"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chartutil"
	helmRelease "helm.sh/helm/v3/pkg/release"
)

//nolint:govet // field alignment doesn't matter in tests
type MockReleaseConfig struct {
	mock.Mock

	t *testing.T
}

func NewMockReleaseConfig(t *testing.T) *MockReleaseConfig {
	t.Helper()

	c := &MockReleaseConfig{}
	c.Mock.Test(t)
	c.t = t

	return c
}

func (r *MockReleaseConfig) HideSecret(_ bool) {
	r.Called()
}

func (r *MockReleaseConfig) SetChartName(_ string) {
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

	u, _ := uniqname.New(r.Name(), r.Namespace(), r.KubeContext())

	return u
}

func (r *MockReleaseConfig) Sync(_ context.Context, _ bool) (*helmRelease.Release, error) {
	args := r.Called()

	return args.Get(0).(*helmRelease.Release), args.Error(1)
}

func (r *MockReleaseConfig) SyncDryRun(ctx context.Context, runHooks bool) (*helmRelease.Release, error) {
	r.DryRun(true)
	defer r.DryRun(false)

	return r.Sync(ctx, runHooks)
}

func (r *MockReleaseConfig) DryRun(bool) {
	r.Called()
}

func (r *MockReleaseConfig) ChartDepsUpd() error {
	return r.Called().Error(0)
}

func (r *MockReleaseConfig) Equal(_ release.Config) bool {
	return r.Called().Bool(0)
}

func (r *MockReleaseConfig) BuildValues(ctx context.Context, dir, templater string) error {
	args := r.Called()
	if errReturn := args.Error(0); errReturn != nil {
		return errReturn
	}

	for i := len(r.Values()) - 1; i >= 0; i-- {
		v := r.Values()[i]
		dst := filepath.Join(dir, Values, filepath.Base(v.Src))
		err := template.Tpl2yml(ctx, v.Src, dst, nil, templater)
		if err != nil {
			return err
		}
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

func (r *MockReleaseConfig) Values() []fileref.Config {
	return r.Called().Get(0).([]fileref.Config)
}

func (r *MockReleaseConfig) Logger() *log.Entry {
	return log.WithField("test", r.t.Name())
}

func (r *MockReleaseConfig) AllowFailure() bool {
	return r.Called().Bool(0)
}

func (r *MockReleaseConfig) HelmWait() bool {
	return true
}

func (r *MockReleaseConfig) DownloadChart(string) error {
	return r.Called().Error(0)
}

func (r *MockReleaseConfig) SetChart(_ string) {}

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

func (r *MockReleaseConfig) Lifecycle() hooks.Lifecycle {
	return r.Called().Get(0).(hooks.Lifecycle)
}
