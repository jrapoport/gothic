package core

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/core/tokens"
	"github.com/jrapoport/gothic/models/account"
	"github.com/jrapoport/gothic/models/token"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/models/types/provider"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/jrapoport/gothic/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPI_GetUser(t *testing.T) {
	t.Parallel()
	a := apiWithTempDB(t)
	u1 := testUser(t, a)
	u2, err := a.GetUser(u1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
	assert.Equal(t, u1.Email, u2.Email)
	assert.Equal(t, u1.Username, u2.Username)
	assert.Equal(t, u1.Data, u2.Data)
	_, err = a.GetUser(user.SystemID)
	assert.Error(t, err)
	_, err = a.GetUser(uuid.New())
	assert.Error(t, err)
}

func TestAPI_GetAuthenticatedUser(t *testing.T) {
	t.Parallel()
	a := apiWithTempDB(t)
	u := testUser(t, a)
	u = confirmUser(t, a, u)
	_, err := a.GetAuthenticatedUser(u.ID)
	assert.Error(t, err)
	bt, err := tokens.GrantBearerToken(a.conn, a.config.JWT, u)
	require.NoError(t, err)
	au, err := a.GetAuthenticatedUser(u.ID)
	assert.NoError(t, err)
	assert.Equal(t, u.ID, au.ID)
	err = a.conn.Delete(u).Error
	assert.NoError(t, err)
	_, err = a.GetAuthenticatedUser(u.ID)
	assert.Error(t, err)
	err = a.conn.Delete(bt.RefreshToken).Error
	assert.NoError(t, err)
	_, err = a.GetAuthenticatedUser(u.ID)
	assert.Error(t, err)
}

func TestAPI_GetUserWithEmail(t *testing.T) {
	t.Parallel()
	a := apiWithTempDB(t)
	u1 := testUser(t, a)
	u2, err := a.GetUserWithEmail(u1.EmailAddress().String())
	assert.NoError(t, err)
	assert.NotNil(t, u2)
	assert.Equal(t, u1.Email, u2.Email)
	assert.Equal(t, u1.Username, u2.Username)
	assert.Equal(t, u1.Data, u2.Data)
	_, err = a.GetUserWithEmail("does-not-exist@example.com")
	assert.Error(t, err)
}

const pageSize = 4

type testUsers struct {
	uid      uuid.UUID
	provider provider.Name
	role     user.Role
	email    string
	name     string
	data     types.Map
}

func testCreateUsers(t *testing.T, a *API) []testUsers {
	var cases = func() []testUsers {
		email := func() string { return tutils.RandomEmail() }
		name := func() string { return utils.RandomUsername() }
		p := a.Provider()
		ru := user.RoleUser
		ra := user.RoleAdmin
		rs := user.RoleSuper
		var tests = []testUsers{
			{uuid.Nil, p, ru, email(), name(), nil},
			{uuid.Nil, p, ru, email(), name(), nil},
			{uuid.Nil, p, ru, email(), name(), nil},
			{uuid.Nil, p, ru, email(), name(), nil},
			{uuid.Nil, p, ru, email(), name(), nil},
			{uuid.Nil, p, ru, email(), name(), nil},
			{uuid.Nil, p, ru, email(), name(), nil},
			{uuid.Nil, p, ra, email(), name(), nil},
			{uuid.Nil, p, ra, email(), name(), nil},
			{uuid.Nil, p, rs, email(), name(), nil},
		}
		for i, test := range tests {
			ext := test
			ext.email = email()
			ext.provider = provider.Google
			if i%2 == 0 {
				ext.provider = provider.Amazon
			}
			tests = append(tests, ext)
		}
		var one sync.Once
		for i, bk := range []interface{}{
			"thing2", testBook, uuid.New().String(),
		} {
			for x, test := range tests {
				test.data = types.Map{
					"dr_suess": "thing1",
					"book":     bk,
				}
				sld := fmt.Sprintf("salad-%d", x+i)
				test.data["extra"] = sld
				one.Do(func() {
					test.data["pepper"] = "spicy"
				})
				test.email = email()
				tests = append(tests, test)
			}
		}
		return tests
	}()
	err := a.conn.Transaction(func(tx *store.Connection) error {
		for i, test := range cases {
			u := user.NewUser(
				test.provider,
				test.role,
				test.email,
				test.name,
				[]byte(testPass),
				test.data,
				nil)
			err := tx.Save(u).Error
			require.NoError(t, err)
			cases[i].uid = u.ID
		}
		return nil
	})
	require.NoError(t, err)
	return cases
}

