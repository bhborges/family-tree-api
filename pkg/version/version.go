// Package version supplies version information collected at build time.
package version

import "runtime"

// When building a image that uses this package, these variables will
// be filled from -ldflags settings.
//   - GitCommit:
//     output of: $(git rev-parse --short HEAD)
//   - BuildDate:
//     ISO8601 format, output of:
//     $(date ${SOURCE_DATE_EPOCH:+"--date=@${SOURCE_DATE_EPOCH}"} -u +'%Y-%m-%dT%H:%M:%SZ')
var (
	gitCommit string
	buildDate string
	goVersion = runtime.Version()
)

// Version contains versioning information.
// It simple wraps all info to avoid global variables.
type Version struct {
	GitCommit string `json:"gitCommit"`
	BuildDate string `json:"buildDate"`
	GoVersion string `json:"goVersion"`
}

// Get returns version info.
func Get() Version {
	return Version{
		GitCommit: gitCommit,
		BuildDate: buildDate,
		GoVersion: goVersion,
	}
}
