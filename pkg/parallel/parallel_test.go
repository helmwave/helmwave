//go:build ignore || unit

package parallel

import (
	"errors"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/suite"
)

type ParallelTestSuite struct {
	suite.Suite
	wg *WaitGroup
}

func (s *ParallelTestSuite) SetupTest() {
	s.wg = NewWaitGroup()
}

func (s *ParallelTestSuite) TestErrors() {
	ch := s.wg.ErrChan()
	s.Require().NotNil(ch)

	err := errors.New("blabla")
	s.wg.Add(1)
	go func(wg *WaitGroup, err error) {
		wg.ErrChan() <- err
		wg.Done()
	}(s.wg, err)

	e := s.wg.Wait()
	s.Require().Error(e)
	s.Require().IsType(&multierror.Error{}, e)
	s.Equal(err, e.(*multierror.Error).Unwrap())
}

func TestParallelTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ParallelTestSuite))
}
