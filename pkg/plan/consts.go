package plan

import "errors"

const (
	// Dir is default directory for generated files.
	Dir = ".helmwave/"

	// File is default file name for planfile.
	File = "planfile"

	// Body is default file name for main config.
	Body = "helmwave.yml"

	// Manifest is default directory under Dir for manifests.
	Manifest = "manifest/"

	// Values is default directory for values.
	Values = "values/"

	// GraphFilename is default filename of markdown graph
	GraphFilename = "graph.md"
)

var (
	// ErrManifestDirNotFound is an error for nonexistent manifest dir.
	ErrManifestDirNotFound = errors.New(Manifest + " dir not found")

	// ErrManifestDirEmpty is an error for empty manifest dir.
	ErrManifestDirEmpty = errors.New(Manifest + " is empty")
)
