package api

import (
	"encoding/json"
	"github.com/jrapoport/gothic/models"
	"github.com/jrapoport/gothic/storage"
	"net/http"

	"github.com/google/uuid"
)

// EmailConfirmationParams holds the parameters for an email confirmation request.
type EmailConfirmationParams struct {
	Id string `json:"id"`
}

// EmailConfirmation sends a confirmation email
func (a *API) EmailConfirmation(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	config := a.getConfig(ctx)
	params := &EmailConfirmationParams{}
	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(params)
	if err != nil {
		return badRequestError("Could not read verification params: %v", err)
	}

	if params.Id == "" {
		return unprocessableEntityError("resending confirmation requires an id")
	}

	uid, err := uuid.Parse(params.Id)
	if err != nil {
		return internalServerError("invalid params").WithInternalError(err)
	}
	user, err := models.FindUserByID(a.db, uid)
	if err != nil {
		if models.IsNotFoundError(err) {
			return notFoundError(err.Error())
		}
		return internalServerError("Database error finding user").WithInternalError(err)
	}

	if user.IsConfirmed() {
		return badRequestError("A user with this email address has already been confirmed")
	}

	err = a.db.Transaction(func(tx *storage.Connection) error {
		if terr := models.NewAuditLogEntry(tx, user, models.UserConfirmationRequestedAction, nil); terr != nil {
			return terr
		}

		mailer := a.Mailer(ctx)
		referrer := a.getReferrer(r)
		return sendConfirmation(tx, user, mailer, config.SMTP.MaxFrequency, referrer)
	})
	if err != nil {
		return internalServerError("Error resending confirmation").WithInternalError(err)
	}

	return sendJSON(w, http.StatusOK, &map[string]string{})
}
