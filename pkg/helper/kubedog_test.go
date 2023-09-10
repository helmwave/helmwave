package helper_test

import (
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
	"github.com/werf/kubedog/pkg/kube"
)

type KubedogTestSuite struct {
	suite.Suite
}

func (s *KubedogTestSuite) TestKubeInit() {
	s.T().Setenv("KUBECONFIG", filepath.Join(tests.Root, "kubeconfig.yaml"))
	err := helper.KubeInit("")

	s.Require().NoError(err)
	s.Require().NotNil(kube.Client)
}

//nolint:paralleltest // can't parallel because of setenv
func TestKubedogTestSuite(t *testing.T) {
	// t.Parallel()
	suite.Run(t, new(KubedogTestSuite))
}
