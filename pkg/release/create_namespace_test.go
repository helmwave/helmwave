//go:build integration

package release_test

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/repo"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type CreateNamespaceTestSuite struct {
	suite.Suite

	ctx       context.Context
	clientSet kubernetes.Interface
}

func TestCreateNamespaceTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(CreateNamespaceTestSuite))
}

func (ts *CreateNamespaceTestSuite) SetupSuite() {
	ts.ctx = tests.GetContext(ts.T())

	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	).ClientConfig()
	ts.Require().NoError(err)

	ts.clientSet, err = kubernetes.NewForConfig(config)
	ts.Require().NoError(err)

	var rs repo.Configs
	str := `
- name: bitnami
  url: https://charts.bitnami.com/bitnami
`
	ts.Require().NoError(yaml.Unmarshal([]byte(str), &rs))
	ts.Require().NoError(plan.SyncRepositories(ts.ctx, rs))
}

// TestSyncIntoExistingNamespace verifies the get-first behavior: when the target
// namespace already exists, a release with create_namespace=true installs without
// attempting to create the namespace.
func (ts *CreateNamespaceTestSuite) TestSyncIntoExistingNamespace() {
	ns := strings.ToLower(strings.ReplaceAll(ts.T().Name(), "/", ""))

	_, err := ts.clientSet.CoreV1().Namespaces().Create(ts.ctx, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: ns},
	}, metav1.CreateOptions{})
	if err != nil && !apierrors.IsAlreadyExists(err) {
		ts.Require().NoError(err)
	}

	rel := release.NewConfig()
	rel.NamespaceF = ns
	rel.CreateNamespace = true
	rel.WaitStrategy = release.WaitStrategyHookOnly
	rel.ChartF.Name = "bitnami/nginx"
	rel.ValuesF = append(rel.ValuesF, release.ValuesReference{
		Dst: filepath.Join(tests.Root, "06_values.yaml"),
	})

	r, err := rel.Sync(ts.ctx, false)
	ts.Require().NoError(err)
	ts.Require().NotNil(r)
}
