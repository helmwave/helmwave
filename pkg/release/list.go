package release

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	helm "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
)

func (rel *Config) List() (*release.Release, error) {
	cfg, err := helper.ActionCfg(rel.Options.Namespace, helm.New())
	if err != nil {
		return nil, err
	}

	client := action.NewList(cfg)
	client.Filter = fmt.Sprintf("^%s$", regexp.QuoteMeta(rel.Name))

	result, err := client.Run()
	if err != nil {
		return nil, err
	}
	switch len(result) {
	case 0:
		return nil, errors.New("release not found")
	case 1:
		return result[0], nil
	default:
		return nil, errors.New("found multiple releases o_0")
	}
}

func List(releases []*Config) error {
	log.Infof("Should be %d releases", len(releases))
	if len(releases) == 0 {
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"name", "namespace", "revision", "updated", "status", "chart", "version"})
	table.SetAutoFormatHeaders(true)
	table.SetAutoMergeCellsByColumnIndex([]int{1, 4})
	table.SetBorder(false)

	for _, rel := range releases {
		r, err := rel.List()
		if err != nil {
			log.Errorf("Failed to list %s release, skipping: %v", rel.UniqName(), err)
			continue
		}

		status := r.Info.Status
		statusColor := tablewriter.Colors{tablewriter.Normal, tablewriter.FgGreenColor}
		if status != release.StatusDeployed {
			statusColor = tablewriter.Color(tablewriter.Bold, tablewriter.BgRedColor)
		}
		r.Chart.Name()

		table.Rich([]string{
			r.Name,
			r.Namespace,
			strconv.Itoa(r.Version),
			r.Info.LastDeployed.String(),
			r.Info.Status.String(),
			r.Chart.Name(),
			r.Chart.Metadata.Version,
		}, []tablewriter.Colors{
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
