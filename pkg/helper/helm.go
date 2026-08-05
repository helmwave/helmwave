package helper

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"k8s.io/apimachinery/pkg/version"

	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v4/pkg/action"
	helm "helm.sh/helm/v4/pkg/cli"
	"helm.sh/helm/v4/pkg/registry"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
)

var (
	// Helm is an instance of helm CLI.
	Helm = helm.New()

	// HelmRegistryClient  is an instance of helm registry client.
	HelmRegistryClient *registry.Client
)

func init() {
	var err error
	HelmRegistryClient, err = registry.NewClient(
		registry.ClientOptDebug(Helm.Debug),
		registry.ClientOptWriter(log.StandardLogger().Writer()),
		registry.ClientOptCredentialsFile(Helm.RegistryConfig),
	)
	if err != nil {
		log.Fatal(err)
	}
}

// NewRegistryClient builds a registry client for a single release whose chart needs a non-default
// OCI transport (plain HTTP or a skipped TLS verification).
//
// The package-global HelmRegistryClient is TLS-only, and helm v4's OCI getter uses an injected
// registry client verbatim -- it only honors plain_http / insecure when it builds its own client
// (see pkg/getter/ocigetter.go). So the shared client can never reach a plain-HTTP or insecure
// registry; a per-release client carrying these options is the only way the per-chart flags reach
// the OCI transport. Mirrors helm's own cmd.newRegistryClient.
func NewRegistryClient(plainHTTP, insecureSkipTLSVerify bool) (*registry.Client, error) {
	opts := []registry.ClientOption{
		registry.ClientOptDebug(Helm.Debug),
		registry.ClientOptWriter(log.StandardLogger().Writer()),
		registry.ClientOptCredentialsFile(Helm.RegistryConfig),
	}

	if insecureSkipTLSVerify {
		opts = append(opts, registry.ClientOptHTTPClient(&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // the chart explicitly asked for it
				Proxy:           http.ProxyFromEnvironment,
			},
		}))
	}

	if plainHTTP {
		opts = append(opts, registry.ClientOptPlainHTTP())
	}

	client, err := registry.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create registry client: %w", err)
	}

	return client, nil
}

func wrapConfigFn(client *rest.Config) *rest.Config {
	client.QPS = 100   // default is 5.0
	client.Burst = 100 // default is 10

	return client
}

// NewCfg creates helm internal configuration for provided namespace and kubecontext.
func NewCfg(ns, kubecontext string) (*action.Configuration, error) {
	cfg := action.NewConfiguration(action.ConfigurationSetLogger(NewSlogHandler()))
	helmDriver := os.Getenv("HELM_DRIVER") // TODO: get rid of getenv in runtime
	config := genericclioptions.NewConfigFlags(true)
	config.WrapConfigFn = wrapConfigFn
	config.Namespace = &ns
	config.Context = &kubecontext

	err := cfg.Init(config, ns, helmDriver)
	if err != nil {
		return nil, fmt.Errorf("failed to create helm configuration for %s namespace: %w", ns, err)
	}

	cfg.RegistryClient = HelmRegistryClient

	return cfg, nil
}

// NewHelm is a hack to create an instance of helm CLI and specifying namespace without environment variables.
func NewHelm(ns string) *helm.EnvSettings {
	env := helm.New()
	env.SetNamespace(ns)

	return env
}

// GetKubernetesVersion returns kubernetes server version.
//
//nolint:wrapcheck
func GetKubernetesVersion(cfg *action.Configuration) (*version.Info, error) {
	clientSet, err := cfg.KubernetesClientSet()
	if err != nil {
		return nil, err
	}

	return clientSet.Discovery().ServerVersion()
}
