package release

import (
	"fmt"
	"regexp"

	"github.com/helmwave/helmwave/pkg/helper"
	"helm.sh/helm/v3/pkg/action"
	helm "helm.sh/helm/v3/pkg/cli"
)

func (rel *config) Cfg() *action.Configuration {
	cfg, err := helper.NewCfg(rel.Namespace(), rel.KubeContext())
	if err != nil {
		rel.Logger().Fatal(err)

		return nil
	}

	return cfg
}

func (rel *config) Helm() *helm.EnvSettings {
	if rel.helm == nil {
		var err error
		rel.helm, err = helper.NewHelm(rel.Namespace())
		if err != nil {
			rel.Logger().Fatal(err)

			return nil
		}

		rel.helm.Debug = helper.Helm.Debug
	}

	return rel.helm
}

func (rel *config) newInstall() *action.Install {
	client := action.NewInstall(rel.Cfg())

	// client.IncludeCRDs = true

	// Only Up
	client.CreateNamespace = rel.CreateNamespace
	client.ReleaseName = rel.Name()

	// Common Part
	client.DryRun = rel.dryRun
	client.Namespace = rel.Namespace()
	client.EnableDNS = rel.EnableDNS
	client.Labels = rel.Labels

	rel.Chart().CopyOptions(&client.ChartPathOptions)

	client.DisableHooks = rel.DisableHooks
	client.SkipCRDs = rel.SkipCRDs
	client.Timeout = rel.Timeout
	client.Wait = rel.Wait
	client.WaitForJobs = rel.WaitForJobs
	client.Atomic = rel.Atomic
	client.DisableOpenAPIValidation = rel.DisableOpenAPIValidation
	client.SubNotes = rel.SubNotes
	client.Description = rel.Description()

	pr, err := rel.PostRenderer()
	if err != nil {
		rel.Logger().WithError(err).Warn("failed to create post-renderer")
	} else {
		client.PostRenderer = pr
	}

	if client.DryRun {
		client.Replace = true
	}

	if client.DryRun && nil != rel.OfflineKubeVersion() {
		client.ClientOnly = true
		client.KubeVersion = rel.OfflineKubeVersion()
	}

	return client
}

func (rel *config) newUpgrade() *action.Upgrade {
	client := action.NewUpgrade(rel.Cfg())

	// Only Upgrade
	client.CleanupOnFail = rel.CleanupOnFail
	client.MaxHistory = rel.MaxHistory
	client.Recreate = rel.Recreate
	client.ReuseValues = rel.ReuseValues
	client.ResetValues = rel.ResetValues
	client.ResetThenReuseValues = rel.ResetThenReuseValues

	// Common Part
	client.DryRun = rel.dryRun
	client.Namespace = rel.Namespace()
	client.EnableDNS = rel.EnableDNS
	client.Labels = rel.Labels

	rel.Chart().CopyOptions(&client.ChartPathOptions)

	client.Force = rel.Force
	client.DisableHooks = rel.DisableHooks
	client.SkipCRDs = rel.SkipCRDs
	client.Timeout = rel.Timeout
	client.Wait = rel.Wait
	client.WaitForJobs = rel.WaitForJobs
	client.Atomic = rel.Atomic
	client.DisableOpenAPIValidation = rel.DisableOpenAPIValidation
	client.SubNotes = rel.SubNotes
	client.Description = rel.Description()

	pr, err := rel.PostRenderer()
	if err != nil {
		rel.Logger().WithError(err).Warn("failed to create post_renderer")
	} else {
		client.PostRenderer = pr
	}

	return client
}

func (rel *config) newUninstall() *action.Uninstall {
	client := action.NewUninstall(rel.Cfg())

	client.KeepHistory = false                // TODO: pass it via flags
	client.DeletionPropagation = "background" // TODO: pass it via flags

	client.IgnoreNotFound = true // make it idempotent

	client.DryRun = rel.dryRun
	client.DisableHooks = rel.DisableHooks
	client.Timeout = rel.Timeout
	client.Wait = rel.Wait
	client.Description = rel.Description()

	return client
}

func (rel *config) newTest() *action.ReleaseTesting {
	client := action.NewReleaseTesting(rel.Cfg())

	client.Timeout = rel.Timeout
	client.Namespace = rel.Namespace()
	client.Filters = rel.Tests.Filters

	return client
}

func (rel *config) newRollback() *action.Rollback {
	client := action.NewRollback(rel.Cfg())

	client.CleanupOnFail = rel.CleanupOnFail
	client.MaxHistory = rel.MaxHistory
	client.Recreate = rel.Recreate
	client.Timeout = rel.Timeout

	client.DisableHooks = rel.DisableHooks
	client.Timeout = rel.Timeout
	client.Wait = rel.Wait
	client.WaitForJobs = rel.WaitForJobs
	client.Force = rel.Force

	return client
}

func (rel *config) newStatus() *action.Status {
	client := action.NewStatus(rel.Cfg())

	client.ShowDescription = true

	return client
}

func (rel *config) newGet() *action.Get {
	client := action.NewGet(rel.Cfg())

	return client
}

func (rel *config) newGetValues() *action.GetValues {
	client := action.NewGetValues(rel.Cfg())

	return client
}

func (rel *config) newList() *action.List {
	client := action.NewList(rel.Cfg())

	client.Filter = fmt.Sprintf("^%s$", regexp.QuoteMeta(rel.Name()))

	return client
}

func (rel *config) newHistory() *action.History {
	client := action.NewHistory(rel.Cfg())

	return client
}
