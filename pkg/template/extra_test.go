//go:build ignore || unit

package template

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ExtraTestSuite struct {
	suite.Suite
}

func (s *ExtraTestSuite) TestToYaml() {
	data := struct {
		Field       interface{}
		nonexported interface{}
	}{
		Field:       "field",
		nonexported: 123,
	}
	yamlData := "field: field"

	y, err := ToYaml(data)
	s.Require().NoError(err)
	s.Require().YAMLEq(yamlData, y)
}

type raw struct{}

func (r raw) MarshalYAML() (interface{}, error) {
	return nil, os.ErrNotExist
}

func (s *ExtraTestSuite) TestToYamlNil() {
	data := raw{}
	y, err := ToYaml(data)
	s.Require().Equal("", y)
	s.Require().ErrorIs(err, os.ErrNotExist)
}

func (s *ExtraTestSuite) TestFromYaml() {
	tests := []struct {
		yaml   string
		result Values
		fails  bool
	}{
		{
			yaml:   "abc: 123",
			result: Values{"abc": 123},
			fails:  false,
		},
		{
			yaml:  "123!!123",
			fails: true,
		},
	}

	for _, test := range tests {
		v, err := FromYaml(test.yaml)
		if test.fails {
			s.Require().Error(err)
			s.Require().Empty(v)
		} else {
			s.Require().NoError(err)
			s.Require().Equal(test.result, v)
		}
	}
}

func (s *ExtraTestSuite) TestExec() {
	res, err := Exec("pwd", []interface{}{})
	s.Require().NoError(err)

	pwd, err := os.Getwd()
	s.Require().NoError(err)

	s.Require().Equal(pwd, strings.TrimSpace(res))
}

func (s *ExtraTestSuite) TestExecInvalidArg() {
	res, err := Exec("pwd", []interface{}{123})
	s.Require().Error(err)
	s.Require().Empty(res)
}

func (s *ExtraTestSuite) TestExecError() {
	res, err := Exec("pwd", []interface{}{"123"})
	expected := &exec.ExitError{}
	s.Require().ErrorAs(err, &expected)
	s.Require().Empty(res)
}

func (s *ExtraTestSuite) TestExecStdin() {
	input := "123"
	res, err := Exec("cat", []interface{}{}, input)
	s.Require().NoError(err)
	s.Require().Equal(input, res)
}

func (s *ExtraTestSuite) TestSetValueAtPath() {
	data := Values{
		"a": map[string]interface{}{
			"b": "123",
		},
		"c": 123,
	}

	tests := []struct {
		path   string
		value  interface{}
		result Values
		fails  bool
	}{
		{
			path:  "c",
			value: 321,
			result: Values{
				"a": map[string]interface{}{"b": "123"},
				"c": 321,
			},
			fails: false,
		},
		{
			path:  "a.b",
			value: "321",
			result: Values{
				"a": map[string]interface{}{"b": "321"},
				"c": 321,
			},
			fails: false,
		},
		{
			path:  "a.c",
			value: "321",
			result: Values{
				"a": map[string]interface{}{"b": "321", "c": "321"},
				"c": 321,
			},
			fails: false,
		},
		{
			path:   "c.a",
			value:  "321",
			result: nil,
			fails:  true,
		},
	}

	for _, test := range tests {
		res, err := SetValueAtPath(test.path, test.value, data)
		if test.fails {
			s.Require().Error(err)
			s.Require().Nil(res)
		} else {
			s.Require().NoError(err)
			s.Require().Equal(test.result, res)
		}
	}
}

func (s *ExtraTestSuite) TestRequiredEnv() {
	name := s.T().Name()

	res, err := RequiredEnv(name)
	s.Require().Error(err)
	s.Require().Empty(res)

	data := "test"
	s.T().Setenv(name, data)

	res, err = RequiredEnv(name)
	s.Require().NoError(err)
	s.Require().Equal(data, res)
}

func (s *ExtraTestSuite) TestRequired() {
	tests := []struct {
		data  interface{}
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
		res, err := Required("blabla", t.data)
		if t.fails {
			s.Require().Error(err)
			s.Require().Nil(res)
		} else {
			s.Require().NoError(err)
			s.Require().Equal(t.data, res)
		}
	}
}

func (s *ExtraTestSuite) TestReadFile() {
	tmpDir := s.T().TempDir()
	tmpFile := filepath.Join(tmpDir, "blablafile")

	res, err := ReadFile(tmpFile)

	s.Require().Equal("", res)
	s.Require().ErrorIs(err, os.ErrNotExist)

	data := s.T().Name()

	s.Require().NoError(os.WriteFile(tmpFile, []byte(data), 0666))
	s.Require().FileExists(tmpFile)

	res, err = ReadFile(tmpFile)

	s.Require().NoError(err)
	s.Require().Equal(data, res)
}

func (s *ExtraTestSuite) TestGet() {
	data := Values{
		"a": map[string]interface{}{
			"b": "123",
		},
		"c": 123,
	}

	tests := []struct {
		path   string
		result interface{}
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
	}

	for _, test := range tests {
		res, err := Get(test.path, data)
		if test.fails {
			s.Require().Error(err)
			s.Require().Nil(res)
		} else {
			s.Require().NoError(err)
			s.Require().Equal(test.result, res)
		}
	}
}

func (s *ExtraTestSuite) TestHasKey() {
	data := Values{
		"a": map[string]interface{}{
			"b": "123",
		},
		"c": 123,
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
	}

	for _, test := range tests {
		res, err := HasKey(test.path, data)
		s.Require().Equal(test.result, res)

		if test.fails {
			s.Require().Error(err)
		} else {
			s.Require().NoError(err)
		}
	}
}

func TestExtraTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ExtraTestSuite))
}
