package store

import (
	"errors"

	"gorm.io/gorm"
)

// Has returns true if the store contains the object, otherwise false.
// If an error besides not found occurs, false and the error are returned.
func Has(tx *gorm.DB, model interface{}, c ...interface{}) (bool, error) {
	err := tx.First(model, c...).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// HasLast returns true if the store contains the object, otherwise false.
// If an error besides not found occurs, false and the error are returned.
func HasLast(tx *gorm.DB, model interface{}, c ...interface{}) (bool, error) {
	err := tx.Last(model, c...).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
