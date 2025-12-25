package release

import (
	"context"
	"crypto"
	_ "crypto/md5" // for crypto.MD5.New to work
	"encoding/hex"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	gotemplate "text/template"
	"time"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/helmwave/helmwave/pkg/template"
	"github.com/invopop/jsonschema"
	log "github.com/sirupsen/logrus"
	"github.com/stoewer/go-strcase"
	"gopkg.in/yaml.v3"
)

// ValuesReference is used to match source values file path and temporary.
//
//nolint:lll
type ValuesReference struct {
	Src            string `yaml:"src" json:"src" jsonschema:"required,description=Source of values. Can be local path or HTTP URL"`
	Dst            string `yaml:"dst" json:"dst" jsonschema:"readOnly"`
	DelimiterLeft  string `yaml:"delimiter_left,omitempty" json:"delimiter_left,omitempty"  jsonschema:"Set left delimiter for template engine,default={{"`
	DelimiterRight string `yaml:"delimiter_right,omitempty" json:"delimiter_right,omitempty" jsonschema:"Set right delimiter for template engine,default=}}"`
	Renderer       string `yaml:"renderer" json:"renderer" jsonschema:"description=How to render the file,enum=sprig,enum=gomplate,enum=copy,enum=sops"`
	Strict         bool   `yaml:"strict" json:"strict" jsonschema:"description=Whether to fail if values is not found,default=false"`
}

//nolint:gocritic
func (v ValuesReference) JSONSchema() *jsonschema.Schema {
	r := &jsonschema.Reflector{
		DoNotReference:             true,
		RequiredFromJSONSchemaTags: true,
		KeyNamer:                   strcase.SnakeCase, // for action.ChartPathOptions
	}

	type values ValuesReference
	schema := r.Reflect(values(v))
	schema.OneOf = []*jsonschema.Schema{
		{
			Type: "string",
		},
		{
			Type: "object",
		},
	}
	schema.Type = ""

	return schema
}

// UnmarshalYAML flexible config.
func (v *ValuesReference) UnmarshalYAML(node *yaml.Node) error {
	type raw ValuesReference
	var err error
	switch node.Kind {
	// single value or reference to another value
	case yaml.ScalarNode, yaml.AliasNode:
		err = node.Decode(&v.Src)
	case yaml.MappingNode:
		err = node.Decode((*raw)(v))
	default:
		err = ErrUnknownFormat
	}

	if err != nil {
		return fmt.Errorf("failed to decode values reference %q from YAML: %w", node.Value, err)
	}

	return nil
}

// MarshalYAML is used to implement Marshaler interface of gopkg.in/yaml.v3.
func (v *ValuesReference) MarshalYAML() (any, error) {
	return struct {
		Src string
		Dst string
	}{
		Src: v.Src,
		Dst: v.Dst,
	}, nil
}

func (v *ValuesReference) isURL() bool {
	return helper.IsURL(v.Src)
}

// Download downloads values by source URL and places to destination path.
func (v *ValuesReference) Download(ctx context.Context) error {
	if err := helper.Download(ctx, v.Dst, v.Src); err != nil {
		return fmt.Errorf("failed to download values %s -> %s: %w", v.Src, v.Dst, err)
	}

	return nil
}

// Get returns destination path of values.
// func (v *ValuesReference) Get() string {
//	return v.Dst
// }

// SetUniq generates unique file path based on provided base directory, release uniqname and sha1 of source path.
func (v *ValuesReference) SetUniq(dir string, name uniqname.UniqName) *ValuesReference {
	h := crypto.MD5.New()
	h.Write([]byte(v.Src))
	hash := h.Sum(nil)
	s := hex.EncodeToString(hash)

	v.Dst = filepath.Join(dir, "values", name.String(), s+".yml")

	return v
}

// ProhibitDst Dst now is public method.
// Dst needs to marshal for export.
// Also, dst needs to unmarshal for import from plan.
func ProhibitDst(values []ValuesReference) error {
	for i := range values {
		v := values[i]
		if v.Dst != "" {
			return fmt.Errorf("dst %q not allowed here, this field reserved", v.Dst)
		}
	}

	return nil
}

