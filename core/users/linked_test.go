package users

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/models/account"
	"github.com/jrapoport/gothic/store/types/provider"
	"github.com/jrapoport/gothic/test/tconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLinkAccount(t *testing.T) {
	t.Parallel()
	conn, c := tconn.TempConn(t)
	u := testUser(t, conn, c.Provider())
	l := account.LinkedAccount{}
	l.Type = -1
	err := LinkAccount(conn, nil, &l)
	assert.Error(t, err)
	err = LinkAccount(conn, u, &l)
	assert.Error(t, err)
	l.Type = account.Auth
	err = LinkAccount(conn, u, &l)
	assert.Error(t, err)
	l.UserID = uuid.New()
	err = LinkAccount(conn, u, &l)
	assert.Error(t, err)
	l.UserID = uuid.Nil
	l.Provider = c.Provider()
	err = LinkAccount(conn, u, &l)
	assert.Error(t, err)
	l.AccountID = uuid.New().String()
	err = LinkAccount(conn, u, &l)
	assert.NoError(t, err)
	err = LinkAccount(conn, u, &l)
	assert.Error(t, err)
	l2 := account.LinkedAccount{
		Provider:  c.Provider(),
		AccountID: l.AccountID,
	}
	err = LinkAccount(conn, u, &l2)
	assert.Error(t, err)
}

func TestGetLinkedUser(t *testing.T) {
	var accountID1 = uuid.New().String()
	conn, c := tconn.TempConn(t)
	u := testUser(t, conn, c.Provider())
	l := account.LinkedAccount{
		Provider:  c.Provider(),
		AccountID: accountID1,
	}
	err := LinkAccount(conn, u, &l)
	assert.NoError(t, err)
	lu, err := HasLinkedUser(conn, c.Provider(), accountID1)
	assert.NoError(t, err)
	assert.Equal(t, u.ID, lu.ID)
	lu, err = HasLinkedUser(conn, c.Provider(), "")
	assert.NoError(t, err)
	assert.Nil(t, lu)
	lu, err = HasLinkedUser(conn, provider.Unknown, accountID1)
	assert.NoError(t, err)
	assert.Nil(t, lu)
	var accountID2 = uuid.New().String()
	badLink := account.LinkedAccount{
		UserID:    uuid.New(),
		Provider:  c.Provider(),
		AccountID: accountID2,
	}
	err = conn.Create(&badLink).Error
	require.NoError(t, err)
	_, err = HasLinkedUser(conn, c.Provider(), accountID2)
	assert.Error(t, err)
}
