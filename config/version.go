package config

import "fmt"

var (
	// Version is the build version
	Version string
	// Build is the build number
	Build string
)

var ver string

func init() {
	ver = version()
}

// BuildVersion the build version string
func BuildVersion() string {
	return ver
}

func version() string {
	var v = "release"
	if debug {
		v = "debug"
	}
	if Version != "" && Build != "" {
		v = fmt.Sprintf("%s (%s)", Version, Build)
	} else if Version != "" {
		v = Version
	} else if Build != "" {
		v = Build
	}
	return v
}
