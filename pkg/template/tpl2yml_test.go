package template_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/template"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
)

type Tpl2YmlTestSuite struct {
	suite.Suite
}

func (s *Tpl2YmlTestSuite) TestNonExistingTemplate() {
	tmpDir := s.T().TempDir()
	tpl := filepath.Join(tmpDir, "values.yaml")
	dst := filepath.Join(tmpDir, "values.yaml")

	err := template.Tpl2yml(tpl, dst, nil, template.TemplaterSprig)
	s.Require().ErrorIs(err, os.ErrNotExist)
}

func (s *Tpl2YmlTestSuite) TestNonExistingDestDir() {
	tmpDir := s.T().TempDir()
	tpl := filepath.Join(tests.Root, "05_values.yaml")
	dst := filepath.Join(tmpDir, "blabla", "values.yaml")

	err := template.Tpl2yml(tpl, dst, nil, template.TemplaterSprig)
	s.Require().NoError(err)
}

func (s *Tpl2YmlTestSuite) TestMissingData() {
	tmpDir := s.T().TempDir()
	tpl := filepath.Join(tests.Root, "08_values.yaml")
	dst := filepath.Join(tmpDir, "values.yaml")

	err := template.Tpl2yml(tpl, dst, nil, template.TemplaterSprig)
	s.Require().EqualError(err, "failed to render template: failed to parse template: template: tpl:1: function \"defineDatasource\" not defined")
}

func (s *Tpl2YmlTestSuite) TestDisabledGomplate() {
	tmpDir := s.T().TempDir()
	tpl := filepath.Join(tests.Root, "09_values.yaml")
	dst := filepath.Join(tmpDir, "values.yaml")

	err := template.Tpl2yml(tpl, dst, nil, template.TemplaterSprig)
	s.Require().Error(err)
}

func (s *Tpl2YmlTestSuite) TestEnabledGomplate() {
	tmpDir := s.T().TempDir()
	tpl := filepath.Join(tests.Root, "09_values.yaml")
	dst := filepath.Join(tmpDir, "values.yaml")

	err := template.Tpl2yml(tpl, dst, nil, template.TemplaterGomplate)
	s.Require().NoError(err)
}

func TestTpl2YmlTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(Tpl2YmlTestSuite))
}
