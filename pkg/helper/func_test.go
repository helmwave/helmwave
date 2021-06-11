// +build ignore unit

package helper

import (
	"os"
	"path"
	"sort"
	"strconv"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type HelperTestSuite struct {
	suite.Suite
}

func (s *HelperTestSuite) TestContains() {
	tests := []struct {
		t   string
		a   []string
		out bool
	}{
		{
			t:   "",
			a:   []string{},
			out: false,
		},
		{
			t:   "123",
			a:   []string{},
			out: false,
		},
		{
			t:   "321",
			a:   []string{"123", "321"},
			out: true,
		},
	}

	for _, t := range tests {
		sort.Strings(t.a)
		s.Equalf(t.out, Contains(t.t, t.a), "checking %v in %v", t.t, t.a)
	}
}

func (s *HelperTestSuite) TestCreateFile() {
	tmpPath := path.Join(os.TempDir(), strconv.Itoa(time.Now().Second()))
	defer func(name string) {
		_ = os.Remove(name)
	}(tmpPath)

	f, err := CreateFile(tmpPath)

	s.Require().FileExists(tmpPath)
	s.Require().NoError(err)
	s.Require().NotNil(f)

	err = f.Close()
	s.Require().NoError(err)

	tmpPath = path.Join(tmpPath, "123")
	_, err = CreateFile(tmpPath)

	s.Error(err)
	s.IsType(&os.PathError{}, err)
	s.Equal(syscall.ENOTDIR, err.(*os.PathError).Err)
}

func TestHelperTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(HelperTestSuite))
}
