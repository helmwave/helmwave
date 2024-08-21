package cache_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type NonParallelConfigTestSuite struct {
	suite.Suite
}

func TestNonParallelConfigTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(NonParallelConfigTestSuite))
}

// func (ts *NonParallelConfigTestSuite) TestInvalidCacheDir() {
// 	oldConfig := cache.Default
// 	cache.Default = cache.Config{
// 		Home:         "/proc/1/bla",
// 	}
// 	defer func() {
// 		cache.Default = oldConfig
// 	}()
//
// 	err := cache.Default.Init()
// 	ts.Require().Error(err)
// }
