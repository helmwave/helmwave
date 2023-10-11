package dependency_test

import (
	"testing"
	"time"

	"github.com/helmwave/helmwave/pkg/release/dependency"
	"github.com/stretchr/testify/suite"
)

type GraphTestSuite struct {
	suite.Suite
}

func (s *GraphTestSuite) TestNewNode() {
	key := "1"
	data := "123"

	graph := dependency.NewGraph[string, string]()
	s.Require().NoError(graph.NewNode(key, data))

	s.Require().Len(graph.Nodes, 1)
	s.Require().Contains(graph.Nodes, key)
	s.Require().Equal(graph.Nodes[key].Data, data)
	s.Require().False(graph.Nodes[key].IsDone())
	s.Require().True(graph.Nodes[key].IsReady())

	s.Require().NoError(graph.Build())
}

func (s *GraphTestSuite) TestDuplicateNode() {
	key := "1"

	graph := dependency.NewGraph[string, string]()

	s.Require().NoError(graph.NewNode(key, "123"))
	s.Require().Error(graph.NewNode(key, "321"))
	s.Require().Len(graph.Nodes, 1)

	s.Require().NoError(graph.Build())
}

// 1 -> 3 -> 4
// 2 /
func (s *GraphTestSuite) TestDependencies() {
	graph := dependency.NewGraph[string, string]()

	s.Require().NoError(graph.NewNode("1", "1"))
	s.Require().NoError(graph.NewNode("2", "2"))
	s.Require().NoError(graph.NewNode("3", "3"))
	s.Require().NoError(graph.NewNode("4", "4"))

	graph.AddDependency("3", "1")
	graph.AddDependency("3", "2")
	graph.AddDependency("4", "3")

	s.Require().NoError(graph.Build())

	ch := graph.Run()

	s.Require().NotNil(ch)
	s.Require().Eventually(func() bool { return len(ch) == 2 }, time.Second, time.Millisecond)
	s.Require().Never(func() bool { return len(ch) > 2 }, time.Second, time.Millisecond)
	d1, d2 := <-ch, <-ch
	s.Require().ElementsMatch([]string{"1", "2"}, []string{d1.Data, d2.Data})
	d1.SetSucceeded()
	d2.SetSucceeded()

	s.Require().Eventually(func() bool { return len(ch) == 1 }, time.Second, time.Millisecond)
	s.Require().Never(func() bool { return len(ch) > 1 }, time.Second, time.Millisecond)
	d := <-ch
	s.Require().Equal("3", d.Data)
	d.SetSucceeded()

	s.Require().Eventually(func() bool { return len(ch) == 1 }, time.Second, time.Millisecond)
	s.Require().Never(func() bool { return len(ch) > 1 }, time.Second, time.Millisecond)
	d = <-ch
	s.Require().Equal("4", d.Data)
	d.SetSucceeded()

	s.Require().Never(func() bool { return len(ch) > 0 }, time.Second, time.Millisecond)
}

func (s *GraphTestSuite) TestFailedDependencies() {
	graph := dependency.NewGraph[string, string]()

	s.Require().NoError(graph.NewNode("1", "1"))
	s.Require().NoError(graph.NewNode("2", "2"))
	s.Require().NoError(graph.NewNode("3", "3"))
	s.Require().NoError(graph.NewNode("4", "4"))

	graph.AddDependency("3", "1")
	graph.AddDependency("3", "2")
	graph.AddDependency("4", "3")

	s.Require().NoError(graph.Build())

	ch := graph.Run()

	s.Require().NotNil(ch)
	s.Require().Eventually(func() bool { return len(ch) == 2 }, time.Second, time.Millisecond)
	s.Require().Never(func() bool { return len(ch) > 2 }, time.Second, time.Millisecond)
	d1, d2 := <-ch, <-ch
	s.Require().ElementsMatch([]string{"1", "2"}, []string{d1.Data, d2.Data})
	d1.SetSucceeded()
	d2.SetFailed()

	s.Require().Never(func() bool { return len(ch) > 0 }, time.Second, time.Millisecond)
}

