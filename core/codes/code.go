package codes

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/models/code"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/utils"
)

// CreateCode creates a signup code.
func CreateCode(conn *store.Connection, userID uuid.UUID, f code.Format, uses int, unique bool) (*code.SignupCode, error) {
	var sc *code.SignupCode
	err := conn.Transaction(func(tx *store.Connection) error {
		for {
			sc = code.NewSignupCode(userID, f, uses)
			if sc == nil {
				return errors.New("invalid code")
			}
			if !unique {
				break
			}
			has, err := sc.HasCode(tx)
			if err != nil {
				return err
			}
			if !has {
				break
			}
		}
		return tx.Create(sc).Error
	})
	if err != nil {
		return nil, err
	}
	return sc, nil
}

// CreateCodes generates a list of unique signup codes.
func CreateCodes(conn *store.Connection, userID uuid.UUID, f code.Format, uses, count int) ([]*code.SignupCode, error) {
	if count < 0 {
		count = 0
	}
	list := make([]*code.SignupCode, count)
	err := conn.Transaction(func(tx *store.Connection) error {
		for i := 0; i < count; i++ {
			c, err := CreateCode(tx, userID, f, uses, true)
			if err != nil {
				return err
			}
			list[i] = c
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}

// GetCode finds the signup code that matches code.
func GetCode(conn *store.Connection, c string) (*code.SignupCode, error) {
	if !utils.IsValidCode(c) {
		return nil, errors.New("invalid code")
	}
	sc := new(code.SignupCode)
	if sc.Format == code.PIN && utils.IsDebugPIN(c) {
		sc.Token = c
		return sc, nil
	}
	err := conn.First(sc, "token = ?", c).Error
	if err != nil {
		return nil, err
	}
	return sc, nil
}

// GetUsableCode returns not nil if a signup code is found that can be used
func GetUsableCode(conn *store.Connection, c string) (*code.SignupCode, error) {
	sc, err := GetCode(conn, c)
	if err != nil {
		return nil, err
	}
	if !sc.Usable() {
		err = fmt.Errorf("unusable code: %s", c)
		return nil, err
	}
	return sc, nil
}

// CodeSent mark a code as sent.
func CodeSent(conn *store.Connection, sc *code.SignupCode) error {
	now := time.Now().UTC()
	sc.SentAt = &now
	return conn.Model(sc).Update("sent_at", sc.SentAt).Error
}

// GetLastSentCode returns the last usable code sent by a user
func GetLastSentCode(conn *store.Connection, userID uuid.UUID) (*code.SignupCode, error) {
	sc := new(code.SignupCode)
	has, err := conn.HasLast(sc, "user_id = ? AND sent_at NOT NULL", userID)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return sc, nil
}
