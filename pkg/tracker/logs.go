package tracker

import (
	"bufio"
	"context"
	"sync"
	"time"

	"github.com/fluxcd/cli-utils/pkg/object"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
)

const streamRetryDelay = 2 * time.Second

const logFieldNamespace = "namespace"

// logStreamer follows the container logs of pods owned by tracked workloads and forwards every
// line into the logrus stream, the way kubedog used to.
type logStreamer struct {
	clientset kubernetes.Interface
	cfg       *Config

	// tracked is the set of workloads whose pods are streamed.
	tracked map[object.ObjMetadata]bool

	// rsOwners caches ReplicaSet -> owning workload lookups: namespace/name -> owner id.
	rsOwners map[string]object.ObjMetadata

	// streams remembers which namespace/pod/container is already being followed.
	streams map[string]bool

	start metav1.Time

	rsOwnersMu sync.Mutex
	streamsMu  sync.Mutex
}

func newLogStreamer(clientset kubernetes.Interface, cfg *Config, ids object.ObjMetadataSet) *logStreamer {
	tracked := make(map[object.ObjMetadata]bool, len(ids))
	for i := range ids {
		if workloadGK[ids[i].GroupKind] {
			tracked[ids[i]] = true
		}
	}

	return &logStreamer{
		clientset: clientset,
		cfg:       cfg,
		start:     metav1.Now(),
		tracked:   tracked,
		rsOwners:  map[string]object.ObjMetadata{},
		streams:   map[string]bool{},
	}
}

// Run watches pods in every tracked namespace until the context is done. It never returns an
// error: log streaming is best effort on top of best-effort tracking.
func (s *logStreamer) Run(ctx context.Context) {
	namespaces := map[string]bool{}
	for id := range s.tracked {
		namespaces[id.Namespace] = true
	}

	var wg sync.WaitGroup
	for ns := range namespaces {
		wg.Add(1)
		go func(ns string) {
			defer wg.Done()
			s.watchPods(ctx, ns)
		}(ns)
	}
	wg.Wait()
}

func (s *logStreamer) watchPods(ctx context.Context, ns string) {
	for ctx.Err() == nil {
		w, err := s.clientset.CoreV1().Pods(ns).Watch(ctx, metav1.ListOptions{})
		if err != nil {
			log.WithError(err).WithField("namespace", ns).Debug("failed to watch pods for log streaming")
			sleepCtx(ctx, streamRetryDelay)

			continue
		}

		for ev := range w.ResultChan() {
			pod, ok := ev.Object.(*corev1.Pod)
			if !ok {
				continue
			}
			s.handlePod(ctx, pod)
		}
		// The watch expired: reconnect unless we are done.
	}
}

func (s *logStreamer) handlePod(ctx context.Context, pod *corev1.Pod) {
	if !s.ownedByTracked(ctx, pod) {
		return
	}

	statuses := make([]corev1.ContainerStatus, 0, len(pod.Status.InitContainerStatuses)+len(pod.Status.ContainerStatuses))
	statuses = append(statuses, pod.Status.InitContainerStatuses...)
	statuses = append(statuses, pod.Status.ContainerStatuses...)

	for i := range statuses {
		st := &statuses[i]
		if st.State.Running == nil && st.State.Terminated == nil {
			continue // there are no logs to read yet
		}

		key := pod.Namespace + "/" + pod.Name + "/" + st.Name
		s.streamsMu.Lock()
		started := s.streams[key]
		s.streams[key] = true
		s.streamsMu.Unlock()

		if started {
			continue
		}

		go s.streamContainer(ctx, pod.Namespace, pod.Name, st.Name)
	}
}

// ownedByTracked reports whether the pod's controller chain leads to a tracked workload.
// The ReplicaSet hop of Deployments is resolved with a cached GET.
func (s *logStreamer) ownedByTracked(ctx context.Context, pod *corev1.Pod) bool {
	owner := metav1.GetControllerOf(pod)
	if owner == nil {
		return false
	}

	id := object.ObjMetadata{
		Namespace: pod.Namespace,
		Name:      owner.Name,
		GroupKind: schema.GroupKind{Group: groupFromAPIVersion(owner.APIVersion), Kind: owner.Kind},
	}

	if owner.Kind == "ReplicaSet" {
		rsOwner, ok := s.replicaSetOwner(ctx, pod.Namespace, owner.Name)
		if !ok {
			return false
		}
		id = rsOwner
	}

	return s.tracked[id]
}

func (s *logStreamer) replicaSetOwner(ctx context.Context, ns, name string) (object.ObjMetadata, bool) {
	key := ns + "/" + name

	s.rsOwnersMu.Lock()
	cached, ok := s.rsOwners[key]
	s.rsOwnersMu.Unlock()
	if ok {
		return cached, cached.Name != ""
	}

	rs, err := s.clientset.AppsV1().ReplicaSets(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return object.ObjMetadata{}, false
	}

	id := object.ObjMetadata{}
	if owner := metav1.GetControllerOf(rs); owner != nil {
		id = object.ObjMetadata{
			Namespace: ns,
			Name:      owner.Name,
			GroupKind: schema.GroupKind{Group: groupFromAPIVersion(owner.APIVersion), Kind: owner.Kind},
		}
	}

	s.rsOwnersMu.Lock()
	s.rsOwners[key] = id
	s.rsOwnersMu.Unlock()

	return id, id.Name != ""
}

// streamContainer follows one container's logs until the context is done or the pod finishes.
// It reconnects after EOF so logs keep flowing across container restarts.
func (s *logStreamer) streamContainer(ctx context.Context, ns, pod, container string) {
	l := log.WithFields(log.Fields{
		logFieldNamespace: ns,
		"pod":             pod,
		"container":       container,
	})

	since := s.start
	for ctx.Err() == nil {
		req := s.clientset.CoreV1().Pods(ns).GetLogs(pod, &corev1.PodLogOptions{
			Container: container,
			Follow:    true,
			SinceTime: &since,
		})

		stream, err := req.Stream(ctx)
		if err != nil {
			if apierrors.IsNotFound(err) || ctx.Err() != nil {
				return
			}
			sleepCtx(ctx, streamRetryDelay)

			continue
		}

		sc := bufio.NewScanner(stream)
		sc.Buffer(make([]byte, 64*1024), 1024*1024)
		for sc.Scan() {
			since = metav1.Now()
			l.Info(trimLine(sc.Text(), s.cfg.LogWidth))
		}
		_ = stream.Close()

		if ctx.Err() != nil || s.podFinished(ctx, ns, pod) {
			return
		}

		// The container terminated or the stream broke: retry to catch restarts.
		sleepCtx(ctx, streamRetryDelay)
	}
}

func (s *logStreamer) podFinished(ctx context.Context, ns, name string) bool {
	pod, err := s.clientset.CoreV1().Pods(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return apierrors.IsNotFound(err)
	}

	return pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed
}

func groupFromAPIVersion(apiVersion string) string {
	gv, err := schema.ParseGroupVersion(apiVersion)
	if err != nil {
		return ""
	}

	return gv.Group
}

func trimLine(line string, width int) string {
	if width > 0 && len(line) > width {
		return line[:width]
	}

	return line
}

func sleepCtx(ctx context.Context, d time.Duration) {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
	case <-t.C:
	}
}
