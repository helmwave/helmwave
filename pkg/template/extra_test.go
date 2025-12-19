package template_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/helmwave/helmwave/pkg/template"
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

	y, err := template.ToYaml(data)
	s.Require().NoError(err)
	s.Require().YAMLEq(yamlData, y)
}

type raw struct{}

func (r raw) MarshalYAML() (any, error) {
	return nil, os.ErrNotExist
}

func (s *ExtraTestSuite) TestToYamlNil() {
	data := raw{}
	y, err := template.ToYaml(data)
	s.Require().Equal("", y)
	s.Require().ErrorIs(err, os.ErrNotExist)
}

func (s *ExtraTestSuite) TestFromYaml() {
	tests := []struct {
		result template.Values
		yaml   string
		fails  bool
	}{
		{
			yaml:   "abc: 123",
			result: template.Values{"abc": 123},
			fails:  false,
		},
		{
			yaml:  "123!!123",
			fails: true,
		},
	}

	for i := range tests {
		v, err := template.FromYaml(tests[i].yaml)
		if tests[i].fails {
			s.Error(err)
			s.Empty(v)
		} else {
			s.NoError(err)
			s.Equal(tests[i].result, v)
		}
	}
}

func (s *ExtraTestSuite) TestFromYamlArray() {
	tests := []struct {
		yaml   string
		result []any
		fails  bool
	}{
		{
			yaml:   "[1, 2, 3]",
			result: []any{1, 2, 3},
			fails:  false,
		},
		{
			yaml:   "- a\n- b\n- c",
			result: []any{"a", "b", "c"},
			fails:  false,
		},
		{
			yaml:  "a: 123",
			fails: true,
		},
	}

	for i := range tests {
		v, err := template.FromYamlArray(tests[i].yaml)
		if tests[i].fails {
			s.Error(err)
			s.Empty(v)
		} else {
			s.NoError(err)
			s.Equal(tests[i].result, v)
		}
	}
}

func (s *ExtraTestSuite) TestFromYamlAll() {
	tests := []struct {
		yaml   string
		result []any
		fails  bool
	}{
		{
			yaml:   "1",
			result: []any{1},
			fails:  false,
		},
		{
			yaml:   "1\n---\na: 123\n---\n[1, 2, 3]",
			result: []any{1, template.Values{"a": 123}, []any{1, 2, 3}},
			fails:  false,
		},
		{
			yaml:   "---\napiVersion: v1\nkind: ConfigMap",
			result: []any{template.Values{"apiVersion": "v1", "kind": "ConfigMap"}},
			fails:  false,
		},
		{
			yaml:   "---\napiVersion: v1\n---\nkind: ConfigMap",
			result: []any{template.Values{"apiVersion": "v1"}, template.Values{"kind": "ConfigMap"}},
			fails:  false,
		},
	}

	for i := range tests {
		v, err := template.FromYamlAll(tests[i].yaml)
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
	res, err := template.Exec("pwd", []any{})
	s.Require().NoError(err)

	pwd, err := os.Getwd()
	s.Require().NoError(err)

	s.Require().Equal(pwd, strings.TrimSpace(res))

	res, err = template.Exec("echo", []any{"-n", "123"})
	s.Require().NoError(err)
	s.Require().Equal("123", res)
}

func (s *ExtraTestSuite) TestExecInvalidArg() {
	res, err := template.Exec("pwd", []any{123})
	s.Require().Error(err)
	s.Require().Empty(res)
}

func (s *ExtraTestSuite) TestExecError() {
	res, err := template.Exec(s.T().Name(), []any{})
	s.Require().Error(err)
	s.Require().Empty(res)
}

func (s *ExtraTestSuite) TestExecStdin() {
	input := "123"
	res, err := template.Exec("cat", []any{}, input)
	s.Require().NoError(err)
	s.Require().Equal(input, res)
}

func (s *ExtraTestSuite) TestExecCommandString() {
	// Test: {{ exec "echo -n 123" }}
	res, err := template.Exec("echo -n 123")
	s.Require().NoError(err)
	s.Require().Equal("123", res)

	// Test: {{ exec "echo -n hello world" }}
	res, err = template.Exec("echo -n hello world")
	s.Require().NoError(err)
	s.Require().Equal("hello world", res)
}

func (s *ExtraTestSuite) TestExecCommandStringWithQuotes() {
	// Test: {{ exec "echo -n 'hello world'" }} - single quotes
	res, err := template.Exec("echo -n 'hello world'")
	s.Require().NoError(err)
	s.Require().Equal("hello world", res)

	// Test: {{ exec `echo -n "hello world"` }} - double quotes
	res, err = template.Exec(`echo -n "hello world"`)
	s.Require().NoError(err)
	s.Require().Equal("hello world", res)
}

func (s *ExtraTestSuite) TestExecCommandStringWithStdin() {
	// Test: {{ "input" | exec "cat" }}
	input := "test input"
	res, err := template.Exec("cat", input)
	s.Require().NoError(err)
	s.Require().Equal(input, res)
}

func (s *ExtraTestSuite) TestExecCommandStringWithNilArgs() {
	// Test: {{ exec "pwd" }} with nil passed (simulating template behavior)
	pwd, err := os.Getwd()
	s.Require().NoError(err)

	res, err := template.Exec("pwd", nil)
	s.Require().NoError(err)
	s.Require().Equal(pwd, strings.TrimSpace(res))
}

func (s *ExtraTestSuite) TestExecUnclosedQuote() {
	// Test unclosed quote error (shlex returns "EOF found when expecting closing quote")
	_, err := template.Exec("echo 'unclosed")
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "expecting closing quote")

	_, err = template.Exec(`echo "unclosed`)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "expecting closing quote")
}

