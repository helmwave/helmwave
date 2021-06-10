package plan

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func (p *Plan) Import() error {
	src, err := ioutil.ReadFile(p.fullPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(src, p.body)
	if err != nil {
		return err
	}

	return nil
}

// Validate that files is existing
func (p *Plan) Validate() error {
	return nil
}
