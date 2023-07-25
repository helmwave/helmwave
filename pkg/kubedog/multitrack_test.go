package kubedog_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/kubedog"
	"github.com/stretchr/testify/suite"
	"github.com/werf/kubedog/pkg/trackers/rollout/multitrack"
	meta1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type MultitrackTestSuite struct {
	suite.Suite
}

func TestMultitrackTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(MultitrackTestSuite))
}

func (s *MultitrackTestSuite) TestNoResources() {
	res := []kubedog.Resource{}
	spec, err := kubedog.MakeSpecs(res, "", false)

	s.Require().NoError(err)

	s.Require().Empty(spec.Canaries)
	s.Require().Empty(spec.DaemonSets)
	s.Require().Empty(spec.Deployments)
	s.Require().Empty(spec.Generics)
	s.Require().Empty(spec.Jobs)
	s.Require().Empty(spec.StatefulSets)
}

func (s *MultitrackTestSuite) TestGenerics() {
	res := []kubedog.Resource{
		{
			TypeMeta: meta1.TypeMeta{
				Kind:       "ServiceAccount",
				APIVersion: "v1",
			},
			ObjectMeta: meta1.ObjectMeta{
				Name: "bla",
			},
		},
	}

	spec, err := kubedog.MakeSpecs(res, "", false)

	s.Require().NoError(err)

	s.Require().Empty(spec.Canaries)
	s.Require().Empty(spec.DaemonSets)
	s.Require().Empty(spec.Deployments)
	s.Require().Empty(spec.Generics)
	s.Require().Empty(spec.Jobs)
	s.Require().Empty(spec.StatefulSets)

	spec, err = kubedog.MakeSpecs(res, "", true)

	s.Require().NoError(err)

	s.Require().Empty(spec.Canaries)
	s.Require().Empty(spec.DaemonSets)
	s.Require().Empty(spec.Deployments)
	s.Require().Len(spec.Generics, 1)
	s.Require().Empty(spec.Jobs)
	s.Require().Empty(spec.StatefulSets)
}

func (s *MultitrackTestSuite) TestCanary() {
	res := kubedog.Resource{
		TypeMeta: meta1.TypeMeta{
			Kind:       "Canary",
			APIVersion: "v1",
		},
		ObjectMeta: meta1.ObjectMeta{
			Name: "bla",
		},
	}

	spec, err := kubedog.MakeSpecs([]kubedog.Resource{res}, "", false)

	s.Require().NoError(err)

	s.Require().Len(spec.Canaries, 1)
	s.Require().Empty(spec.DaemonSets)
	s.Require().Empty(spec.Deployments)
	s.Require().Empty(spec.Generics)
	s.Require().Empty(spec.Jobs)
	s.Require().Empty(spec.StatefulSets)
}

func (s *MultitrackTestSuite) TestDaemonSet() {
	res := kubedog.Resource{
		TypeMeta: meta1.TypeMeta{
			Kind:       "DaemonSet",
			APIVersion: "v1",
		},
		ObjectMeta: meta1.ObjectMeta{
			Name: "bla",
		},
	}

	spec, err := kubedog.MakeSpecs([]kubedog.Resource{res}, "", false)

	s.Require().NoError(err)

	s.Require().Empty(spec.Canaries)
	s.Require().Len(spec.DaemonSets, 1)
	s.Require().Empty(spec.Deployments)
	s.Require().Empty(spec.Generics)
	s.Require().Empty(spec.Jobs)
	s.Require().Empty(spec.StatefulSets)
}

func (s *MultitrackTestSuite) TestDeployment() {
	res := kubedog.Resource{
		TypeMeta: meta1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "v1",
		},
		ObjectMeta: meta1.ObjectMeta{
			Name: "bla",
		},
	}

	spec, err := kubedog.MakeSpecs([]kubedog.Resource{res}, "", false)

	s.Require().NoError(err)

	s.Require().Empty(spec.Canaries)
	s.Require().Empty(spec.DaemonSets)
	s.Require().Len(spec.Deployments, 1)
	s.Require().Empty(spec.Generics)
	s.Require().Empty(spec.Jobs)
	s.Require().Empty(spec.StatefulSets)
}

func (s *MultitrackTestSuite) TestJob() {
	res := kubedog.Resource{
		TypeMeta: meta1.TypeMeta{
			Kind:       "Job",
			APIVersion: "v1",
		},
		ObjectMeta: meta1.ObjectMeta{
			Name: "bla",
		},
	}

	spec, err := kubedog.MakeSpecs([]kubedog.Resource{res}, "", false)

	s.Require().NoError(err)

	s.Require().Empty(spec.Canaries)
	s.Require().Empty(spec.DaemonSets)
	s.Require().Empty(spec.Deployments)
	s.Require().Empty(spec.Generics)
	s.Require().Len(spec.Jobs, 1)
	s.Require().Empty(spec.StatefulSets)
}

func (s *MultitrackTestSuite) TestStatefulSet() {
	res := kubedog.Resource{
		TypeMeta: meta1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "v1",
		},
		ObjectMeta: meta1.ObjectMeta{
			Name: "bla",
		},
	}

	spec, err := kubedog.MakeSpecs([]kubedog.Resource{res}, "", false)

	s.Require().NoError(err)

	s.Require().Empty(spec.Canaries)
	s.Require().Empty(spec.DaemonSets)
	s.Require().Empty(spec.Deployments)
	s.Require().Empty(spec.Generics)
	s.Require().Empty(spec.Jobs)
	s.Require().Len(spec.StatefulSets, 1)
}

func (s *MultitrackTestSuite) TestAnnotations() {
	res := kubedog.Resource{
		TypeMeta: meta1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "v1",
		},
		ObjectMeta: meta1.ObjectMeta{
			Name: "bla",
			Annotations: map[string]string{
				kubedog.SkipLogsAnnoName:                  "true",
				kubedog.ShowEventsAnnoName:                "true",
				kubedog.LogRegexAnnoName:                  "true",
				kubedog.FailuresAllowedPerReplicaAnnoName: "100",
				kubedog.TrackTerminationModeAnnoName:      string(multitrack.NonBlocking),
				kubedog.FailModeAnnoName:                  string(multitrack.HopeUntilEndOfDeployProcess),
				kubedog.SkipLogsForContainersAnnoName:     "true,false",
				kubedog.ShowLogsOnlyForContainersAnnoName: "blabla",
				kubedog.LogRegexForAnnoPrefix + "bla":     "true",
			},
		},
	}

	spec, err := kubedog.MakeSpecs([]kubedog.Resource{res}, "", false)

	s.Require().NoError(err)

	s.Require().Len(spec.Deployments, 1)
	d := spec.Deployments[0]
	s.Require().Equal(d.ResourceName, res.ObjectMeta.Name)
	s.Require().True(d.SkipLogs)
	s.Require().True(d.ShowServiceMessages)
	s.Require().Equal("true", d.LogRegex.String())
	s.Require().Equal(100, *d.AllowFailuresCount)
	s.Require().Equal(multitrack.NonBlocking, d.TrackTerminationMode)
	s.Require().Equal(multitrack.HopeUntilEndOfDeployProcess, d.FailMode)
	s.Require().Equal([]string{"true", "false"}, d.SkipLogsForContainers)
	s.Require().Equal([]string{"blabla"}, d.ShowLogsOnlyForContainers)
	s.Require().Len(d.LogRegexByContainerName, 1)
	s.Require().Contains(d.LogRegexByContainerName, "bla")
	s.Require().Equal("true", d.LogRegexByContainerName["bla"].String())
}
