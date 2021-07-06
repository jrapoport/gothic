// +build !release

package utils

// IsDebugPIN returns true if the pin is a debug pin.
func IsDebugPIN(_ string) bool {
	return false
}
