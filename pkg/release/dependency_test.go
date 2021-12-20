package release_test

import (
	"testing"
	"time"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/stretchr/testify/suite"
)

type DependencyTestSuite struct {
	suite.Suite
}

func (s *DependencyTestSuite) TestSingleRelease() {
	rel := release.NewConfig()
	rel.HandleDependencies([]release.Config{rel})

	s.Require().Empty(rel.GetDependencies())
}

func (s *DependencyTestSuite) TestSingleDependency() {
	rel2 := release.NewConfig()

	rel1 := release.NewConfig()
	rel1.DependsOnF = []string{string(rel2.Uniq())}

	releases := []release.Config{rel1, rel2}
	rel1.HandleDependencies(releases)
	rel2.HandleDependencies(releases)

	s.Require().NotEmpty(rel1.GetDependencies())
	s.Require().Contains(rel1.GetDependencies(), rel2.Uniq())

	s.Require().Empty(rel2.GetDependencies())
}

func (s *DependencyTestSuite) TestMissingDependency() {
	rel := release.NewConfig()
	rel.DependsOnF = []string{"blabla@blabla"}

	releases := []release.Config{rel}
	rel.HandleDependencies(releases)

	s.Require().Empty(rel.GetDependencies())
}

func (s *DependencyTestSuite) TestWaitForDependenciesDryRun() {
	rel := release.NewConfig()
	rel.DryRun(true)

	rel.HandleDependencies([]release.Config{rel})

	s.Require().Eventually(func() bool {
		return rel.WaitForDependencies() == nil
	}, 5*time.Second, time.Second)
}

func (s *DependencyTestSuite) TestHangWaitForDependencies() {
	rel2 := release.NewConfig()

	relHang := release.NewConfig()
	relHang.DependsOnF = []string{string(rel2.Uniq())}

	releases := []release.Config{relHang, rel2}
	relHang.HandleDependencies(releases)
	rel2.HandleDependencies(releases)

	s.Require().Never(func() bool {
		return relHang.WaitForDependencies() == nil
	}, 5*time.Second, time.Second)
}

func (s *DependencyTestSuite) TestDependencyFailed() {
	rel2 := release.NewConfig()

	rel1 := release.NewConfig()
	rel1.DependsOnF = []string{string(rel2.Uniq())}

	releases := []release.Config{rel1, rel2}
	rel1.HandleDependencies(releases)
	rel2.HandleDependencies(releases)

	rel2.NotifyFailed()

	s.Require().ErrorIs(rel1.WaitForDependencies(), release.ErrDepFailed)
}

func (s *DependencyTestSuite) TestDependencyAllowedToFail() {
	rel2 := release.NewConfig()
	rel2.AllowFailure = true

	rel1 := release.NewConfig()
	rel1.DependsOnF = []string{string(rel2.Uniq())}

	releases := []release.Config{rel1, rel2}
	rel1.HandleDependencies(releases)
	rel2.HandleDependencies(releases)

	rel2.NotifyFailed()

	s.Require().NoError(rel1.WaitForDependencies())
}

func (s *DependencyTestSuite) TestDependencySucceed() {
	rel2 := release.NewConfig()

	rel1 := release.NewConfig()
	rel1.DependsOnF = []string{string(rel2.Uniq())}

	releases := []release.Config{rel1, rel2}
	rel1.HandleDependencies(releases)
	rel2.HandleDependencies(releases)

	rel2.NotifySuccess()

	s.Require().NoError(rel1.WaitForDependencies())
}

func TestDependencyTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(DependencyTestSuite))
}
