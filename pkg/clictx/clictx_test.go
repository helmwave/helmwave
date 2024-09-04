package clictx_test

import (
	"context"
	"testing"

	"github.com/helmwave/helmwave/pkg/clictx"
	"github.com/stretchr/testify/suite"
)

type CliCtxTestSuite struct {
	suite.Suite
}

func TestCliCtxTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(CliCtxTestSuite))
}

func (ts *CliCtxTestSuite) TestAddFlagToContext() {
	ctx := context.Background()
	ctx = clictx.AddFlagToContext(ctx, "testFlag", "testValue")

	value := clictx.GetFlagFromContext(ctx, "testFlag")
	ts.Require().Equal("testValue", value)
}

// func (ts *CliCtxTestSuite) TestAddAndRetrieveCLIContext(t *testing.T) {
//	ts.Require().Error(nil)
//}
//
// func (ts *CliCtxTestSuite) TestCLIContextToContext(t *testing.T) {
//	ts.Require().Error(nil)
//}
