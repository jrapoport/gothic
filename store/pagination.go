package store

import (
	"math"

	"github.com/vcraescu/go-paginator/v2"
	"github.com/vcraescu/go-paginator/v2/adapter"
	"gorm.io/gorm"
)

// MaxPageSize is the default max per page.
const MaxPageSize = 50

// Pagination holds paged results.
type Pagination struct {
	// Index is the current page number.
	Index int
	// Size is the max page size.
	Size int
	// Prev is the index of the previous page, or 0 if there is no previous page.
	Prev int
	// Next is the index of the next page, or 0 if there is no next page.
	Next int
	// Count is the total pages of pages.
	Count int
	// Items holds the result slice.
	Items interface{}
	// Length is the actual number of items.
	Length int
	// Total is the total number of items across all pages.
	Total uint64
}

// Sort is the sort order for the results.
type Sort uint8

// Sort orders
const (
	Ascending Sort = iota
	Descending
)

func (s Sort) String() string {
	switch s {
	case Descending:
		return "DESC"
	case Ascending:
		fallthrough
	default:
		return "ASC"
	}
}

// NextPage fetches the next page.
func NextPage(query *gorm.DB, page *Pagination, s Sort) (*Pagination, error) {
	if page.Size == 0 {
		page.Size = MaxPageSize
	}
	q := query.Order("created_at " + s.String())
	p := paginator.New(adapter.NewGORMAdapter(q), page.Size)
	p.SetPage(page.Index)
	err := p.Results(page.Items)
	if err != nil {
		return nil, err
	}
	total, err := p.Nums()
	if err != nil {
		return nil, err
	}
	page.Total = uint64(total)
	page.Count = int(math.Ceil(float64(page.Total) / float64(page.Size)))
	page.Prev = 0
	if page.Index > 1 {
		page.Prev = page.Index - 1
	}
	page.Next = 0
	if page.Index < page.Count {
		page.Next = page.Index + 1
	}
	page.Length = int(page.Total - uint64(page.Size*(page.Index-1)))
	if page.Length > page.Size {
		page.Length = page.Size
	}
	return page, nil
}
