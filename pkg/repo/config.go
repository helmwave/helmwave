package repo

import (
	"errors"
	"github.com/invopop/jsonschema"
	"helm.sh/helm/v3/pkg/repo"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// Configs type of array Config.
type Configs []Config

// UnmarshalYAML parse Config.
func (r *Configs) UnmarshalYAML(node *yaml.Node) error {
	var err error
	*r, err = UnmarshalYAML(node)

	return err
}

type Repository struct {
	log   *log.Entry
	entry *repo.Entry

	NameF                 string `json:"name"`
	URLF                  string `json:"url"`
	Username              string `json:"username,omitempty"`
	Password              string `json:"password,omitempty"`
	CertFile              string `json:"cert_file,omitempty"`
	KeyFile               string `json:"key_file,omitempty"`
	CAFile                string `json:"ca_file,omitempty"`
	InsecureSkipTLSverify bool   `json:"insecure_skip_tls_verify,omitempty"`
	PassCredentialsAll    bool   `json:"pass_credentials_all,omitempty"`
	Force                 bool   `json:"force,omitempty"`
}

func (c *Repository) Name() string {
	return c.NameF
}

func (c *Repository) URL() string {
	return c.URLF
}

func (c *Repository) Logger() *log.Entry {
	if c.log == nil {
		c.log = log.WithField("repository", c.Name())
	}

	return c.log
}

func (c *Repository) JSONSchema() *jsonschema.Schema {
	return jsonschema.Reflect(c)
}

// ErrNotFound is an error for not declared repository name.
var ErrNotFound = errors.New("repository not found")
