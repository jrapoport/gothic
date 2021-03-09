package auditlog_test

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
	"github.com/jrapoport/gothic/core/audit"
	"github.com/jrapoport/gothic/hosts/rest"
	audit_log "github.com/jrapoport/gothic/hosts/rest/admin/auditlog"
	"github.com/jrapoport/gothic/models/auditlog"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/store/types/key"
	"github.com/jrapoport/gothic/test/thttp"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/segmentio/encoding/json"
	"github.com/stretchr/testify/suite"
)

const endpoint = audit_log.Endpoint

var (
	testUID  = uuid.New()
	testBook = uuid.New()
)

type AuditTestSuite struct {
	suite.Suite
	srv   *rest.Host
	web   *httptest.Server
	conn  *store.Connection
	tests []audit_log.TestCase
	uid   uuid.UUID
}

func TestAudit(t *testing.T) {
	ts := &AuditTestSuite{}
	suite.Run(t, ts)
}

func (ts *AuditTestSuite) SetupSuite() {
	ts.srv, ts.web, _ = tsrv.RESTHost(ts.T(), []rest.RegisterServer{
		audit_log.RegisterServer,
	}, false)
	conn, err := store.Dial(ts.srv.Config(), nil)
	ts.Require().NoError(err)
	ts.conn = conn
	ts.tests = audit_log.SetupTestLogs(ts.T(), ts.srv.Config(),
		testUID, testBook)
}

func (ts *AuditTestSuite) TestPageHeaders() {
	res, err := thttp.Do(ts.T(), ts.web, http.MethodGet, endpoint, nil, nil)
	ts.NoError(err)
	var logs []interface{}
	err = json.NewDecoder(res.Body).Decode(&logs)
	ts.NoError(err)
	ts.Len(logs, store.MaxPerPage)
	e := logs[0].(map[string]interface{})
	f := e[key.Fields].(map[string]interface{})
	id := uint(e["ID"].(float64))
	le, err := audit.GetLogEntry(ts.conn, id)
	ts.Require().NoError(err)
	ts.Equal(id, le.ID)
	ts.Equal(f["dr_suess"], le.Fields["dr_suess"])
	pn := res.Header.Get(rest.PageNumber)
	ts.Equal("1", pn)
	pc := res.Header.Get(rest.PageCount)
	cnt := int(math.Ceil(float64(len(ts.tests)) / float64(store.MaxPerPage)))
	testCount := strconv.Itoa(cnt)
	ts.Equal(testCount, pc)
	pl := res.Header.Get(rest.PageLength)
	testLen := strconv.Itoa(store.MaxPerPage)
	ts.Equal(testLen, pl)
	tot := res.Header.Get(rest.PageTotal)
	// +1 because of audit.LogStartup
	testTotal := strconv.Itoa(len(ts.tests) + 1)
	ts.Equal(testTotal, tot)
}

func (ts *AuditTestSuite) TestPageLinks() {
	startLink := func() string {
		return fmt.Sprintf("%s%s?%s=1&%s=%d",
			ts.web.URL, endpoint, key.Page, key.PerPage, store.MaxPerPage)
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
		res, err := thttp.Do(ts.T(), ts.web, http.MethodGet, uri, nil, nil)
		ts.Require().NoError(err)
		var logs []interface{}
		err = json.NewDecoder(res.Body).Decode(&logs)
		ts.Require().NoError(err)
		pc := res.Header.Get(rest.PageLength)
		cnt, err := strconv.Atoi(pc)
		ts.Require().NoError(err)
		ts.Len(logs, cnt)
		l := res.Header.Get(rest.Link)
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

func (ts *AuditTestSuite) TestSearchFilters() {
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
		res, err := thttp.Do(ts.T(), ts.web, http.MethodGet, endpoint, test.v, nil)
		ts.NoError(err)
		var logs []interface{}
		err = json.NewDecoder(res.Body).Decode(&logs)
		ts.NoError(err)
		ts.Greater(len(logs), 0)
		for _, log := range logs {
			e := log.(map[string]interface{})
			test.comp(e)
		}
	}
}

func (ts *AuditTestSuite) TestSearchSort() {
	// search Ascending
	v := url.Values{
		key.Sort:   []string{string(store.Ascending)},
		"dr_suess": []string{"thing1"},
		"sorted":   []string{"yes"},
	}
	var logs []interface{}
	res, err := thttp.Do(ts.T(), ts.web, http.MethodGet, endpoint, v, nil)
	ts.NoError(err)
	err = json.NewDecoder(res.Body).Decode(&logs)
	ts.NoError(err)
	ts.Greater(len(logs), 0)
	// reverse the indexes
	testIdx := make([]int, len(logs))
	for i := len(logs) - 1; i >= 0; i-- {
		log := logs[i]
		e := log.(map[string]interface{})
		f := e[key.Fields].(map[string]interface{})
		ts.Equal("thing1", f["dr_suess"])
		ts.Equal("yes", f["sorted"])
		testIdx[i] = int(e["ID"].(float64))
	}
	// search Descending
	v[key.Sort] = []string{string(store.Descending)}
	res, err = thttp.Do(ts.T(), ts.web, http.MethodGet, endpoint, v, nil)
	ts.NoError(err)
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
