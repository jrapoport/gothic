package storage

import "gorm.io/gorm"

type errNotFound struct {
	error
}

// ErrNotFound is returned when something is not found in the database.
var ErrNotFound = errNotFound{gorm.ErrRecordNotFound}

func wrapError(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
		return ErrNotFound
	}
	return err
}
