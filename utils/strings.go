package utils

import (
	"fmt"
	"strings"
)

// Namespaced returns a namespaced formatted string.
func Namespaced(ns, name string) string {
	if ns != "" {
		ns += "_"
	}
	return ns + name
}

// MaskString replaces a section of the string with mask character '*'.
func MaskString(s string, n int) string {
	const maskToken = "*"
	maskLen := n
	maskLen = len(s) - maskLen
	if maskLen <= 0 {
		maskLen = len(s)
	}
	mask := strings.Repeat(maskToken, maskLen)
	return strings.Replace(s, s[len(s)-maskLen:], mask, 1)
}

// MaskEmail replaces sections of the email with mask characters '*'.
func MaskEmail(e string) string {
	parts := strings.Split(e, "@")
	if len(parts) < 2 {
		return ""
	}
	domain := strings.Split(parts[1], ".")
	if len(domain) < 2 {
		return ""
	}
	parts[0] = MaskString(parts[0], 2)
	domain[0] = MaskString(domain[0], 1)
	return fmt.Sprintf("%s@%s.%s", parts[0], domain[0], domain[1])
}
