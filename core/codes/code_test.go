package codes

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/models/code"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/test/tconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testFormats = []code.Format{code.Invite, code.PIN}

func testName(f code.Format) string {
	switch f {
	case code.Invite:
		return "Invite"
	case code.PIN:
		return "PIN"
	default:
		return ""
	}
}

func TestCreateCode(t *testing.T) {
	t.Parallel()
	for _, f := range testFormats {
		t.Run(testName(f), func(t *testing.T) {
			t.Parallel()
			testCreateCode(t, f)
		})
	}
}

func testCreateCode(t *testing.T, f code.Format) {
	conn, _ := tconn.TempConn(t)
	tests := []struct {
		uses   int
		unique bool
		max    int
	}{
		{-2, false, code.InfiniteUse},
		{-2, true, code.InfiniteUse},
		{code.SingleUse, false, 1},
		{code.SingleUse, true, 1},
		{2, false, 2},
		{2, true, 2},
		{code.InfiniteUse, false, code.InfiniteUse},
		{code.InfiniteUse, true, code.InfiniteUse},
	}
	err := conn.Transaction(func(tx *store.Connection) error {
		for _, test := range tests {
			sc, err := CreateSignupCode(tx, uuid.Nil, f, test.uses, test.unique)
			assert.NoError(t, err)
			assert.False(t, sc.CreatedAt.IsZero())
			assert.Equal(t, user.SystemID, sc.UserID)
			assert.Equal(t, f, sc.Format)
			assert.NotEmpty(t, sc.Code())
			assert.Equal(t, test.max, sc.MaxUses)
			assert.Equal(t, 0, sc.Used)
		}
		return nil
	})
	require.NoError(t, err)
}

func TestCreateCodes(t *testing.T) {
	t.Parallel()
	for _, f := range testFormats {
		t.Run(testName(f), func(t *testing.T) {
			t.Parallel()
			testCreateCodes(t, f)
		})
	}
}

func testCreateCodes(t *testing.T, f code.Format) {
	conn, _ := tconn.TempConn(t)
	tests := []struct {
		uses  int
		count int
		max   int
		len   int
	}{
		{-2, -1, code.InfiniteUse, 0},
		{code.SingleUse, 0, 1, 0},
		{2, 1, 2, 1},
		{code.InfiniteUse, 2, code.InfiniteUse, 2},
		{code.InfiniteUse, 48, code.InfiniteUse, 48},
	}
	err := conn.Transaction(func(tx *store.Connection) error {
		for _, test := range tests {
			list, err := CreateSignupCodes(tx, uuid.Nil, f, test.uses, test.count)
			assert.NoError(t, err)
			assert.Len(t, list, test.len)
			if len(list) <= 0 {
				continue
			}
			sc := list[0]
			assert.NotEmpty(t, sc.Code())
			assert.Equal(t, test.max, sc.MaxUses)
		}
		return nil
	})
	require.NoError(t, err)
}

func TestGetCode(t *testing.T) {
	t.Parallel()
	for _, f := range testFormats {
		t.Run(testName(f), func(t *testing.T) {
			t.Parallel()
			testGetCode(t, f)
		})
	}
}

func testGetCode(t *testing.T, f code.Format) {
	conn, _ := tconn.TempConn(t)
	uid := uuid.New()
	testCode := func() *code.SignupCode {
		sc, err := CreateSignupCode(conn, uid, f, code.SingleUse, true)
		require.NoError(t, err)
		return sc
	}
	deletedCode := testCode()
	err := conn.Delete(deletedCode).Error
	require.NoError(t, err)
	tests := []struct {
		sc  *code.SignupCode
		Err assert.ErrorAssertionFunc
	}{
		{&code.SignupCode{}, assert.Error},
		{code.NewSignupCode(uuid.Nil, f, code.SingleUse), assert.Error},
		{testCode(), assert.NoError},
		{deletedCode, assert.Error},
	}
	err = conn.Transaction(func(tx *store.Connection) error {
		for _, test := range tests {
			sc, err := GetSignupCode(tx, test.sc.Code())
			test.Err(t, err)
			if sc != nil {
				assert.Equal(t, f, sc.Format)
			}
		}
		return nil
	})
	require.NoError(t, err)
}

