package hooks

import (
	"github.com/google/shlex"
	"github.com/helmwave/helmwave/pkg/helper"
	"gopkg.in/yaml.v3"
)

// UnmarshalYAML is an unmarshaller for gopkg.in/yaml.v3 to parse YAML into `Hook` interface.
func (h *Hooks) UnmarshalYAML(node *yaml.Node) error {
	rr := make([]*hook, 0)
	err := node.Decode(&rr)
	if err != nil {
		return NewYAMLDecodeError(err)
	}

	*h = helper.SlicesMap(rr, func(h *hook) Hook { return h })

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
		words, err := shlex.Split(script)
		if err != nil {
			return NewYAMLDecodeError(err)
		}

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
		err = ErrUnknownFormat
	}

	if err != nil {
		return NewYAMLDecodeError(err)
	}

	return nil
}
