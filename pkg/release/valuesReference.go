package release

import (
	"context"
	"crypto"
	_ "crypto/sha1"
	"fmt"
	"github.com/hashicorp/go-getter/v2"
	"os"
)

type ValuesReference struct {
	srcURI    string
	localPath string
}

func (v *ValuesReference) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return unmarshal(&v.srcURI)
}

func (v *ValuesReference) GetPath() string {
	if v.localPath != "" {
		return v.localPath
	}
	return v.srcURI
}

func (v *ValuesReference) ManifestPath() (string, error) {
	s := v.GetPath()
	if v.IsLocal() {
		return s, nil
	}

	hasher := crypto.SHA1.New()
	_, err := hasher.Write([]byte(s))
	return fmt.Sprintf("values/%x", hasher.Sum(nil)), err
}

func (v *ValuesReference) IsLocal() bool {
	stat, err := os.Stat(v.srcURI)
	return err == nil && !stat.IsDir()
}

func (v *ValuesReference) Download() error {
	if v.IsLocal() {
		v.SetProcessedPath(v.srcURI)
		return nil
	}

	f, err := os.CreateTemp(os.TempDir(), "helmwave-*")
	if err != nil {
		return err
	}
	tmpPath := f.Name()

	err = f.Close()
	if err != nil {
		return err
	}

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	req := &getter.Request{
		Src:     v.srcURI,
		Dst:     tmpPath,
		Pwd:     pwd,
		GetMode: getter.ModeFile,
	}
	_, err = getter.DefaultClient.Get(context.TODO(), req)
	if err == nil {
		v.SetProcessedPath(tmpPath)
	}

	return err
}

func (v *ValuesReference) SetProcessedPath(path string) {
	v.localPath = path
}

func (v *ValuesReference) UnlinkProcessed() {
	if !v.IsLocal() && v.localPath != "" {
		os.Remove(v.localPath)
	}
}
