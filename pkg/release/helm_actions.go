package release

import (
	"fmt"
	"regexp"

	"github.com/helmwave/helmwave/pkg/helper"
	"helm.sh/helm/v4/pkg/action"
	helm "helm.sh/helm/v4/pkg/cli"
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
		rel.helm = helper.NewHelm(rel.Namespace())

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
	client.Namespace = rel.Namespace()
	client.EnableDNS = rel.EnableDNS
	client.TakeOwnership = rel.TakeOwnership
	client.SkipSchemaValidation = rel.SkipSchemaValidation
	client.Labels = rel.Labels
	client.HideSecret = rel.hideSecret

	rel.Chart().CopyOptions(&client.ChartPathOptions)
	rel.applyOCIRegistryClient(client.SetRegistryClient)

	client.DisableHooks = rel.DisableHooks
	client.SkipCRDs = rel.SkipCRDs
	client.Timeout = rel.Timeout
	client.WaitStrategy = rel.WaitStrategy.Helm()
	client.WaitForJobs = rel.WaitForJobs
	client.RollbackOnFailure = rel.RollbackOnFailure
	client.ForceConflicts = rel.ForceConflicts
	// install takes a boolean where upgrade takes true/false/auto, the same way helm's own CLI
	// narrows it when it turns an upgrade into an install.
	client.ServerSideApply = rel.serverSideApply() != "false"
	client.DisableOpenAPIValidation = rel.DisableOpenAPIValidation
	client.SubNotes = rel.SubNotes
	client.Description = rel.Description()

	pr, err := rel.PostRenderer()
	if err != nil {
		rel.Logger().WithError(err).Warn("failed to create post-renderer")
	} else {
		client.PostRenderer = pr
		client.PostRenderStrategy = rel.postRenderStrategy()
	}

	if rel.dryRun {
		client.Replace = true

		// The strategy is what helm actually looks at, and "server" makes it render through a
		// live REST client. That is what we want for a normal build, so `lookup` in a chart sees
		// the cluster. With offline_kube_version the user asked for the opposite, so say
		// "client" — otherwise helm builds a REST config and renders against the cluster.
		if rel.OfflineKubeVersion() != nil {
			client.DryRunStrategy = action.DryRunClient
			client.KubeVersion = rel.OfflineKubeVersion()
		} else {
			client.DryRunStrategy = action.DryRunServer
		}
	}

	return client
}

func (rel *config) newUpgrade() *action.Upgrade {
	client := action.NewUpgrade(rel.Cfg())

	// Only Upgrade
	client.CleanupOnFail = rel.CleanupOnFail
	client.MaxHistory = rel.MaxHistory
	client.ReuseValues = rel.ReuseValues
	client.ResetValues = rel.ResetValues
	client.ResetThenReuseValues = rel.ResetThenReuseValues

	// Common Part
	client.Namespace = rel.Namespace()
	client.EnableDNS = rel.EnableDNS
	client.TakeOwnership = rel.TakeOwnership
	client.SkipSchemaValidation = rel.SkipSchemaValidation
	client.Labels = rel.Labels
	client.HideSecret = rel.hideSecret

	rel.Chart().CopyOptions(&client.ChartPathOptions)
	rel.applyOCIRegistryClient(client.SetRegistryClient)

	client.ForceReplace = rel.ForceReplace
	client.ForceConflicts = rel.ForceConflicts
	client.ServerSideApply = rel.serverSideApply()
	client.DisableHooks = rel.DisableHooks
	client.SkipCRDs = rel.SkipCRDs
	client.Timeout = rel.Timeout
	client.WaitStrategy = rel.WaitStrategy.Helm()
	client.WaitForJobs = rel.WaitForJobs
	client.RollbackOnFailure = rel.RollbackOnFailure
	client.DisableOpenAPIValidation = rel.DisableOpenAPIValidation
	client.SubNotes = rel.SubNotes
	client.Description = rel.Description()

	// A dry run always goes through the install action (see config.upgrade), so this only needs
	// to render without touching the cluster.
	if rel.dryRun {
		client.DryRunStrategy = action.DryRunClient
	}

	pr, err := rel.PostRenderer()
	if err != nil {
		rel.Logger().WithError(err).Warn("failed to create post_renderer")
	} else {
		client.PostRenderer = pr
		client.PostRenderStrategy = rel.postRenderStrategy()
	}

	return client
}

func (rel *config) newUninstall() *action.Uninstall {
	client := action.NewUninstall(rel.Cfg())

	client.KeepHistory = false // TODO: pass it via flags

	client.IgnoreNotFound = true // make it idempotent

	client.DryRun = rel.dryRun
	client.DisableHooks = rel.DisableHooks
	client.Timeout = rel.Timeout
	client.WaitStrategy = rel.WaitStrategy.Helm()
	client.DeletionPropagation = rel.DeletePropagation
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
	client.Timeout = rel.Timeout

	client.DisableHooks = rel.DisableHooks
	client.Timeout = rel.Timeout
	client.WaitStrategy = rel.WaitStrategy.Helm()
	client.WaitForJobs = rel.WaitForJobs
	client.ForceReplace = rel.ForceReplace
	client.ForceConflicts = rel.ForceConflicts
	client.ServerSideApply = rel.serverSideApply()

	return client
}

func (rel *config) newStatus() *action.Status {
	return action.NewStatus(rel.Cfg())
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
