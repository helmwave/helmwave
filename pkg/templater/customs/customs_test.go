package customs_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/helmwave/helmwave/pkg/templater/customs"
	"github.com/stretchr/testify/suite"
)

type ExtraTestSuite struct {
	suite.Suite
}

func (s *ExtraTestSuite) TestToYaml() {
	data := struct {
		Field       any
		nonexported any
	}{
		Field:       "field",
		nonexported: 123,
	}
	yamlData := "field: field"

	y, err := customs.ToYaml(data)
	s.Require().NoError(err)
	s.Require().YAMLEq(yamlData, y)
}

type raw struct{}

func (r raw) MarshalYAML() (any, error) {
	return nil, os.ErrNotExist
}

func (s *ExtraTestSuite) TestToYamlNil() {
	data := raw{}
	y, err := customs.ToYaml(data)
	s.Require().Equal("", y)
	s.Require().ErrorIs(err, os.ErrNotExist)
}

func (s *ExtraTestSuite) TestFromYaml() {
	tests := []struct {
		result customs.Values
		yaml   string
		fails  bool
	}{
		{
			yaml:   "abc: 123",
			result: customs.Values{"abc": 123},
			fails:  false,
		},
		{
			yaml:  "123!!123",
			fails: true,
		},
	}

	for i := range tests {
		v, err := customs.FromYaml(tests[i].yaml)
		if tests[i].fails {
			s.Error(err)
			s.Empty(v)
		} else {
			s.NoError(err)
			s.Equal(tests[i].result, v)
		}
	}
}

func (s *ExtraTestSuite) TestExec() {
	res, err := customs.Exec("pwd", []any{})
	s.Require().NoError(err)

	pwd, err := os.Getwd()
	s.Require().NoError(err)

	s.Require().Equal(pwd, strings.TrimSpace(res))

	res, err = customs.Exec("echo", []any{"-n", "123"})
	s.Require().NoError(err)
	s.Require().Equal("123", res)
}

func (s *ExtraTestSuite) TestExecInvalidArg() {
	res, err := customs.Exec("pwd", []any{123})
	s.Require().Error(err)
	s.Require().Empty(res)
}

func (s *ExtraTestSuite) TestExecError() {
	res, err := customs.Exec(s.T().Name(), []any{})
	s.Require().Error(err)
	s.Require().Empty(res)
}

func (s *ExtraTestSuite) TestExecStdin() {
	input := "123"
	res, err := customs.Exec("cat", []any{}, input)
	s.Require().NoError(err)
	s.Require().Equal(input, res)
}

func (s *ExtraTestSuite) TestSetValueAtPath() {
	data := customs.Values{
		"a": map[string]any{
			"b": "123",
		},
		"c": 123,
		"d": map[any]any{
			"e": "f",
		},
	}

	tests := []struct {
		result customs.Values
		value  any
		path   string
		fails  bool
	}{
		{
			path:  "c",
			value: 321,
			result: customs.Values{
				"a": map[string]any{"b": "123"},
				"c": 321,
				"d": map[any]any{
					"e": "f",
				},
			},
			fails: false,
		},
		{
			path:  "a.b",
			value: "321",
			result: customs.Values{
				"a": map[string]any{"b": "321"},
				"c": 321,
				"d": map[any]any{
					"e": "f",
				},
			},
			fails: false,
		},
		{
			path:  "a.c",
			value: "321",
			result: customs.Values{
				"a": map[string]any{"b": "321", "c": "321"},
				"c": 321,
				"d": map[any]any{
					"e": "f",
				},
			},
			fails: false,
		},
		{
			path:   "c.a",
			value:  "321",
			result: nil,
			fails:  true,
		},
		{
			path:  "d.e",
			value: "321",
			result: customs.Values{
				"a": map[string]any{"b": "321", "c": "321"},
				"c": 321,
				"d": map[any]any{
					"e": "321",
				},
			},
			fails: false,
		},
	}

	for i := range tests {
		res, err := customs.SetValueAtPath(tests[i].path, tests[i].value, data)
		if tests[i].fails {
			s.Error(err)
			s.Nil(res)
		} else {
			s.NoError(err)
			s.Equal(tests[i].result, res)
		}
	}
}

