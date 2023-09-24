package plan

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type BuildReleasesTestSuite struct {
	suite.Suite
}

func TestBuildReleasesTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(BuildReleasesTestSuite))
}

func (ts *BuildReleasesTestSuite) TestCheckTagInclusion() {
	cases := []struct {
		targetTags  []string
		releaseTags []string
		matchAll    bool
		result      bool
	}{
		{
			targetTags:  []string{},
			releaseTags: []string{"bla"},
			matchAll:    false,
			result:      true,
		},
		{
			targetTags:  []string{},
			releaseTags: []string{"bla"},
			matchAll:    true,
			result:      true,
		},
		{
			targetTags:  []string{"bla"},
			releaseTags: []string{},
			matchAll:    false,
			result:      false,
		},
		{
			targetTags:  []string{"bla"},
			releaseTags: []string{},
			matchAll:    true,
			result:      false,
		},
		{
			targetTags:  []string{"bla"},
			releaseTags: []string{"abc"},
			matchAll:    false,
			result:      false,
		},
		{
			targetTags:  []string{"bla"},
			releaseTags: []string{"abc"},
			matchAll:    true,
			result:      false,
		},
		{
			targetTags:  []string{"bla"},
			releaseTags: []string{"bla"},
			matchAll:    false,
			result:      true,
		},
		{
			targetTags:  []string{"1", "2", "3"},
			releaseTags: []string{"3", "2"},
			matchAll:    false,
			result:      true,
		},
		{
			targetTags:  []string{"1", "2", "3"},
			releaseTags: []string{"3", "2"},
			matchAll:    true,
			result:      false,
		},
		{
			targetTags:  []string{"1", "2", "3"},
			releaseTags: []string{"2"},
			matchAll:    false,
			result:      true,
		},
		{
			targetTags:  []string{"1", "2", "3"},
			releaseTags: []string{"2"},
			matchAll:    true,
			result:      false,
		},
		{
			targetTags:  []string{"1", "2", "3"},
			releaseTags: []string{"3", "2", "1"},
			matchAll:    false,
			result:      true,
		},
		{
			targetTags:  []string{"1", "2", "3"},
			releaseTags: []string{"3", "2", "1"},
			matchAll:    true,
			result:      true,
		},
		{
			targetTags:  []string{"1", "2", "3"},
			releaseTags: []string{"3", "4", "1", "2"},
			matchAll:    false,
			result:      true,
		},
		{
			targetTags:  []string{"1", "2", "3"},
			releaseTags: []string{"3", "4", "1", "2"},
			matchAll:    true,
			result:      true,
		},
	}

	for i := range cases {
		c := cases[i]
		res := checkTagInclusion(c.targetTags, c.releaseTags, c.matchAll)
		ts.Equal(c.result, res, c)
	}
}

func (ts *BuildReleasesTestSuite) TestNoReleases() {
	p := New()
	p.NewBody()

	releases, err := p.buildReleases([]string{}, false)

	ts.Require().NoError(err)
	ts.Require().Empty(releases)
}

func (ts *BuildReleasesTestSuite) TestNoMatchingReleases() {
	p := New()

	mockedRelease := &MockReleaseConfig{}
	mockedRelease.On("Tags").Return([]string{"bla"})

	p.SetReleases(mockedRelease)

	releases, err := p.buildReleases([]string{"abc"}, true)
	ts.Require().NoError(err)
	ts.Require().Empty(releases)

	mockedRelease.AssertExpectations(ts.T())
}

func (ts *BuildReleasesTestSuite) TestDuplicateReleases() {
	p := New()

	tags := []string{"bla"}
	u := uniqname.UniqName(ts.T().Name())

	rel1 := &MockReleaseConfig{}
	rel1.On("Tags").Return(tags)
	rel1.On("Uniq").Return(u)
	rel1.On("DependsOn").Return([]*release.DependsOnReference{})
	rel1.On("SetDependsOn", []*release.DependsOnReference{}).Return()

	rel2 := &MockReleaseConfig{}
	rel2.On("Tags").Return(tags)
	rel2.On("Uniq").Return(u)

	p.SetReleases(rel1, rel2)

	releases, err := p.buildReleases(tags, true)
	ts.Require().ErrorIs(err, release.DuplicateError{})
	ts.Require().Empty(releases)

	rel1.AssertExpectations(ts.T())
	rel2.AssertExpectations(ts.T())
}

func (ts *BuildReleasesTestSuite) TestMissingRequiredDependency() {
	p := New()

	tags := []string{"bla"}
	u := uniqname.UniqName(ts.T().Name())

	rel := &MockReleaseConfig{}
	rel.On("Tags").Return(tags)
	rel.On("Uniq").Return(u)
	rel.On("DependsOn").Return([]*release.DependsOnReference{{Name: "blabla", Optional: false}})
	rel.On("Logger").Return(log.WithField("test", ts.T().Name()))

	p.SetReleases(rel)

	releases, err := p.buildReleases(tags, true)
	ts.Require().ErrorIs(err, release.ErrDepFailed)
	ts.Require().Empty(releases)

	rel.AssertExpectations(ts.T())
}

func (ts *BuildReleasesTestSuite) TestMissingOptionalDependency() {
	p := New()

	tags := []string{"bla"}
	u := uniqname.UniqName(ts.T().Name())

	rel := &MockReleaseConfig{}
	rel.On("Tags").Return(tags)
	rel.On("Uniq").Return(u)
	rel.On("DependsOn").Return([]*release.DependsOnReference{{Name: "blabla", Optional: true}})
	rel.On("SetDependsOn", []*release.DependsOnReference{}).Return()
	rel.On("Logger").Return(log.WithField("test", ts.T().Name()))

	p.SetReleases(rel)

	releases, err := p.buildReleases(tags, true)
	ts.Require().NoError(err)
	ts.Require().Len(releases, 1)
	ts.Require().Contains(releases, rel)

	rel.AssertExpectations(ts.T())
}

func (ts *BuildReleasesTestSuite) TestUnmatchedDependency() {
	p := New()

	tags := []string{"bla"}
	u1 := uniqname.UniqName(ts.T().Name())
	u2 := uniqname.UniqName("blabla")
	deps := []*release.DependsOnReference{{Name: u2.String()}}

	rel1 := &MockReleaseConfig{}
	rel1.On("Tags").Return(tags)
	rel1.On("Uniq").Return(u1)
	rel1.On("DependsOn").Return(deps)
	rel1.On("SetDependsOn", deps).Return()
	rel1.On("Logger").Return(log.WithField("test", ts.T().Name()))

	rel2 := &MockReleaseConfig{}
	rel2.On("Tags").Return([]string{})
	rel2.On("Uniq").Return(u2)
	rel2.On("DependsOn").Return([]*release.DependsOnReference{})
	rel2.On("SetDependsOn", []*release.DependsOnReference{}).Return()

	p.SetReleases(rel1, rel2)

	releases, err := p.buildReleases(tags, true)
	ts.Require().NoError(err)
	ts.Require().Len(releases, 2)
	ts.Require().Contains(releases, rel1)
	ts.Require().Contains(releases, rel2)

	rel1.AssertExpectations(ts.T())
	rel2.AssertExpectations(ts.T())
}
