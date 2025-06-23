package gsrv

import (
	"fmt"
	"github.com/maloquacious/semver"
	"github.com/spf13/cobra"
)

var (
	version = semver.Version{Major: 0, Minor: 1, Patch: 0, PreRelease: "alpha"}
)

func Version() semver.Version {
	return version
}
