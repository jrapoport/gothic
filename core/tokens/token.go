package tokens

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/models/token"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
)

func grantToken(conn *store.Connection, userID uuid.UUID, issueToken func() token.Token) (token.Token, error) {
	if userID == user.SystemID {
		return nil, errors.New("system user")
	}
	t := issueToken()
	err := conn.FirstOrCreate(t, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return t, nil
}

func revokeAll(conn *store.Connection, t token.Token, userID uuid.UUID) error {
	return conn.Unscoped().Where("user_id = ?", userID).Delete(t).Error
}
