package helper_test

import (
	"context"
	"testing"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
)

type ContextTestSuite struct {
	suite.Suite

	ctx context.Context
}

func TestContextTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ContextTestSuite))
}

func (ts *ContextTestSuite) SetupTest() {
	ts.ctx = tests.GetContext(ts.T())
}

func (ts *ContextTestSuite) TestReleaseUniq() {
	uniqBase, _ := uniqname.NewFromString("test")
	ctx := helper.ContextWithReleaseUniq(ts.ctx, uniqBase)

	uniq, exists := helper.ContextGetReleaseUniq(ctx)
	ts.Require().True(exists)
	ts.Require().Equal(uniqBase, uniq)
}

func (ts *ContextTestSuite) TestNoReleaseUniq() {
	_, exists := helper.ContextGetReleaseUniq(ts.ctx)
	ts.Require().False(exists)
}

func (ts *ContextTestSuite) TestLifecycleType() {
	typBase := "test"
	ctx := helper.ContextWithLifecycleType(ts.ctx, typBase)

	typ, exists := helper.ContextGetLifecycleType(ctx)
	ts.Require().True(exists)
	ts.Require().Equal(typBase, typ)
}

func (ts *ContextTestSuite) TestNoLifecycleType() {
	_, exists := helper.ContextGetLifecycleType(ts.ctx)
	ts.Require().False(exists)
}
