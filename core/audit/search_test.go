package audit

import (
	"math"
	"testing"

	"github.com/jrapoport/gothic/models/auditlog"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/store/types/key"
	"github.com/jrapoport/gothic/test/tconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const pageSize = 10

func createEntries(t *testing.T, conn *store.Connection) []testCase {
	tests := testCases()
	for i, test := range tests {
		le, err := CreateLogEntry(test.ctx, conn, test.act, test.uid, test.fields)
		require.NoError(t, err)
		require.NotNil(t, le)
		tests[i].logID = le.ID
	}
	return tests
}

func TestSearchEntries(t *testing.T) {
	conn, _ := tconn.TempConn(t)
	var tests []testCase
	t.Run("TestEntries", func(t *testing.T) {
		tests = createEntries(t, conn)
	})
	// find all no page
	logs, err := SearchEntries(conn, store.Descending, nil, nil)
	assert.NoError(t, err)
	assert.Len(t, logs, len(tests))
	for i := 0; i < 100; i += 10 {
		test := tests[i]
		log := logs[i]
		assert.Equal(t, test.act, log.Action)
		assert.Equal(t, test.uid, log.UserID)
		for k, v := range test.fields {
			assert.EqualValues(t, v, log.Fields[k])
		}
	}
	// find all page
	page := &store.Pagination{
		Page: 1,
		Size: 40,
	}
	logs, err = SearchEntries(conn, store.Descending, nil, page)
	assert.NoError(t, err)
	assert.Len(t, logs, 40)
	log := logs[20]
	for _, test := range tests {
		if test.logID == log.ID {
			assert.Equal(t, test.act, log.Action)
			assert.Equal(t, test.uid, log.UserID)
			for k, v := range test.fields {
				assert.EqualValues(t, v, log.Fields[k])
			}
		}
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
			key.Type: auditlog.Account,
		}, tests, paged)
		assert.Equal(t, test(48, paged), cnt)
		cnt = filtersTest(t, conn, store.Filters{
			key.Type:   auditlog.Signup.Type(),
			key.Action: auditlog.Signup,
		}, tests, paged)
		assert.Equal(t, test(16, paged), cnt)
		cnt = filtersTest(t, conn, store.Filters{
			key.Action: auditlog.Signup,
			key.Type:   auditlog.Signup.Type(),
			key.UserID: testUID,
		}, tests, paged)
		assert.Equal(t, 8, cnt)
		// re-run the last test with type conversion
		cnt = filtersTest(t, conn, store.Filters{
			key.Action: auditlog.Signup,
			key.Type:   auditlog.Signup.Type().String(),
			key.UserID: testUID.String(),
		}, tests, paged)
		assert.Equal(t, 8, cnt)
		cnt = filtersTest(t, conn, store.Filters{
			"book": "thing2",
		}, tests, paged)
		assert.Equal(t, test(26, paged), cnt)
		cnt = filtersTest(t, conn, store.Filters{
			key.Action: auditlog.Signup,
			"book":     "thing2",
		}, tests, paged)
		assert.Equal(t, 2, cnt)
		cnt = filtersTest(t, conn, store.Filters{
			"book":  testBook,
			"extra": "salad-0",
		}, tests, paged)
		assert.Equal(t, 0, cnt)
		var tc testCase
		for _, tst := range tests {
			v, ok := tst.fields["pepper"]
			if ok && v == "spicy" {
				tc = tst
				break
			}
		}
		cnt = filtersTest(t, conn, store.Filters{
			"book":   tc.fields["book"],
			"extra":  tc.fields["extra"],
			"pepper": tc.fields["pepper"],
		}, tests, paged)
		assert.Equal(t, 1, cnt)
	}
}

func filtersTest(t *testing.T, conn *store.Connection, f store.Filters, tests []testCase, paged bool) int {
	var page *store.Pagination
	if paged {
		page = &store.Pagination{
			Page: 0,
			Size: pageSize,
		}
	}
	list, err := SearchEntries(conn, store.Descending, f, page)
	assert.NoError(t, err)
	cnt := 0
	for _, test := range tests {
		if v, ok := f[key.Type]; ok {
			if test.act.Type() != v {
				continue
			}
		}
		if v, ok := f[key.Action]; ok {
			if test.act != v {
				continue
			}
		}
		if v, ok := f[key.UserID]; ok {
			if test.uid != v {
				continue
			}
		}
		if v, ok := f["book"]; ok {
			if test.fields["book"] != v {
				continue
			}
		}
		if v, ok := f["extra"]; ok {
			if test.fields["extra"] != v {
				continue
			}
		}
		if v, ok := f["pepper"]; ok {
			if test.fields["pepper"] != v {
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
			case key.Type:
				assert.Equal(t, v, l.Type)
			case key.Action:
				assert.Equal(t, v, l.Action)
			case key.UserID:
				assert.Equal(t, v, l.UserID)
			case "book":
				assert.Equal(t, v, l.Fields["book"])
			case "extra":
				assert.Equal(t, v, l.Fields["extra"])
			case "pepper":
				assert.Equal(t, v, l.Fields["pepper"])
			}
		}
	}
	return cnt
}
