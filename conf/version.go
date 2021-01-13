package conf

import "fmt"

var Version string
var Build string

func CurrentVersion() string {
	return fmt.Sprintf("%s (%s)", Version, Build)
}
