package api

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/google/uuid"
	"github.com/jrapoport/gothic/models"
	"github.com/jrapoport/gothic/storage"
)

// UserUpdateParams parameters for updating a user
type UserUpdateParams struct {
	Email            string                 `json:"email"`
	Password         string                 `json:"password"`
	EmailChangeToken string                 `json:"email_change_token"`
	Data             map[string]interface{} `json:"data"`
	AppData          map[string]interface{} `json:"app_metadata,omitempty"`
}

// UserGet returns a user
func (a *API) UserGet(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	claims := getClaims(ctx)
	if claims == nil {
		return badRequestError("could not read claims")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return badRequestError("could not read username id claim")
	}

	if err = claims.VerifyAudience(jwt.DefaultValidationHelper, a.config.JWT.Aud); err != nil {
		return badRequestError("token audience doesn't match request audience")
	}

	user, err := models.FindUserByID(a.db, userID)
	if err != nil {
		if models.IsNotFoundError(err) {
			return notFoundError(err.Error())
		}
		return internalServerError("name error finding user").WithInternalError(err)
	}

	return sendJSON(w, http.StatusOK, user)
}

// UserUpdate updates fields on a user
func (a *API) UserUpdate(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	params := &UserUpdateParams{}
	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(params)
	if err != nil {
		return badRequestError("could not read Username Update params: %v", err)
	}

	claims := getClaims(ctx)
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return badRequestError("could not read username id claim")
	}

	user, err := models.FindUserByID(a.db, userID)
	if err != nil {
		if models.IsNotFoundError(err) {
			return notFoundError(err.Error())
		}
		return internalServerError("name error finding user").WithInternalError(err)
	}

	log := getLogEntry(r)
	log.Debugf("checking params for token %v", params)

	err = a.db.Transaction(func(tx *storage.Connection) error {
		var terr error
		if params.Password != "" {
			if err = a.validatePassword(params.Password); err != nil {
				return err
			}
			if terr = user.UpdatePassword(tx, params.Password); terr != nil {
				return internalServerError("error during password storage").WithInternalError(terr)
			}
		}

		if params.Data != nil {
			if terr = user.UpdateUserMetaData(tx, params.Data); terr != nil {
				return internalServerError("error updating user").WithInternalError(terr)
			}
		}

		if params.AppData != nil {
			if !a.isAdmin(ctx, user) {
				return unauthorizedError("Updating app_metadata requires admin privileges")
			}

			if terr = user.UpdateAppMetaData(tx, params.AppData); terr != nil {
				return internalServerError("Error updating user").WithInternalError(terr)
			}
		}

		if params.EmailChangeToken != "" {
			log.Debugf("Got change token %v", params.EmailChangeToken)

			if params.EmailChangeToken != user.EmailChangeToken {
				return unauthorizedError("Email Change Token didn't match token on file")
			}

			if terr = user.ConfirmEmailChange(tx); terr != nil {
				return internalServerError("Error updating user").WithInternalError(terr)
			}
		} else if params.Email != "" && params.Email != user.Email {
			if terr = a.validateEmail(ctx, params.Email); terr != nil {
				return terr
			}

			var exists bool
			if exists, terr = models.IsDuplicatedEmail(tx, params.Email); terr != nil {
				return internalServerError("Name error checking email").WithInternalError(terr)
			} else if exists {
				return unprocessableEntityError("Email address already registered by another user")
			}

			mailer := a.Mailer(ctx)
			referrer := a.getReferrer(r)
			if terr = a.sendEmailChange(tx, user, mailer, params.Email, referrer); terr != nil {
				return internalServerError("Error sending change email").WithInternalError(terr)
			}
		}

		if terr = models.NewAuditLogEntry(tx, user, models.UserModifiedAction, nil); terr != nil {
			return internalServerError("Error recording audit log entry").WithInternalError(terr)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return sendJSON(w, http.StatusOK, user)
}

func (a *API) validatePassword(password string) error {
	if password == "" {
		return unprocessableEntityError("password is required")
	}
	if a.config.Validation.PasswordRegex == "" {
		return nil
	}
	rex, err := regexp.Compile(a.config.Validation.PasswordRegex)
	if err != nil {
		return err
	}
	if !rex.MatchString(password) {
		return unprocessableEntityError("invalid password")
	}
	return nil
}

func (a *API) validateUsername(username string) error {
	if username == "" {
		return unprocessableEntityError("username is required")
	}
	if a.config.Validation.UsernameRegex == "" {
		return nil
	}
	rex, err := regexp.Compile(a.config.Validation.PasswordRegex)
	if err != nil {
		return err
	}
	if !rex.MatchString(username) {
		return unprocessableEntityError("invalid username")
	}
	return nil
}
