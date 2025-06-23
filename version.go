package gsrv

import (
	"runtime/debug"

	"github.com/maloquacious/semver"
)

// version defines the semantic version of the gsrv package.
var (
	version = semver.Version{Major: 1, Minor: 0, Patch: 1}
)

// Version returns build and version information for the gsrv package.
// It combines the defined semantic version with build-time information
// extracted from the Go module system, including the package version,
// VCS modification status, and revision hash.
func Version() (result struct {
	Version        semver.Version // Semantic version defined in this package
	PackageVersion string         // Version from go.mod or build system
	Modified       string         // VCS modification status ("true" if dirty, "false" if clean)
	Revision       string         // VCS revision hash
}) {
	result.Version = version
	if info, ok := debug.ReadBuildInfo(); ok {
		result.PackageVersion = info.Main.Version
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.modified":
				result.Modified = setting.Value
			case "vcs.revision":
				result.Revision = setting.Value
			}
		}
	}
	return result
}
