package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// IsLocalPath returns true if the string is not an http url.
func IsLocalPath(s string) bool {
	if s == "" {
		return false
	}
	return !strings.HasPrefix(s, "http")
}

// HasExt returns true if the path has a file ext.
func HasExt(file string) bool {
	return filepath.Ext(file) != ""
}

// IsDirectory returns true if the path is a directory.
func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
