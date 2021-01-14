package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jrapoport/gothic/models"
	"github.com/jrapoport/gothic/storage"
)

const (
	signupConfirmation   = "signup"
	recoveryConfirmation = "recovery"
)

// ConfirmParams are the parameters the Confirm endpoint accepts
type ConfirmParams struct {
	Type     string `json:"type"`
	Token    string `json:"token"`
	Password string `json:"password"`
}

// Confirm exchanges a confirmation or recovery token to a refresh token
func (a *API) Confirm(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	config := a.getConfig(ctx)

	params := &ConfirmParams{}
	cookie := r.Header.Get(useCookieHeader)
	jsonDecoder := json.NewDecoder(r.Body)
	if err := jsonDecoder.Decode(params); err != nil {
		return badRequestError("Could not read verification params: %v", err)
	}

	if params.Token == "" {
		return unprocessableEntityError("Confirm requires a token")
	}

	var (
		user  *models.User
		err   error
		token *AccessTokenResponse
	)

	err = a.db.Transaction(func(tx *storage.Connection) error {
		var terr error
		switch params.Type {
		case signupConfirmation:
			user, terr = a.signupConfirm(ctx, tx, params)
		case recoveryConfirmation:
			user, terr = a.recoverConfirm(ctx, tx, params)
		default:
			return unprocessableEntityError("Confirm requires a verification type")
		}

		if terr != nil {
			return terr
		}

		token, terr = a.issueRefreshToken(ctx, tx, user)
		if terr != nil {
			return terr
		}

		if cookie != "" && config.Cookies.Duration > 0 {
			if terr = a.setCookieToken(config, token.Token, cookie == useSessionCookie, w); terr != nil {
				return internalServerError("Failed to set JWT cookie. %s", terr)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return sendJSON(w, http.StatusOK, token)
}

func (a *API) signupConfirm(ctx context.Context, conn *storage.Connection, params *ConfirmParams) (*models.User, error) {
	config := a.getConfig(ctx)

	user, err := models.FindUserByConfirmationToken(conn, params.Token)
	if err != nil {
		if models.IsNotFoundError(err) {
			return nil, notFoundError(err.Error())
		}
		return nil, internalServerError("Name error finding user").WithInternalError(err)
	}

	err = conn.Transaction(func(tx *storage.Connection) error {
		var terr error
		if user.EncryptedPassword == "" {
			if user.InvitedAt != nil {
				if err = a.validatePassword(params.Password); err != nil {
					return err
				}
				if terr = user.UpdatePassword(tx, params.Password); terr != nil {
					return internalServerError("error storing password").WithInternalError(terr)
				}
			}
		}

		if terr = models.NewAuditLogEntry(tx, user, models.UserSignedUpAction, nil); terr != nil {
			return terr
		}

		if terr = triggerEventHooks(ctx, tx, SignupEvent, user, config); terr != nil {
			return terr
		}

		if terr = user.Confirm(tx); terr != nil {
			return internalServerError("Error confirming user").WithInternalError(terr)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (a *API) recoverConfirm(ctx context.Context, conn *storage.Connection, params *ConfirmParams) (*models.User, error) {
	config := a.getConfig(ctx)
	user, err := models.FindUserByRecoveryToken(conn, params.Token)
	if err != nil {
		if models.IsNotFoundError(err) {
			return nil, notFoundError(err.Error())
		}
		return nil, internalServerError("Name error finding user").WithInternalError(err)
	}

	err = conn.Transaction(func(tx *storage.Connection) error {
		var terr error
		if terr = user.Recover(tx); terr != nil {
			return terr
		}
		if !user.IsConfirmed() {
			if terr = models.NewAuditLogEntry(tx, user, models.UserSignedUpAction, nil); terr != nil {
				return terr
			}

			if terr = triggerEventHooks(ctx, tx, SignupEvent, user, config); terr != nil {
				return terr
			}
			if terr = user.Confirm(tx); terr != nil {
				return terr
			}
		}
		return nil
	})

	if err != nil {
		return nil, internalServerError("Name error updating user").WithInternalError(err)
	}
	return user, nil
}