func TestSearchUsers(t *testing.T) {
	t.Parallel()
	a := apiWithTempDB(t)
	var tests []testUsers
	t.Run("CreateUsers", func(t *testing.T) {
		tests = testCreateUsers(t, a)
	})
	// find all no page
	ctx := context.Background()
	ctx.SetSort(store.Descending)
	list, err := a.SearchUsers(ctx, nil, nil)
	assert.NoError(t, err)
	assert.Len(t, list, len(tests))
	for _, idx := range []int{0, 5, 10, 20} {
		test := tests[idx]
		u := list[idx]
		assert.Equal(t, test.email, u.Email)
		assert.Equal(t, test.name, u.Username)
		assert.Equal(t, test.provider, u.Provider)
		assert.Equal(t, test.role, u.Role)
		assert.Equal(t, test.data, u.Data)
	}
	// find all page
	var size = len(tests) / 2
	page := &store.Pagination{
		Page: 1,
		Size: size,
	}
	list, err = a.SearchUsers(ctx, nil, page)
	assert.NoError(t, err)
	assert.Len(t, list, size)
	u := list[10]
	var found bool
	for _, test := range tests {
		found = test.email == u.Email
		if found {
			assert.Equal(t, test.name, u.Username)
			assert.Equal(t, test.provider, u.Provider)
			assert.Equal(t, test.role, u.Role)
			assert.Equal(t, test.data, u.Data)
			break
		}
	}
	assert.True(t, found)
	for _, filter := range []string{
		key.Email,
		key.Provider,
		key.Role,
		key.UserID,
		key.Username,
	} {
		name := strings.Title(string(filter) + "Filter")
		t.Run(name, func(t *testing.T) {
			filterTest(t, a, filter, tests)
		})
	}
	test := func(i int, paged bool) int {
		if !paged {
			return i
		}
		i = int(math.Min(float64(i), float64(pageSize)))
		return i
	}
	for _, paged := range []bool{false, true} {
		cnt := filtersTest(t, a, store.Filters{
			key.Provider: provider.Google,
			key.Role:     user.RoleAdmin,
		}, tests, paged)
		assert.Equal(t, test(8, paged), cnt)
		cnt = filtersTest(t, a, store.Filters{
			key.Provider: u.Provider,
			key.Role:     u.Role,
			key.UserID:   u.ID,
		}, tests, paged)
		assert.Equal(t, 1, cnt)
		cnt = filtersTest(t, a, store.Filters{
			key.Provider: provider.Apple,
			key.Role:     u.Role,
			key.UserID:   uuid.Nil,
		}, tests, paged)
		assert.Equal(t, 0, cnt)
		cnt = filtersTest(t, a, store.Filters{
			"book": "thing2",
		}, tests, paged)
		assert.Equal(t, test(20, paged), cnt)
		cnt = filtersTest(t, a, store.Filters{
			key.Provider: provider.Google,
			"book":       "thing2",
		}, tests, paged)
		assert.Equal(t, test(5, paged), cnt)
		cnt = filtersTest(t, a, store.Filters{
			key.Provider: u.Provider,
			"book":       u.Data["book"],
			"extra":      u.Data["extra"],
		}, tests, paged)
		assert.Equal(t, 1, cnt)
		cnt = filtersTest(t, a, store.Filters{
			"book":  "thing2",
			"extra": "caesar",
		}, tests, paged)
		assert.Equal(t, 0, cnt)
		var tc testUsers
		for _, tst := range tests {
			v, ok := tst.data["pepper"]
			if ok && v == "spicy" {
				tc = tst
				break
			}
		}
		cnt = filtersTest(t, a, store.Filters{
			"book":   tc.data["book"],
			"extra":  tc.data["extra"],
			"pepper": tc.data["pepper"],
		}, tests, paged)
		assert.Equal(t, 1, cnt)
	}
}

