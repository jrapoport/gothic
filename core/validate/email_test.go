package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmailTemplate_Valid(t *testing.T) {
	t.Parallel()
	tests := []struct {
		email    string
		expected string
		Err      assert.ErrorAssertionFunc
	}{
		// valid addresses
		{"email@example.com", "email@example.com", assert.NoError},
		{"firstname.lastname@example.com", "firstname.lastname@example.com", assert.NoError},
		{"email@subdomain.example.com", "email@subdomain.example.com", assert.NoError},
		{"firstname+lastname@example.com", "firstname+lastname@example.com", assert.NoError},
		{"email@123.123.123.123", "email@123.123.123.123", assert.NoError},
		{"1234567890@example.com", "1234567890@example.com", assert.NoError},
		{"email@example-one.com", "email@example-one.com", assert.NoError},
		{"_______@example.com", "_______@example.com", assert.NoError},
		{"email@example.name", "email@example.name", assert.NoError},
		{"email@example.museum", "email@example.museum", assert.NoError},
		{"email@example.co.jp", "email@example.co.jp", assert.NoError},
		{"firstname-lastname@example.com", "firstname-lastname@example.com", assert.NoError},
		{"Joe Smith <email@example.com>", "email@example.com", assert.NoError},
		{"email@example.com (Joe Smith)", "email@example.com", assert.NoError},
		// technically the following might be invalid
		{"email@example", "email@example", assert.NoError},
		{"email@111.222.333.44444", "email@111.222.333.44444", assert.NoError},
		{"\"email\"@example.com", "email@example.com", assert.NoError},
		// invalid addresses
		{"", "", assert.Error},
		{"plainaddress", "", assert.Error},
		{"#@%^%#$@#$@#.com", "", assert.Error},
		{"@example.com", "", assert.Error},
		{"email.example.com", "", assert.Error},
		{"email@example@example.com", "", assert.Error},
		{"あいうえお@example.com", "", assert.Error},
		{"email@-example.com", "", assert.Error},
		{"email@example..com", "", assert.Error},
		{`”(),:;<>[\]@example.com`, "", assert.Error},
		{`just”not”right@example.com`, "", assert.Error},
		{`this\ is"really"not\allowed@example.com`, "", assert.Error},
		// technically the following might be valid
		{".email@example.com", "", assert.Error},
		{"email.@example.com", "", assert.Error},
		{"email..email@example.com", "", assert.Error},
		{"Abc..123@example.com", "", assert.Error},
		{"email@[123.123.123.123]", "", assert.Error},
		{`much.”more\ unusual”@example.com`, "", assert.Error},
		{`very.unusual.”@”.unusual.com@example.com`, "", assert.Error},
		{`very.”(),:;<>[]”.VERY.”very@\\ "very”.unusual@strange.example.com`, "", assert.Error},
	}
	for _, test := range tests {
		e, err := Email(test.email)
		test.Err(t, err, test.email)
		assert.Equal(t, test.expected, e)
	}
}
