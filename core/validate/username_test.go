package validate

import (
	"testing"

	"github.com/jrapoport/gothic/test/tconf"
	"github.com/stretchr/testify/assert"
)

func TestUsername(t *testing.T) {
	t.Parallel()
	const tooLongName = "nzYgP4vScpzAG2wyXX1MaWNsRwlGof9m1HAy2tX05" +
		"54WGWgyEyokhSkHCkQFjnEsbx6Nk6CetjU0H24pmQw9kPu0LnlNzh2" +
		"KdHNkCJMHSDBie9mhKHinN75Ot2Z5oaqlobNp0wdFl7PxFj4guA7uX" +
		"fT01DglJnjIoYRZm1TmNFCaFNdEWUmcRQeEAoqSILDumIXDclFG6o0" +
		"bEsOG9vh8uRj7lhaFA4SKHpKnwxDMKkpiqjBJz5t5TV0u7umSy3HD"
	const randomName = "random_name"
	c := tconf.Config(t)
	regex := c.Validation.UsernameRegex // default: ^[a-zA-Z0-9_]{2,255}$
	tests := []struct {
		username string
		regex    string
		Err      assert.ErrorAssertionFunc
	}{
		{"", "", assert.NoError},
		{"a", "", assert.NoError},
		{"Z", "", assert.NoError},
		{"1", "", assert.NoError},
		{"_", "", assert.NoError},
		{".", "", assert.NoError},
		{" ", "", assert.NoError},
		{"aa", "", assert.NoError},
		{"aZ", "", assert.NoError},
		{"a1", "", assert.NoError},
		{"a_", "", assert.NoError},
		{"a.", "", assert.NoError},
		{"a ", "", assert.NoError},
		{"aa" + randomName, "", assert.NoError},
		{"aZ" + randomName, "", assert.NoError},
		{"a1" + randomName, "", assert.NoError},
		{"a_" + randomName, "", assert.NoError},
		{"a." + randomName, "", assert.NoError},
		{"a " + randomName, "", assert.NoError},
		{tooLongName, "", assert.NoError},
		{"", regex, assert.Error},
		{"a", regex, assert.Error},
		{"Z", regex, assert.Error},
		{"1", regex, assert.Error},
		{"_", regex, assert.Error},
		{".", regex, assert.Error},
		{" ", regex, assert.Error},
		{"aa", regex, assert.NoError},
		{"aZ", regex, assert.NoError},
		{"a1", regex, assert.NoError},
		{"a_", regex, assert.NoError},
		{"a.", regex, assert.Error},
		{"a ", regex, assert.Error},
		{" aa", regex, assert.Error},
		{"aaa", regex, assert.NoError},
		{"aZa", regex, assert.NoError},
		{"a1a", regex, assert.NoError},
		{"a_a", regex, assert.NoError},
		{"a.a", regex, assert.Error},
		{"a a", regex, assert.Error},
		{"aa" + randomName, regex, assert.NoError},
		{"aZ" + randomName, regex, assert.NoError},
		{"a1" + randomName, regex, assert.NoError},
		{"a_" + randomName, regex, assert.NoError},
		{"a." + randomName, regex, assert.Error},
		{"a " + randomName, regex, assert.Error},
		{tooLongName, regex, assert.Error},
	}
	for i, test := range tests {
		c.Validation.UsernameRegex = test.regex
		err := Username(c, test.username)
		test.Err(t, err, "(%d) username: %s regex: %s",
			i, test.username, test.regex)
	}
}
