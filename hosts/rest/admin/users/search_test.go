package users

import (
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/models/types/provider"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/test/tconn"
	"github.com/jrapoport/gothic/test/thttp"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/jrapoport/gothic/utils"
	"github.com/segmentio/encoding/json"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestUsers struct {
	uid      uuid.UUID
	provider provider.Name
	role     user.Role
	email    string
	name     string
	data     types.Map
}

var testBook = uuid.New().String()

func CreateTestUsers(t *testing.T, srv *usersServer) []TestUsers {
	var cases = func() []TestUsers {
		email := func() string { return tutils.RandomEmail() }
		name := func() string { return utils.RandomUsername() }
		p := srv.Provider()
		ru := user.RoleUser
		ra := user.RoleAdmin
		var tests = []TestUsers{
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
		var idx = 0
		for i, bk := range []interface{}{
			"thing2", testBook, uuid.New().String(),
		} {
			for x, test := range tests {
				test.data = types.Map{
					"dr_suess": "thing1",
					"book":     bk,
				}
				if idx < 50 {
					test.data["sorted"] = "yes"
				}
				idx++
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
	ctx := context.Background()
	admin, _ := testUser(t, srv, true)
	ctx.SetAdminID(admin.ID)
	ctx.SetProvider(srv.Provider())
	srv.Config().Signup.AutoConfirm = true
	srv.Config().Signup.Default.Username = false
	srv.Config().Signup.Default.Color = false
	srv.Config().Mail.SpamProtection = false
	conn := tconn.Conn(t, srv.Config())
	for i, test := range cases {
		u, err := srv.API.AdminCreateUser(ctx, test.email, test.name, testPass, test.data, false)
		require.NoError(t, err)
		u, err = srv.ChangeRole(ctx, u.ID, test.role)
		require.NoError(t, err)
		u.Provider = test.provider
		err = conn.Save(u).Error
		require.NoError(t, err)
		cases[i].uid = u.ID
	}
	return cases
}

type UserServerTestSuite struct {
	suite.Suite
	srv   *usersServer
	conn  *store.Connection
	tests []TestUsers
	uid   uuid.UUID
}

func TestUserServer_Search(t *testing.T) {
	t.Parallel()
	ts := &UserServerTestSuite{}
	suite.Run(t, ts)
}

func (ts *UserServerTestSuite) SetupSuite() {
	s, _ := tsrv.RESTServer(ts.T(), false)
	ts.srv = newUserServer(s)
	conn, err := store.Dial(ts.srv.Config(), nil)
	ts.Require().NoError(err)
	ts.conn = conn
	ts.tests = CreateTestUsers(ts.T(), ts.srv)
}

func (ts *UserServerTestSuite) searchUsers(ep string, v url.Values) *httptest.ResponseRecorder {
	r := thttp.Request(ts.T(), http.MethodGet, ep, "", v, nil)
	w := httptest.NewRecorder()
	ts.srv.SearchUsers(w, r)
	return w
}

func (ts *UserServerTestSuite) TestErrors() {
	// invalid req
	r := thttp.Request(ts.T(), http.MethodGet, Endpoint, "", nil, []byte("\n"))
	w := httptest.NewRecorder()
	ts.srv.SearchUsers(w, r)
	ts.NotEqual(http.StatusOK, w.Code)
	// bad paging
	r = thttp.Request(ts.T(), http.MethodGet, Endpoint, "", url.Values{
		key.Page: []string{"\n"},
	}, nil)
	w = httptest.NewRecorder()
	ts.srv.SearchUsers(w, r)
	ts.NotEqual(http.StatusOK, w.Code)
}

func (ts *UserServerTestSuite) TestPageHeaders() {
	res := ts.searchUsers(Endpoint, nil)
	ts.Equal(http.StatusOK, res.Code)
	var list []interface{}
	err := json.NewDecoder(res.Body).Decode(&list)
	ts.NoError(err)
	ts.Len(list, store.MaxPerPage)
	e := list[0].(map[string]interface{})
	f := e[key.Data].(map[string]interface{})
	uid := uuid.MustParse(e[key.ID].(string))
	u, err := ts.srv.API.GetUser(uid)
	ts.Require().NoError(err)
	ts.Equal(uid, u.ID)
	ts.Equal(f["dr_suess"], u.Data["dr_suess"])
	pn := res.Header().Get(rest.PageNumber)
	ts.Equal("1", pn)
	pc := res.Header().Get(rest.PageCount)
	cnt := int(math.Ceil(float64(len(ts.tests)) / float64(store.MaxPerPage)))
	testCount := strconv.Itoa(cnt)
	ts.Equal(testCount, pc)
	pl := res.Header().Get(rest.PageLength)
	testLen := strconv.Itoa(store.MaxPerPage)
	ts.Equal(testLen, pl)
	tot := res.Header().Get(rest.PageTotal)
	// +1 because of audit.LogStartup
	testTotal := strconv.Itoa(len(ts.tests) + 1)
	ts.Equal(testTotal, tot)
}

func (ts *UserServerTestSuite) TestPageLinks() {
	startLink := func() string {
		return fmt.Sprintf("%s?%s=1&%s=%d",
			Endpoint, key.Page, key.PerPage, store.MaxPerPage)
	}
	var nextLink = startLink()
	for {
		if nextLink == "" {
			break
		}
		u, err := url.Parse(nextLink)
		ts.Require().NoError(err)
		nextLink = ""
		u.Scheme = ""
		u.Host = ""
		uri := u.String()
		res := ts.searchUsers(uri, nil)
		ts.Equal(http.StatusOK, res.Code)
		var logs []interface{}
		err = json.NewDecoder(res.Body).Decode(&logs)
		ts.Require().NoError(err)
		pc := res.Header().Get(rest.PageLength)
		cnt, err := strconv.Atoi(pc)
		ts.Require().NoError(err)
		ts.Len(logs, cnt)
		l := res.Header().Get(rest.Link)
		links := strings.Split(l, ",")
		if len(links) <= 0 {
			break
		}
		for _, lnk := range links {
			next := `rel="next"`
			if strings.HasSuffix(lnk, next) {
				nextLink = strings.ReplaceAll(lnk, next, "")
				nextLink = strings.Trim(nextLink, " <>;")
				break
			}
		}
	}
}

func (ts *UserServerTestSuite) TestSearchFilters() {
	tests := []struct {
		v    url.Values
		comp func(e map[string]interface{})
	}{
		{
			url.Values{
				key.UserID: []string{ts.tests[0].uid.String()},
			},
			func(e map[string]interface{}) {
				uid := e[key.ID].(string)
				ts.Equal(ts.tests[0].uid.String(), uid)
			},
		},
		{
			url.Values{
				key.ID: []string{ts.tests[0].uid.String()},
			},
			func(e map[string]interface{}) {
				uid := e[key.ID].(string)
				ts.Equal(ts.tests[0].uid.String(), uid)
			},
		},
		{
			url.Values{
				key.Username: []string{ts.tests[0].name},
			},
			func(e map[string]interface{}) {
				name := e[key.Username].(string)
				ts.Equal(ts.tests[0].name, name)
			},
		},
		{
			url.Values{
				key.Email: []string{ts.tests[0].email},
			},
			func(e map[string]interface{}) {
				em := e[key.Email].(string)
				ts.Equal(ts.tests[0].email, em)
			},
		},
		{
			url.Values{
				key.Provider: []string{provider.Google.String()},
			},
			func(e map[string]interface{}) {
				p := e[key.Provider].(string)
				ts.Equal(provider.Google.String(), p)
			},
		},
		{
			url.Values{
				key.Role: []string{user.RoleUser.String()},
			},
			func(e map[string]interface{}) {
				r := e[key.Role].(string)
				ts.Equal(user.RoleUser.String(), r)
			},
		},
		{
			url.Values{
				"dr_suess": []string{"thing1"},
			},
			func(e map[string]interface{}) {
				f := e[key.Data].(map[string]interface{})
				ts.Equal("thing1", f["dr_suess"])
			},
		},
		{
			url.Values{
				key.Role:   []string{user.RoleUser.String()},
				"dr_suess": []string{"thing1"},
			},
			func(e map[string]interface{}) {
				r := e[key.Role].(string)
				ts.Equal(user.RoleUser.String(), r)
				f := e[key.Data].(map[string]interface{})
				ts.Equal("thing1", f["dr_suess"])
			},
		},
		{
			url.Values{
				"dr_suess": []string{"thing1"},
				"book":     []string{testBook},
			},
			func(e map[string]interface{}) {
				f := e[key.Data].(map[string]interface{})
				ts.Equal("thing1", f["dr_suess"])
				ts.Equal(testBook, f["book"])
			},
		},
	}
	for _, test := range tests {
		res := ts.searchUsers(Endpoint, test.v)
		ts.Equal(http.StatusOK, res.Code)
		var list []interface{}
		err := json.NewDecoder(res.Body).Decode(&list)
		ts.NoError(err)
		ts.Greater(len(list), 0)
		for _, log := range list {
			e := log.(map[string]interface{})
			test.comp(e)
		}
	}
}

func (ts *UserServerTestSuite) TestSearchSort() {
	// search Ascending
	v := url.Values{
		key.Sort:   []string{string(store.Ascending)},
		"dr_suess": []string{"thing1"},
		"sorted":   []string{"yes"},
	}
	var list []interface{}
	res := ts.searchUsers(Endpoint, v)
	ts.Equal(http.StatusOK, res.Code)
	err := json.NewDecoder(res.Body).Decode(&list)
	ts.NoError(err)
	ts.Greater(len(list), 0)
	testIDs := make([]string, len(list))
	for i, log := range list {
		e := log.(map[string]interface{})
		f := e[key.Data].(map[string]interface{})
		ts.Equal("thing1", f["dr_suess"])
		ts.Equal("yes", f["sorted"])
		testIDs[i] = e[key.ID].(string)
	}
	// search Descending
	v[key.Sort] = []string{string(store.Descending)}
	res = ts.searchUsers(Endpoint, v)
	ts.Equal(http.StatusOK, res.Code)
	err = json.NewDecoder(res.Body).Decode(&list)
	ts.NoError(err)
	ts.Greater(len(list), 0)
	ts.Require().Len(list, len(testIDs))
	descIDs := make([]string, len(list))
	for i, log := range list {
		e := log.(map[string]interface{})
		f := e[key.Data].(map[string]interface{})
		ts.Equal("thing1", f["dr_suess"])
		ts.Equal("yes", f["sorted"])
		descIDs[i] = e[key.ID].(string)
	}
	// reverse the ids
	reverse := func(ids []string) {
		for i, j := 0, len(ids)-1; i < j; i, j = i+1, j-1 {
			ids[i], ids[j] = ids[j], ids[i]
		}
	}
	reverse(descIDs)
	ts.Equal(testIDs, descIDs)
}
