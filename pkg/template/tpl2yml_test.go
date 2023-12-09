package template_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/template"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
)

type Tpl2YmlTestSuite struct {
	suite.Suite

	ctx context.Context
}

func TestTpl2YmlTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(Tpl2YmlTestSuite))
}

func (ts *Tpl2YmlTestSuite) SetupTest() {
	ts.ctx = tests.GetContext(ts.T())
}

func (ts *Tpl2YmlTestSuite) TestNonExistingTemplate() {
	tmpDir := ts.T().TempDir()
	tpl := filepath.Join(tmpDir, "values.yaml")
	dst := filepath.Join(tmpDir, "values.yaml")

	err := template.Tpl2yml(ts.ctx, tpl, dst, nil, template.TemplaterSprig)
	ts.Require().ErrorIs(err, os.ErrNotExist)
}

func (ts *Tpl2YmlTestSuite) TestNonExistingDestDir() {
	tmpDir := ts.T().TempDir()
	tpl := filepath.Join(tests.Root, "05_values.yaml")
	dst := filepath.Join(tmpDir, "blabla", "values.yaml")

	err := template.Tpl2yml(ts.ctx, tpl, dst, nil, template.TemplaterSprig)
	ts.Require().NoError(err)
}

func (ts *Tpl2YmlTestSuite) TestMissingData() {
	tmpDir := ts.T().TempDir()
	tpl := filepath.Join(tests.Root, "08_values.yaml")
	dst := filepath.Join(tmpDir, "values.yaml")

	err := template.Tpl2yml(ts.ctx, tpl, dst, nil, template.TemplaterSprig)
	ts.Require().EqualError(err, "failed to render template: failed to parse template: template: tpl:1: function \"defineDatasource\" not defined")
}

func (ts *Tpl2YmlTestSuite) TestDisabledGomplate() {
	tmpDir := ts.T().TempDir()
	tpl := filepath.Join(tests.Root, "09_values.yaml")
	dst := filepath.Join(tmpDir, "values.yaml")

	err := template.Tpl2yml(ts.ctx, tpl, dst, nil, template.TemplaterSprig)
	ts.Require().Error(err)
}

func (ts *Tpl2YmlTestSuite) TestEnabledGomplate() {
	tmpDir := ts.T().TempDir()
	tpl := filepath.Join(tests.Root, "09_values.yaml")
	dst := filepath.Join(tmpDir, "values.yaml")

	err := template.Tpl2yml(ts.ctx, tpl, dst, nil, template.TemplaterGomplate)
	ts.Require().NoError(err)
}
