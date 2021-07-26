package plan

import (
	"crypto/sha1"
	"encoding/hex"
	"os"
	"path/filepath"

	"github.com/helmwave/helmwave/pkg/helper"
)

// Export allows save plan to file
func (p *Plan) Export() error {
	if err := os.RemoveAll(p.dir); err != nil {
		return err
	}

	if err := p.exportManifest(); err != nil {
		return err
	}

	// TODO make it better later
	if err := p.buildValues(p.dir); err != nil {
		return err
	}

	if err := helper.SaveInterface(p.fullPath, p.body); err != nil {
		return err
	}

	if err := p.exportGraphMD(); err != nil {
		return err
	}

	return nil
}

func (p *Plan) exportManifest() error {
	if len(p.body.Releases) == 0 {
		return nil
	}

	for k, v := range p.manifests {
		m := filepath.Join(p.dir, Manifest, string(k))

		f, err := helper.CreateFile(m)
		if err != nil {
			return err
		}

		_, err = f.WriteString(v)
		if err != nil {
			return err
		}

		err = f.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Plan) exportGraphMD() error {
	f, err := helper.CreateFile(filepath.Join(p.dir, "graph.md"))
	if err != nil {
		return err
	}

	_, err = f.WriteString(p.graphMD)
	if err != nil {
		return err
	}

	return f.Close()
}

func (p *Plan) exportValues() error {
	if len(p.body.Releases) == 0 {
		return nil
	}

	found := false
	h := sha1.New() // nolint:gosec

	for i, rel := range p.body.Releases {
		for j := range p.body.Releases[i].Values {
			found = true
			h.Write([]byte(p.body.Releases[i].Values[j].Src))
			hash := h.Sum(nil)
			hs := hex.EncodeToString(hash)
			p.body.Releases[i].Values[j].Set(filepath.Join(p.dir, "values", string(rel.Uniq()), hs+".yml"))
		}
	}

	if !found {
		return nil
	}

	return os.Rename(
		filepath.Join(p.tmpDir, Values),
		filepath.Join(p.dir, Values),
	)
}

// IsExist returns true if planfile exists
func (p *Plan) IsExist() bool {
	return helper.IsExists(p.fullPath)
}

// IsManifestExist returns true if planfile exists
func (p *Plan) IsManifestExist() bool {
	return helper.IsExists(filepath.Join(p.dir, Manifest))
}
