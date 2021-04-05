package users

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/models/types/provider"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/test/tconn"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/jrapoport/gothic/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const pageSize = 4

type testCase struct {
	uid      uuid.UUID
	provider provider.Name
	role     user.Role
	email    string
	name     string
	data     types.Map
}

var testBook = utils.RandomUsername()

func testCreateUsers(t *testing.T, conn *store.Connection, c *config.Config) []testCase {
	p := c.Provider()
	var cases = func() []testCase {
		email := func() string { return tutils.RandomEmail() }
		name := func() string { return utils.RandomUsername() }
		ru := user.RoleUser
		ra := user.RoleAdmin
		var tests = []testCase{
			{uuid.Nil, p, ru, email(), name(), nil},
			{uuid.Nil, p, ru, email(), name(), nil},
			{uuid.Nil, p, ru, email(), name(), nil},
			{uuid.Nil, p, ru, email(), name(), nil},
			{uuid.Nil, p, ru, email(), name(), nil},
			{uuid.Nil, p, ru, email(), name(), nil},
			{uuid.Nil, p, ru, email(), name(), nil},
			{uuid.Nil, p, ra, email(), name(), nil},
			{uuid.Nil, p, ra, email(), name(), nil},
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
		var once sync.Once
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
				once.Do(func() {
					test.data["pepper"] = "spicy"
				})
				test.email = email()
				tests = append(tests, test)
			}
		}
		return tests
	}()
	err := conn.Transaction(func(tx *store.Connection) error {
		for i, test := range cases {
			u, err := createUser(tx, test.provider, test.email, test.name, "", test.data, nil)
			require.NoError(t, err)
			u.Role = test.role
			err = tx.Save(u).Error
			require.NoError(t, err)
			cases[i].uid = u.ID
		}
		return nil
	})
	require.NoError(t, err)
	su := user.NewSuperAdmin("")
	su.Provider = p
	err = conn.Create(su).Error
	require.NoError(t, err)
	return cases
}

func TestSearchUsers(t *testing.T) {
	t.Parallel()
	conn, c := tconn.TempConn(t)
	var tests []testCase
	t.Run("Populate Users", func(t *testing.T) {
		tests = testCreateUsers(t, conn, c)
	})
	// find all no page
	list, err := SearchUsers(conn, store.Descending, nil, nil)
	assert.NoError(t, err)
	require.Len(t, list, len(tests))
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
		Index: 1,
		Size:  size,
	}
	list, err = SearchUsers(conn, store.Descending, nil, page)
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
		name := strings.Title(filter + "Filter")
		t.Run(name, func(t *testing.T) {
			filterTest(t, conn, filter, tests)
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
		cnt := filtersTest(t, conn, store.Filters{
			key.Provider: provider.Google,
			key.Role:     user.RoleAdmin,
		}, tests, paged)
		assert.Equal(t, test(8, paged), cnt)
		cnt = filtersTest(t, conn, store.Filters{
			key.Provider: u.Provider,
			key.Role:     u.Role,
			key.UserID:   u.ID,
		}, tests, paged)
		assert.Equal(t, 1, cnt)
		cnt = filtersTest(t, conn, store.Filters{
			key.Provider: provider.Apple,
			key.Role:     u.Role,
			key.UserID:   uuid.Nil,
		}, tests, paged)
		assert.Equal(t, 0, cnt)
		cnt = filtersTest(t, conn, store.Filters{
			"book": "thing2",
		}, tests, paged)
		assert.Equal(t, test(18, paged), cnt)
		cnt = filtersTest(t, conn, store.Filters{
			key.Provider: provider.Google,
			"book":       "thing2",
		}, tests, paged)
		assert.Equal(t, test(4, paged), cnt)
		cnt = filtersTest(t, conn, store.Filters{
			key.Provider: u.Provider,
			"book":       u.Data["book"],
			"extra":      u.Data["extra"],
		}, tests, paged)
		assert.Equal(t, 1, cnt)
		cnt = filtersTest(t, conn, store.Filters{
			"book":  "thing2",
			"extra": "caesar",
		}, tests, paged)
		assert.Equal(t, 0, cnt)
		var tc testCase
		for _, tst := range tests {
			v, ok := tst.data["pepper"]
			if ok && v == "spicy" {
				tc = tst
				break
			}
		}
		cnt = filtersTest(t, conn, store.Filters{
			"book":   tc.data["book"],
			"extra":  tc.data["extra"],
			"pepper": tc.data["pepper"],
		}, tests, paged)
		assert.Equal(t, 1, cnt)
	}
}

func filterTest(t *testing.T, conn *store.Connection, f string, tests []testCase) {
	filters := store.Filters{}
	tc := tests[20]
	switch f {
	case key.Email:
		filters[f] = tc.email
		cnt := filtersTest(t, conn, filters, tests, false)
		assert.Equal(t, 1, cnt)
	case key.Provider:
		filters[f] = tc.provider
		cnt := filtersTest(t, conn, filters, tests, false)
		assert.Equal(t, 72, cnt)
	case key.Role:
		filters[f] = tc.role
		cnt := filtersTest(t, conn, filters, tests, false)
		assert.Equal(t, 112, cnt)
	case key.UserID:
		filters[f] = tc.uid
		cnt := filtersTest(t, conn, filters, tests, false)
		assert.Equal(t, 1, cnt)
	case key.Username:
		filters[f] = tc.name
		cnt := filtersTest(t, conn, filters, tests, false)
		assert.Equal(t, 16, cnt)
	}
}

func filtersTest(t *testing.T, conn *store.Connection, f store.Filters, tests []testCase, paged bool) int {
	var page *store.Pagination
	if paged {
		page = &store.Pagination{
			Index: 0,
			Size:  pageSize,
		}
	}
	list, err := SearchUsers(conn, store.Descending, f, page)
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
