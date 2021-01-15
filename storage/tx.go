package storage

import "gorm.io/gorm"

func First(tx *gorm.DB, v interface{}) error {
	return wrapError(tx.First(v).Error)
}
