package plan

import (
	"errors"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/helmwave/helmwave/pkg/repo"
	"github.com/helmwave/helmwave/pkg/version"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const (
	Dir      = ".helmwave/"
	File     = "planfile"
	Body     = "helmwave.yml"
	Manifest = ".manifest/"
)

var (
	ErrManifestDirNotFound = errors.New(Manifest + " dir not found")
	ErrManifestDirEmpty    = errors.New(Manifest + " is empty")
)

type Plan struct {
	body     *planBody
	dir      string
	fullPath string

	manifests map[uniqname.UniqName]string
}

type planBody struct {
	Project      string
	Version      string
	Repositories []*repo.Config
	Releases     []*release.Config
}

func NewBody(file string) (*planBody, error) {
	b := &planBody{
		Version: version.Version,
	}

	src, err := ioutil.ReadFile(file)
	if err != nil {
		return b, err
	}

	err = yaml.Unmarshal(src, b)
	if err != nil {
		return b, err
	}

	// Setup dev version
	//if b.Version == "" {
	//	b.Version = version.Version
	//}

	return b, err

}

func New(dir string) *Plan {
	if dir[len(dir)-1:] != "/" {
		dir += "/"
	}

	plan := &Plan{
		dir:       dir,
		fullPath:  dir + File,
		manifests: make(map[uniqname.UniqName]string),
	}

	return plan
}