// SetViaRelease downloads and templates values file.
// Returns ErrValuesNotExist if values can't be downloaded or doesn't exist in local FS.
func (v *ValuesReference) SetViaRelease(
	ctx context.Context,
	rel Config,
	dir, templater string,
	templateFuncs gotemplate.FuncMap,
	renderedFiles *renderedValuesFiles,
) error {
	if v.Renderer == "" {
		v.Renderer = templater
	}

	if v.Renderer == template.TemplaterGomplate && v.DelimiterLeft == "" && v.DelimiterRight == "" {
		v.DelimiterLeft = "[["
		v.DelimiterRight = "]]"
	}

	v.SetUniq(dir, rel.Uniq())

	l := rel.Logger().WithField("values src", v.Src).WithField("values Dst", v.Dst)

	l.Trace("Building values reference")

	data := struct {
		Release Config
	}{
		Release: rel,
	}

	err := v.fetch(ctx, l)
	if err != nil {
		return err
	}

	opts := []template.TemplaterOptions{template.SetDelimiters(v.DelimiterLeft, v.DelimiterRight)}
	for name, value := range templateFuncs {
		opts = append(opts, template.AddFunc(name, value))
	}

	if renderedFiles != nil {
		buf := &strings.Builder{}
		defer func() {
			renderedFiles.Add(v.Src, buf)
		}()

		opts = append(opts,
			template.CopyOutput(buf),
			template.AddFunc("getValues", func(filename string) (any, error) {
				s := renderedFiles.Get(filename).String()

				var res any
				err := yaml.Unmarshal([]byte(s), &res)

				//nolint:wrapcheck
				return res, err
			},
			))
	}
	if v.isURL() {
		err = template.Tpl2yml(ctx, v.Dst, v.Dst, data, v.Renderer, opts...)
	} else {
		err = template.Tpl2yml(ctx, v.Src, v.Dst, data, v.Renderer, opts...)
	}

	if err != nil {
		return fmt.Errorf("failed to render %q values: %w", v.Src, err)
	}

	return nil
}

func (v *ValuesReference) fetch(ctx context.Context, l *log.Entry) error {
	if v.isURL() {
		err := v.Download(ctx)
		if err != nil {
			l.WithError(err).Warnf("%q skipping: cant download", v.Src)

			return ErrValuesNotExist
		}
	} else if !helper.IsExists(v.Src) {
		l.Warn("skipping: local file not found")

		return ErrValuesNotExist
	}

	return nil
}

func (rel *config) BuildValues(
	ctx context.Context,
	dir, templater string,
	templateFuncs gotemplate.FuncMap,
) error {
	mu := &sync.Mutex{}
	vals := rel.Values()

	wg := parallel.NewWaitGroup()
	wg.Add(len(vals))

	// we need to keep rendered values in memory to use them in other values
	renderedValuesMap := newRenderedValuesFiles()
	// we need to keep track of values that we need to delete (e.g. non-existent files)
	toDeleteMap := make(map[*ValuesReference]bool)

	// just in case of dependency cycle or long http requests
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	l := rel.Logger()

	for i := range vals {
		go func(v *ValuesReference) {
			defer wg.Done()

			l := l.WithField("values", v)

			err := v.SetViaRelease(ctx, rel, dir, templater, templateFuncs, renderedValuesMap)
			switch {
			case !v.Strict && errors.Is(ErrValuesNotExist, err):
				l.WithError(err).Warn("skipping values...")
				mu.Lock()
				toDeleteMap[v] = true
				mu.Unlock()
			case err != nil:
				l.WithError(err).Error("failed to build values")

				wg.ErrChan() <- err
			}
		}(&vals[i])
	}

	err := wg.WaitWithContext(ctx)
	if err != nil {
		return err
	}

	for i := len(vals) - 1; i >= 0; i-- {
		if toDeleteMap[&vals[i]] {
			vals = slices.Delete(vals, i, i+1)
		}
	}
	rel.ValuesF = vals

	return nil
}

type renderedValuesFiles struct {
	bufs map[string]fmt.Stringer
	cond *sync.Cond
}

func newRenderedValuesFiles() *renderedValuesFiles {
	r := &renderedValuesFiles{
		bufs: make(map[string]fmt.Stringer),
	}
	r.cond = sync.NewCond(&sync.Mutex{})

	return r
}

func (r *renderedValuesFiles) Add(name string, buf fmt.Stringer) {
	r.cond.L.Lock()
	r.bufs[name] = buf
	r.cond.Broadcast()
	r.cond.L.Unlock()
}

func (r *renderedValuesFiles) Get(name string) fmt.Stringer {
	r.cond.L.Lock()
	for r.bufs[name] == nil {
		r.cond.Wait()
	}
	r.cond.L.Unlock()

	return r.bufs[name]
}
