package dependency_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/release/dependency"
	"github.com/stretchr/testify/suite"
)

type NodeTestSuite struct {
	suite.Suite
}

func (s *NodeTestSuite) TestNewNode() {
	data := "123"
	node := dependency.NewNode(data)

	s.Require().IsType(data, node.Data)
}

func (s *NodeTestSuite) TestPending() {
	node := dependency.NewNode("")

	s.Require().False(node.IsDone())
}

func (s *NodeTestSuite) TestSucceeded() {
	node := dependency.NewNode("")

	s.Require().False(node.IsDone())
	s.Require().False(node.IsFailed())

	node.SetSucceeded()
	s.Require().True(node.IsDone())
	s.Require().False(node.IsFailed())
}

func (s *NodeTestSuite) TestFailed() {
	node := dependency.NewNode("")

	s.Require().False(node.IsDone())
	s.Require().False(node.IsFailed())

	node.SetFailed()
	s.Require().True(node.IsDone())
	s.Require().True(node.IsFailed())
}

func (s *NodeTestSuite) TestReadyWithDependencies() {
	node := dependency.NewNode("")
	nodeDep := dependency.NewNode("")

	s.Require().True(node.IsReady())

	node.AddDependency(nodeDep)

	s.Require().False(node.IsReady())

	nodeDep.SetSucceeded()

	s.Require().True(node.IsReady())

	nodeDep.SetFailed()

	s.Require().False(node.IsReady())
	s.Require().True(node.IsFailed())
}

func TestNodeTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(NodeTestSuite))
}
