package build

import (
	"runtime/debug"
)

// Define public variables here which you wish to be configurable at build time
var (
	Version           = "dev"
	DefaultPageSize   = 100
	DefaultPageNumber = 1
	ConsoleURL        = "https://console.redhat.com"
	AuthURL           = "https://auth.redhat.com/auth/realms/EmployeeIDP"
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
