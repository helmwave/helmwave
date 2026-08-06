package helper_test

import (
	"context"
	"testing"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/stretchr/testify/suite"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

type HelmTestSuite struct {
	suite.Suite
}

func (s *HelmTestSuite) TestNewCfg() {
	ns := s.T().Name()
	cfg, err := helper.NewCfg(ns, "")

	s.Require().NoError(err)
	s.Require().NotNil(cfg)
}

func (s *HelmTestSuite) TestNewHelmNS() {
	ns := s.T().Name()
	h1 := helper.NewHelm(ns)

	s.Require().NotNil(h1)
	s.Require().Equal(ns, h1.Namespace())
}

func (s *HelmTestSuite) TestNamespaceExists() {
	ns := "existing-ns"
	clientSet := fake.NewSimpleClientset(&corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: ns},
	})

	exists, err := helper.NamespaceExists(context.Background(), clientSet, ns)
	s.Require().NoError(err)
	s.Require().True(exists)
}

func (s *HelmTestSuite) TestNamespaceExistsMissing() {
	clientSet := fake.NewSimpleClientset()

	exists, err := helper.NamespaceExists(context.Background(), clientSet, "missing-ns")
	s.Require().NoError(err)
	s.Require().False(exists)
}

func TestHelmTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(HelmTestSuite))
}
