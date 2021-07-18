package release

import (
	"testing"

	"gopkg.in/yaml.v2"
)

func TestValuesReferenceParse(t *testing.T) {
	type config struct {
		Values []ValuesReference
	}

	src := `
values:
- a
- b
`
	c := &config{}

	err := yaml.Unmarshal([]byte(src), c)
	if err != nil {
		t.Error(err)
	}

	if "a" != c.Values[0].Src || "b" != c.Values[1].Src {
		t.Log(c.Values)
		t.Error("error parsed ValuesReference List")
	}

	src = `
values:
- src: 1
  dst: a
- src: 2
  dst: b
`

	err = yaml.Unmarshal([]byte(src), c)
	if err != nil {
		t.Error(err)
	}

	if "1" != c.Values[0].Src || "a" != c.Values[0].dst || "2" != c.Values[1].Src || "b" != c.Values[1].dst {
		t.Log(c.Values)
		t.Error("error parsed ValuesReference MAP")
	}
}
