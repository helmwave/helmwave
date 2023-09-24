package plan

import (
	"io/fs"
	"path"

	"github.com/helmwave/go-fsimpl"
)

func (p *Plan) exportCharts(srcFS fsimpl.CurrentPathFS, plandirFS fsimpl.WriteableFS) error {
	for i, rel := range p.body.Releases {
		l := p.Logger().WithField("release", rel.Uniq())

		if !rel.Chart().IsRemote(plandirFS) {
			l.Info("chart is local, skipping exporting it")

			continue
		}

		dst := path.Join(Charts, rel.Uniq().String())
		err := rel.DownloadChart(srcFS, plandirFS, dst)
		if err != nil {
			return err
		}

		// Chart is places as an archive under this directory.
		// So we need to find it and use.
		entries, err := plandirFS.(fs.ReadDirFS).ReadDir(dst)
		if err != nil {
			l.WithError(err).Warn("failed to read directory with downloaded chart, skipping")

			continue
		}

		if len(entries) != 1 {
			l.WithField("entries", entries).Warn("don't know which file is downloaded chart, skipping")

			continue
		}

		chart := entries[0]
		p.body.Releases[i].SetChartName(path.Join(dst, chart.Name()))
	}

	return nil
}
