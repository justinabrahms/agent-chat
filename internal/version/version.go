// Package version provides build-time version information.
// Variables are set via ldflags during build.
package version

import "fmt"

// Build-time variables set via ldflags
var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

// String returns a formatted version string.
func String() string {
	return fmt.Sprintf("agent-chat %s\ncommit: %s\nbuilt: %s", Version, GitCommit, BuildDate)
}
