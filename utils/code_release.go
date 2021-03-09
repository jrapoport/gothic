// +build release

package utils

// IsDebugPIN returns true if the pin is a debug pin.
func IsDebugPIN(code string) bool {
	return code == DebugPIN
}
