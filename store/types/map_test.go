package types

import (
	"testing"

	"github.com/jrapoport/gothic/test/tdb"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func TestData(t *testing.T) {
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

func TestDataFromMap(t *testing.T) {
	m := map[string]interface{}{
		"foo":   "bar",
		"quack": 1,
	}
	test := DataFromMap(m)
	assert.NotEmpty(t, test)
	assert.Equal(t, "bar", test.Get("foo"))
	assert.Equal(t, 1, test["quack"])
	assert.Equal(t, "", test.Get("bar"))
}
