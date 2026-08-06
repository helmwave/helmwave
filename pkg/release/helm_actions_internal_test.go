package release

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"helm.sh/helm/v4/pkg/action"
)

type HelmActionsInternalTestSuite struct {
	suite.Suite
}

func TestHelmActionsInternalTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(HelmActionsInternalTestSuite))
}

// isDryRun and interactWithRemote mirror helm's own predicates in action
// (helm.sh/helm/v4 pkg/action/action.go). They are copied here because they are
// unexported, and because the options newInstall sets only matter through them:
// when interactWithRemote is true helm renders the chart through a live REST
// client, so it needs a reachable cluster and `lookup` queries it.
func isDryRun(c *action.Install) bool {
	return c.DryRunStrategy == action.DryRunClient || c.DryRunStrategy == action.DryRunServer
}

func interactWithRemote(c *action.Install) bool {
	return c.DryRunStrategy == action.DryRunNone || c.DryRunStrategy == action.DryRunServer
}

// A real install must not be turned into a dry run by a stray strategy.
func (ts *HelmActionsInternalTestSuite) TestNewInstallWithoutDryRun() {
	rel := NewConfig()

	client := rel.newInstall()

	ts.Equal(action.DryRunNone, client.DryRunStrategy, "any other strategy puts helm into dry-run mode")
	ts.False(isDryRun(client))
	ts.True(interactWithRemote(client))
}

func (ts *HelmActionsInternalTestSuite) TestNewInstallDryRun() {
	rel := NewConfig()
	rel.DryRun(true)

	client := rel.newInstall()

	ts.True(client.Replace)
	ts.Equal(action.DryRunServer, client.DryRunStrategy)
	ts.True(interactWithRemote(client), "a normal build renders against the cluster so `lookup` works")
}

func (ts *HelmActionsInternalTestSuite) TestNewInstallDryRunOfflineKubeVersion() {
	rel := NewConfig()
	rel.DryRun(true)
	rel.OfflineKubeVersionF = "v1.29.0"

	client := rel.newInstall()

	ts.Require().NotNil(client.KubeVersion)
	ts.Equal("v1.29.0", client.KubeVersion.Version)

	ts.Equal(action.DryRunClient, client.DryRunStrategy)
	ts.False(interactWithRemote(client), "an offline build must not touch the cluster")
}
