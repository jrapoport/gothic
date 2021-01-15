package api

import (
	"errors"
	"math/rand"
	"time"

	"github.com/jrapoport/gothic/models"
)

// CodeFormat is the format of the access code.
type CodeFormat int

const (
	// CodeFormatPIN is a PIN style access code.
	CodeFormatPIN CodeFormat = iota
)

// CodeType is the type of code (e.g. single use)
type CodeType int

const (
	// CodeTypeSingleUse is a code that can be used once.
	CodeTypeSingleUse CodeType = iota

	// CodeTypeMultiUse is a code that can be used repeatedly.
	CodeTypeMultiUse
)

const (
	// PINCodeMaxLen is the max length of a PIN (6).
	PINCodeMaxLen = 6

	// DebugPINCode is bypass code that is only valid in debug builds.
	DebugPINCode = "000000"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// NewRandomAccessCode generates a new & unique access code of the format and type.
func (a *API) NewRandomAccessCode(fmt CodeFormat, t CodeType) (*models.AccessCode, error) {
	var code string
	switch fmt {
	case CodeFormatPIN:
		for {
			code = randomPINCodeN(PINCodeMaxLen)
			if code != DebugPINCode {
				break
			}
		}
	default:
		return nil, errors.New("unknown code format")
	}
	exists, err := models.IsDuplicatedCode(a.db, code)
	if err != nil {
		return nil, err
	}
	if exists {
		return a.NewRandomAccessCode(fmt, t)
	}
	return &models.AccessCode{Code: code, Format: int(fmt), Type: int(t)}, nil
}

// NewRandomAccessCodes generates a slice of new & unique access code of the format and type.
func (a *API) NewRandomAccessCodes(fmt CodeFormat, t CodeType, count int) ([]*models.AccessCode, error) {
	codes := make([]*models.AccessCode, count)
	for i := 0; i < count; i++ {
		code, err := a.NewRandomAccessCode(fmt, t)
		if err != nil {
			return nil, err
		}
		codes[i] = code
	}
	return codes, nil
}

func randomPINCodeN(length int) string {
	const pool = "1234567890"
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = pool[rand.Intn(len(pool))]
	}
	return string(bytes)
}
