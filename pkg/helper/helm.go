package helper

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"helm.sh/helm/v3/pkg/action"
	helm "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/registry"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

//nolint:gochecknoglobals // TODO: get rid of globals
var (
	// Helm is an instance of helm CLI.
	Helm = helm.New()

	// Default logLevel for helm logs.
	helmLogLevel = log.Debugf

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

// NewCfg creates helm internal configuration for provided namespace.
func NewCfg(ns string) (*action.Configuration, error) {
	cfg := new(action.Configuration)
	helmDriver := os.Getenv("HELM_DRIVER") // TODO: get rid of getenv in runtime
	config := genericclioptions.NewConfigFlags(false)
	config.Namespace = &ns
	config.Context = &Helm.KubeContext

	if Helm.Debug {
		helmLogLevel = log.Infof
	}
	err := cfg.Init(config, ns, helmDriver, helmLogLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to create helm configuration for %s namespace: %w", ns, err)
	}

	cfg.RegistryClient = HelmRegistryClient

	return cfg, nil
}

// NewHelm is a hack to create an instance of helm CLI and specifying namespace without environment variables.
func NewHelm(ns string) (*helm.EnvSettings, error) {
	env := helm.New()
	fs := &pflag.FlagSet{}
	env.AddFlags(fs)
	flag := fs.Lookup("namespace")

	if err := flag.Value.Set(ns); err != nil {
		return nil, fmt.Errorf("failed to set namespace %s for helm: %w", ns, err)
	}

	return env, nil
}
