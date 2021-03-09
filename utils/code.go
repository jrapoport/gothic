package utils

// maxPIN is the max length of a PINCode (6).
const maxPIN = 6

// DebugPIN is bypass code that is only valid in debug builds.
const DebugPIN = "000000"

// PINCode returns a new random pin.
func PINCode() string {
	code := ""
	for {
		code = RandomPIN(maxPIN)
		// make sure this isn't the debug code
		if code != DebugPIN {
			break
		}
	}
	return code
}

// IsValidCode returns true if rhe code is valid.
func IsValidCode(code string) bool {
	return len(code) > 0
}
