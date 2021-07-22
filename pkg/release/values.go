package release

import (
	"crypto/sha1"
	"encoding/hex"
	"os"
	"path/filepath"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/template"
	log "github.com/sirupsen/logrus"
)

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

func (v *ValuesReference) SetViaRelease(rel *Config, dir string) error {
	h := sha1.New() // nolint:gosec
	h.Write([]byte(v.Src))
	hash := h.Sum(nil)
	hs := hex.EncodeToString(hash)
	// b64 := base64.URLEncoding.EncodeToString(hash)

	v.dst = filepath.Join(dir, "values", string(rel.Uniq()), hs+".yml")

	log.WithFields(log.Fields{
		"release": rel.Uniq(),
		"src":     v.Src,
		"dst":     v.dst,
	}).Trace("Building values reference")

	if v.isURL() {
		err := v.Download()
		if err != nil {
			log.Warn(v.Src, "skipping: cant download ", err)
			return nil
		}
		return template.Tpl2yml(v.dst, v.dst, struct{ Release *Config }{rel})
	} else if !helper.IsExists(v.Src) {
		log.Warn(v.Src, "skipping: local not found")
		return nil
	}

	return template.Tpl2yml(v.Src, v.dst, struct{ Release *Config }{rel})
}

func (rel *Config) BuildValues(dir string) error {
	wg := parallel.NewWaitGroup()
	wg.Add(len(rel.Values))

	for i := range rel.Values {
		go func(wg *parallel.WaitGroup, i int) {
			defer wg.Done()
			err := rel.Values[i].SetViaRelease(rel, dir)
			if err != nil {
				log.WithFields(log.Fields{
					"release": rel.Uniq(),
					"err":     err,
					"values":  rel.Values[i],
				}).Fatal("Values failed")
			}

			// log.WithField("values", rel.Values).Info(rel.Uniq(), " values are ok ")
			log.Info(rel.Uniq(), " values are ok ")
		}(wg, i)
	}

	return wg.Wait()
}
