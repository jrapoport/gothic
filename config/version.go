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
	var v = "release"
	if debug {
		v = "debug"
	}
	if Version != "" {
		v = fmt.Sprintf("%s (%s)", Version, Build)
	}
	ver = v
}

// BuildVersion the build version string
func BuildVersion() string {
	return ver
}
