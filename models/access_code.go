package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/storage"
	"gorm.io/gorm"
)

func init() {
	storage.AddMigration(&AccessCode{})
}

// AccessCode can be used for signup access codes.
type AccessCode struct {
	gorm.Model
	Code      string     `json:"code" gorm:"index:access_code_idx"`
	Format    int        `json:"format"`
	Type      int        `json:"type"`
	UserID    *uuid.UUID `json:"user_id"`
	Invalid   bool       `json:"invalid"`
	InvalidAt *time.Time `json:"invalid_at"`
}

// IsDuplicatedCode returns whether a user exists with a matching email and audience.
func IsDuplicatedCode(tx *storage.Connection, code string) (bool, error) {
	_, err := FindAccessCodeByCode(tx, code)
	if errors.Is(err, storage.ErrNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// FindUserByEmail finds a user with the matching email and audience.
func FindAccessCodeByCode(tx *storage.Connection, code string) (*AccessCode, error) {
	q := tx.Where("code = ?", code)
	c := &AccessCode{}
	return c, storage.First(q, c)
}
