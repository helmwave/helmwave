package release

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/helmwave/helmwave/pkg/template"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// ErrSkipValues is returned when values cannot be used and are skipped.
var ErrSkipValues = errors.New("values have been skipped")

// ValuesReference is used to match source values file path and temporary.
type ValuesReference struct {
	Src string
	dst string
}

// UnmarshalYAML is used to implement Unmarshaler interface of gopkg.in/yaml.v3.
func (v *ValuesReference) UnmarshalYAML(node *yaml.Node) error {
	m := make(map[string]string)
	if err := node.Decode(&m); err != nil {
		if err := node.Decode(&v.Src); err != nil {
			return fmt.Errorf("failed to decode values reference from YAML: %w", err)
		}
	}

	v.Src = m["src"]
	v.dst = m["dst"]

	return nil
}

// MarshalYAML is used to implement Marshaler interface of gopkg.in/yaml.v3.
func (v ValuesReference) MarshalYAML() (interface{}, error) {
	return struct {
		Src string
		Dst string
	}{
		Src: v.Src,
		Dst: v.dst,
	}, nil
}

func (v *ValuesReference) isURL() bool {
	return helper.IsURL(v.Src)
}

// Download downloads values by source URL and places to destination path.
func (v *ValuesReference) Download() error {
	if err := helper.Download(v.dst, v.Src); err != nil {
		return fmt.Errorf("failed to download values %s -> %s: %w", v.Src, v.dst, err)
	}

	return nil
}

// Get returns destination path of values.
func (v *ValuesReference) Get() string {
	return v.dst
}

// SetUniq generates unique file path based on provided base directory, release uniqname and sha1 of source path.
func (v *ValuesReference) SetUniq(dir string, name uniqname.UniqName) *ValuesReference {
	h := sha1.New() // nolint:gosec
	h.Write([]byte(v.Src))
	hash := h.Sum(nil)
	s := hex.EncodeToString(hash)

	v.dst = filepath.Join(dir, "values", string(name), s+".yml")

	return v
}

// func (v *ValuesReference) Set(dst string) *ValuesReference {
//	v.dst = dst
//	return v
// }

// SetViaRelease downloads and templates values file.
// Returns ErrSkipValues if values cannot be downloaded or doesn't exist in local FS.
func (v *ValuesReference) SetViaRelease(rel Config, dir, templater string) error {
	v.SetUniq(dir, rel.Uniq())

	log.WithFields(log.Fields{
		"release": rel.Uniq(),
		"src":     v.Src,
		"dst":     v.dst,
	}).Trace("Building values reference")

	data := struct {
		Release Config
	}{
		Release: rel,
	}

	if v.isURL() {
		err := v.Download()
		if err != nil {
			log.WithError(err).Warnf("%s skipping: cant download", v.Src)

			return ErrSkipValues
		}

		if err := template.Tpl2yml(v.dst, v.dst, data, templater); err != nil {
			return fmt.Errorf("failed to render %s values: %w", v.Src, err)
		}

		return nil
	} else if !helper.IsExists(v.Src) {
		log.Warnf("%s skipping: local not found", v.Src)

		return ErrSkipValues
	}

	if err := template.Tpl2yml(v.Src, v.dst, data, templater); err != nil {
		return fmt.Errorf("failed to render %s values: %w", v.Src, err)
	}

	return nil
}

func (rel *config) BuildValues(dir, templater string) error {
	for i := len(rel.Values()) - 1; i >= 0; i-- {
		err := rel.Values()[i].SetViaRelease(rel, dir, templater)
		if errors.Is(ErrSkipValues, err) {
			rel.ValuesF = append(rel.ValuesF[:i], rel.ValuesF[i+1:]...)
		} else if err != nil {
			log.WithFields(log.Fields{
				"release": rel.Uniq(),
				"err":     err,
				"values":  rel.Values()[i],
			}).Fatal("Values failed")

			return err
		}
	}

	return nil
}
