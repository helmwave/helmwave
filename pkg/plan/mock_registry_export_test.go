package plan

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

type MockRegistryConfig struct {
	mock.Mock
}

func (r *MockRegistryConfig) Install() error {
	return r.Called().Error(0)
}

func (r *MockRegistryConfig) Host() string {
	return r.Called().String(0)
}

func (r *MockRegistryConfig) Logger() *logrus.Entry {
	return r.Called().Get(0).(*logrus.Entry)
}

func (r *MockRegistryConfig) Validate() error {
	return r.Called().Error(0)
}