func (s *GraphTestSuite) TestCycles() {
	graph := dependency.NewGraph[string, string]()

	s.Require().NoError(graph.NewNode("1", "1"))
	s.Require().NoError(graph.NewNode("2", "2"))

	graph.AddDependency("1", "2")
	graph.AddDependency("2", "1")

	s.Require().Error(graph.Build())
}

// 1 -> 3 -> 4
// 2 /
// Reversed
// 4 -> 3 -> 1
// --------> 2
func (s *GraphTestSuite) TestReverse() {
	graph := dependency.NewGraph[string, string]()

	s.Require().NoError(graph.NewNode("1", "1"))
	s.Require().NoError(graph.NewNode("2", "2"))
	s.Require().NoError(graph.NewNode("3", "3"))
	s.Require().NoError(graph.NewNode("4", "4"))

	graph.AddDependency("3", "1")
	graph.AddDependency("3", "2")
	graph.AddDependency("4", "3")

	s.Require().NoError(graph.Build())
	graph, err := graph.Reverse()
	s.Require().Nil(err)

	ch := graph.Run()

	s.Require().NotNil(ch)
	s.Require().Eventually(func() bool { return len(ch) == 1 }, time.Second, time.Millisecond)
	s.Require().Never(func() bool { return len(ch) > 1 }, time.Second, time.Millisecond)
	d := <-ch
	s.Require().Equal("4", d.Data)
	d.SetSucceeded()

	s.Require().Eventually(func() bool { return len(ch) == 1 }, time.Second, time.Millisecond)
	s.Require().Never(func() bool { return len(ch) > 1 }, time.Second, time.Millisecond)
	d = <-ch
	s.Require().Equal("3", d.Data)
	d.SetSucceeded()

	s.Require().Eventually(func() bool { return len(ch) == 2 }, time.Second, time.Millisecond)
	s.Require().Never(func() bool { return len(ch) > 2 }, time.Second, time.Millisecond)
	d1, d2 := <-ch, <-ch
	s.Require().ElementsMatch([]string{"1", "2"}, []string{d1.Data, d2.Data})
	d1.SetSucceeded()
	d2.SetSucceeded()

	s.Require().Never(func() bool { return len(ch) > 0 }, time.Second, time.Millisecond)
}

func (s *GraphTestSuite) TestReverseFailedDependencies() {
	graph := dependency.NewGraph[string, string]()

	s.Require().NoError(graph.NewNode("1", "1"))
	s.Require().NoError(graph.NewNode("2", "2"))
	s.Require().NoError(graph.NewNode("3", "3"))
	s.Require().NoError(graph.NewNode("4", "4"))

	graph.AddDependency("3", "1")
	graph.AddDependency("3", "2")
	graph.AddDependency("4", "3")

	s.Require().NoError(graph.Build())
	graph, err := graph.Reverse()
	s.Require().Nil(err)

	ch := graph.Run()

	s.Require().NotNil(ch)
	s.Require().Eventually(func() bool { return len(ch) == 1 }, time.Second, time.Millisecond)
	s.Require().Never(func() bool { return len(ch) > 1 }, time.Second, time.Millisecond)
	d := <-ch
	s.Require().Equal("4", d.Data)
	d.SetSucceeded()

	s.Require().Eventually(func() bool { return len(ch) == 1 }, time.Second, time.Millisecond)
	s.Require().Never(func() bool { return len(ch) > 1 }, time.Second, time.Millisecond)
	d = <-ch
	s.Require().Equal("3", d.Data)
	d.SetFailed()

	s.Require().Never(func() bool { return len(ch) > 0 }, time.Second, time.Millisecond)
}

func TestGraphTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(GraphTestSuite))
}
