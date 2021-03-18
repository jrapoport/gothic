package validate

import (
	"testing"

	"github.com/jrapoport/gothic/test/tconf"
	"github.com/stretchr/testify/assert"
)

func TestPassword(t *testing.T) {
	t.Parallel()
	const randomPass8 = "7r/M3Z&F"
	const randomPass8AN = "isBSfVB9"
	const randomPass8NO = "01234569"
	const randomPass8MC = "dnGfemqX"
	const randomPass8UC = "PSKBRQUX"
	const randomPass8LC = "pskbrqux"
	const randomPass28 = "u)}H@j!Up*}dH~9(aH$K5S{x*,@?"
	const randomPass40 = "Pd5ajK7?XH-We@k^Jp$/Df=x{eQxXW:m'CqQhPf{"
	const tooLongPass = randomPass40 + randomPass8
	c := tconf.Config(t)
	regex := c.Validation.PasswordRegex // default: ^[a-zA-Z0-9[:punct:]]{8,40}$
	tests := []struct {
		pw    string
		regex string
		Err   assert.ErrorAssertionFunc
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
		{"ZZ", "", assert.NoError},
		{"a1", "", assert.NoError},
		{"Z1", "", assert.NoError},
		{"11", "", assert.NoError},
		{"_a", "", assert.NoError},
		{"_1", "", assert.NoError},
		{" a", "", assert.NoError},
		{randomPass8, "", assert.NoError},
		{randomPass8AN, "", assert.NoError},
		{randomPass8NO, "", assert.NoError},
		{randomPass8MC, "", assert.NoError},
		{randomPass8UC, "", assert.NoError},
		{randomPass8LC, "", assert.NoError},
		{randomPass28, "", assert.NoError},
		{randomPass40, "", assert.NoError},
		{tooLongPass, "", assert.NoError},
		{"", regex, assert.Error},
		{"a", regex, assert.Error},
		{"Z", regex, assert.Error},
		{"1", regex, assert.Error},
		{"_", regex, assert.Error},
		{".", regex, assert.Error},
		{" ", regex, assert.Error},
		{"aa", regex, assert.Error},
		{"aZ", regex, assert.Error},
		{"ZZ", regex, assert.Error},
		{"a1", regex, assert.Error},
		{"Z1", regex, assert.Error},
		{"11", regex, assert.Error},
		{"_a", regex, assert.Error},
		{"_1", regex, assert.Error},
		{" a", regex, assert.Error},
		{randomPass8, regex, assert.NoError},
		{randomPass8AN, regex, assert.NoError},
		{randomPass8NO, regex, assert.NoError},
		{randomPass8MC, regex, assert.NoError},
		{randomPass8UC, regex, assert.NoError},
		{randomPass8LC, regex, assert.NoError},
		{randomPass28, regex, assert.NoError},
		{randomPass40, regex, assert.NoError},
		{tooLongPass, regex, assert.Error},
	}
	for _, test := range tests {
		c.Validation.PasswordRegex = test.regex
		err := Password(c, test.pw)
		test.Err(t, err, "pw: %s regex: %s", test.pw, test.regex)
	}
}
