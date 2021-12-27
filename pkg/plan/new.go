package plan

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/helmwave/helmwave/pkg/repo"
	"github.com/helmwave/helmwave/pkg/version"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const (
	// Dir is default directory for generated files.
	Dir = ".helmwave/"

	// File is default file name for planfile.
	File = "planfile"

	// Body is default file name for main config.
	Body = "helmwave.yml"

	// Manifest is default directory under Dir for manifests.
	Manifest = "manifest/"

	// Values is default directory for values.
	Values = "values/"
)

var (
	// ErrManifestDirNotFound is an error for nonexisting manifest dir.
	ErrManifestDirNotFound = errors.New(Manifest + " dir not found")

	// ErrManifestDirEmpty is an error for empty manifest dir.
	ErrManifestDirEmpty = errors.New(Manifest + " is empty")
)

// Plan contains full helmwave state.
type Plan struct {
	body     *planBody
	dir      string
	fullPath string

	tmpDir string

	manifests map[uniqname.UniqName]string

	graphMD string

	templater string
}

type repoConfigs []repo.Config

func (r *repoConfigs) UnmarshalYAML(node *yaml.Node) error {
	if r == nil {
		r = new(repoConfigs)
	}
	var err error

	*r, err = repo.UnmarshalYAML(node)

	return err
}

type releaseConfigs []release.Config

func (r *releaseConfigs) UnmarshalYAML(node *yaml.Node) error {
	if r == nil {
		r = new(releaseConfigs)
	}
	var err error

	*r, err = release.UnmarshalYAML(node)

	return err
}

type planBody struct {
	Project      string
	Version      string
	Repositories repoConfigs
	Releases     releaseConfigs
}

func NewBody(file string) (*planBody, error) { // nolint:revive
	b := &planBody{
		Version: version.Version,
	}

	src, err := os.ReadFile(file)
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

	if err := b.Validate(); err != nil {
		return nil, err
	}

	return b, nil
}

// New returns empty *Plan for provided directory.
func New(dir string) *Plan {
	// if dir[len(dir)-1:] != "/" {
	//	dir += "/"
	// }

	plan := &Plan{
		tmpDir:    os.TempDir(),
		dir:       dir,
		fullPath:  filepath.Join(dir, File),
		manifests: make(map[uniqname.UniqName]string),
	}

	return plan
}

// PrettyPlan logs releases and repositories names.
func (p *Plan) PrettyPlan() {
	a := make([]string, 0, len(p.body.Releases))
	for _, r := range p.body.Releases {
		a = append(a, string(r.Uniq()))
	}

	b := make([]string, 0, len(p.body.Repositories))
	for _, r := range p.body.Repositories {
		b = append(b, r.Name())
	}

	log.WithFields(log.Fields{
		"releases":     a,
		"repositories": b,
	}).Info("🏗 Plan")
}
