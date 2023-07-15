package hooks

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// UnmarshalYAML is an unmarshaller for gopkg.in/yaml.v3 to parse YAML into `Hook` interface.
func (h *Hooks) UnmarshalYAML(node *yaml.Node) error {
	rr := make([]*hook, 0)
	err := node.Decode(&rr)
	if err != nil {
		return fmt.Errorf("failed to decode lifecycle config from YAML: %w", err)
	}

	*h = make([]Hook, len(rr))
	for i := range rr {
		(*h)[i] = rr[i]
	}

	return nil
}

func (h *hook) UnmarshalYAML(node *yaml.Node) error {
	type raw hook
	var err error

	// show by default
	h.Show = true

	switch node.Kind {
	// single value or reference to another value
	case yaml.ScalarNode, yaml.AliasNode:
		var script string
		err = node.Decode(&script)
		if err != nil {
			break
		}

		// Short name
		words := strings.Fields(script)
		if len(words) > 1 {
			h.Cmd = words[0]
			h.Args = words[1:]
		} else {
			h.Cmd = script
			h.Args = []string{}
		}

	case yaml.MappingNode:
		err = node.Decode((*raw)(h))
	default:
		err = fmt.Errorf("unknown format")
	}

	if err != nil {
		return fmt.Errorf("failed to decode values reference %q from YAML: %w", node.Value, err)
	}

	return nil
}
