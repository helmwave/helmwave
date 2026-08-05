package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/helmwave/helmwave/pkg/helper"
)

// ClusterEnv opts cluster tests into a real kubernetes API server.
// Its value is a kubeconfig path; it becomes KUBECONFIG for the whole test process.
const ClusterEnv = "HELMWAVE_TEST_CLUSTER"

// init pins KUBECONFIG so tests can never reach whatever cluster the developer happens to be
// logged into. It runs on import because KUBECONFIG is process-wide and t.Setenv cannot be
// called from a parallel test. Without ClusterEnv set, KUBECONFIG points at a kubeconfig that
// goes nowhere.
//
// helper.Helm snapshots KUBECONFIG when its package initializes; importing helper here makes
// that snapshot happen first, so it can be overwritten along with the environment.
func init() {
	kubeconfig := deadEndKubeconfig()

	if cluster := os.Getenv(ClusterEnv); cluster != "" {
		abs, err := filepath.Abs(cluster)
		if err != nil {
			panic(fmt.Sprintf("%s=%q cannot be resolved: %s", ClusterEnv, cluster, err))
		}
		if _, err := os.Stat(abs); err != nil {
			// A missing kubeconfig must be a loud failure, not a silent fallback:
			// opting into one cluster must never opt you into another.
			panic(fmt.Sprintf("%s=%q is not a readable kubeconfig: %s", ClusterEnv, cluster, err))
		}
		kubeconfig = abs
	}

	if err := os.Setenv("KUBECONFIG", kubeconfig); err != nil {
		panic(fmt.Sprintf("failed to pin KUBECONFIG: %s", err))
	}
	helper.Helm.KubeConfig = kubeconfig
}

// deadEndKubeconfig returns the fixture pointing at a server that does not exist.
func deadEndKubeconfig() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("failed to locate the dead-end kubeconfig fixture")
	}

	return filepath.Join(filepath.Dir(file), "kubeconfig.yaml")
}

// RequireCluster marks a test as needing a kubernetes API server. Without ClusterEnv it skips,
// so a laptop run stays green instead of failing with a cluster-unreachable error.
func RequireCluster(t *testing.T) {
	t.Helper()

	if os.Getenv(ClusterEnv) == "" {
		t.Skipf("test needs a kubernetes API server: set %s to a kubeconfig to run it", ClusterEnv)
	}
}
