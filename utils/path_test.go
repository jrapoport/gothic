package utils

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsLocalPath(t *testing.T) {
	const (
		testURL1 = "http://exmaple.com"
		testURL2 = "https://exmaple.com"
	)
	dir := t.TempDir()
	is := IsLocalPath(dir)
	assert.True(t, is)
	is = IsLocalPath(testURL1)
	assert.False(t, is)
	is = IsLocalPath(testURL2)
	assert.False(t, is)
	is = IsLocalPath("")
	assert.False(t, is)
}

func TestHasExt(t *testing.T) {
	dir := t.TempDir()
	has := HasExt(dir)
	assert.False(t, has)
	test := filepath.Join(dir, "test.text")
	has = HasExt(test)
	assert.True(t, has)
}

func TestIsDirectory(t *testing.T) {
	dir := t.TempDir()
	is := IsDirectory(dir)
	assert.True(t, is)
	test := filepath.Join(dir, "test.text")
	is = IsDirectory(test)
	assert.False(t, is)
	testFile, err := ioutil.TempFile(dir, "pre-*.txt")
	require.NoError(t, err)
	defer testFile.Close()
	is = IsDirectory(testFile.Name())
	assert.False(t, is)
}

func TestExecutableName(t *testing.T) {
	name, err := ExecutableName()
	assert.NoError(t, err)
	assert.NotEmpty(t, name)
}
