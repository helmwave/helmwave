package tracker

import (
	"context"
	"testing"
	"time"

	"github.com/fluxcd/cli-utils/pkg/object"
	"github.com/stretchr/testify/suite"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/fake"
)

type LogsTestSuite struct {
	suite.Suite
}

func TestLogsTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(LogsTestSuite))
}

func deploymentID(ns, name string) object.ObjMetadata {
	return object.ObjMetadata{
		Namespace: ns,
		Name:      name,
		GroupKind: schema.GroupKind{Group: groupApps, Kind: kindDeployment},
	}
}

func (s *LogsTestSuite) newStreamer(clientset *fake.Clientset, ids ...object.ObjMetadata) *logStreamer {
	return newLogStreamer(clientset, &Config{Logs: true}, ids)
}

func (s *LogsTestSuite) TestOwnedByTrackedJob() {
	clientset := fake.NewSimpleClientset()
	str := s.newStreamer(clientset, object.ObjMetadata{
		Namespace: "ns",
		Name:      "job",
		GroupKind: schema.GroupKind{Group: "batch", Kind: "Job"},
	})

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "job-abc",
			Namespace: "ns",
			OwnerReferences: []metav1.OwnerReference{{
				APIVersion: "batch/v1",
				Kind:       kindJob,
				Name:       "job",
				Controller: new(true),
			}},
		},
	}

	s.True(str.ownedByTracked(context.Background(), pod))
}

func (s *LogsTestSuite) TestOwnedByTrackedDeploymentViaReplicaSet() {
	rs := &appsv1.ReplicaSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "web-12345",
			Namespace: "ns",
			OwnerReferences: []metav1.OwnerReference{{
				APIVersion: "apps/v1",
				Kind:       kindDeployment,
				Name:       "web",
				Controller: new(true),
			}},
		},
	}
	clientset := fake.NewSimpleClientset(rs)
	str := s.newStreamer(clientset, deploymentID("ns", "web"))

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "web-12345-xyz",
			Namespace: "ns",
			OwnerReferences: []metav1.OwnerReference{{
				APIVersion: "apps/v1",
				Kind:       "ReplicaSet",
				Name:       "web-12345",
				Controller: new(true),
			}},
		},
	}

	s.True(str.ownedByTracked(context.Background(), pod))
	// The second call must hit the cache and still agree.
	s.True(str.ownedByTracked(context.Background(), pod))
}

func (s *LogsTestSuite) TestForeignPodIgnored() {
	clientset := fake.NewSimpleClientset()
	str := s.newStreamer(clientset, deploymentID("ns", "web"))

	orphan := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "orphan", Namespace: "ns"}}
	s.False(str.ownedByTracked(context.Background(), orphan))

	foreign := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "other-abc",
			Namespace: "ns",
			OwnerReferences: []metav1.OwnerReference{{
				APIVersion: "batch/v1",
				Kind:       kindJob,
				Name:       "other",
				Controller: new(true),
			}},
		},
	}
	s.False(str.ownedByTracked(context.Background(), foreign))
}

func (s *LogsTestSuite) TestOnlyWorkloadsAreStreamed() {
	clientset := fake.NewSimpleClientset()
	str := s.newStreamer(clientset,
		deploymentID("ns", "web"),
		object.ObjMetadata{
			Namespace: "ns",
			Name:      "canary",
			GroupKind: schema.GroupKind{Group: "flagger.app", Kind: "Canary"},
		},
	)

	s.Len(str.tracked, 1)
}

func (s *LogsTestSuite) TestHandlePodStartsOneStreamPerContainer() {
	clientset := fake.NewSimpleClientset()
	str := s.newStreamer(clientset, object.ObjMetadata{
		Namespace: "ns",
		Name:      "job",
		GroupKind: schema.GroupKind{Group: "batch", Kind: "Job"},
	})

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "job-abc",
			Namespace: "ns",
			OwnerReferences: []metav1.OwnerReference{{
				APIVersion: "batch/v1",
				Kind:       kindJob,
				Name:       "job",
				Controller: new(true),
			}},
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodSucceeded, // finished: the stream goroutine exits after one pass
			ContainerStatuses: []corev1.ContainerStatus{
				{Name: "main", State: corev1.ContainerState{Terminated: &corev1.ContainerStateTerminated{}}},
				{Name: "waiting", State: corev1.ContainerState{Waiting: &corev1.ContainerStateWaiting{}}},
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	str.handlePod(ctx, pod)
	str.handlePod(ctx, pod) // the same pod again must not double-register

	str.streamsMu.Lock()
	defer str.streamsMu.Unlock()
	s.Len(str.streams, 1)
	s.True(str.streams["ns/job-abc/main"])
}

func (s *LogsTestSuite) TestTrimLine() {
	s.Equal("abc", trimLine("abcdef", 3))
	s.Equal("abcdef", trimLine("abcdef", 0))
	s.Equal("abcdef", trimLine("abcdef", 100))
}
