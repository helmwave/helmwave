package plan

import (
	"go.beyondstorage.io/v5/services"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"

	"github.com/helmwave/helmwave/pkg/release/uniqname"
	log "github.com/sirupsen/logrus"
	_ "go.beyondstorage.io/services/fs/v4"
	storage "go.beyondstorage.io/v5/types"
)

// Plan contains full helmwave state.
type Plan struct {
	body *planBody

	fsys   fs.FS
	store  storage.Storager
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

// New create Plan
// src is can be:
// fs://./<workdir>
// fs://<workdir>
// fs:///<workdir>
func New(plandir string) (p *Plan, err error) {
	URL, err := url.Parse(plandir)
	if err != nil {
		plandir = LocalScheme + plandir
		URL, err = url.Parse(plandir)
		if err != nil {
			return nil, err
		}
	}

	store, err := services.NewStoragerFromString(plandir)
	if err != nil {
		return nil, err
	}

	return &Plan{
		store:     store,
		url:       URL,
		tmpDir:    os.TempDir(),
		manifests: make(map[uniqname.UniqName]string),
	}, nil
}

// File is path to planfile.
func (p *Plan) File() string {
	return filepath.Join(p.Dir(), File)
}

// GraphPath is path to graph.md.
func (p *Plan) GraphPath() string {
	return filepath.Join(p.Dir(), GraphFilename)
}

// Dir is path to plandir.
func (p *Plan) Dir() string {
	return filepath.Dir(filepath.Join(p.url.Host, p.url.Path))
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
