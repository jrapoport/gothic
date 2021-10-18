package config

import "fmt"

var (
	// Version is the build version
	Version string
	// Build is the build number
	Build string
	// ExeName is the exe name
	ExeName string
)

var ver string

func init() {
	ver = version()
}

// BuildVersion the build version string
func BuildVersion() string {
	return ver
}

// BuildName is the build name of the executable
func BuildName() string {
	return ExeName
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
