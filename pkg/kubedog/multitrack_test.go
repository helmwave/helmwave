package kubedog

import (
	"github.com/stretchr/testify/suite"
	"github.com/werf/kubedog/pkg/trackers/rollout/multitrack"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
	"testing"
)

var (
	deployment = Resource{
		TypeMeta: metav1.TypeMeta{
			Kind: "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "deploy",
		},
		Spec: Spec{},
	}
	statefulSet = Resource{
		TypeMeta: metav1.TypeMeta{
			Kind: "StatefulSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "ss",
		},
		Spec: Spec{},
	}
	job = Resource{
		TypeMeta: metav1.TypeMeta{
			Kind: "Job",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "job",
		},
		Spec: Spec{},
	}
	daemonset = Resource{
		TypeMeta: metav1.TypeMeta{
			Kind: "DaemonSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "daemonset",
		},
		Spec: Spec{},
	}
)

type MultitrackTestSuite struct {
	suite.Suite
}

func (s *MultitrackTestSuite) TestSplitContainers() {
	_, err1 := splitContainers("")
	s.Error(err1)

	_, err2 := splitContainers(" ")
	s.Error(err2)

	c, err := splitContainers("a2,s4,d5")
	s.NoError(err)
	s.Equal([]string{"a2", "s4", "d5"}, c)
}

func (s *MultitrackTestSuite) TestSkipLogsAnno() {
	resources := []Resource{deployment, statefulSet, job, daemonset}
	ns := "test"

	resources[0].SetAnnotations(map[string]string{
		SkipLogsAnnoName: "blabla",
	})

	_, err := MakeSpecs(resources, ns)
	s.Error(err)

	resources[0].SetAnnotations(map[string]string{
		SkipLogsAnnoName: "false",
	})

	specs, err := MakeSpecs(resources, ns)
	s.Require().NoError(err)
	s.Require().NotNil(specs)
	s.Require().Len(specs.Deployments, 1)
	s.False(specs.Deployments[0].SkipLogs)
}

func (s *MultitrackTestSuite) TestShowEventsAnno() {
	resources := []Resource{deployment, statefulSet, job, daemonset}
	ns := "test"

	resources[0].SetAnnotations(map[string]string{
		ShowEventsAnnoName: "blabla",
	})

	_, err := MakeSpecs(resources, ns)
	s.Error(err)

	resources[0].SetAnnotations(map[string]string{
		ShowEventsAnnoName: "false",
	})

	specs, err := MakeSpecs(resources, ns)
	s.Require().NoError(err)
	s.Require().NotNil(specs)
	s.Require().Len(specs.Deployments, 1)
	s.False(specs.Deployments[0].ShowServiceMessages)
}

func (s *MultitrackTestSuite) TestLogRegexAnno() {
	resources := []Resource{deployment, statefulSet, job, daemonset}
	ns := "test"

	resources[1].SetAnnotations(map[string]string{
		LogRegexAnnoName: "/(((/",
	})

	_, err := MakeSpecs(resources, ns)
	s.Error(err)

	regex := "/123/"

	resources[1].SetAnnotations(map[string]string{
		LogRegexAnnoName: regex,
	})

	specs, err := MakeSpecs(resources, ns)
	s.Require().NoError(err)
	s.Require().NotNil(specs)
	s.Require().Len(specs.StatefulSets, 1)
	s.Require().NotNil(specs.StatefulSets[0].LogRegex)
	s.Equal(regex, specs.StatefulSets[0].LogRegex.String())
}

func (s *MultitrackTestSuite) TestFailuresAllowedPerReplicaAnno() {
	resources := []Resource{deployment, statefulSet, job, daemonset}
	ns := "test"

	resources[3].SetAnnotations(map[string]string{
		FailuresAllowedPerReplicaAnnoName: "blabla",
	})

	_, err := MakeSpecs(resources, ns)
	s.Error(err)

	resources[3].Spec.Replicas = new(uint32)
	*resources[3].Spec.Replicas = 2
	resources[3].SetAnnotations(map[string]string{
		FailuresAllowedPerReplicaAnnoName: "3",
	})

	specs, err := MakeSpecs(resources, ns)
	s.Require().NoError(err)
	s.Require().NotNil(specs)
	s.Require().Len(specs.DaemonSets, 1)
	s.Require().NotNil(specs.DaemonSets[0].AllowFailuresCount)
	s.Equal(6, *specs.DaemonSets[0].AllowFailuresCount)
}

func (s *MultitrackTestSuite) TestJobFailuresAllowedPerReplicaAnno() {
	resources := []Resource{job}
	ns := "test"

	resources[0].SetAnnotations(map[string]string{
		FailuresAllowedPerReplicaAnnoName: "blabla",
	})

	_, err := MakeSpecs(resources, ns)
	s.Error(err)
}

