package cache_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/helmwave/go-fsimpl"
	"github.com/helmwave/go-fsimpl/filefs"
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

func (ts *ChartsTestSuite) TestNonWriteableFS() {
	cfg := cache.Config{}

	cacheFS := os.DirFS(ts.T().TempDir())
	err := cfg.Init(cacheFS, ".")

	ts.Require().ErrorIs(err, cache.ErrNotWriteableFS)
}

func (ts *ChartsTestSuite) TestEnabled() {
	cfg := cache.Config{}

	ts.Require().False(cfg.IsEnabled())

	_, err := cfg.FindInCache("", "")
	ts.Require().ErrorIs(err, cache.ErrCacheDisabled)

	mux := fsimpl.NewMux()
	mux.Add(filefs.FS)
	cacheFS, err := mux.Lookup(fmt.Sprintf("file://%s", ts.T().TempDir()))

	ts.Require().NoError(err)

	err = cfg.Init(cacheFS, ".")

	ts.Require().NoError(err)
	ts.Require().True(cfg.IsEnabled())
}

func (ts *ChartsTestSuite) TestFindNonexisting() {
	cfg := cache.Config{}
	mux := fsimpl.NewMux()
	mux.Add(filefs.FS)
	cacheFS, err := mux.Lookup(fmt.Sprintf("file://%s", ts.T().TempDir()))

	ts.Require().NoError(err)
	err = cfg.Init(cacheFS, ".")

	ts.Require().NoError(err)

	_, err = cfg.FindInCache("bla", "bla")

	ts.Require().ErrorIs(err, cache.ErrChartNotFound)
}
