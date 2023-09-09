package kubedog_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/kubedog"
	"github.com/stretchr/testify/suite"
)

type ManifestTestSuite struct {
	suite.Suite
}

func TestManifestTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ManifestTestSuite))
}

func (ts *ManifestTestSuite) TestEmptyYAML() {
	res := kubedog.Parse([]byte{})
	ts.Require().Empty(res)
}

func (ts *ManifestTestSuite) TestInvalidYAML() {
	res := kubedog.Parse([]byte("a: {[]}"))
	ts.Require().Empty(res)

	res = kubedog.Parse([]byte("----"))
	ts.Require().Empty(res)
}

func (ts *ManifestTestSuite) TestEmptyYAMLDocument() {
}
