//go:build !release
// +build !release

package utils

// IsDebugPIN returns true if the pin is a debug pin.
func IsDebugPIN(pin string) bool {
	return pin == DebugPIN
}
