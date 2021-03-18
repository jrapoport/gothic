package audit

import (
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/core/audit"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/models/auditlog"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/store/types"
	"github.com/jrapoport/gothic/store/types/key"
	"github.com/jrapoport/gothic/test/thttp"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/segmentio/encoding/json"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestCase struct {
	act    auditlog.Action
	uid    uuid.UUID
	fields types.Map
	logID  uint
}

var (
	testUID  = uuid.New()
	testBook = uuid.New()
)

func SetupTestLogs(t *testing.T, c *config.Config, uid, bid uuid.UUID) []TestCase {
	var tests = []TestCase{
		{auditlog.Startup, uuid.New(), nil, 0},
		{auditlog.Shutdown, uuid.New(), nil, 0},
		{auditlog.Signup, uuid.New(), nil, 0},
		{auditlog.ConfirmSent, uuid.New(), nil, 0},
		{auditlog.Confirmed, uuid.New(), nil, 0},
		{auditlog.Granted, uuid.New(), nil, 0},
		{auditlog.Revoked, uuid.New(), nil, 0},
		{auditlog.RevokedAll, uuid.New(), nil, 0},
		{auditlog.Login, uuid.New(), nil, 0},
		{auditlog.Logout, uuid.New(), nil, 0},
		{auditlog.Password, uuid.New(), nil, 0},
		{auditlog.Email, uuid.New(), nil, 0},
		{auditlog.Updated, uuid.New(), nil, 0},
	}
	for _, test := range tests {
		tst := test
		tst.uid = uid
		tests = append(tests, tst)
	}
	var idx int
	for _, bk := range []interface{}{
		"thing2", bid.String(), uuid.New().String(),
	} {
		for _, test := range tests {
			test.fields = types.Map{
				"dr_suess": "thing1",
				"book":     bk,
			}
			if idx < 50 {
				test.fields["sorted"] = "yes"
			}
			idx++
			tests = append(tests, test)
		}
	}
	ctx := context.Background()
	ctx.SetProvider(c.Provider())
	ctx.SetIPAddress("127.0.0.1")
	conn, err := store.Dial(c, nil)
	require.NoError(t, err)
	err = conn.Transaction(func(tx *store.Connection) error {
		for i, test := range tests {
			le, err := audit.CreateLogEntry(ctx, tx, test.act, test.uid, test.fields)
			require.NoError(t, err)
			tests[i].logID = le.ID
		}
		return nil
	})
	require.NoError(t, err)
	return tests
}

type AuditServerTestSuite struct {
	suite.Suite
	srv   *auditServer
	conn  *store.Connection
	tests []TestCase
	uid   uuid.UUID
}

func TestAuditLogs(t *testing.T) {
	t.Parallel()
	ts := &AuditServerTestSuite{}
	suite.Run(t, ts)
}

func (ts *AuditServerTestSuite) SetupSuite() {
	s, _ := tsrv.RESTServer(ts.T(), false)
	ts.srv = newAuditServer(s)
	conn, err := store.Dial(ts.srv.Config(), nil)
	ts.Require().NoError(err)
	ts.conn = conn
	ts.tests = SetupTestLogs(ts.T(), ts.srv.Config(),
		testUID, testBook)
}

func (ts *AuditServerTestSuite) searchAuditLogs(ep string, v url.Values) *httptest.ResponseRecorder {
	r := thttp.Request(ts.T(), http.MethodGet, ep, "", v, nil)
	w := httptest.NewRecorder()
	ts.srv.SearchAuditLogs(w, r)
	return w
}

func (ts *AuditServerTestSuite) TestErrors() {
	// invalid req
	r := thttp.Request(ts.T(), http.MethodGet, Endpoint, "", nil, []byte("\n"))
	w := httptest.NewRecorder()
	ts.srv.SearchAuditLogs(w, r)
	ts.NotEqual(http.StatusOK, w.Code)
	// bad paging
	r = thttp.Request(ts.T(), http.MethodGet, Endpoint, "", url.Values{
		key.Page: []string{"\n"},
	}, nil)
	w = httptest.NewRecorder()
	ts.srv.SearchAuditLogs(w, r)
	ts.NotEqual(http.StatusOK, w.Code)
}

