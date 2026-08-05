package release

import (
	releaseiface "helm.sh/helm/v4/pkg/release"
	release "helm.sh/helm/v4/pkg/release/v1"
)

// asRelease unwraps the release interface helm actions return.
//
// helm's actions return the release.Releaser interface (an `any`), but every storage driver
// still decodes into the v1 struct, so unwrap it here rather than spreading the interface
// through helmwave, and fail loudly if helm ever starts handing out something else.
func asRelease(r releaseiface.Releaser) (*release.Release, error) {
	switch v := r.(type) {
	case nil:
		return nil, ErrNilRelease
	case *release.Release:
		if v == nil {
			// A typed nil pointer is still "no release": callers dereference the
			// result the moment the error is nil, so it has to be an error too.
			return nil, ErrNilRelease
		}

		return v, nil
	case release.Release:
		return &v, nil
	default:
		return nil, NewUnexpectedReleaseTypeError(r)
	}
}

// unwrapRelease adapts the (Releaser, error) pair a helm action returns. helm hands back a
// partially built release alongside the error when a sync fails, and callers here forward both, so
// the action's own error always wins over a conversion failure.
func unwrapRelease(r releaseiface.Releaser, err error) (*release.Release, error) {
	rel, convErr := asRelease(r)
	if err != nil {
		return rel, err
	}

	return rel, convErr
}
