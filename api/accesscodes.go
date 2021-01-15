package api

import (
	"errors"
	"math/rand"
	"time"

	"github.com/jrapoport/gothic/models"
)

type CodeFormat int

const (
	CodeFormatPIN CodeFormat = iota
)

type CodeType int

const (
	CodeTypeSingleUse CodeType = iota
	CodeTypeMultiUse
)

const (
	PINCodeMaxLen = 6

	DebugPINCode = "000000"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

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