func (ts *AuditServerTestSuite) TestPageHeaders() {
	res := ts.searchAuditLogs(Endpoint, nil)
	ts.Equal(http.StatusOK, res.Code)
	var logs []interface{}
	err := json.NewDecoder(res.Body).Decode(&logs)
	ts.NoError(err)
	ts.Len(logs, store.MaxPerPage)
	e := logs[0].(map[string]interface{})
	f := e[key.Fields].(map[string]interface{})
	id := uint(e["ID"].(float64))
	le, err := audit.GetLogEntry(ts.conn, id)
	ts.Require().NoError(err)
	ts.Equal(id, le.ID)
	ts.Equal(f["dr_suess"], le.Fields["dr_suess"])
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

func (ts *AuditServerTestSuite) TestPageLinks() {
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
		res := ts.searchAuditLogs(uri, nil)
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

func (ts *AuditServerTestSuite) TestSearchFilters() {
	tests := []struct {
		v    url.Values
		comp func(e map[string]interface{})
	}{
		{
			url.Values{
				key.UserID: []string{testUID.String()},
			},
			func(e map[string]interface{}) {
				uid := e[key.UserID].(string)
				ts.Equal(testUID.String(), uid)
			},
		},
		{
			url.Values{
				key.Action: []string{auditlog.Startup.String()},
			},
			func(e map[string]interface{}) {
				act := e[key.Action].(string)
				ts.Equal(auditlog.Startup.String(), act)
			},
		},
		{
			url.Values{
				"dr_suess": []string{"thing1"},
			},
			func(e map[string]interface{}) {
				f := e[key.Fields].(map[string]interface{})
				ts.Equal("thing1", f["dr_suess"])
			},
		},
		{
			url.Values{
				key.Type:   []string{auditlog.Account.String()},
				"dr_suess": []string{"thing1"},
			},
			func(e map[string]interface{}) {
				typ := auditlog.Type(e[key.Type].(float64))
				ts.Equal(auditlog.Account, typ)
				f := e[key.Fields].(map[string]interface{})
				ts.Equal("thing1", f["dr_suess"])
			},
		},
		{
			url.Values{
				"dr_suess": []string{"thing1"},
				"book":     []string{testBook.String()},
			},
			func(e map[string]interface{}) {
				f := e[key.Fields].(map[string]interface{})
				ts.Equal("thing1", f["dr_suess"])
				ts.Equal(testBook.String(), f["book"])
			},
		},
	}
	for _, test := range tests {
		res := ts.searchAuditLogs(Endpoint, test.v)
		ts.Equal(http.StatusOK, res.Code)
		var logs []interface{}
		err := json.NewDecoder(res.Body).Decode(&logs)
		ts.NoError(err)
		ts.Greater(len(logs), 0)
		for _, log := range logs {
			e := log.(map[string]interface{})
			test.comp(e)
		}
	}
}

func (ts *AuditServerTestSuite) TestSearchSort() {
	// search Ascending
	v := url.Values{
		key.Sort:   []string{string(store.Ascending)},
		"dr_suess": []string{"thing1"},
		"sorted":   []string{"yes"},
	}
	var logs []interface{}
	res := ts.searchAuditLogs(Endpoint, v)
	ts.Equal(http.StatusOK, res.Code)
	err := json.NewDecoder(res.Body).Decode(&logs)
	ts.NoError(err)
	ts.Greater(len(logs), 0)
	testIdx := make([]int, len(logs))
	for i, log := range logs {
		e := log.(map[string]interface{})
		f := e[key.Fields].(map[string]interface{})
		ts.Equal("thing1", f["dr_suess"])
		ts.Equal("yes", f["sorted"])
		testIdx[i] = int(e["ID"].(float64))
	}
	// search Descending
	v[key.Sort] = []string{string(store.Descending)}
	res = ts.searchAuditLogs(Endpoint, v)
	ts.Equal(http.StatusOK, res.Code)
	err = json.NewDecoder(res.Body).Decode(&logs)
	ts.NoError(err)
	ts.Greater(len(logs), 0)
	ts.Require().Len(logs, len(testIdx))
	descIdx := make([]int, len(logs))
	for i, log := range logs {
		e := log.(map[string]interface{})
		f := e[key.Fields].(map[string]interface{})
		ts.Equal("thing1", f["dr_suess"])
		ts.Equal("yes", f["sorted"])
		descIdx[i] = int(e["ID"].(float64))
	}
	// reverse the indexes
	sort.Ints(descIdx)
	ts.Equal(testIdx, descIdx)
}
