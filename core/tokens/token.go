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
	err := conn.Transaction(func(tx *store.Connection) error {
		err := tx.FirstOrCreate(t, "user_id = ?", userID).Error
		if err != nil {
			return err
		}
		// do we need to reissue this token?
		if t.Usable() {
			return nil
		}
		err = tx.Delete(t).Error
		if err != nil {
			return err
		}
		t, err = grantToken(tx, userID, issueToken)
		return err
	})
	if err != nil {
		return nil, err
	}
	return t, nil
}

func revokeAll(conn *store.Connection, t token.Token, userID uuid.UUID) error {
	return conn.Unscoped().Where("user_id = ?", userID).Delete(t).Error
}

// UseToken burns a usable token
func UseToken(conn *store.Connection, t token.Token) error {
	if !t.Usable() {
		return errors.New("invalid")
	}
	return conn.Transaction(func(tx *store.Connection) error {
		t.Use()
		err := tx.Save(t).Error
		if err != nil {
			return err
		}
		if t.Usable() {
			return nil
		}
		return tx.Delete(t).Error
	})
}
