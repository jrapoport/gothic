package store

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/store/migration"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/jrapoport/gothic/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

const pageSize = 4

type SearchModel struct {
	gorm.Model
	AccountID uuid.UUID
	Name      string
	Data      types.Map
}

type testCase struct {
	account uuid.UUID
	name    string
	data    types.Map
}

var testName = utils.RandomUsername()
var testBook = utils.RandomUsername()

func testSearchModels(t *testing.T, conn *Connection) []testCase {
	var cases = func() []testCase {
		name := func() string { return utils.RandomUsername() }
		data := types.Map{
			"dr_suess": "thing1",
			"book":     "thing2",
		}
		var tests = []testCase{
			{uuid.New(), name(), data},
			{uuid.New(), name(), data},
			{uuid.New(), name(), data},
			{uuid.New(), name(), data},
			{uuid.New(), name(), data},
			{uuid.New(), name(), data},
			{uuid.New(), name(), data},
			{uuid.New(), name(), data},
			{uuid.New(), name(), data},
		}
		for _, test := range tests {
			ext := test
			ext.name = testName
			test.account = uuid.New()
			tests = append(tests, ext)
		}
		var once sync.Once
		for i, bk := range []interface{}{
			testBook, uuid.New().String(),
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
				test.account = uuid.New()
				tests = append(tests, test)
			}
		}
		return tests
	}()
	err := conn.Transaction(func(tx *Connection) error {
		for _, test := range cases {
			m := &SearchModel{
				AccountID: test.account,
				Name:      test.name,
				Data:      test.data,
			}
			err := tx.Create(m).Error
			require.NoError(t, err)
		}
		return nil
	})
	require.NoError(t, err)
	return cases
}

func doSearch(t *testing.T, conn *Connection, s Sort, f Filters, p *Pagination) []*SearchModel {
	tx := conn.Model(new(SearchModel))
	flt := Filter{
		Filters:   f,
		DataField: key.Data,
		Fields: []string{
			key.AccountID,
			key.Name,
		},
	}
	var models []*SearchModel
	err := Search(tx, &models, s, flt, p)
	require.NoError(t, err)
	return models
}

func TestSearch(t *testing.T) {
	t.Parallel()
	c := tconf.TempDB(t)
	conn, err := Dial(c, nil)
	require.NoError(t, err)
	require.NotNil(t, conn)
	p := migration.Plan{
		migration.NewMigration("1", SearchModel{}),
	}
	err = p.Run(conn.DB, false)
	require.NoError(t, err)
	var tests []testCase
	t.Run("Populate Search Models", func(t *testing.T) {
		tests = testSearchModels(t, conn)
	})
	// find all no page
	list := doSearch(t, conn, Descending, nil, nil)
	assert.Len(t, list, len(tests))
	for _, idx := range []int{0, 5, 10, 20} {
		test := tests[idx]
		item := list[idx]
		assert.Equal(t, test.account, item.AccountID)
		assert.Equal(t, test.name, item.Name)
		assert.Equal(t, test.data, item.Data)
	}
	// find all page
	var size = len(tests) / 2
	page := &Pagination{
		Page: 1,
		Size: size,
	}
	list = doSearch(t, conn, Descending, nil, page)
	assert.NoError(t, err)
	assert.Len(t, list, size)
	item := list[10]
	var found bool
	for _, test := range tests {
		found = test.account == item.AccountID
		if found {
			assert.Equal(t, test.name, item.Name)
			assert.Equal(t, test.data, item.Data)
			break
		}
	}
	assert.True(t, found)
	for _, filter := range []string{
		key.AccountID,
		key.Name,
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
	var tc testCase
	for _, paged := range []bool{false, true} {
		cnt := filtersTest(t, conn, Filters{
			key.Name: testName,
		}, tests, paged)
		assert.Equal(t, test(36, paged), cnt)
		cnt = filtersTest(t, conn, Filters{
			key.AccountID: item.AccountID,
			key.Name:      item.Name,
		}, tests, paged)
		assert.Equal(t, 1, cnt)
		cnt = filtersTest(t, conn, Filters{
			"book": "thing2",
		}, tests, paged)
		assert.Equal(t, test(18, paged), cnt)
		cnt = filtersTest(t, conn, Filters{
			key.Name: testName,
			"book":   "thing2",
		}, tests, paged)
		assert.Equal(t, test(9, paged), cnt)
		cnt = filtersTest(t, conn, Filters{
			"book":  item.Data["book"],
			"extra": item.Data["extra"],
		}, tests, paged)
		assert.Equal(t, 1, cnt)
		cnt = filtersTest(t, conn, Filters{
			"book":  "thing2",
			"extra": "caesar",
		}, tests, paged)
		assert.Equal(t, 0, cnt)
		for _, tst := range tests {
			v, ok := tst.data["pepper"]
			if ok && v == "spicy" {
				tc = tst
				break
			}
		}
		cnt = filtersTest(t, conn, Filters{
			"book":   tc.data["book"],
			"extra":  tc.data["extra"],
			"pepper": tc.data["pepper"],
		}, tests, paged)
		assert.Equal(t, 1, cnt)
	}
	// NOT name
	list = doSearch(t, conn, Descending, Filters{
		key.Name + "!": testName,
		"book":         "thing2",
	}, nil)
	assert.Len(t, list, 9)
	for _, item = range list {
		assert.NotEqual(t, testName, item.Name)
	}
	// NOT name paged
	page = &Pagination{
		Page: 1,
		Size: pageSize,
	}
	list = doSearch(t, conn, Descending, Filters{
		key.Name + "!": testName,
		"book":         "thing2",
	}, page)
	assert.Len(t, list, 4)
}

func filterTest(t *testing.T, conn *Connection, f string, tests []testCase) {
	filters := Filters{}
	tc := tests[20]
	switch f {
	case key.AccountID:
		filters[f] = tc.account
		cnt := filtersTest(t, conn, filters, tests, false)
		assert.Equal(t, 1, cnt)
	case key.Name:
		filters[f] = testName
		cnt := filtersTest(t, conn, filters, tests, false)
		assert.Equal(t, 36, cnt)
	}
}

func filtersTest(t *testing.T, conn *Connection, f Filters, tests []testCase, paged bool) int {
	var page *Pagination
	if paged {
		page = &Pagination{
			Page: 0,
			Size: pageSize,
		}
	}
	list := doSearch(t, conn, Descending, f, page)
	cnt := 0
	for _, test := range tests {
		if v, ok := f[key.AccountID]; ok {
			if test.account != v {
				continue
			}
		}
		if v, ok := f[key.Name]; ok {
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
			case key.AccountID:
				assert.Equal(t, v, l.AccountID)
			case key.Name:
				assert.Equal(t, v, l.Name)
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
