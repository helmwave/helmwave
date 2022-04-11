package plan

import (
	"testing"

	"github.com/hairyhenderson/go-fsimpl"
	"github.com/hairyhenderson/go-fsimpl/blobfs"
	"github.com/hairyhenderson/go-fsimpl/filefs"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	src := ".helmwave"
	// Allowed FS
	mux := fsimpl.NewMux()
	mux.Add(filefs.FS)
	mux.Add(blobfs.FS)

	_, err := mux.Lookup(src)
	if err != nil {
		src = "file://" + src
		_, err = mux.Lookup(src)
		require.NoError(t, err)
	}

	require.NoError(t, err)

}
