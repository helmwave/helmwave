package cache_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/cache"
	"github.com/stretchr/testify/suite"
)

type ChartsTestSuite struct {
	suite.Suite
}

func TestChartsTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ChartsTestSuite))
}

func (ts *ChartsTestSuite) TestEnabled() {
	cfg := cache.Config{}

	ts.Require().False(cfg.IsEnabled())

	_, err := cfg.FindInCache("", "")
	ts.Require().ErrorIs(err, cache.ErrCacheDisabled)

	err = cfg.Init(ts.T().TempDir())

	ts.Require().NoError(err)
	ts.Require().True(cfg.IsEnabled())
}

func (ts *ChartsTestSuite) TestFindNonexisting() {
	cfg := cache.Config{}
	err := cfg.Init(ts.T().TempDir())

	ts.Require().NoError(err)

	_, err = cfg.FindInCache("bla", "bla")

	ts.Require().ErrorIs(err, cache.ErrChartNotFound)
}
