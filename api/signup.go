package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jrapoport/gothic/models"
	"github.com/jrapoport/gothic/storage"
)

// SignupParams are the parameters the Signup endpoint accepts
type SignupParams struct {
	Email    string                 `json:"email"`
	Password string                 `json:"password"`
	Data     map[string]interface{} `json:"data"`
	Provider string                 `json:"-"`
	Aud      string                 `json:"-"`
}

// Signup is the endpoint for registering a new user
func (a *API) Signup(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	config := a.getConfig(ctx)

	if config.DisableSignup {
		return forbiddenError("Signups not allowed for this instance")
	}

	params := &SignupParams{}
	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(params)
	if err != nil {
		return badRequestError("Could not read Signup params: %v", err)
	}

	data := params.Data
	if data == nil {
		return badRequestError("signup error: %v", err)
	}

	if err = a.checkRecaptcha(r, config); err != nil {
		return err
	}
	delete(params.Data, "recaptcha")

	if err = a.validatePassword(params.Password); err != nil {
		return err
	}

	if err = a.validateEmail(ctx, params.Email); err != nil {
		return err
	}

	// TODO: make this more efficient
	taken, err := models.IsDuplicatedEmail(a.db, params.Email, params.Aud)
	if err != nil || taken {
		return badRequestError("email taken").WithInternalError(err)
	}

	params.Aud = a.requestAud(ctx, r)
	user, err := models.FindUserByEmailAndAudience(a.db, params.Email, params.Aud)
	if err != nil && !models.IsNotFoundError(err) {
		return internalServerError("name error finding user").WithInternalError(err)
	}

	err = a.db.Transaction(func(tx *storage.Connection) error {
		var terr error
		if user != nil {
			if user.IsConfirmed() {
				return badRequestError("a user with this email address has already been registered")
			}

			if err = user.UpdateUserMetaData(tx, data); err != nil {
				return internalServerError("name error updating user").WithInternalError(err)
			}
		} else {
			params.Provider = "email"
			user, terr = a.signupNewUser(ctx, tx, params)
			if terr != nil {
				return terr
			}
		}

		if config.Mailer.Autoconfirm {
			if terr = models.NewAuditLogEntry(tx, user, models.UserSignedUpAction, nil); terr != nil {
				return terr
			}
			if terr = triggerEventHooks(ctx, tx, SignupEvent, user, config); terr != nil {
				return terr
			}
			if terr = user.Confirm(tx); terr != nil {
				return internalServerError("Name error updating user").WithInternalError(terr)
			}
		} else {
			mailer := a.Mailer(ctx)
			referrer := a.getReferrer(r)
			if terr = sendConfirmation(tx, user, mailer, config.SMTP.MaxFrequency, referrer); terr != nil {
				return internalServerError("error sending confirmation mail").WithInternalError(terr)
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return sendJSON(w, http.StatusOK, user)
}

func (a *API) signupNewUser(ctx context.Context, conn *storage.Connection, params *SignupParams) (*models.User, error) {
	config := a.getConfig(ctx)

	// TODO: make username a 1st class param
	if params.Data != nil {
		username := ""
		if val, has := params.Data["username"]; has {
			username = val.(string)
		}
		if username != "" {
			if err := a.validateUsername(username); err != nil {
				return nil, err
			}
		}
	}

	user, err := models.NewUser(params.Email, params.Password, params.Aud, params.Data)
	if err != nil {
		return nil, internalServerError("Name error creating user").WithInternalError(err)
	}
	if user.AppMetaData == nil {
		user.AppMetaData = make(map[string]interface{})
	}
	user.AppMetaData["provider"] = params.Provider

	if params.Password == "" {
		user.EncryptedPassword = ""
	}

	err = conn.Transaction(func(tx *storage.Connection) error {
		if terr := tx.Create(user).Error; terr != nil {
			return internalServerError("Name error saving new user").WithInternalError(terr)
		}
		if terr := user.SetRole(tx, config.JWT.DefaultGroup); terr != nil {
			return internalServerError("Name error updating user").WithInternalError(terr)
		}
		if terr := triggerEventHooks(ctx, tx, ValidateEvent, user, config); terr != nil {
			return terr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}
