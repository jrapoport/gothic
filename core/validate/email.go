package validate

import (
	"net/mail"

	"github.com/badoux/checkmail"
)

// Email validate an email offline
func Email(address string) (string, error) {
	e, err := mail.ParseAddress(address)
	if err != nil {
		return "", err
	}
	address = e.Address
	err = checkmail.ValidateFormat(address)
	if err != nil {
		return "", err
	}
	return address, nil
}

// EmailAccount validate an email online with an smtp server
func EmailAccount(host, from, address string) (string, error) {
	var err error
	address, err = Email(address)
	if err != nil {
		return "", err
	}
	err = checkmail.ValidateHostAndUser(host, from, address)
	if err != nil {
		return "", err
	}
	return address, nil
}
