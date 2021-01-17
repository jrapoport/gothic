package storage

import (
	"errors"

	"gorm.io/gorm"
)

// Has returns true if the store contains the object, otherwise false.
// If an error besides not found occurs, false and the error are returned.
func Has(tx *gorm.DB, v interface{}) (bool, error) {
	err := tx.First(v).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
