package plan

import (
	"path/filepath"
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
	tmpDir := ts.T().TempDir()
	p := New(filepath.Join(tmpDir, Dir))
	p.NewBody()

	releases, err := p.buildReleases(BuildOptions{})

	ts.Require().NoError(err)
	ts.Empty(releases)
}

func (ts *BuildReleasesTestSuite) TestNoMatchingReleases() {
	tmpDir := ts.T().TempDir()
	p := New(filepath.Join(tmpDir, Dir))

	mockedRelease := &MockReleaseConfig{}
	mockedRelease.On("Tags").Return([]string{"bla"})

	p.SetReleases(mockedRelease)

	releases, err := p.buildReleases(BuildOptions{Tags: []string{"abc"}, MatchAll: true})
	ts.Require().NoError(err)
	ts.Empty(releases)

	mockedRelease.AssertExpectations(ts.T())
}

func (ts *BuildReleasesTestSuite) TestDuplicateReleases() {
	tmpDir := ts.T().TempDir()
	p := New(filepath.Join(tmpDir, Dir))

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

	releases, err := p.buildReleases(BuildOptions{Tags: tags, MatchAll: true, EnableDependencies: true})

	var e *release.DuplicateError
	ts.Require().ErrorAs(err, &e)
	ts.Equal(u, e.Uniq)

	ts.Empty(releases)

	rel1.AssertExpectations(ts.T())
	rel2.AssertExpectations(ts.T())
}

func (ts *BuildReleasesTestSuite) TestMissingRequiredDependency() {
	tmpDir := ts.T().TempDir()
	p := New(filepath.Join(tmpDir, Dir))

	tags := []string{"bla"}
	u := uniqname.UniqName(ts.T().Name())

	rel := &MockReleaseConfig{}
	rel.On("Tags").Return(tags)
	rel.On("Uniq").Return(u)
	rel.On("DependsOn").Return([]*release.DependsOnReference{{Name: "blabla", Optional: false}})
	rel.On("Logger").Return(log.WithField("test", ts.T().Name()))

	p.SetReleases(rel)

	releases, err := p.buildReleases(BuildOptions{Tags: tags, MatchAll: true, EnableDependencies: true})
	ts.ErrorIs(err, release.ErrDepFailed)
	ts.Empty(releases)

	rel.AssertExpectations(ts.T())
}

func (ts *BuildReleasesTestSuite) TestMissingOptionalDependency() {
	tmpDir := ts.T().TempDir()
	p := New(filepath.Join(tmpDir, Dir))

	tags := []string{"bla"}
	u := uniqname.UniqName(ts.T().Name())

	rel := &MockReleaseConfig{}
	rel.On("Tags").Return(tags)
	rel.On("Uniq").Return(u)
	rel.On("DependsOn").Return([]*release.DependsOnReference{{Name: "blabla", Optional: true}})
	rel.On("SetDependsOn", []*release.DependsOnReference{}).Return()
	rel.On("Logger").Return(log.WithField("test", ts.T().Name()))

	p.SetReleases(rel)

	releases, err := p.buildReleases(BuildOptions{Tags: tags, MatchAll: true, EnableDependencies: true})
	ts.Require().NoError(err)
	ts.Len(releases, 1)
	ts.Contains(releases, rel)

	rel.AssertExpectations(ts.T())
}

func (ts *BuildReleasesTestSuite) TestUnmatchedDependency() {
	tmpDir := ts.T().TempDir()
	p := New(filepath.Join(tmpDir, Dir))

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

	releases, err := p.buildReleases(BuildOptions{Tags: tags, MatchAll: true, EnableDependencies: true})
	ts.Require().NoError(err)
	ts.Len(releases, 2)
	ts.Contains(releases, rel1)
	ts.Contains(releases, rel2)

	rel1.AssertExpectations(ts.T())
	rel2.AssertExpectations(ts.T())
}

func (ts *BuildReleasesTestSuite) TestDisabledDependencies() {
	tmpDir := ts.T().TempDir()
	p := New(filepath.Join(tmpDir, Dir))

	tags := []string{"bla"}
	u1 := uniqname.UniqName(ts.T().Name())

	rel1 := &MockReleaseConfig{}
	rel1.On("Tags").Return(tags)
	rel1.On("Uniq").Return(u1)
	rel1.On("SetDependsOn", []*release.DependsOnReference{}).Return()

	rel2 := &MockReleaseConfig{}
	rel2.On("Tags").Return([]string{})

	p.SetReleases(rel1, rel2)

	releases, err := p.buildReleases(BuildOptions{Tags: tags, MatchAll: true, EnableDependencies: false})
	ts.Require().NoError(err)
	ts.Len(releases, 1)
	ts.Contains(releases, rel1)

	rel1.AssertExpectations(ts.T())
	rel2.AssertExpectations(ts.T())
}
