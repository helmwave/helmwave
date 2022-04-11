package registry

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// Config is an interface to manage particular helm reg.
type Config interface {
	Install() error
	Logger() *log.Entry
	Host() string
	// Username() string
	// Password() string
}

// UnmarshalYAML is an unmarshaller for gopkg.in/yaml.v3 to parse YAML into `Config` interface.
func UnmarshalYAML(node *yaml.Node) ([]Config, error) {
	r := make([]*config, 0)
	if err := node.Decode(&r); err != nil {
		return nil, fmt.Errorf("failed to decode reg config from YAML: %w", err)
	}

	res := make([]Config, len(r))
	for i := range r {
		res[i] = r[i]
	}

	return res, nil
}
