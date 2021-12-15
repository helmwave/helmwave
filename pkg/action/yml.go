package action

import (
	"github.com/helmwave/helmwave/pkg/template"
	log "github.com/sirupsen/logrus"
)

type Yml struct {
	tpl, file string
}

func (i *Yml) Run() error {
	err := template.Tpl2yml(i.tpl, i.file, nil, &template.GomplateConfig{Enabled: false})
	if err != nil {
		return err
	}

	log.WithField(
		"build plan with next command",
		"helmwave build -f "+i.file,
	).Info("📄 YML is ready!")

	return nil
}
