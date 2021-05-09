package release

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/helmwave/helmwave/pkg/feature"
	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	helm "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
	"os"
	"regexp"
	"sort"
	"strconv"
)

var emptyReleases = errors.New("releases are empty")

func (rel *Config) Sync(manifestPath string) error {
	log.Infof("🛥 %s", rel.UniqName())

	if err := rel.waitForDependencies(); err != nil {
		return err
	}

	// I hate Private
	_ = os.Setenv("HELM_NAMESPACE", rel.Options.Namespace)
	settings := helm.New()
	cfg, err := helper.ActionCfg(rel.Options.Namespace, settings)
	if err != nil {
		return err
	}

	install, err := rel.Install(cfg, settings)
	if install != nil {
		log.Trace(install.Manifest)
	}

	if err != nil {
		return err
	}

	m := manifestPath + rel.UniqName() + ".yml"
	f, err := helper.CreateFile(m)
	if err != nil {
		return err
	}
	_, err = f.WriteString(install.Manifest)
	if err != nil {
		return err
	}

	return f.Close()
}

func (rel *Config) SyncWithFails(fails *[]*Config, manifestPath string) {
	err := rel.Sync(manifestPath)
	if err != nil {
		log.Error("❌ ", err)
		rel.NotifyFailed()
		*fails = append(*fails, rel)
	} else {
		rel.NotifySuccess()
		log.Infof("✅ %s", rel.UniqName())
	}
}

func (rel *Config) Status() (*release.Release, error) {
	cfg, err := helper.ActionCfg(rel.Options.Namespace, helm.New())
	if err != nil {
		return nil, err
	}

	client := action.NewStatus(cfg)
	client.ShowDescription = true

	return client.Run(rel.Name)
}

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

func Sync(releases []*Config, manifestPath string) (err error) {
	if len(releases) == 0 {
		return emptyReleases
	}

	log.Info("🛥 Sync releases")
	var fails []*Config

	if feature.Parallel {
		wg := parallel.NewWaitGroup()
		log.Debug("🐞 Run in parallel mode")
		wg.Add(len(releases))
		for i := range releases {
			go func(wg *parallel.WaitGroup, release *Config, fails *[]*Config, manifestPath string) {
				defer wg.Done()
				release.SyncWithFails(fails, manifestPath)
			}(wg, releases[i], &fails, manifestPath)
		}
		err := wg.Wait()
		if err != nil {
			return err
		}
	} else {
		for _, r := range releases {
			r.SyncWithFails(&fails, manifestPath)
		}
	}

	n := len(releases)
	k := len(fails)

	log.Infof("Success %d / %d", n-k, n)
	if k > 0 {
		for _, rel := range fails {
			log.Errorf("%q was not deployed to %q", rel.Name, rel.Options.Namespace)
		}

		return errors.New("deploy failed")
	}
	return nil
}

func Status(allReleases []*Config, releasesNames []string) error {
	r := allReleases

	if len(releasesNames) > 0 {
		sort.Strings(releasesNames)
		r = make([]*Config, 0, len(allReleases))

		for _, rel := range allReleases {
			if helper.Contains(rel.UniqName(), releasesNames) {
				r = append(r, rel)
			}
		}
	}

	for _, rel := range r {
		status, err := rel.Status()
		if err != nil {
			log.Errorf("Failed to get status of %s: %v", rel.UniqName(), err)
			continue
		}

		labels, _ := json.Marshal(status.Labels)
		values, _ := json.Marshal(status.Config)

		log.WithFields(log.Fields{
			"name":          status.Name,
			"namespace":     status.Namespace,
			"chart":         fmt.Sprintf("%s-%s", status.Chart.Name(), status.Chart.Metadata.Version),
			"last deployed": status.Info.LastDeployed,
			"status":        status.Info.Status,
			"revision":      status.Version,
		}).Infof("General status of %s", rel.UniqName())

		log.WithFields(log.Fields{
			"notes":         status.Info.Notes,
			"labels":        string(labels),
			"chart sources": status.Chart.Metadata.Sources,
			"values":        string(values),
		}).Debugf("Debug status of %s", rel.UniqName())

		log.WithFields(log.Fields{
			"hooks":    status.Hooks,
			"manifest": status.Manifest,
		}).Tracef("Superdebug status of %s", rel.UniqName())
	}

	return nil
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