func filterTest(t *testing.T, a *API, f string, tests []testUsers) {
	filters := store.Filters{}
	tc := tests[20]
	switch f {
	case key.Email:
		filters[f] = tc.email
		cnt := filtersTest(t, a, filters, tests, false)
		assert.Equal(t, 1, cnt)
	case key.Provider:
		filters[f] = tc.provider
		cnt := filtersTest(t, a, filters, tests, false)
		assert.Equal(t, 80, cnt)
	case key.Role:
		filters[f] = tc.role
		cnt := filtersTest(t, a, filters, tests, false)
		assert.Equal(t, 112, cnt)
	case key.UserID:
		filters[f] = tc.uid
		cnt := filtersTest(t, a, filters, tests, false)
		assert.Equal(t, 1, cnt)
	case key.Username:
		filters[f] = tc.name
		cnt := filtersTest(t, a, filters, tests, false)
		assert.Equal(t, 16, cnt)
	}
}

func filtersTest(t *testing.T, a *API, f store.Filters, tests []testUsers, paged bool) int {
	var page *store.Pagination
	if paged {
		page = &store.Pagination{
			Page: 0,
			Size: pageSize,
		}
	}
	ctx := context.Background()
	ctx.SetSort(store.Descending)
	list, err := a.SearchUsers(ctx, f, page)
	assert.NoError(t, err)
	cnt := 0
	for _, test := range tests {
		if v, ok := f[key.Email]; ok {
			if test.email != v {
				continue
			}
		}
		if v, ok := f[key.Provider]; ok {
			if test.provider != v {
				continue
			}
		}
		if v, ok := f[key.Role]; ok {
			if test.role != v {
				continue
			}
		}
		if v, ok := f[key.UserID]; ok {
			if test.uid != v {
				continue
			}
		}
		if v, ok := f[key.Username]; ok {
			if test.name != v {
				continue
			}
		}
		if v, ok := f["book"]; ok {
			if test.data["book"] != v {
				continue
			}
		}
		if v, ok := f["extra"]; ok {
			if test.data["extra"] != v {
				continue
			}
		}
		if v, ok := f["pepper"]; ok {
			if test.data["pepper"] != v {
				continue
			}
		}
		cnt++
		if paged && cnt >= pageSize {
			break
		}
	}
	assert.Len(t, list, cnt)
	for _, l := range list {
		for k, v := range f {
			switch k {
			case key.Email:
				assert.Equal(t, v, l.Email)
			case key.Provider:
				assert.Equal(t, v, l.Provider)
			case key.Role:
				assert.Equal(t, v, l.Role)
			case key.UserID:
				assert.Equal(t, v, l.ID)
			case key.Username:
				assert.Equal(t, v, l.Username)
			case "book":
				assert.Equal(t, v, l.Data["book"])
			case "extra":
				assert.Equal(t, v, l.Data["extra"])
			case "pepper":
				assert.Equal(t, v, l.Data["pepper"])
			}
		}
	}
	return cnt
}

func TestAPI_ChangePassword(t *testing.T) {
	t.Parallel()
	var newPass = utils.SecureToken()
	a := apiWithTempDB(t)
	ctx := testContext(a)
	a.config.Validation.PasswordRegex = ""
	u := testUser(t, a)
	_, err := a.ChangePassword(nil, uuid.Nil, testPass, newPass)
	assert.Error(t, err)
	_, err = a.ChangePassword(ctx, u.ID, "", newPass)
	assert.Error(t, err)
	_, err = a.ChangePassword(ctx, u.ID, testPass, newPass)
	assert.NoError(t, err)
	banUser(t, a, u)
	_, err = a.ChangePassword(ctx, u.ID, newPass, testPass)
	assert.Error(t, err)
	a.config.Validation.PasswordRegex = "!"
	_, err = a.ChangePassword(ctx, u.ID, "", "")
	assert.Error(t, err)
}

