package main

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

// helper.Dotenv replaces helper.Helm, so --progress must be bound to the instance
// that exists when the app is created.
func TestProgressFlagIsBoundToCurrentHelmSettings(t *testing.T) { //nolint:paralleltest // replaces a global
	helper.Dotenv()

	up := CreateApp().Command("up")
	require.NotNil(t, up)

	for _, flag := range up.Flags {
		if progress, ok := flag.(*cli.BoolFlag); ok && progress.Name == "progress" {
			require.Same(t, &helper.Helm.Debug, progress.Destination)

			return
		}
	}

	t.Fatal("progress flag not found")
}
