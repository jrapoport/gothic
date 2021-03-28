package store

import (
	"github.com/jrapoport/gothic/models/types"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Filters hols the search filters.
type Filters types.Map

// Copy returns a copy of the filter
func (f Filters) Copy() Filters {
	return Filters(types.Map(f).Copy())
}

// FiltersFromMap returns a set of filters from a string map.
func FiltersFromMap(m map[string]string) Filters {
	d := make(Filters, len(m))
	for k, v := range m {
		d[k] = v
	}
	return d
}

// Filter to apply
type Filter struct {
	Filters   Filters
	DataField string
	Fields    []string
}

// Search searches a table for hits
func Search(tx *gorm.DB, models interface{}, s Sort, f Filter, p *Pagination) error {
	if s == "" {
		s = Descending
	}
	filters := f.Filters.Copy()
	for k, v := range filters {
		if v == "" {
			delete(filters, k)
		}
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