func TestAPI_CmdChangeUserRole(t *testing.T) {
	t.Parallel()
	a := apiWithTempDB(t)
	ctx := testContext(a)
	u := testUser(t, a)
	_, err := a.ChangeRole(nil, uuid.Nil, user.RoleAdmin)
	assert.Error(t, err)
	_, err = a.ChangeRole(ctx, uuid.New(), user.RoleAdmin)
	assert.Error(t, err)
	// promote
	u, err = a.ChangeRole(ctx, u.ID, user.RoleAdmin)
	assert.NoError(t, err)
	assert.True(t, u.IsAdmin())
	// re-promote
	u, err = a.ChangeRole(ctx, u.ID, user.RoleAdmin)
	assert.NoError(t, err)
	assert.True(t, u.IsAdmin())
	// super-promote
	_, err = a.ChangeRole(ctx, u.ID, user.RoleSuper)
	assert.Error(t, err)
	// sneak promote
	u.Role = user.RoleSuper
	err = a.conn.Save(u).Error
	assert.NoError(t, err)
	u, err = a.ChangeRole(ctx, u.ID, user.RoleAdmin)
	assert.NoError(t, err)
	assert.True(t, u.IsAdmin())
	// demote
	u, err = a.ChangeRole(ctx, u.ID, user.RoleUser)
	assert.NoError(t, err)
	assert.False(t, u.IsAdmin())
	banUser(t, a, u)
	_, err = a.ChangeRole(ctx, u.ID, user.RoleAdmin)
	assert.Error(t, err)
	_, err = a.ChangeRole(ctx, u.ID, user.RoleUser)
	assert.Error(t, err)
}

func TestAPI_ConfirmUser(t *testing.T) {
	t.Parallel()
	a := apiWithTempDB(t)
	u := testUser(t, a)
	ctx := testContext(a)
	assert.False(t, u.IsConfirmed())
	ct, err := tokens.GrantConfirmToken(a.conn, u.ID, token.NoExpiration)
	assert.NoError(t, err)
	// bad token
	u, err = a.ConfirmUser(nil, "")
	assert.Error(t, err)
	// good token
	u, err = a.ConfirmUser(ctx, ct.String())
	assert.NoError(t, err)
	assert.True(t, u.IsConfirmed())
	ct, err = tokens.GrantConfirmToken(a.conn, u.ID, token.NoExpiration)
	assert.NoError(t, err)
	_, err = a.ConfirmUser(ctx, ct.String())
	assert.NoError(t, err)

}

func TestAPI_ConfirmPassword(t *testing.T) {
	t.Parallel()
	const (
		empty   = ""
		badPass = "pass"
		passRx  = "^[a-zA-Z0-9[:punct:]]{8,40}$"
	)
	a := apiWithTempDB(t)
	u := testUser(t, a)
	assert.False(t, u.IsConfirmed())
	ct, err := tokens.GrantConfirmToken(a.conn, u.ID, token.NoExpiration)
	assert.NoError(t, err)
	// password validation
	a.config.Validation.PasswordRegex = passRx
	ctx := testContext(a)
	// bad token
	u, err = a.ConfirmResetPassword(nil, "", testPass)
	assert.Error(t, err)
	_, err = a.ConfirmResetPassword(nil, ct.String(), empty)
	assert.Error(t, err)
	_, err = a.ConfirmResetPassword(ctx, ct.String(), badPass)
	assert.Error(t, err)
	u, err = a.ConfirmResetPassword(ctx, ct.String(), testPass)
	assert.NoError(t, err)
	assert.True(t, u.IsConfirmed())
	err = u.Authenticate(testPass)
	assert.NoError(t, err)
	// no password validation
	a.config.Validation.PasswordRegex = empty
	ct, err = tokens.GrantConfirmToken(a.conn, u.ID, token.NoExpiration)
	assert.NoError(t, err)
	u, err = a.ConfirmResetPassword(ctx, ct.String(), empty)
	assert.NoError(t, err)
	err = u.Authenticate(empty)
	assert.NoError(t, err)
	// can't reuse token
	_, err = a.ConfirmResetPassword(ctx, ct.String(), empty)
	assert.Error(t, err)
}

