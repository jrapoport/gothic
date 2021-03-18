package audit

import (
	"fmt"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/models/auditlog"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/store/types"
	"github.com/jrapoport/gothic/store/types/key"
	"github.com/jrapoport/gothic/test/tconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	ctx    context.Context
	act    auditlog.Action
	uid    uuid.UUID
	fields types.Map
	logID  uint
}

const testProvider = "test_provider"
const testIPAddress = "127.0.0.1"

var (
	testUID  = uuid.New()
	testBook = uuid.New().String()
)

func testContext() context.Context {
	ctx := context.Background()
	ctx.SetIPAddress(testIPAddress)
	ctx.SetProvider(testProvider)
	ctx.SetUserID(testUID)
	ctx.SetAdminID(testUID)
	return ctx
}

func testCases() []testCase {
	ctx := testContext()
	var tests = []testCase{
		{ctx, auditlog.Startup, uuid.New(), nil, 0},
		{ctx, auditlog.Shutdown, uuid.New(), nil, 0},
		{ctx, auditlog.Signup, uuid.New(), nil, 0},
		{ctx, auditlog.ConfirmSent, uuid.New(), nil, 0},
		{ctx, auditlog.Confirmed, uuid.New(), nil, 0},
		{ctx, auditlog.Granted, uuid.New(), nil, 0},
		{ctx, auditlog.Revoked, uuid.New(), nil, 0},
		{ctx, auditlog.RevokedAll, uuid.New(), nil, 0},
		{ctx, auditlog.Login, uuid.New(), nil, 0},
		{ctx, auditlog.Logout, uuid.New(), nil, 0},
		{ctx, auditlog.Password, uuid.New(), nil, 0},
		{ctx, auditlog.Email, uuid.New(), nil, 0},
		{ctx, auditlog.Updated, uuid.New(), nil, 0},
	}
	for _, test := range tests {
		tst := test
		tst.uid = testUID
		tests = append(tests, tst)
	}
	var once sync.Once
	for i, bk := range []interface{}{
		"thing2", testBook, uuid.New().String(),
	} {
		for x, test := range tests {
			test.fields = types.Map{
				"dr_suess":    "thing1",
				"book":        bk,
				key.IPAddress: test.ctx.GetIPAddress(),
				key.Provider:  test.ctx.GetProvider(),
				key.UserID:    test.ctx.GetUserID().String(),
				key.AdminID:   test.ctx.GetAdminID().String(),
			}
			sld := fmt.Sprintf("salad-%d", x+i)
			test.fields["extra"] = sld
			once.Do(func() {
				test.fields["pepper"] = "spicy"
			})
			tests = append(tests, test)
		}
	}
	return tests
}

func TestCreateLogEntry(t *testing.T) {
	t.Parallel()
	conn, _ := tconn.TempConn(t)
	fields := types.Map{
		key.IPAddress: "test",
		key.Provider:  "test",
		key.UserID:    "test",
		key.AdminID:   "test",
	}
	ctx := testContext()
	le, err := CreateLogEntry(ctx, conn, auditlog.Login, testUID, fields)
	assert.NoError(t, err)
	assert.Equal(t, fields, le.Fields)
	ctx = context.Background()
	ctx.SetIPAddress(testIPAddress)
	ctx.SetProvider(testProvider)
	ctx.SetUserID(testUID)
	le, err = CreateLogEntry(ctx, conn, auditlog.Login, testUID, nil)
	assert.NoError(t, err)
	fields = types.Map{
		key.IPAddress: ctx.GetIPAddress(),
		key.Provider:  ctx.GetProvider(),
		key.UserID:    ctx.GetUserID().String(),
	}
	assert.Equal(t, fields, le.Fields)
	assert.Nil(t, le.Fields[key.AdminID])
	for _, test := range testCases() {
		testCreate(t, conn,
			createTest{
				test.act,
				test.uid,
				test.fields,
				func(ctx context.Context, conn *store.Connection, uid uuid.UUID, fields types.Map) error {
					_, err = CreateLogEntry(test.ctx, conn, test.act, uid, fields)
					return err
				},
			})
	}
}

type logFunc func(ctx context.Context, conn *store.Connection, uid uuid.UUID, fields types.Map) error

type createTest struct {
	a    auditlog.Action
	uid  uuid.UUID
	data types.Map
	fn   logFunc
}

func testCreate(t *testing.T, conn *store.Connection, test createTest) *auditlog.AuditLog {
	require.NotEqual(t, auditlog.Unknown, test.a.Type())
	ctx := testContext()
	err := test.fn(ctx, conn, test.uid, test.data)
	require.NoError(t, err)
	le := getLast(t, conn)
	require.NotNil(t, le)
	assert.Equal(t, test.uid, le.UserID)
	assert.Equal(t, test.a, le.Action)
	assert.Equal(t, test.a.Type(), le.Type)
	te, err := GetLogEntry(conn, le.ID)
	assert.NoError(t, err)
	assert.Equal(t, le.UserID, te.UserID)
	for k, v := range test.data {
		assert.EqualValues(t, v, le.Fields[k])
	}
	return le
}

func getLast(t *testing.T, conn *store.Connection) *auditlog.AuditLog {
	le := &auditlog.AuditLog{}
	err := conn.Last(le).Error
	require.NoError(t, err)
	return le
}

func testLogEntry(t *testing.T, a auditlog.Action, uid uuid.UUID, data types.Map, fn logFunc) *auditlog.AuditLog {
	conn, _ := tconn.TempConn(t)
	return testCreate(t, conn, createTest{a, uid, data, fn})
}
