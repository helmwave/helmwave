package release

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/helmwave/helmwave/pkg/template"
	log "github.com/sirupsen/logrus"
)

var ErrSkipValues = errors.New("values has been skip")

type ValuesReference struct {
	Src string
	dst string
}

func (v *ValuesReference) UnmarshalYAML(unmarshal func(interface{}) error) error {
	m := make(map[string]string)
	if err := unmarshal(&m); err != nil {
		return unmarshal(&v.Src)
	}

	v.Src = m["src"]
	v.dst = m["dst"]

	return nil
}

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

func (v *ValuesReference) IsLocal() bool {
	stat, err := os.Stat(v.Src)
	return err == nil && !stat.IsDir()
}

func (v *ValuesReference) Download() error {
	return helper.Download(v.dst, v.Src)
}

func (v *ValuesReference) Get() string {
	return v.dst
}

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

func (v *ValuesReference) SetViaRelease(rel *Config, dir string) error {
	v.SetUniq(dir, rel.Uniq())

	log.WithFields(log.Fields{
		"release": rel.Uniq(),
		"src":     v.Src,
		"dst":     v.dst,
	}).Trace("Building values reference")

	if v.isURL() {
		err := v.Download()
		if err != nil {
			log.Warnf("%s skipping: cant download %v", v.Src, err)
			return ErrSkipValues
		}
		return template.Tpl2yml(v.dst, v.dst, struct{ Release *Config }{rel})
	} else if !helper.IsExists(v.Src) {
		log.Warnf("%s skipping: local not found", v.Src)
		return ErrSkipValues
	}

	return template.Tpl2yml(v.Src, v.dst, struct{ Release *Config }{rel})
}

func (rel *Config) BuildValues(dir string) error {
	for i := len(rel.Values) - 1; i >= 0; i-- {
		err := rel.Values[i].SetViaRelease(rel, dir)
		if errors.Is(ErrSkipValues, err) {
			rel.Values = append(rel.Values[:i], rel.Values[i+1:]...)
		} else if err != nil {
			log.WithFields(log.Fields{
				"release": rel.Uniq(),
				"err":     err,
				"values":  rel.Values[i],
			}).Fatal("Values failed")
			return err
		}
	}

	return nil
}
