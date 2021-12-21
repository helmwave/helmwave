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
	gomplateConfig := &template.GomplateConfig{
		Enabled: false,
	}

	tmpDir := s.T().TempDir()
	tpl := filepath.Join(tmpDir, "values.yaml")
	dst := filepath.Join(tmpDir, "values.yaml")

	err := template.Tpl2yml(tpl, dst, nil, gomplateConfig)
	s.Require().ErrorIs(err, os.ErrNotExist)
}

func (s *Tpl2YmlTestSuite) TestNonExistingDestDir() {
	gomplateConfig := &template.GomplateConfig{
		Enabled: false,
	}

	tmpDir := s.T().TempDir()
	tpl := filepath.Join(tests.Root, "05_values.yaml")
	dst := filepath.Join(tmpDir, "blabla", "values.yaml")

	err := template.Tpl2yml(tpl, dst, nil, gomplateConfig)
	s.Require().NoError(err)
}

func (s *Tpl2YmlTestSuite) TestMissingData() {
	gomplateConfig := &template.GomplateConfig{
		Enabled: false,
	}

	tmpDir := s.T().TempDir()
	tpl := filepath.Join(tests.Root, "08_values.yaml")
	dst := filepath.Join(tmpDir, "values.yaml")

	err := template.Tpl2yml(tpl, dst, nil, gomplateConfig)
	s.Require().EqualError(err, "template: tpl:1: function \"include\" not defined")
}

func (s *Tpl2YmlTestSuite) TestDisabledGomplate() {
	gomplateConfig := &template.GomplateConfig{
		Enabled: false,
	}

	tmpDir := s.T().TempDir()
	tpl := filepath.Join(tests.Root, "09_values.yaml")
	dst := filepath.Join(tmpDir, "values.yaml")

	err := template.Tpl2yml(tpl, dst, nil, gomplateConfig)
	s.Require().Error(err)
}

func (s *Tpl2YmlTestSuite) TestEnabledGomplate() {
	gomplateConfig := &template.GomplateConfig{
		Enabled: true,
	}

	tmpDir := s.T().TempDir()
	tpl := filepath.Join(tests.Root, "09_values.yaml")
	dst := filepath.Join(tmpDir, "values.yaml")

	err := template.Tpl2yml(tpl, dst, nil, gomplateConfig)
	s.Require().NoError(err)
}

func TestTpl2YmlTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(Tpl2YmlTestSuite))
}
