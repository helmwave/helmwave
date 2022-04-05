package plan

import (
	"errors"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"

	"github.com/hairyhenderson/go-fsimpl"
	"github.com/hairyhenderson/go-fsimpl/blobfs"
	"github.com/hairyhenderson/go-fsimpl/filefs"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var (
	// ErrManifestDirNotFound is an error for nonexistent manifest dir.
	ErrManifestDirNotFound = errors.New(Manifest + " dir not found")

	// ErrManifestDirEmpty is an error for empty manifest dir.
	ErrManifestDirEmpty = errors.New(Manifest + " is empty")
)

// Plan contains full helmwave state.
type Plan struct {
	body *planBody

	fsys   fs.FS
	url    *url.URL
	tmpDir string

	manifests map[uniqname.UniqName]string

	graphMD string

	templater string
}

// NewAndImport wrapper for New and Import in one
func NewAndImport(src string) (p *Plan, err error) {
	p, err = New(src)
	if err != nil {
		return nil, err
	}

	if err = p.Import(); err != nil {
		return p, err
	}

	return p, nil
}

func New(src string) (*Plan, error) {

	// Allowed FS
	mux := fsimpl.NewMux()
	mux.Add(filefs.FS)
	mux.Add(blobfs.FS)

	// Looking for FS
	fsys, err := mux.Lookup(src)
	if err != nil {
		src = "file://" + src
		fsys, err = mux.Lookup(src)
		if err != nil {
			return nil, err
		}
	}

	URL, _ := url.Parse(src)

	return &Plan{
		fsys:      fsys,
		url:       URL,
		tmpDir:    os.TempDir(),
		manifests: make(map[uniqname.UniqName]string),
	}, nil
}

func (p *Plan) File() string {
	return filepath.Join(p.Dir(), File)
}

func (p *Plan) GraphPath() string {
	return filepath.Join(p.Dir(), GraphFilename)
}

func (p *Plan) Dir() string {
	return p.url.Path
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

	c := make([]string, 0, len(p.body.Registries))
	for _, r := range p.body.Registries {
		c = append(c, r.Host())
	}

	log.WithFields(log.Fields{
		"releases":     a,
		"repositories": b,
		"registries":   c,
	}).Info("🏗 Plan")
}