func (s *MultitrackTestSuite) TestTrackTerminationModeAnno() {
	resources := []Resource{deployment, statefulSet, job, daemonset}
	ns := "test"

	resources[0].SetAnnotations(map[string]string{
		TrackTerminationModeAnnoName: "blabla",
	})

	_, err := MakeSpecs(resources, ns)
	s.Error(err)

	value := multitrack.WaitUntilResourceReady
	resources[0].SetAnnotations(map[string]string{
		TrackTerminationModeAnnoName: string(value),
	})

	specs, err := MakeSpecs(resources, ns)
	s.Require().NoError(err)
	s.Require().NotNil(specs)
	s.Require().Len(specs.Deployments, 1)
	s.Equal(value, specs.Deployments[0].TrackTerminationMode)
}

func (s *MultitrackTestSuite) TestFailModeAnno() {
	resources := []Resource{deployment, statefulSet, job, daemonset}
	ns := "test"

	resources[0].SetAnnotations(map[string]string{
		FailModeAnnoName: "blabla",
	})

	_, err := MakeSpecs(resources, ns)
	s.Error(err)

	value := multitrack.HopeUntilEndOfDeployProcess
	resources[0].SetAnnotations(map[string]string{
		FailModeAnnoName: string(value),
	})

	specs, err := MakeSpecs(resources, ns)
	s.Require().NoError(err)
	s.Require().NotNil(specs)
	s.Require().Len(specs.Deployments, 1)
	s.Equal(value, specs.Deployments[0].FailMode)
}

func (s *MultitrackTestSuite) TestSkipLogsForContainersAnno() {
	resources := []Resource{deployment, statefulSet, job, daemonset}
	ns := "test"

	resources[0].SetAnnotations(map[string]string{
		SkipLogsForContainersAnnoName: "",
	})

	_, err := MakeSpecs(resources, ns)
	s.Error(err)

	value := []string{"1", "2"}
	resources[0].SetAnnotations(map[string]string{
		SkipLogsForContainersAnnoName: strings.Join(value, ","),
	})

	specs, err := MakeSpecs(resources, ns)
	s.Require().NoError(err)
	s.Require().NotNil(specs)
	s.Require().Len(specs.Deployments, 1)
	s.Equal(value, specs.Deployments[0].SkipLogsForContainers)
}

func (s *MultitrackTestSuite) TestShowLogsOnlyForContainersAnno() {
	resources := []Resource{deployment, statefulSet, job, daemonset}
	ns := "test"

	resources[0].SetAnnotations(map[string]string{
		ShowLogsOnlyForContainersAnnoName: "",
	})

	_, err := MakeSpecs(resources, ns)
	s.Error(err)

	value := []string{"1", "2"}
	resources[0].SetAnnotations(map[string]string{
		ShowLogsOnlyForContainersAnnoName: strings.Join(value, ","),
	})

	specs, err := MakeSpecs(resources, ns)
	s.Require().NoError(err)
	s.Require().NotNil(specs)
	s.Require().Len(specs.Deployments, 1)
	s.Equal(value, specs.Deployments[0].ShowLogsOnlyForContainers)
}

func (s *MultitrackTestSuite) TestLogRegexForAnno() {
	resources := []Resource{deployment, statefulSet, job, daemonset}
	ns := "test"

	resources[0].SetAnnotations(map[string]string{
		LogRegexForAnnoPrefix + "123": "/((((/",
	})

	_, err := MakeSpecs(resources, ns)
	s.Error(err)

	regex := "/123/"
	container := "123"

	resources[0].SetAnnotations(map[string]string{
		LogRegexForAnnoPrefix + container: regex,
	})

	specs, err := MakeSpecs(resources, ns)
	s.Require().NoError(err)
	s.Require().NotNil(specs)
	s.Require().Len(specs.Deployments, 1)
	s.Require().NotNil(specs.Deployments[0].LogRegexByContainerName)
	s.Require().NotNil(specs.Deployments[0].LogRegexByContainerName[container])
	s.Equal(regex, specs.Deployments[0].LogRegexByContainerName[container].String())
}

func (s *MultitrackTestSuite) TestMakeSpecs() {
	resources := []Resource{deployment, statefulSet, job, daemonset}
	ns := "test"

	resources[0].SetAnnotations(map[string]string{
		SkipLogsAnnoName: "false",
	})

	specs, err := MakeSpecs(resources, ns)

	s.NoError(err)
	s.NotNil(specs)

	s.Require().Len(specs.Deployments, 1)
	s.Require().Len(specs.StatefulSets, 1)
	s.Require().Len(specs.Jobs, 1)
	s.Require().Len(specs.DaemonSets, 1)
	s.Equal(ns, specs.Deployments[0].Namespace)
	s.Equal(ns, specs.StatefulSets[0].Namespace)
	s.Equal(ns, specs.Jobs[0].Namespace)
	s.Equal(ns, specs.DaemonSets[0].Namespace)
}

func TestMultitrackTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(MultitrackTestSuite))
}
