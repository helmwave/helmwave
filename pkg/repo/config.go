package repo

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/repo"
)

type config struct {
	log        *log.Entry `yaml:"-" json:"-"`
	repo.Entry `yaml:",inline" json:",inline"`
	Force      bool `yaml:"force" json:"force" jsonschema:"title=force flag,description=force update helm repo list and download dependencies,default=false"`
}

func (c *config) Name() string {
	return c.Entry.Name
}

func (c *config) URL() string {
	return c.Entry.URL
}

func (c *config) Logger() *log.Entry {
	if c.log == nil {
		c.log = log.WithField("repository", c.Name())
	}

	return c.log
}

// UnmarshalYAML is an unmarshaller for gopkg.in/yaml.v3 to parse YAML into `Config` interface.
func (c *config) UnmarshalYAML(value *yaml.Node) error {
	// Step 1: Convert YAML to map[string]any
	tempMap := make(map[string]any)
	if err := value.Decode(&tempMap); err != nil {
		return err
	}

	// Step 2: Marshal the map to JSON
	jsonData, err := json.Marshal(tempMap)
	if err != nil {
		return err
	}

	// Step 3: Unmarshal JSON into the config struct
	if err := json.Unmarshal(jsonData, c); err != nil {
		return err
	}

	return nil
}

func (c *config) MarshalYAML() (any, error) {
	// Step 1: Marshal the config struct to JSON
	jsonData, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	// Step 2: Unmarshal JSON back into a map[string]any
	tempMap := make(map[string]any)
	if err := json.Unmarshal(jsonData, &tempMap); err != nil {
		return nil, err
	}

	// Step 3: Encode the map into YAML
	return tempMap, nil
}
