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

// CreateSignupCode creates a signup code.
func CreateSignupCode(conn *store.Connection, userID uuid.UUID, f code.Format, uses int, unique bool) (*code.SignupCode, error) {
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

// CreateSignupCodes generates a list of unique signup codes.
func CreateSignupCodes(conn *store.Connection, userID uuid.UUID, f code.Format, uses, count int) ([]*code.SignupCode, error) {
	if count < 0 {
		count = 0
	}
	list := make([]*code.SignupCode, count)
	err := conn.Transaction(func(tx *store.Connection) error {
		for i := 0; i < count; i++ {
			c, err := CreateSignupCode(tx, userID, f, uses, true)
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

// GetSignupCode finds the signup code that matches code.
func GetSignupCode(conn *store.Connection, tok string) (*code.SignupCode, error) {
	if !utils.IsValidCode(tok) {
		return nil, errors.New("invalid code")
	}
	sc := new(code.SignupCode)
	if utils.IsDebugPIN(tok) {
		sc.Format = code.PIN
		sc.Token = tok
		return sc, nil
	}
	err := conn.First(sc, "token = ?", tok).Error
	if err != nil {
		return nil, err
	}
	return sc, nil
}

// GetUsableSignupCode returns not nil if a signup code is found that can be used
func GetUsableSignupCode(conn *store.Connection, tok string) (*code.SignupCode, error) {
	sc, err := GetSignupCode(conn, tok)
	if err != nil {
		return nil, err
	}
	if !sc.Usable() {
		err = fmt.Errorf("%w: %s", code.ErrUnusableCode, tok)
		return nil, err
	}
	return sc, nil
}

// SignupCodeSent mark a code as sent.
func SignupCodeSent(conn *store.Connection, sc *code.SignupCode) error {
	now := time.Now().UTC()
	sc.SentAt = &now
	return conn.Model(sc).Update("sent_at", sc.SentAt).Error
}

// GetLastSentSignupCode returns the last usable code sent by a user
func GetLastSentSignupCode(conn *store.Connection, userID uuid.UUID) (*code.SignupCode, error) {
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

// VoidSignupCode removes a signup code from use
func VoidSignupCode(conn *store.Connection, tok string) error {
	return conn.Transaction(func(tx *store.Connection) error {
		sc, err := GetUsableSignupCode(tx, tok)
		if err != nil {
			return err
		}
		return tx.Delete(sc).Error
	})
}
