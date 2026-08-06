package release

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/helmwave/helmwave/pkg/helper"
	"helm.sh/helm/v4/pkg/action"
	chart "helm.sh/helm/v4/pkg/chart/v2"
	"helm.sh/helm/v4/pkg/cli/values"
	"helm.sh/helm/v4/pkg/getter"
	release "helm.sh/helm/v4/pkg/release/v1"
)

// Helm wraps a lot of meta.NoKindMatchError into fmt.Errorf which makes errors.Is unusable.
// So we have to find this substring in error string.
const errMissingCRD = "unable to build kubernetes objects from release manifest:"

func (rel *config) upgrade(ctx context.Context) (*release.Release, error) {
	ch, err := rel.GetChart()
	if err != nil {
		return nil, err
	}

	// Values
	valuesFiles := helper.SlicesMap(rel.Values(), func(v ValuesReference) string {
		return v.Dst
	})

	valOpts := &values.Options{ValueFiles: valuesFiles}
	vals, err := valOpts.MergeValues(getter.All(rel.Helm()))
	if err != nil {
		return nil, fmt.Errorf("failed to merge values %v: %w", valuesFiles, err)
	}

	// Install or Template
	if rel.dryRun {
		rel.Logger().Debug("I'll dry-run.")
		r, err := rel.installWithRetry(ctx, ch, vals)
		if err != nil {
			return nil, fmt.Errorf("failed with dry-run %q: %w", rel.Uniq(), err)
		}

		return r, nil
	} else if !rel.dryRun && !rel.isInstalled() {
		rel.Logger().Debug("🧐 Release does not exist. Installing it now.")
		r, err := rel.installWithRetry(ctx, ch, vals)
		if err != nil {
			return nil, fmt.Errorf("failed to install %q: %w", rel.Uniq(), err)
		}

		return r, nil
	}

	pending, err := rel.isPending()
	if err != nil {
		return nil, fmt.Errorf("failed to check %q for pending status: %w", rel.Uniq(), err)
	}
	if pending {
		err := rel.fixPending(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to fix %q pending status: %w", rel.Uniq(), err)
		}
	}

	// Upgrade
	r, err := rel.upgradeWithRetry(ctx, ch, vals)
	if err != nil {
		return nil, fmt.Errorf("failed to upgrade %s: %w", rel.Uniq(), err)
	}

	return r, nil
}

//nolint:wrapcheck // we wrap it later
func (rel *config) installWithRetry(
	ctx context.Context,
	ch *chart.Chart,
	vals map[string]any,
) (*release.Release, error) {
	client := rel.newInstall()
	rel.adjustCreateNamespace(ctx, client)

	r, err := client.RunWithContext(ctx, ch, vals)

	if err != nil && strings.Contains(err.Error(), errMissingCRD) && rel.dryRun {
		er := rel.forceOfflineKubeVersion()
		// return original error if we can't get kubernetes version
		if er != nil {
			return unwrapRelease(r, err)
		}

		client = rel.newInstall()
		rel.adjustCreateNamespace(ctx, client)

		return unwrapRelease(client.RunWithContext(ctx, ch, vals))
	}

	return unwrapRelease(r, err)
}

// adjustCreateNamespace implements a get-first namespace creation strategy.
//
// Helm's CreateNamespace issues an unconditional `create`, which the Kubernetes API
// rejects at the authorization layer before the object is checked for existence. As a
// result, an identity that may only `get/list/watch` namespaces (but is admin inside
// its own, already provisioned namespace) fails with a 403 even though the namespace
// already exists.
//
// To avoid requiring cluster-scoped `create` on namespaces for that common case, we
// first check whether the namespace exists and, if so, disable helm's creation. When
// the namespace is missing, or the existence check cannot be performed, we keep the
// original behavior and let helm create it.
func (rel *config) adjustCreateNamespace(ctx context.Context, client *action.Install) {
	if !client.CreateNamespace || rel.dryRun {
		return
	}

	exists, err := rel.namespaceExists(ctx)
	if err != nil {
		rel.Logger().WithError(err).Warnf(
			"failed to check if namespace %q exists, will attempt to create it",
			rel.Namespace(),
		)

		return
	}

	if exists {
		rel.Logger().Debugf("namespace %q already exists, skipping namespace creation", rel.Namespace())
		client.CreateNamespace = false
	}
}

func (rel *config) namespaceExists(ctx context.Context) (bool, error) {
	clientSet, err := rel.Cfg().KubernetesClientSet()
	if err != nil {
		return false, fmt.Errorf("failed to create kubernetes clientset: %w", err)
	}

	return helper.NamespaceExists(ctx, clientSet, rel.Namespace())
}

//nolint:wrapcheck // we wrap it later
func (rel *config) upgradeWithRetry(
	ctx context.Context,
	ch *chart.Chart,
	vals map[string]any,
) (*release.Release, error) {
	r, err := rel.newUpgrade().RunWithContext(ctx, rel.Name(), ch, vals)

	if err != nil && strings.Contains(err.Error(), errMissingCRD) && rel.dryRun {
		er := rel.forceOfflineKubeVersion()
		// return original error if we can't get kubernetes version
		if er != nil {
			return unwrapRelease(r, err)
		}

		return unwrapRelease(rel.newUpgrade().RunWithContext(ctx, rel.Name(), ch, vals))
	}

	return unwrapRelease(r, err)
}

func (rel *config) forceOfflineKubeVersion() error {
	rel.Logger().Warn("🤔hmm, it looks like some required CRDs are not installed, setting offline_kube_version and trying again")

	v, err := helper.GetKubernetesVersion(rel.Cfg())
	if err != nil {
		rel.Logger().WithError(err).Error("cannot get current kubernetes version, you need to set it manually")

		return err
	}

	rel.OfflineKubeVersionF = v.GitVersion
	rel.Logger().WithField("version", rel.OfflineKubeVersionF).Info("discovered kubernetes version")

	return nil
}

func (rel *config) test() (err error) {
	rel.Logger().Info("running helm tests")

	client := rel.newTest()

	// helm defers deleting the test pods to the returned shutdown func, so that the logs below
	// can still be read. It has to run whether or not the tests passed.
	res, shutdown, err := client.Run(rel.Name())
	defer func() {
		if shutdownErr := shutdown(); shutdownErr != nil && err == nil {
			err = shutdownErr
		}
	}()

	if (err != nil) || rel.Tests.ForceShowLogs {
		r, convErr := asRelease(res)
		if convErr != nil {
			return convErr
		}

		var buf bytes.Buffer
		if r != nil {
			_ = client.GetPodLogs(&buf, r)
		}

		if err != nil {
			rel.Logger().WithError(err).WithField("output", buf.String()).Error("helm tests failed")

			return NewHelmTestsError(err)
		}

		rel.Logger().WithField("output", buf.String()).Info("helm tests output")
	}

	return nil
}
