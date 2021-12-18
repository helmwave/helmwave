package release

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type DependencyTestSuite struct {
	suite.Suite
}

func (s *DependencyTestSuite) TestSingleRelease() {
	rel := &config{
		NameF:      "blabla" + strconv.Itoa(rand.Int()),
		NamespaceF: "blabla",
		DependsOnF: []string{},
	}
	rel.HandleDependencies([]Config{rel})

	s.Require().Empty(rel.dependencies)
}

func (s *DependencyTestSuite) TestSingleDependency() {
	rel2 := &config{
		NameF:      "blabla" + strconv.Itoa(rand.Int()),
		NamespaceF: "blabla",
	}

	rel1 := &config{
		NameF:      "blabla" + strconv.Itoa(rand.Int()),
		NamespaceF: "blabla",
		DependsOnF: []string{string(rel2.Uniq())},
	}

	releases := []Config{rel1, rel2}
	rel1.HandleDependencies(releases)
	rel2.HandleDependencies(releases)

	s.Require().NotEmpty(rel1.dependencies)
	s.Require().Contains(rel1.dependencies, rel2.Uniq())

	s.Require().Empty(rel2.dependencies)
}

func (s *DependencyTestSuite) TestMissingDependency() {
	rel := &config{
		NameF:      "blabla" + strconv.Itoa(rand.Int()),
		NamespaceF: "blabla",
		DependsOnF: []string{"blabla@blabla"},
	}

	releases := []Config{rel}
	rel.HandleDependencies(releases)

	s.Require().Empty(rel.dependencies)
}

func (s *DependencyTestSuite) TestWaitForDependenciesDryRun() {
	rel := &config{
		NameF:      "blabla" + strconv.Itoa(rand.Int()),
		NamespaceF: "blabla",
		DependsOnF: []string{},
	}
	rel.DryRun(true)

	rel.HandleDependencies([]Config{rel})

	s.Require().Eventually(func() bool {
		return rel.waitForDependencies() == nil
	}, 5*time.Second, time.Second)
}

func (s *DependencyTestSuite) TestHangWaitForDependencies() {
	rel2 := &config{
		NameF:      "blabla" + strconv.Itoa(rand.Int()),
		NamespaceF: "blabla",
	}

	relHang := &config{
		NameF:      "blabla" + strconv.Itoa(rand.Int()),
		NamespaceF: "blabla",
		DependsOnF: []string{string(rel2.Uniq())},
	}

	releases := []Config{relHang, rel2}
	relHang.HandleDependencies(releases)
	rel2.HandleDependencies(releases)

	s.Require().Never(func() bool {
		return relHang.waitForDependencies() == nil
	}, 5*time.Second, time.Second)
}

func (s *DependencyTestSuite) TestDependencyFailed() {
	rel2 := &config{
		NameF:      "blabla" + strconv.Itoa(rand.Int()),
		NamespaceF: "blabla",
	}

	rel1 := &config{
		NameF:      "blabla" + strconv.Itoa(rand.Int()),
		NamespaceF: "blabla",
		DependsOnF: []string{string(rel2.Uniq())},
	}

	releases := []Config{rel1, rel2}
	rel1.HandleDependencies(releases)
	rel2.HandleDependencies(releases)

	rel2.NotifyFailed()

	s.Require().ErrorIs(rel1.waitForDependencies(), ErrDepFailed)
}

func (s *DependencyTestSuite) TestDependencyAllowedToFail() {
	rel2 := &config{
		NameF:        "blabla" + strconv.Itoa(rand.Int()),
		NamespaceF:   "blabla",
		AllowFailure: true,
	}

	rel1 := &config{
		NameF:      "blabla" + strconv.Itoa(rand.Int()),
		NamespaceF: "blabla",
		DependsOnF: []string{string(rel2.Uniq())},
	}

	releases := []Config{rel1, rel2}
	rel1.HandleDependencies(releases)
	rel2.HandleDependencies(releases)

	rel2.NotifyFailed()

	s.Require().NoError(rel1.waitForDependencies())
}

func (s *DependencyTestSuite) TestDependencySucceed() {
	rel2 := &config{
		NameF:      "blabla" + strconv.Itoa(rand.Int()),
		NamespaceF: "blabla",
	}

	rel1 := &config{
		NameF:      "blabla" + strconv.Itoa(rand.Int()),
		NamespaceF: "blabla",
		DependsOnF: []string{string(rel2.Uniq())},
	}

	releases := []Config{rel1, rel2}
	rel1.HandleDependencies(releases)
	rel2.HandleDependencies(releases)

	rel2.NotifySuccess()

	s.Require().NoError(rel1.waitForDependencies())
}

func TestDependencyTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(DependencyTestSuite))
}
