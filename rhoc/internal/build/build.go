package build

import (
	"runtime/debug"
)

// Define public variables here which you wish to be configurable at build time
var (
	// Version is dynamically set by the toolchain or overridden by the Makefile.
	Version = "dev"

	// Language used, can be overridden by Makefile or CI
	Language = "en"

	// DefaultPageSize is the default number of items per page when using list commands
	DefaultPageSize = 100

	// DefaultPageNumber is the default page number when using list commands
	DefaultPageNumber = 1
)

func init() {
	if isDevBuild() {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "(devel)" {
			Version = info.Main.Version
		}
	}
}

// isDevBuild returns true if the current build is "dev" (dev build)
func isDevBuild() bool {
	return Version == "dev"
}
