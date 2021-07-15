package release

import (
	"github.com/helmwave/helmwave/pkg/helper"
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

//func (v *ValuesReference) SetViaRelease(rel *Config, dir string) error {
//	dst := dir + ".values/" + string(rel.Uniq()) + "/" + string(rune(i)) + ".yml"
//	v.Set(dst)
//}
