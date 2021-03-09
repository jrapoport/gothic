package tutils

import (
	"net/mail"

	"github.com/jrapoport/gothic/utils"
)

// RandomEmail returns a random email address for tests.
func RandomEmail() string {
	account := utils.RandomUsername()
	account += utils.RandomPIN(4)
	return account + "@example.com"
}

// RandomAddress returns a random email address for tests.
func RandomAddress() string {
	a := &mail.Address{
		Name:    utils.RandomUsername(),
		Address: RandomEmail(),
	}
	return a.String()
}