func TestAPI_ConfirmEmail(t *testing.T) {
	t.Parallel()
	const (
		empty    = ""
		badEmail = "@"
	)
	var testEmail = tutils.RandomEmail()
	a := apiWithTempDB(t)
	ctx := testContext(a)
	u := testUser(t, a)
	assert.False(t, u.IsConfirmed())
	ct, err := tokens.GrantConfirmToken(a.conn, u.ID, token.NoExpiration)
	assert.NoError(t, err)
	// bad token
	_, err = a.ConfirmChangeEmail(nil, "", testEmail)
	assert.Error(t, err)
	_, err = a.ConfirmChangeEmail(ctx, ct.String(), empty)
	assert.Error(t, err)
	_, err = a.ConfirmChangeEmail(ctx, ct.String(), badEmail)
	assert.Error(t, err)
	u, err = a.ConfirmChangeEmail(ctx, ct.String(), testEmail)
	assert.NoError(t, err)
	assert.True(t, u.IsConfirmed())
	// can't reuse token
	_, err = a.ConfirmChangeEmail(ctx, ct.String(), testEmail)
	assert.Error(t, err)
}

func TestAPI_UpdateUser(t *testing.T) {
	t.Parallel()
	var testName = "peaches"
	data := types.Map{
		"foo":   "bar",
		"tasty": "salad",
	}
	a := apiWithTempDB(t)
	u := testUser(t, a)
	confirmUser(t, a, u)
	ctx := testContext(a)
	// system user
	_, err := a.UpdateUser(nil, uuid.Nil, nil, nil)
	assert.Error(t, err)
	// user not found
	_, err = a.UpdateUser(ctx, uuid.New(), nil, nil)
	assert.Error(t, err)
	a.config.Validation.UsernameRegex = "0"
	_, err = a.UpdateUser(ctx, u.ID, &testName, nil)
	assert.Error(t, err)
	a.config.Validation.UsernameRegex = ""
	u, err = a.UpdateUser(ctx, u.ID, &testName, nil)
	assert.NoError(t, err)
	require.NotNil(t, u)
	assert.Equal(t, testName, u.Username)
	var fooName = "foo"
	u, err = a.UpdateUser(ctx, u.ID, &fooName, data)
	assert.NoError(t, err)
	assert.Equal(t, fooName, u.Username)
	assert.EqualValues(t, data, u.Data)
}

func TestAPI_BanUser(t *testing.T) {
	t.Parallel()
	a := apiWithTempDB(t)
	u := testUser(t, a)
	assert.False(t, u.IsBanned())
	// no user id
	_, err := a.BanUser(nil, uuid.Nil)
	assert.Error(t, err)
	// bad user id
	_, err = a.BanUser(nil, uuid.New())
	assert.Error(t, err)
	// ban user
	u, err = a.BanUser(nil, u.ID)
	assert.NoError(t, err)
	require.NotNil(t, u)
	assert.True(t, u.IsBanned())
	// "re" ban user
	u, err = a.BanUser(nil, u.ID)
	assert.NoError(t, err)
	require.NotNil(t, u)
	assert.True(t, u.IsBanned())
}

func TestAPI_DeleteUser(t *testing.T) {
	t.Parallel()
	a := apiWithTempDB(t)
	u := testUser(t, a)
	assert.True(t, u.Valid())
	assert.False(t, u.DeletedAt.Valid)
	// no user id
	err := a.DeleteUser(nil, uuid.Nil)
	assert.Error(t, err)
	// bad user id
	err = a.DeleteUser(nil, uuid.New())
	assert.Error(t, err)
	// delete user
	err = a.DeleteUser(nil, u.ID)
	assert.NoError(t, err)
	_, err = a.GetUser(u.ID)
	assert.Error(t, err)
}

