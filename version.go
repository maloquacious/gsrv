package gsrv

import (
	"runtime/debug"

	"github.com/maloquacious/semver"
)

var (
	version = semver.Version{Major: 1, Minor: 0, Patch: 0}
)

func Version() (result struct {
	Version        semver.Version
	PackageVersion string
	Modified       string
	Revision       string
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
