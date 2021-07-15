package release

import (
	"crypto/sha1"
	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/template"
	log "github.com/sirupsen/logrus"
	"os"
)

type ValuesReference struct {
	Src string
	dst string
}

//func (v *ValuesReference) UnmarshalYAML(unmarshal func(interface{}) error) error {
//	return unmarshal(&v.Src)
//}

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
	return helper.IsUrl(v.Src)
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
	h := sha1.New()
	h.Write([]byte(v.Src))
	bs := h.Sum(nil)

	// Todo: fmt.Sprintf
	dst := dir + ".values/" + string(rel.Uniq()) + "/" + string(bs) + ".yml"

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

	for _, v := range rel.Values {
		go func(wg *parallel.WaitGroup, v ValuesReference) {
			defer wg.Done()
			err := v.SetViaRelease(rel, dir)
			if err != nil {
				log.Fatal(err)
			}

		}(wg, v)
	}

	return wg.Wait()
}
