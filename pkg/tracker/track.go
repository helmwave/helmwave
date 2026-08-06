package tracker

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/fluxcd/cli-utils/pkg/kstatus/polling/aggregator"
	"github.com/fluxcd/cli-utils/pkg/kstatus/polling/collector"
	"github.com/fluxcd/cli-utils/pkg/kstatus/polling/event"
	"github.com/fluxcd/cli-utils/pkg/kstatus/status"
	"github.com/fluxcd/cli-utils/pkg/kstatus/watcher"
	"github.com/fluxcd/cli-utils/pkg/object"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

// Track watches the given objects and logs their status transitions until every one of them is
// ready, the context is done or the timeout expires. It uses the same kstatus engine as helm's
// `wait: watcher` strategy.
func Track(ctx context.Context, cfg *Config, ids object.ObjMetadataSet, kubecontext string) error {
	if len(ids) == 0 {
		log.Info("🐶 no resources to track")

		return nil
	}

	dyn, mapper, clientset, err := newClients(kubecontext)
	if err != nil {
		return err
	}

	var trackCtx context.Context
	var cancel context.CancelFunc
	if cfg.Timeout > 0 {
		trackCtx, cancel = context.WithTimeout(ctx, cfg.Timeout)
	} else {
		trackCtx, cancel = context.WithCancel(ctx)
	}
	defer cancel()

	if cfg.Logs {
		go newLogStreamer(clientset, cfg, ids).Run(trackCtx)
	}

	sw := watcher.NewDefaultStatusWatcher(dyn, mapper)
	eventCh := sw.Watch(trackCtx, ids, watcher.Options{
		RESTScopeStrategy: watcher.RESTScopeNamespace,
	})

	coll := collector.NewResourceStatusCollector(ids)
	done := coll.ListenWithObserver(eventCh, newLoggingObserver(cancel))
	<-done

	if coll.Error != nil {
		return fmt.Errorf("failed to track resources: %w", coll.Error)
	}

	reportNotReady(coll)

	return nil
}

// newLoggingObserver logs every status transition and stops the watch once everything is ready.
func newLoggingObserver(cancel context.CancelFunc) collector.ObserverFunc {
	var mu sync.Mutex
	last := map[object.ObjMetadata]status.Status{}

	return func(c *collector.ResourceStatusCollector, _ event.Event) {
		mu.Lock()
		defer mu.Unlock()

		rss := make([]*event.ResourceStatus, 0, len(c.ResourceStatuses))
		for _, rs := range c.ResourceStatuses {
			if rs == nil {
				continue
			}
			rss = append(rss, rs)
			logTransition(last, rs)
		}

		if aggregator.AggregateStatus(rss, status.CurrentStatus) == status.CurrentStatus {
			log.Info("🐶 all tracked resources are ready")
			cancel()
		}
	}
}

func logTransition(last map[object.ObjMetadata]status.Status, rs *event.ResourceStatus) {
	if last[rs.Identifier] == rs.Status {
		return
	}
	last[rs.Identifier] = rs.Status

	l := log.WithFields(log.Fields{
		"kind":            rs.Identifier.GroupKind.Kind,
		"name":            rs.Identifier.Name,
		logFieldNamespace: rs.Identifier.Namespace,
	})

	switch rs.Status {
	case status.CurrentStatus:
		l.Infof("🐶 ready: %s", rs.Message)
	case status.FailedStatus:
		l.Errorf("🐶 failed: %s", rs.Message)
	case status.InProgressStatus, status.TerminatingStatus:
		l.Infof("🐶 %s: %s", rs.Status, rs.Message)
	case status.NotFoundStatus, status.UnknownStatus:
		l.Debugf("🐶 %s: %s", rs.Status, rs.Message)
	}
}

// reportNotReady logs what was still not ready when tracking stopped. Tracking is advisory, so
// this reports instead of failing the deploy.
func reportNotReady(coll *collector.ResourceStatusCollector) {
	for id, rs := range coll.ResourceStatuses {
		if rs == nil || rs.Status == status.CurrentStatus {
			continue
		}

		log.WithFields(log.Fields{
			"kind":            id.GroupKind.Kind,
			"name":            id.Name,
			logFieldNamespace: id.Namespace,
		}).Warnf("🐶 still not ready: %s, %s", rs.Status, rs.Message)
	}
}

func newClients(kubecontext string) (dynamic.Interface, meta.RESTMapper, kubernetes.Interface, error) {
	flags := genericclioptions.NewConfigFlags(true)
	if kubecontext != "" {
		flags.Context = &kubecontext
	}
	if kubeconfig, ok := os.LookupEnv("KUBECONFIG"); ok {
		flags.KubeConfig = &kubeconfig
	}

	restCfg, err := flags.ToRESTConfig()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to build kubernetes config for tracking: %w", err)
	}

	mapper, err := flags.ToRESTMapper()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to build REST mapper for tracking: %w", err)
	}

	dyn, err := dynamic.NewForConfig(restCfg)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to build dynamic client for tracking: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to build kubernetes clientset for tracking: %w", err)
	}

	return dyn, mapper, clientset, nil
}
