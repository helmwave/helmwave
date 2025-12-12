package plan

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"
	gotemplate "text/template"

	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/helmwave/helmwave/pkg/template"
	"github.com/stretchr/testify/suite"
)

type TemplateFuncsTestSuite struct {
	suite.Suite
}

func TestTemplateFuncsTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(TemplateFuncsTestSuite))
}

func (ts *TemplateFuncsTestSuite) renderTemplate(ctx context.Context, tpl string, data any, templateFuncs gotemplate.FuncMap) (string, error) {
	tmpDir := ts.T().TempDir()
	tplFile := filepath.Join(tmpDir, "test.tpl")
	ymlFile := filepath.Join(tmpDir, "test.yml")

	err := os.WriteFile(tplFile, []byte(tpl), 0o600)
	ts.Require().NoError(err)

	opts := []template.TemplaterOptions{}
	for name, value := range templateFuncs {
		opts = append(opts, template.AddFunc(name, value))
	}

	err = template.Tpl2yml(ctx, tplFile, ymlFile, data, template.TemplaterSprig, opts...)
	if err != nil {
		return "", err
	}

	content, err := os.ReadFile(ymlFile)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func (ts *TemplateFuncsTestSuite) TestGetPlan() {
	p := New(".")
	body := p.NewBody()
	body.Project = "my-project"
	body.Version = "1.0.0"

	mu := &sync.Mutex{}
	templateFuncs := p.templateFuncs(mu)

	ctx := context.Background()
	tpl := `{{ $plan := getPlan }}{{ $plan.project }},{{ $plan.version }}`

	data := struct {
		Release struct {
			Name string
		}
	}{
		Release: struct {
			Name string
		}{
			Name: "test-release",
		},
	}

	rendered, err := ts.renderTemplate(ctx, tpl, data, templateFuncs)
	ts.Require().NoError(err)

	ts.Require().Equal("my-project,1.0.0", rendered)
}

func (ts *TemplateFuncsTestSuite) TestGetManifestsEmpty() {
	p := New(".")
	p.NewBody()

	mu := &sync.Mutex{}
	templateFuncs := p.templateFuncs(mu)

	ts.Require().NotEmpty(templateFuncs)
}

func (ts *TemplateFuncsTestSuite) TestGetManifestsSingleDocument() {
	p := New(".")
	p.NewBody()

	uniq, _ := uniqname.NewFromString("redis@default")
	p.manifests[uniq] = `---
apiVersion: v1
kind: ConfigMap
metadata:
  name: test`

	mu := &sync.Mutex{}
	templateFuncs := p.templateFuncs(mu)

	ctx := context.Background()
	tpl := `{{ $manifests := getManifests "redis@default" }}{{ len $manifests }}`
	rendered, err := ts.renderTemplate(ctx, tpl, nil, templateFuncs)
	ts.Require().NoError(err)

	ts.Require().Equal("1", rendered)
}

func (ts *TemplateFuncsTestSuite) TestGetManifestsMultipleReleases() {
	p := New(".")
	p.NewBody()

	uniq1, _ := uniqname.NewFromString("redis@default")
	p.manifests[uniq1] = `---
apiVersion: v1
kind: ConfigMap
metadata:
  name: redis-config
---
apiVersion: v1
kind: Service
metadata:
  name: redis-svc`

	uniq2, _ := uniqname.NewFromString("nginx@default")
	p.manifests[uniq2] = `---
apiVersion: v1
kind: ConfigMap
metadata:
  name: nginx-config
---
apiVersion: v1
kind: Service
metadata:
  name: nginx-svc`

	mu := &sync.Mutex{}
	templateFuncs := p.templateFuncs(mu)

	ctx := context.Background()
	tpl := `{{ $m1 := getManifests "redis@default" }}{{ $m2 := getManifests "nginx@default" }}{{ len $m1 }},{{ len $m2 }}`
	rendered, err := ts.renderTemplate(ctx, tpl, nil, templateFuncs)
	ts.Require().NoError(err)

	ts.Require().Equal("2,2", rendered)
}

func (ts *TemplateFuncsTestSuite) TestGetManifestsContent() {
	p := New(".")
	p.NewBody()

	uniq, _ := uniqname.NewFromString("redis@default")
	p.manifests[uniq] = `---
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-config
data:
  key: value`

	mu := &sync.Mutex{}
	templateFuncs := p.templateFuncs(mu)

	ctx := context.Background()
	tpl := `{{ $manifests := getManifests "redis@default" }}{{ range $manifests }}{{ .kind }}{{ end }}`
	rendered, err := ts.renderTemplate(ctx, tpl, nil, templateFuncs)
	ts.Require().NoError(err)

	ts.Require().Equal("ConfigMap", rendered)
}

func (ts *TemplateFuncsTestSuite) TestGetManifestsNotFound() {
	p := New(".")
	p.NewBody()

	mu := &sync.Mutex{}
	templateFuncs := p.templateFuncs(mu)

	ctx := context.Background()
	tpl := `{{ $manifests := getManifests "nonexistent@default" }}`

	_, err := ts.renderTemplate(ctx, tpl, nil, templateFuncs)
	ts.Require().Error(err)
	ts.Require().ErrorContains(err, "not found")
}
