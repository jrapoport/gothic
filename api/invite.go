package api

import (
	"encoding/json"
	"net/http"

	"github.com/jrapoport/gothic/models"
	"github.com/jrapoport/gothic/storage"
)

// InviteParams are the parameters the Signup endpoint accepts
type InviteParams struct {
	Email string                 `json:"email"`
	Data  map[string]interface{} `json:"data"`
}

// Invite is the endpoint for inviting a new user
func (a *API) Invite(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	adminUser := getAdminUser(ctx)
	params := &InviteParams{}

	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(params)
	if err != nil {
		return badRequestError("Could not read Invite params: %v", err)
	}

	if err = a.validateEmail(ctx, params.Email); err != nil {
		return err
	}

	user, err := models.FindUserByEmail(a.db, params.Email)
	if err != nil && !models.IsNotFoundError(err) {
		return internalServerError("Name error finding user").WithInternalError(err)
	}
	if user != nil {
		return unprocessableEntityError("Email address already registered by another user")
	}

	err = a.db.Transaction(func(tx *storage.Connection) error {
		signupParams := SignupParams{
			Email:    params.Email,
			Data:     params.Data,
			Provider: "email",
		}
		user, err = a.signupNewUser(ctx, tx, &signupParams)
		if err != nil {
			return err
		}

		if terr := models.NewAuditLogEntry(tx, adminUser, models.UserInvitedAction, map[string]interface{}{
			"user_id":    user.ID,
			"user_email": user.Email,
		}); terr != nil {
			return terr
		}

		mailer := a.Mailer(ctx)
		referrer := a.getReferrer(r)
		if err = sendInvite(tx, user, mailer, referrer); err != nil {
			return internalServerError("error inviting user").WithInternalError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return sendJSON(w, http.StatusOK, user)
}
