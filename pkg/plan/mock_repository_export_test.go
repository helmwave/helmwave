package plan

import (
	"context"
	"testing"

	"github.com/helmwave/helmwave/pkg/repo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"helm.sh/helm/v3/pkg/cli"
	repo2 "helm.sh/helm/v3/pkg/repo"
)

type MockRepositoryConfig struct {
	mock.Mock
}

func NewMockRepositoryConfig(t *testing.T) *MockRepositoryConfig {
	t.Helper()

	c := &MockRepositoryConfig{}
	c.Mock.Test(t)

	return c
}

func (r *MockRepositoryConfig) Equal(_ repo.Config) bool {
	return r.Called().Bool(0)
}

func (r *MockRepositoryConfig) Install(context.Context, *cli.EnvSettings, *repo2.File) error {
	return r.Called().Error(0)
}

func (r *MockRepositoryConfig) Name() string {
	return r.Called().String(0)
}

func (r *MockRepositoryConfig) URL() string {
	return r.Called().String(0)
}

func (r *MockRepositoryConfig) Logger() *logrus.Entry {
	return r.Called().Get(0).(*logrus.Entry)
}

func (r *MockRepositoryConfig) Validate() error {
	return r.Called().Error(0)
}