func (s *ExtraTestSuite) TestSetValueAtPath() {
	data := template.Values{
		"a": map[string]any{
			"b": "123",
		},
		"c": 123,
		"d": map[any]any{
			"e": "f",
		},
	}

	tests := []struct {
		result template.Values
		value  any
		path   string
		fails  bool
	}{
		{
			path:  "c",
			value: 321,
			result: template.Values{
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
			result: template.Values{
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
			result: template.Values{
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
			result: template.Values{
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
		res, err := template.SetValueAtPath(tests[i].path, tests[i].value, data)
		if tests[i].fails {
			s.Error(err)
			s.Nil(res)
		} else {
			s.NoError(err)
			s.Equal(tests[i].result, res)
		}
	}
}

func (s *ExtraTestSuite) TestSetValueAtPathWithArrayIndex() {
	data := template.Values{
		"items": []any{"a", "b", "c"},
		"nested": map[string]any{
			"list": []any{
				map[string]any{"name": "first"},
				map[string]any{"name": "second"},
			},
		},
	}

	// Test setting array element with dot notation: items.1
	res, err := template.SetValueAtPath("items.1", "updated", data)
	s.NoError(err)
	s.Equal("updated", res["items"].([]any)[1])

	// Test setting nested array element property: nested.list.0.name
	res, err = template.SetValueAtPath("nested.list.0.name", "changed", data)
	s.NoError(err)
	s.Equal("changed", res["nested"].(map[string]any)["list"].([]any)[0].(map[string]any)["name"])

	// Test out of bounds error
	_, err = template.SetValueAtPath("items.10", "x", data)
	s.Error(err)

	// Numeric key on map sets key "0"
	res, err = template.SetValueAtPath("nested.0", "x", data)
	s.NoError(err)
	s.Equal("x", res["nested"].(map[string]any)["0"])
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
		res, err := template.Required("blabla", t.data)
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

	res, err := template.ReadFile(tmpFile)

	s.Require().Equal("", res)
	s.Require().ErrorIs(err, os.ErrNotExist)

	data := s.T().Name()

	s.Require().NoError(os.WriteFile(tmpFile, []byte(data), 0o600))
	s.Require().FileExists(tmpFile)

	res, err = template.ReadFile(tmpFile)

	s.Require().NoError(err)
	s.Require().Equal(data, res)
}

func (s *ExtraTestSuite) TestGetValueAtPath() {
	data := template.Values{
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
		res, err := template.GetValueAtPath(tests[i].path, data)
		if tests[i].fails {
			s.Error(err)
			s.Nil(res)
		} else {
			s.NoError(err)
			s.Equal(tests[i].result, res)
		}
	}
}

func (s *ExtraTestSuite) TestGetValueAtPathWithArrayIndex() {
	data := template.Values{
		"items": []any{"a", "b", "c"},
		"nested": map[string]any{
			"list": []any{
				map[string]any{"name": "first"},
				map[string]any{"name": "second"},
			},
		},
	}

	// Test getting array element with dot notation: items.1
	res, err := template.GetValueAtPath("items.1", data)
	s.NoError(err)
	s.Equal("b", res)

	// Test getting nested array element property: nested.list.0.name
	res, err = template.GetValueAtPath("nested.list.0.name", data)
	s.NoError(err)
	s.Equal("first", res)

	// Test out of bounds error
	_, err = template.GetValueAtPath("items.10", data)
	s.Error(err)

	// Test invalid index error
	_, err = template.GetValueAtPath("items.abc", data)
	s.Error(err)

	// Test with default value for out of bounds
	res, err = template.GetValueAtPath("items.10", "default", data)
	s.NoError(err)
	s.Equal("default", res)
}

func (s *ExtraTestSuite) TestHasValueAtPath() {
	data := template.Values{
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
		res, err := template.HasValueAtPath(test.path, data)
		s.Equal(test.result, res)

		if test.fails {
			s.Error(err)
		} else {
			s.NoError(err)
		}
	}
}

func (s *ExtraTestSuite) TestHasValueAtPathWithArrayIndex() {
	data := template.Values{
		"items": []any{"a", "b", "c"},
		"nested": map[string]any{
			"list": []any{
				map[string]any{"name": "first"},
				map[string]any{"name": "second"},
			},
		},
	}

	// Test checking array element exists: items.1
	res, err := template.HasValueAtPath("items.1", data)
	s.NoError(err)
	s.True(res)

	// Test checking nested array element property exists: nested.list.0.name
	res, err = template.HasValueAtPath("nested.list.0.name", data)
	s.NoError(err)
	s.True(res)

	// Test out of bounds returns false
	res, err = template.HasValueAtPath("items.10", data)
	s.NoError(err)
	s.False(res)

	// Test invalid index returns false
	res, err = template.HasValueAtPath("items.abc", data)
	s.NoError(err)
	s.False(res)

	// Test with default (defSet=true) for out of bounds
	res, err = template.HasValueAtPath("items.10", true, data)
	s.NoError(err)
	s.True(res)
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

	res, err := template.RequiredEnv(name)
	s.Require().Error(err)
	s.Require().Empty(res)

	data := "test"
	s.T().Setenv(name, data)

	res, err = template.RequiredEnv(name)
	s.Require().NoError(err)
	s.Require().Equal(data, res)
}

//nolint:paralleltest // can't parallel because of setenv
func TestNonParallelExtraTestSuite(t *testing.T) {
	// t.Parallel()
	suite.Run(t, new(NonParallelExtraTestSuite))
}
