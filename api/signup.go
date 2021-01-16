package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jrapoport/gothic/models"
	"github.com/jrapoport/gothic/storage"
	"github.com/jrapoport/gothic/utils"
)

// SignupParams are the parameters the Signup endpoint accepts
type SignupParams struct {
	Email    string                 `json:"email"`
	Password string                 `json:"password"`
	Data     map[string]interface{} `json:"data"`
	Provider string                 `json:"-"`
}

// Signup is the endpoint for registering a new user
func (a *API) Signup(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	config := a.getConfig(ctx)

	if config.DisableSignup {
		return forbiddenError("signup disabled")
	}

	params := &SignupParams{}
	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(params)
	if err != nil {
		return badRequestError("invalid params: %v", err)
	}

	if recaptcha, ok := params.Data["recaptcha"].(string); ok {
		ipaddr := a.getClientIP(r)
		err = a.checkRecaptcha(ipaddr, recaptcha)
		if err != nil {
			return err
		}
		delete(params.Data, "recaptcha")
	}

	code, ok := params.Data["code"].(string)
	err = a.CheckSignupCode(code)
	if err != nil {
		return err
	}
	if ok {
		delete(params.Data, "code")
	}

	if err = a.validatePassword(params.Password); err != nil {
		return err
	}

	if err = a.validateEmail(ctx, params.Email); err != nil {
		return err
	}

	// This user has already signed up
	taken, err := models.IsDuplicatedEmail(a.db, params.Email)
	if err != nil || taken {
		return badRequestError("email taken").WithInternalError(err)
	}

	var user *models.User
	err = a.db.Transaction(func(tx *storage.Connection) error {
		params.Provider = "email"
		user, err = a.signupNewUser(ctx, tx, params)
		if err != nil {
			return err
		}

		if config.Mailer.Autoconfirm {
			if err = models.NewAuditLogEntry(tx, user, models.UserSignedUpAction, nil); err != nil {
				return err
			}
			if err = triggerEventHooks(ctx, tx, SignupEvent, user, config); err != nil {
				return err
			}
			if err = user.Confirm(tx); err != nil {
				return internalServerError("Name error updating user").WithInternalError(err)
			}
		} else {
			mailer := a.Mailer(ctx)
			referrer := a.getReferrer(r)
			if err = sendConfirmation(tx, user, mailer, config.SMTP.MaxFrequency, referrer); err != nil {
				return internalServerError("error sending confirmation mail").WithInternalError(err)
			}
		}
		if code == "" {
			return nil
		}
		if err = models.UseSignupCode(tx, code, user); err != nil {
			return internalServerError("error saving signup code").WithInternalError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	token, err := a.issueAccessToken(ctx, w, r, user)
	if err != nil {
		return internalServerError("failed to issue access token").WithInternalError(err)
	}

	return sendJSON(w, http.StatusOK, token)
}

func (a *API) signupNewUser(ctx context.Context, conn *storage.Connection, params *SignupParams) (*models.User, error) {
	config := a.getConfig(ctx)

	if a.config.Signup.Defaults {
		if _, ok := params.Data["username"]; !ok {
			// we didn't find a user name so we'll make one up that's random and unique
			params.Data["username"] = a.randomUsername()
		}
		if _, ok := params.Data["color"]; !ok {
			params.Data["color"] = utils.RandomColor()
		}
	}

	username, ok := params.Data["username"].(string)
	if ok {
		delete(params.Data, "username")
	}

	user, err := models.NewUser(params.Email, params.Password, params.Data)
	if err != nil {
		return nil, internalServerError("Name error creating user").WithInternalError(err)
	}

	user.Username = username

	// if there is no username this will not return
	// an error since a username is not required.
	if err = a.validateUsername(user.Username); err != nil {
		return nil, err
	}

	if user.AppMetaData == nil {
		user.AppMetaData = map[string]interface{}{}
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

func (a *API) randomUsername() string {
	username := ""
	for {
		username = utils.RandomUsername(28)
		taken, err := models.IsDuplicatedUsername(a.db, username)
		if err != nil {
			a.config.Log.Warn(err)
			break
		}
		if !taken {
			break
		}
	}
	return username
}
