package plan

import (
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"os"
)

func (p *Plan) List() error {
	log.Infof("Should be %d releases", len(p.body.Releases))
	if len(p.body.Releases) == 0 {
		return nil
	}

	return nil
}


func newListTable() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"name", "namespace", "revision", "updated", "status", "chart", "version"})
	table.SetAutoFormatHeaders(true)
	table.SetAutoMergeCellsByColumnIndex([]int{1, 4})
	table.SetBorder(false)

	return table
}

