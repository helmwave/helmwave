package release

import "helm.sh/helm/v3/pkg/action"

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

	client.DryRun = rel.dryRun
	client.DisableHooks = rel.DisableHooks
	client.Timeout = rel.Timeout
	client.Wait = rel.Wait
	client.Description = rel.Description()

	return client
}
