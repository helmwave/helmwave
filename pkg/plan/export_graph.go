package plan

import (
	"fmt"

	"github.com/helmwave/go-fsimpl"
	"github.com/helmwave/helmwave/pkg/helper"
)

func (p *Plan) exportGraphMD(plandirFS fsimpl.WriteableFS) error {
	found := false
	for _, rel := range p.body.Releases {
		if len(rel.DependsOn()) > 0 {
			found = true

			break
		}
	}

	if !found {
		return nil
	}

	const filename = "graph.md"
	f, err := helper.CreateFile(plandirFS, filename)
	if err != nil {
		return err
	}

	_, err = f.Write([]byte(p.graphMD))
	if err != nil {
		return fmt.Errorf("failed to write graph file %s: %w", filename, err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close graph file %s: %w", filename, err)
	}

	return nil
}