func (s *ExtraTestSuite) TestRequired() {
	tests := []struct {
		data  any
		fails bool
	}{
		{
			data:  nil,
			fails: true,
		},
		{
			data:  4,
			fails: false,
		},
		{
			data:  "",
			fails: true,
		},
		{
			data:  "123",
			fails: false,
		},
	}

	for _, t := range tests {
		res, err := customs.Required("blabla", t.data)
		if t.fails {
			s.Error(err)
			s.Nil(res)
		} else {
			s.NoError(err)
			s.Equal(t.data, res)
		}
	}
}

func (s *ExtraTestSuite) TestReadFile() {
	tmpDir := s.T().TempDir()
	tmpFile := filepath.Join(tmpDir, "blablafile")

	res, err := customs.ReadFile(tmpFile)

	s.Require().Equal("", res)
	s.Require().ErrorIs(err, os.ErrNotExist)

	data := s.T().Name()

	s.Require().NoError(os.WriteFile(tmpFile, []byte(data), 0o600))
	s.Require().FileExists(tmpFile)

	res, err = customs.ReadFile(tmpFile)

	s.Require().NoError(err)
	s.Require().Equal(data, res)
}

func (s *ExtraTestSuite) TestGet() {
	data := customs.Values{
		"a": map[string]any{
			"b": "123",
		},
		"c": 123,
		"d": map[any]any{
			"e": "f",
		},
	}

	tests := []struct {
		result any
		path   string
		fails  bool
	}{
		{
			path:   "c",
			result: 123,
			fails:  false,
		},
		{
			path:   "a.b",
			result: "123",
			fails:  false,
		},
		{
			path:   "a.c",
			result: nil,
			fails:  true,
		},
		{
			path:   "c.a",
			result: nil,
			fails:  true,
		},
		{
			path:   "d.e",
			result: "f",
			fails:  false,
		},
	}

	for i := range tests {
		res, err := customs.Get(tests[i].path, data)
		if tests[i].fails {
			s.Error(err)
			s.Nil(res)
		} else {
			s.NoError(err)
			s.Equal(tests[i].result, res)
		}
	}
}

func (s *ExtraTestSuite) TestHasKey() {
	data := customs.Values{
		"a": map[string]any{
			"b": "123",
		},
		"c": 123,
		"d": map[any]any{
			"e": "f",
		},
	}

	tests := []struct {
		path   string
		result bool
		fails  bool
	}{
		{
			path:   "c",
			result: true,
			fails:  false,
		},
		{
			path:   "a.b",
			result: true,
			fails:  false,
		},
		{
			path:   "a.c",
			result: false,
			fails:  false,
		},
		{
			path:   "c.a",
			result: false,
			fails:  true,
		},
		{
			path:   "d.e",
			result: true,
			fails:  false,
		},
	}

	for _, test := range tests {
		res, err := customs.HasKey(test.path, data)
		s.Equal(test.result, res)

		if test.fails {
			s.Error(err)
		} else {
			s.NoError(err)
		}
	}
}

func TestExtraTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ExtraTestSuite))
}

type NonParallelExtraTestSuite struct {
	suite.Suite
}

func (s *NonParallelExtraTestSuite) TestRequiredEnv() {
	name := s.T().Name()

	res, err := customs.RequiredEnv(name)
	s.Require().Error(err)
	s.Require().Empty(res)

	data := "test"
	s.T().Setenv(name, data)

	res, err = customs.RequiredEnv(name)
	s.Require().NoError(err)
	s.Require().Equal(data, res)
}

//nolint:paralleltest // can't parallel because of setenv
func TestNonParallelExtraTestSuite(t *testing.T) {
	// t.Parallel()
	suite.Run(t, new(NonParallelExtraTestSuite))
}
