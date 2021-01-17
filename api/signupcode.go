package api

import (
	"errors"

	"github.com/jrapoport/gothic/models"
	"github.com/jrapoport/gothic/storage"
)

// NewSignupCodes generates a list of unique signup codes.
func (a *API) NewSignupCodes(f models.Format, t models.Type, count int) ([]*models.SignupCode, error) {
	codes := make([]*models.SignupCode, count)
	err := a.db.Transaction(func(tx *storage.Connection) error {
		for i := 0; i < count; i++ {
			code, err := models.CreateSignupCode(tx, f, t, true)
			if err != nil {
				return err
			}
			codes[i] = code
		}
		return nil
	})
	return codes, err
}

// NewSignupCode generates a unique signup code.
func (a *API) NewSignupCode(f models.Format, t models.Type) (code *models.SignupCode, err error) {
	return models.CreateSignupCode(a.db, f, t, true)
}

func (a *API) CheckSignupCode(code string) error {
	if !a.config.Signup.Code {
		return nil
	} else if code == "" {
		return badRequestError("signup code required")
	}
	available, err := models.CanUseSignupCode(a.db, code)
	if err != nil {
		return err
	} else if !available {
		return errors.New("invalid signup code")
	}
	return nil
}
