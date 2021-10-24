package admin

import (
	"sort"
	"testing"

	"github.com/google/uuid"
	api "github.com/jrapoport/gothic/api/grpc/rpc"
	"github.com/jrapoport/gothic/api/grpc/rpc/admin"
	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/core/audit"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/models/auditlog"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/structpb"
)

type Page struct {
	Number   int64
	PageSize int64
}

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

type AdminAuditServerTestSuite struct {
	suite.Suite
	srv   *server
	conn  *store.Connection
	tests []TestCase
	uid   uuid.UUID
}

func TestAuditLogs(t *testing.T) {
	t.Parallel()
	ts := &AdminAuditServerTestSuite{}
	suite.Run(t, ts)
}

func (ts *AdminAuditServerTestSuite) SetupSuite() {
	s, _ := tsrv.RPCServer(ts.T(), false)
	ts.srv = newAdminServer(s)
	conn, err := store.Dial(ts.srv.Config(), nil)
	ts.Require().NoError(err)
	ts.conn = conn
	ts.tests = SetupTestLogs(ts.T(), ts.srv.Config(),
		testUID, testBook)
}

func (ts *AdminAuditServerTestSuite) searchAuditLogs(s store.Sort, f store.Filters, p *Page) (*admin.AuditLogsResult, error) {
	var err error
	ctx := rootContext(ts.srv.Config())
	srt := api.Sort(s)
	req := &api.SearchRequest{}
	req.Sort = &srt
	if p != nil {
		req.Page = p.Number
		req.PageSize = &p.PageSize
	}
	if f != nil {
		req.Filters, err = structpb.NewStruct(f)
		ts.Require().NoError(err)
	}
	return ts.srv.SearchAuditLogs(ctx, req)
}

func (ts *AdminAuditServerTestSuite) TestSearchFilters() {
	tests := []struct {
		filters store.Filters
		comp    func(test *admin.AuditLog)
	}{
		{
			store.Filters{
				key.UserID: testUID.String(),
			},
			func(res *admin.AuditLog) {
				ts.Equal(testUID.String(), res.UserId)
			},
		},
		{
			store.Filters{
				key.Action: auditlog.Startup.String(),
			},
			func(res *admin.AuditLog) {
				act := res.Action
				ts.Equal(auditlog.Startup.String(), act)
			},
		},
		{
			store.Filters{
				"dr_suess": "thing1",
			},
			func(res *admin.AuditLog) {
				f := res.Fields.AsMap()
				ts.Equal("thing1", f["dr_suess"])
			},
		},
		{
			store.Filters{
				key.Type:   auditlog.Account.String(),
				"dr_suess": "thing1",
			},
			func(res *admin.AuditLog) {
				typ := auditlog.Type(res.Type)
				ts.Equal(auditlog.Account, typ)
				f := res.Fields.AsMap()
				ts.Equal("thing1", f["dr_suess"])
			},
		},
		{
			store.Filters{
				"dr_suess": "thing1",
				"book":     testBook.String(),
			},
			func(res *admin.AuditLog) {
				f := res.Fields.AsMap()
				ts.Equal("thing1", f["dr_suess"])
				ts.Equal(testBook.String(), f["book"])
			},
		},
	}
	for _, test := range tests {
		res, err := ts.searchAuditLogs(store.Ascending, test.filters, nil)
		ts.Require().NoError(err)
		ts.Require().NotNil(res)
		logs := res.Logs
		ts.Greater(len(logs), 0)
		for _, log := range logs {
			test.comp(log)
		}
	}
}

func (ts *AdminAuditServerTestSuite) TestSearchSort() {
	// search Ascending
	filters := store.Filters{
		"dr_suess": "thing1",
		"sorted":   "yes",
	}
	res, err := ts.searchAuditLogs(store.Ascending, filters, nil)
	ts.Require().NoError(err)
	ts.Require().NotNil(res)
	logs := res.Logs
	ts.Greater(len(logs), 0)
	testIdx := make([]int, len(logs))
	for i, log := range logs {
		f := log.Fields.AsMap()
		ts.Equal("thing1", f["dr_suess"])
		ts.Equal("yes", f["sorted"])
		testIdx[i] = int(log.Id)
	}
	// search Descending
	res, err = ts.searchAuditLogs(store.Descending, filters, nil)
	ts.Require().NoError(err)
	ts.Require().NotNil(res)
	logs = res.Logs
	ts.Greater(len(logs), 0)
	ts.Require().Len(logs, len(testIdx))
	descIdx := make([]int, len(logs))
	for i, log := range logs {
		f := log.Fields.AsMap()
		ts.Equal("thing1", f["dr_suess"])
		ts.Equal("yes", f["sorted"])
		descIdx[i] = int(log.Id)
	}
	// reverse the indexes
	sort.Ints(descIdx)
	ts.Equal(testIdx, descIdx)
}

func (ts *AdminAuditServerTestSuite) TestPaging() {
	res, err := ts.searchAuditLogs(store.Ascending, nil, nil)
	ts.Require().NoError(err)
	ts.Require().NotNil(res)
	ts.Len(res.Logs, store.MaxPageSize)
	log := res.Logs[0]
	f := log.Fields.AsMap()
	id := uint(log.Id)
	le, err := audit.GetLogEntry(ts.conn, id)
	ts.Require().NoError(err)
	ts.Equal(id, le.ID)
	ts.Equal(f["dr_suess"], le.Fields["dr_suess"])
	pn := int(res.Page.Index)
	ts.Equal(1, pn)
	pc := int(res.Page.Count)
	cnt := len(ts.tests) / store.MaxPageSize
	ts.Equal(cnt+1, pc)
	sz := int(res.Page.Size)
	ts.Equal(store.MaxPageSize, sz)
	tot := int(res.Page.Total)
	// +1 because of audit.LogStartup
	ts.Equal(len(ts.tests)+1, tot)
	page := &Page{Number: 1}
	var last uint64
	for i := 0; i < pc; i++ {
		res, err = ts.searchAuditLogs(store.Ascending, nil, page)
		ts.Require().NoError(err)
		ts.Require().NotNil(res)
		ts.Equal(len(res.Logs), int(res.Page.Size))
		ts.Equal(i+1, int(res.Page.Index))
		ts.NotEqual(last, res.Logs[0].Id)
		last = res.Logs[0].Id
		page.Number = res.Page.Index + 1
	}
}

func (ts *AdminAuditServerTestSuite) TestErrors() {
	// invalid req
	_, err := ts.srv.SearchAuditLogs(nil, nil)
	ts.Error(err)
	// invalid req
	_, err = ts.searchAuditLogs(store.Ascending, store.Filters{
		key.UserID: "1",
	}, nil)
	ts.Error(err)
}
