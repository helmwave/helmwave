package plan

import (
	"bytes"
	"fmt"
	"github.com/helmwave/helmwave/pkg/registry"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/repo"
	"github.com/helmwave/helmwave/pkg/version"

	"gopkg.in/yaml.v3"
	"io/fs"
)

// planBody contains only yaml fields
type planBody struct {
	Project      string
	Version      string
	Repositories repo.Configs
	Registries   registry.Configs
	Releases     release.Configs
}

func N() *planBody {
	return &planBody{}
}

// NewBody creates planBody
func NewBody(fsys fs.FS, file string) (*planBody, error) { // nolint:revive
	b := &planBody{
		Version: version.Version,
	}

	src, err := fs.ReadFile(fsys, file)
	if err != nil {
		return b, fmt.Errorf("failed to read plan file %s: %w", file, err)
	}

	err = yaml.Unmarshal(src, b)
	if err != nil {
		return b, fmt.Errorf("failed to unmarshal YAML plan %s: %w", file, err)
	}

	if err := b.Validate(); err != nil {
		return nil, err
	}

	return b, nil
}

func (p *Plan) NewBody() error {
	b := &planBody{
		Version: version.Version,
	}

	buf := new(bytes.Buffer)
	_, err := p.store.Read(p.File(), buf)
	err = yaml.Unmarshal(buf.Bytes(), b)

	if err != nil {
		return fmt.Errorf("failed to read plan file %s: %w", p.File(), err)
	}

	return nil
}
