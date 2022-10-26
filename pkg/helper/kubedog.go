package helper

import (
	"fmt"
	"os"

	"github.com/werf/kubedog/pkg/kube"
)

// KubeInit init kubeconfig for kubedog.
func KubeInit(kubecontext string) (err error) {
	opts := kube.InitOptions{}
	kubeconfigPath, IsExists := os.LookupEnv("KUBECONFIG")
	if IsExists {
		opts.ConfigPath = kubeconfigPath
	}
	if kubecontext != "" {
		opts.Context = kubecontext
	} else {
		opts.Context = Helm.KubeContext
	}

	err = kube.Init(opts)
	if err != nil {
		return fmt.Errorf("failed to initialize kubernetes config: %w", err)
	}

	return nil
}
