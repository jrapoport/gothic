package types

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/jrapoport/gothic/store/drivers"
	"github.com/markbates/goth"
	"github.com/segmentio/encoding/json"
	"gorm.io/gorm"
	"gorm.io/gorm/migrator"
	"gorm.io/gorm/schema"
)

// Map is a map that satisfies various db and other interfaces.
type Map map[string]interface{}

var (
	_ goth.Params                    = (*Map)(nil)
	_ sql.Scanner                    = (*Map)(nil)
	_ driver.Valuer                  = (*Map)(nil)
	_ schema.GormDataTypeInterface   = (*Map)(nil)
	_ migrator.GormDataTypeInterface = (*Map)(nil)
	_ fmt.Stringer                   = (*Map)(nil)
)

// Copy returns a copy
func (m Map) Copy() Map {
	cpy := make(Map, len(m))
	for k, v := range m {
		cpy[k] = v
	}
	return cpy
}

// Get satisfies the goth.Params interface
func (m Map) Get(k string) string {
	v, ok := m[k].(string)
	if !ok {
		return ""
	}
	return v
}

// JSON returns a json representation of the map
func (m Map) JSON() ([]byte, error) {
	return json.Marshal(m)
}

// Value satisfies the driver.Valuer interface
func (m Map) Value() (driver.Value, error) {
	data, _ := m.JSON()
	return driver.Value(string(data)), nil
}

// Scan satisfies the sql.Scanner interface
func (m *Map) Scan(src interface{}) error {
	var source []byte
	switch v := src.(type) {
	case string:
		source = []byte(v)
	case []byte:
		source = v
	default:
		return errors.New("invalid Map data")
	}
	if len(source) == 0 {
		source = []byte("{}")
	}
	return json.Unmarshal(source, m)
}

// GormDataType gorm common data type for schema.GormDataTypeInterface
func (Map) GormDataType() string {
	return "json"
}

// GormDBDataType gorm db data type for migrator.GormDataTypeInterface
func (Map) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	switch db.Name() {
	case drivers.MySQL:
		return "JSON"
	case drivers.Postgres:
		return "JSONB"
	case drivers.SQLServer:
		return "NVARCHAR(MAX)"
	case drivers.SQLite, drivers.SQLite3:
		return "JSON"
	}
	return "JSON"
}

// String satisfies the fmt.Stringer interface.
func (m Map) String() string {
	b, _ := m.JSON()
	return string(b)
}

// DataFromMap returns a generic map as a Map.
func DataFromMap(m map[string]interface{}) Map {
	d := make(Map, len(m))
	for k, v := range m {
		d[k] = v
	}
	return d
}