func TestAPI_LinkAccount(t *testing.T) {
	t.Parallel()
	var (
		p     = provider.Google
		aid   = uuid.New().String()
		email = tutils.RandomEmail()
		data  = types.Map{
			"hello":       "world",
			key.IPAddress: "127.0.0.1",
		}
	)
	a := apiWithTempDB(t)
	u := testUser(t, a)
	ctx := context.Background()
	// errors
	err := a.LinkAccount(ctx, uuid.Nil, nil)
	assert.Error(t, err)
	bad := account.NewAccount(p, aid, email, data)
	err = a.LinkAccount(ctx, uuid.Nil, bad)
	assert.Error(t, err)
	bad = account.NewAccount(provider.Unknown, aid, email, data)
	err = a.LinkAccount(ctx, u.ID, bad)
	assert.Error(t, err)
	bad = account.NewAccount(p, "", email, data)
	err = a.LinkAccount(ctx, u.ID, bad)
	assert.Error(t, err)
	// success
	la1 := account.NewAccount(p, aid, email, data)
	err = a.LinkAccount(ctx, u.ID, la1)
	assert.NoError(t, err)
	require.NotNil(t, la1)
	assert.Equal(t, aid, la1.AccountID)
	assert.Equal(t, email, la1.Email)
	assert.Equal(t, data, la1.Data)
	assert.Equal(t, u.ID, la1.UserID)
	cnt := a.conn.Model(u).Association("Linked").Count()
	assert.Equal(t, int64(1), cnt)
	// relink error
	la1 = account.NewAccount(p, aid, email, data)
	err = a.LinkAccount(ctx, u.ID, la1)
	assert.Error(t, err)
	// 2nd account
	la2 := account.NewAccount(provider.GitHub, aid, email, data)
	err = a.LinkAccount(ctx, u.ID, la2)
	assert.NoError(t, err)
	require.NotNil(t, la2)
	assert.NotEqual(t, la1.ID, la2.ID)
	cnt = a.conn.Model(u).Association("Linked").Count()
	assert.Equal(t, int64(2), cnt)
}

func TestAPI_GetLinkedAccounts(t *testing.T) {
	t.Parallel()
	testAccount := func(name provider.Name) *account.Account {
		link := account.NewAccount(name,
			uuid.New().String(),
			tutils.RandomEmail(),
			types.Map{
				"hello":   "world",
				key.Token: name.ID().String(),
			})
		return link
	}
	links := []*account.Account{
		testAccount(provider.Google),
		testAccount(provider.GitHub),
		testAccount(provider.Stripe),
		testAccount(provider.PayPal),
	}
	a := apiWithTempDB(t)
	u := testUser(t, a)
	ctx := context.Background()
	for _, link := range links {
		err := a.LinkAccount(ctx, u.ID, link)
		require.NoError(t, err)
	}
	cnt := a.conn.Model(u).Association("Linked").Count()
	assert.Equal(t, int64(len(links)), cnt)
	tests := []struct {
		typ      account.Type
		filters  store.Filters
		expected []int
	}{
		{account.Any, nil, []int{0, 1, 2, 3}},
		{account.Payment, nil, []int{2, 3}},
		{account.Wallet, nil, []int{3}},
		{account.Payment | account.Wallet, nil, []int{2, 3}},
		{account.Any, store.Filters{
			key.Provider: links[0].Provider,
		}, []int{0}},
		{account.Any, store.Filters{
			key.Email: links[2].Email,
		}, []int{2}},
		{account.Any, store.Filters{
			key.Token: links[1].Provider.ID().String(),
		}, []int{1}},
		{account.Any, store.Filters{
			"hello": "world",
		}, []int{0, 1, 2, 3}},
		{account.Payment, store.Filters{
			"hello": "world",
		}, []int{2, 3}},
		{account.Any, store.Filters{
			"bad": "unknown",
		}, []int{}},
	}
	for _, test := range tests {
		res, err := a.GetLinkedAccounts(ctx, u.ID, test.typ, test.filters)
		require.NoError(t, err)
		require.Len(t, res, len(test.expected))
		for i, r := range res {
			idx := test.expected[i]
			link := links[idx]
			assert.Equal(t, link.ID, r.ID)
			assert.True(t, r.Type.Has(test.typ))
			assert.Equal(t, link.Email, r.Email)
			assert.Equal(t, link.Provider, r.Provider)
			assert.EqualValues(t, link.Data, r.Data)
		}
	}
	// bad user id
	_, err := a.GetLinkedAccounts(ctx, uuid.Nil, account.Any, nil)
	assert.Error(t, err)
	// user not found
	_, err = a.GetLinkedAccounts(ctx, uuid.New(), account.Any, nil)
	assert.Error(t, err)
}
