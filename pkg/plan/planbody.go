package plan

import (
	"fmt"

	"github.com/helmwave/helmwave/pkg/registry"
)

type planBody struct {
	Project      string
	Version      string
	Repositories repoConfigs
	Registries   registry.Configs
	Releases     releaseConfigs
}

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

	// Setup dev version
	// if b.Version == "" {
	// 	 b.Version = version.Version
	// }

	if err := b.Validate(); err != nil {
		return nil, err
	}

	return b, nil
}