func TestGetUsableCode(t *testing.T) {
	t.Parallel()
	for _, f := range testFormats {
		t.Run(testName(f), func(t *testing.T) {
			t.Parallel()
			testGetUsableCode(t, f)
		})
	}
}

func testGetUsableCode(t *testing.T, f code.Format) {
	conn, _ := tconn.TempConn(t)
	uses := []int{
		-2, code.SingleUse, 2, code.InfiniteUse,
	}
	type testCase struct {
		code string
		Err  assert.ErrorAssertionFunc
		Nil  assert.ValueAssertionFunc
	}
	usedTest := func() testCase {
		used, err := CreateSignupCode(conn, uuid.Nil, f, code.SingleUse, true)
		assert.NoError(t, err)
		used.Used = code.SingleUse
		err = conn.Save(used).Error
		assert.NoError(t, err)
		return testCase{
			used.Code(),
			assert.Error,
			assert.Nil,
		}
	}
	deleteTest := func() testCase {
		deleted, err := CreateSignupCode(conn, uuid.Nil, f, code.SingleUse, true)
		assert.NoError(t, err)
		err = conn.Delete(deleted).Error
		assert.NoError(t, err)
		return testCase{
			deleted.Code(),
			assert.Error,
			assert.Nil,
		}
	}
	tests := []testCase{
		{"", assert.Error, assert.Nil},
		{"0xFFFFFF", assert.Error, assert.Nil},
		usedTest(),
		deleteTest(),
	}
	err := conn.Transaction(func(tx *store.Connection) error {
		for _, use := range uses {
			sc, err := CreateSignupCode(tx, uuid.Nil, f, use, true)
			assert.NoError(t, err)
			tests = append(tests, testCase{
				sc.Code(),
				assert.NoError,
				assert.NotNil,
			})
		}
		for _, test := range tests {
			sc, err := GetUsableSignupCode(tx, test.code)
			test.Err(t, err)
			test.Nil(t, sc)
		}
		return nil
	})
	require.NoError(t, err)
}

func TestCodeSent(t *testing.T) {
	t.Parallel()
	conn, _ := tconn.TempConn(t)
	uid := uuid.New()
	sc, err := CreateSignupCode(conn, uid, code.Invite, code.SingleUse, true)
	require.NoError(t, err)
	require.NotNil(t, sc)
	assert.Nil(t, sc.SentAt)
	err = SignupCodeSent(conn, sc)
	assert.NoError(t, err)
	assert.NotNil(t, sc.SentAt)
	tc, err := GetSignupCode(conn, sc.Code())
	require.NoError(t, err)
	require.NotNil(t, sc)
	assert.Equal(t, sc.SentAt.String(), tc.SentAt.String())
}

func TestGetLastSentCode(t *testing.T) {
	t.Parallel()
	conn, _ := tconn.TempConn(t)
	uid := uuid.New()
	codes := make([]*code.SignupCode, 0)
	err := conn.Transaction(func(tx *store.Connection) error {
		for i := 0; i < 3; i++ {
			sc, err := CreateSignupCode(tx, uid, code.Invite, code.SingleUse, true)
			require.NoError(t, err)
			require.NotNil(t, sc)
			if i > 0 {
				err = SignupCodeSent(tx, sc)
				assert.NoError(t, err)
			}
			for _, c := range codes {
				assert.NotEqual(t, c.ID, sc.ID)
			}
			codes = append(codes, sc)
		}
		lc, err := GetLastSentSignupCode(tx, uid)
		assert.NoError(t, err)
		assert.Equal(t, codes[len(codes)-1].ID, lc.ID)
		lc, err = GetLastSentSignupCode(tx, uuid.Nil)
		assert.NoError(t, err)
		assert.Nil(t, lc)
		return nil
	})
	require.NoError(t, err)
}
