package plan

import (
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/helmwave/helmwave/pkg/repo"
	"github.com/helmwave/helmwave/pkg/version"
	"gopkg.in/yaml.v2"
)

const (
	Dir      = ".helmwave/"
	File     = "planfile"
	Body     = "helmwave.yml"
	Manifest = "manifest/"
	Values   = "values/"
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

	graphMD string
}

type planBody struct {
	Project      string
	Version      string
	Repositories []*repo.Config
	Releases     []*release.Config
}

func NewBody(file string) (*planBody, error) { // nolint:revive
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
	// if b.Version == "" {
	// 	 b.Version = version.Version
	// }

	return b, err
}

func New(dir string) *Plan {
	//if dir[len(dir)-1:] != "/" {
	//	dir += "/"
	//}

	plan := &Plan{
		dir:       dir,
		fullPath:  filepath.Join(dir, File),
		manifests: make(map[uniqname.UniqName]string),
	}

	return plan
}
