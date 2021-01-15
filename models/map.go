package models

import (
	"database/sql/driver"
	"errors"

	"github.com/segmentio/encoding/json"
)

type Map map[string]interface{}

func (m Map) Value() (driver.Value, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return driver.Value(""), err
	}
	return driver.Value(string(data)), nil
}

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
