package types

import (
	"github.com/jrapoport/gothic/store/drivers"
	"testing"

	"github.com/jrapoport/gothic/test/tdb"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

func TestData(t *testing.T) {
	t.Parallel()
	type Tests struct {
		gorm.Model
		Name       string
		Attributes Map
	}
	db := tdb.DB(t)
	err := db.Migrator().AutoMigrate(&Tests{})
	assert.NoError(t, err)
	test1 := Map{
		"age":  18,
		"name": "data-1",
		"orgs": Map{"orga": "orga"},
		"tags": []string{"tag1", "tag2"},
	}
	test2 := Map{
		"name": "data-2",
		"age":  28,
		"tags": []string{"tag1", "tag3"},
		"role": "admin",
		"orgs": Map{"orgb": "orgb"},
	}
	tests := []Tests{
		{
			Name:       "data-1",
			Attributes: test1,
		}, {
			Name:       "data-2",
			Attributes: test2,
		}}
	err = db.Create(&tests).Error
	assert.NoError(t, err)
	var result Tests
	q := datatypes.JSONQuery("attributes").HasKey("role")
	err = db.First(&result, q).Error
	assert.NoError(t, err)
	assert.Equal(t, result.Name, tests[1].Name)
	var result2 Tests
	q = datatypes.JSONQuery("attributes").HasKey("orgs", "orga")
	err = db.First(&result2, q).Error
	assert.NoError(t, err)
	assert.Equal(t, result2.Name, tests[0].Name)
	var result3 Tests
	q = datatypes.JSONQuery("attributes").Equals("data-1", "name")
	err = db.First(&result3, q).Error
	assert.NoError(t, err)
	assert.Equal(t, result3.Name, tests[0].Name)
	var result4 Tests
	q = datatypes.JSONQuery("attributes").Equals("orgb", "orgs", "orgb")
	err = db.First(&result4, q).Error
	assert.NoError(t, err)
	assert.Equal(t, result4.Name, tests[1].Name)
	test1["age"] = 19
	b, err := test1.JSON()
	assert.NoError(t, err)
	jsonMap := map[string]interface{}{"Attributes": string(b)}
	err = db.Where(&Tests{Name: "data-1"}).Assign(jsonMap).FirstOrCreate(&Tests{}).Error
	assert.NoError(t, err)
	var result5 Tests
	q = datatypes.JSONQuery("attributes").Equals(19, "age")
	err = db.First(&result5, q).Error
	assert.NoError(t, err)
}

func TestMap_Scan(t *testing.T) {
	m := Map{}
	i := 0
	err := m.Scan(i)
	assert.Error(t, err)
	err = m.Scan([]byte("\n"))
	assert.Error(t, err)
	err = m.Scan([]byte(""))
	assert.NoError(t, err)
}

func TestDataFromMap(t *testing.T) {
	t.Parallel()
	const testString = "{\"foo\":\"bar\",\"quack\":1}"
	m := map[string]interface{}{
		"foo":   "bar",
		"quack": 1,
	}
	test := DataFromMap(m)
	assert.NotEmpty(t, test)
	assert.Equal(t, "bar", test.Get("foo"))
	assert.Equal(t, 1, test["quack"])
	assert.Equal(t, "", test.Get("bar"))
	s := test.String()
	assert.Equal(t, testString, s)
}

func TestMap_GormDBDataType(t *testing.T) {
	m := Map{}
	tests := []struct {
		d string
		t string
	}{
		{drivers.MySQL, "JSON"},
		{drivers.Postgres, "JSONB"},
		{drivers.SQLServer, "NVARCHAR(MAX)"},
		{drivers.SQLite, "JSON"},
		{drivers.SQLite3, "JSON"},
		{"fake", "JSON"},
	}
	for _, test := range tests {
		typ := m.GormDBDataType(dummyDB(test.d), nil)
		assert.Equal(t, test.t, typ)
	}
}

func dummyDB(name string) *gorm.DB {
	d := &dummy{name}
	c := &gorm.Config{Dialector: d}
	return &gorm.DB{Config: c}
}

type dummy struct {
	name string
}

func (m dummy) Name() string {
	return m.name
}

func (m dummy) Initialize(*gorm.DB) error {
	panic("implement me")
}

func (m dummy) Migrator(*gorm.DB) gorm.Migrator {
	panic("implement me")
}

func (m dummy) DataTypeOf(*schema.Field) string {
	panic("implement me")
}

func (m dummy) DefaultValueOf(*schema.Field) clause.Expression {
	panic("implement me")
}

func (m dummy) BindVarTo(clause.Writer, *gorm.Statement, interface{}) {
	panic("implement me")
}

func (m dummy) QuoteTo(clause.Writer, string) {
	panic("implement me")
}

func (m dummy) Explain(string, ...interface{}) string {
	panic("implement me")
}

var _ gorm.Dialector = (*dummy)(nil)
