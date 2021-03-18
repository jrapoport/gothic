package core

import (
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/core/audit"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/models/auditlog"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testLog struct {
	act    auditlog.Action
	uid    uuid.UUID
	fields types.Map
	logID  uint
}

var (
	testUID  = uuid.New()
	testBook = uuid.New().String()
)

var _testLogs []testLog
var _logAPI *API
var once sync.Once

func setupLogs(t *testing.T) (*API, []testLog) {
	once.Do(func() {
		a := apiWithTempDB(t)
		var tests = []testLog{
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
			tst.uid = testUID
			tests = append(tests, tst)
		}
		for _, bk := range []interface{}{
			"thing2", testBook, uuid.New().String(),
		} {
			for _, test := range tests {
				test.fields = types.Map{
					"dr_suess": "thing1",
					"book":     bk,
				}
				tests = append(tests, test)
			}
		}
		ctx := context.Background()
		err := a.conn.Transaction(func(tx *store.Connection) error {
			for i, test := range tests {
				le, err := audit.CreateLogEntry(ctx, tx, test.act, test.uid, test.fields)
				require.NoError(t, err)
				tests[i].logID = le.ID
			}
			return nil
		})
		require.NoError(t, err)
		_testLogs = tests
		_logAPI = a
	})
	require.NotNil(t, _logAPI)
	require.NotEmpty(t, _testLogs)
	return _logAPI, _testLogs
}

func TestAPI_GetAuditLog(t *testing.T) {
	t.Parallel()
	a, tests := setupLogs(t)
	_, err := a.GetAuditLog(nil, 9999)
	assert.Error(t, err)
	for _, test := range tests {
		le, err := a.GetAuditLog(nil, test.logID)
		assert.NoError(t, err)
		assert.EqualValues(t, test.uid, le.UserID)
		assert.EqualValues(t, test.act, le.Action)
		for k, v := range test.fields {
			assert.EqualValues(t, v, le.Fields[k])
		}
	}
}

func TestAPI_SearchAuditLogs(t *testing.T) {
	t.Parallel()
	a, _ := setupLogs(t)
	tests := []struct {
		filters store.Filters
		comp    func(log *auditlog.AuditLog)
	}{
		{
			store.Filters{
				key.UserID: testUID.String(),
			},
			func(log *auditlog.AuditLog) {
				assert.Equal(t, testUID, log.UserID)
			},
		},
		{
			store.Filters{
				key.Action: auditlog.Startup.String(),
			},
			func(log *auditlog.AuditLog) {
				assert.Equal(t, auditlog.Startup, log.Action)
			},
		},
		{
			store.Filters{
				"dr_suess": "thing1",
			},
			func(log *auditlog.AuditLog) {
				assert.Equal(t, "thing1", log.Fields["dr_suess"])
			},
		},
		{
			store.Filters{
				key.Type:   auditlog.Account.String(),
				"dr_suess": "thing1",
			},
			func(log *auditlog.AuditLog) {
				assert.Equal(t, auditlog.Account, log.Type)
				assert.Equal(t, "thing1", log.Fields["dr_suess"])
			},
		},
		{
			store.Filters{
				"dr_suess": "thing1",
				"book":     testBook,
			},
			func(log *auditlog.AuditLog) {
				assert.Equal(t, "thing1", log.Fields["dr_suess"])
				assert.Equal(t, testBook, log.Fields["book"])
			},
		},
	}
	for _, test := range tests {
		logs, err := a.SearchAuditLogs(nil, test.filters, nil)
		assert.NoError(t, err)
		assert.Greater(t, len(logs), 0)
		for _, log := range logs {
			test.comp(log)
		}
	}
}

func TestAPI_SearchAuditLogs_Sort(t *testing.T) {
	t.Parallel()
	filters := store.Filters{
		"dr_suess": []string{"thing1"},
		"book":     []string{testBook},
	}
	a, _ := setupLogs(t)
	ctx := testContext(a)
	ctx.SetSort(store.Descending)
	logs, err := a.SearchAuditLogs(ctx, filters, nil)
	assert.NoError(t, err)
	assert.Greater(t, len(logs), 0)
	// reverse the indexes
	testIdx := make([]uint, len(logs))
	for i := len(logs) - 1; i >= 0; i-- {
		log := logs[i]
		assert.Equal(t, "thing1", log.Fields["dr_suess"])
		assert.Equal(t, testBook, log.Fields["book"])
		testIdx[i] = log.ID
	}
	// reverse the sort (and the indexes)
	ctx.SetSort(store.Ascending)
	logs, err = a.SearchAuditLogs(ctx, filters, nil)
	assert.NoError(t, err)
	assert.Greater(t, len(logs), 0)
	require.Len(t, logs, len(testIdx))
	descIdx := make([]uint, len(logs))
	for i, log := range logs {
		assert.Equal(t, "thing1", log.Fields["dr_suess"])
		assert.Equal(t, testBook, log.Fields["book"])
		descIdx[i] = log.ID
	}
	assert.Equal(t, testIdx, descIdx)
}
