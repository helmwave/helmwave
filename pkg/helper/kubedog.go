package helper

import (
	"fmt"
	"os"

	"github.com/werf/kubedog/pkg/kube"
)

func KubeInit() error {
	opts := kube.InitOptions{}
	kubeconfigPath, IsExist := os.LookupEnv("KUBECONFIG")
	if IsExist {
		opts.ConfigPath = kubeconfigPath
	}

	err := kube.Init(opts)
	if err != nil {
		return fmt.Errorf("failed to initialize kubernetes config: %w", err)
	}

	return nil
}
