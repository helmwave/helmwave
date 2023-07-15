package helper_test

import (
	"context"
	"testing"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/stretchr/testify/suite"
)

type ContextTestSuite struct {
	suite.Suite
}

func TestContextTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ContextTestSuite))
}

func (s *ContextTestSuite) TestReleaseUniq() {
	ctxBase := context.Background()
	uniqBase := uniqname.UniqName("test")

	ctx := helper.ContextWithReleaseUniq(ctxBase, uniqBase)

	uniq, exists := helper.ContextGetReleaseUniq(ctx)
	s.Require().True(exists)
	s.Require().Equal(uniqBase, uniq)
}

func (s *ContextTestSuite) TestNoReleaseUniq() {
	ctx := context.Background()

	_, exists := helper.ContextGetReleaseUniq(ctx)
	s.Require().False(exists)
}

func (s *ContextTestSuite) TestLifecycleType() {
	ctxBase := context.Background()
	typBase := "test"

	ctx := helper.ContextWithLifecycleType(ctxBase, typBase)

	typ, exists := helper.ContextGetLifecycleType(ctx)
	s.Require().True(exists)
	s.Require().Equal(typBase, typ)
}

func (s *ContextTestSuite) TestNoLifecycleType() {
	ctx := context.Background()

	_, exists := helper.ContextGetLifecycleType(ctx)
	s.Require().False(exists)
}
