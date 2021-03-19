package store

import (
	"github.com/jrapoport/gothic/models/types"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Filters hols the search filters.
type Filters types.Map

// Filter to apply
type Filter struct {
	Filters   Filters
	DataField string
	Fields    []string
}

// Search searches a table for hits
func Search(tx *gorm.DB, models interface{}, s Sort, f Filter, p *Pagination, cond ...string) error {
	if s == "" {
		s = Descending
	}
	filters := make(Filters, len(f.Filters))
	for k, v := range f.Filters {
		filters[k] = v
	}
	for _, field := range f.Fields {
		if v, ok := filters[field]; ok {
			tx = tx.Where(field+" = ?", v)
			delete(filters, field)
		}
	}
	for _, field := range f.Fields {
		if v, ok := filters[field+"!"]; ok {
			tx = tx.Where(field+" <> ?", v)
			delete(filters, field+"!")
		}
	}
	if len(cond) > 0 {
		for _, c := range cond {
			tx = tx.Where(c)
		}
	}
	for k, v := range filters {
		tx = tx.Where(datatypes.JSONQuery(f.DataField).Equals(v, k))
	}
	var err error
	if p == nil {
		err = tx.Find(models).Error
	} else {
		p.Items = models
		p, err = NextPage(tx, p, s)
	}
	return err
}
