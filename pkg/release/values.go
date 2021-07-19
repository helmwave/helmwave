package release

import (
	"crypto/sha1"
	"encoding/base64"
	"os"

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

func (v *ValuesReference) Set(dst string) error {
	v.dst = dst

	if v.isURL() {
		err := v.Download()
		if err != nil {
			return err
		}
	}

	return nil
}

func (v *ValuesReference) SetViaRelease(rel *Config, dir string) error {
	h := sha1.New() // nolint:gosec
	h.Write([]byte(v.Src))
	sha := base64.URLEncoding.EncodeToString(h.Sum(nil))

	// Todo: fmt.Sprintf
	dst := dir + "values/" + string(rel.Uniq()) + "/" + sha + ".yml"

	log.WithFields(log.Fields{
		"release": rel.Uniq(),
		"src":     v.Src,
		"dst":     dst,
	}).Debug("Building values reference")

	err := v.Set(dst)
	if err != nil {
		log.Warn(v.Src, " skipping: ", err)
		return nil
	}

	return template.Tpl2yml(dst, dst, struct{ Release *Config }{rel})
}

func (rel *Config) BuildValues(dir string) error {
	wg := parallel.NewWaitGroup()
	wg.Add(len(rel.Values))

	for i := range rel.Values {
		go func(wg *parallel.WaitGroup, i int) {
			defer wg.Done()
			err := rel.Values[i].SetViaRelease(rel, dir)
			if err != nil {
				log.Fatal(err)
			}

			// log.WithField("values", rel.Values).Info(rel.Uniq(), " values are ok ")
		}(wg, i)
	}

	return wg.Wait()
}
