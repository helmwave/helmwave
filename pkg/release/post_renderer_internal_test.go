package release

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/suite"
)

type PostRendererInternalTestSuite struct {
	suite.Suite
}

func TestPostRendererInternalTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(PostRendererInternalTestSuite))
}

func (ts *PostRendererInternalTestSuite) TestNotFound() {
	// No path separator, so it is looked up in $PATH rather than resolved as a path.
	pr, err := newExecPostRenderer("helmwave-no-such-post-renderer")

	ts.Require().ErrorIs(err, exec.ErrNotFound)

	var e *PostRendererNotFoundError
	ts.Require().ErrorAs(err, &e)
	ts.Nil(pr)
}

// The manifests must reach the command's stdin and come back from its stdout untouched.
func (ts *PostRendererInternalTestSuite) TestManifestsRoundTrip() {
	pr, err := newExecPostRenderer("cat")
	ts.Require().NoError(err)

	manifests := "apiVersion: v1\nkind: ConfigMap\n"

	out, err := pr.Run(bytes.NewBufferString(manifests))

	ts.Require().NoError(err)
	ts.Equal(manifests, out.String())
}

func (ts *PostRendererInternalTestSuite) TestCommandFails() {
	pr, err := newExecPostRenderer("false")
	ts.Require().NoError(err)

	out, err := pr.Run(bytes.NewBufferString("kind: ConfigMap\n"))

	var e *PostRendererError
	ts.Require().ErrorAs(err, &e)
	ts.Nil(out)
}
