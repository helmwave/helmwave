package plan

import (
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/release"
)

//nolint:gochecknoglobals // can't make these colors const
var (
	// FailStatusColor is tablewriter color for failed releases.
	FailStatusColor = tablewriter.Color(tablewriter.Bold, tablewriter.BgRedColor)

	// SuccessStatusColor is tablewriter color for succeeded releases.
	SuccessStatusColor = tablewriter.Color(tablewriter.Bold, tablewriter.BgRedColor)
)

// List renders releases status table and writes it to stdout.
// Listed only releases that exist in plan.
func (p *Plan) List() error {
	log.Infof("Should be %d releases", len(p.body.Releases))
	if len(p.body.Releases) == 0 {
		return nil
	}

	table := newListTable()

	for _, rel := range p.body.Releases {
		r, err := rel.List()
		if err != nil {
			rel.Logger().WithError(err).Error("Failed to list. Skipping.")

			continue
		}

		status := r.Info.Status

		statusColor := SuccessStatusColor
		if status != release.StatusDeployed {
			statusColor = FailStatusColor
		}

		row := []string{
			r.Name,
			r.Namespace,
			strconv.Itoa(r.Version),
			r.Info.LastDeployed.String(),
			string(r.Info.Status),
			r.Chart.Name(),
			r.Chart.Metadata.Version,
		}

		table.Rich(row, []tablewriter.Colors{
			{},
			{},
			{},
			{},
			statusColor,
			{},
			{},
		})
	}

	table.Render()

	return nil
}

func newListTable() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"name", "namespace", "revision", "updated", "status", "chart", "version"})
	table.SetAutoFormatHeaders(true)
	table.SetBorder(false)

	return table
}
