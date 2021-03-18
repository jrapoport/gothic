package store

import (
	"context"
	"database/sql/driver"
	"math"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/store/drivers"
	"github.com/jrapoport/gothic/store/migration"
	"github.com/jrapoport/gothic/store/types"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ModelA is a model for tests.
type ModelA struct {
	gorm.Model
	Value string
}

// ModelB is a model with indexes for tests.
type ModelB struct {
	gorm.Model
	Value string `gorm:"index:idx_value"`
}

// ModelBIndex is the name of the index for ModelB
const ModelBIndex = "idx_value"

func TestDial(t *testing.T) {
	t.Parallel()
	// bad url
	url := &config.Config{}
	url.DB.DSN = "\n"
	// bad host
	host := &config.Config{}
	host.DB.Driver = drivers.MySQL
	host.DB.DSN = "root@tcp(255.255.255.255:3306)/test" //"255.255.255.255"
	// bad retry
	retry := &config.Config{}
	retry.DB.Driver = drivers.MySQL
	retry.DB.DSN = "root@tcp(255.255.255.255:3306)/test"
	retry.DB.MaxRetries = 1
	// good mock
	mctx, mocked := tconf.MockDB(t)
	mock := tconf.MockFromContext(mctx)
	create := "CREATE DATABASE IF NOT EXISTS " // + mocked.DB.Name
	mock.ExpectExec(create).WillReturnResult(sqlmock.NewResult(1, 1))
	use := "USE " // + mocked.DB.Name
	mock.ExpectExec(use).WillReturnResult(driver.ResultNoRows)
	type dialTest struct {
		ctx context.Context
		c   *config.Config
		l   logrus.FieldLogger
		Err assert.ErrorAssertionFunc
	}
	log := logrus.New()
	tests := []dialTest{
		{nil, nil, nil, assert.Error},
		{nil, nil, log, assert.Error},
		{nil, url, nil, assert.Error},
		{nil, url, log, assert.Error},
		{nil, host, log, assert.Error},
		{nil, retry, log, assert.Error},
		{mctx, mocked, log, assert.NoError},
		{nil, tconf.TempDB(t), log, assert.NoError},
	}
	for _, test := range tests {
		if test.c != nil {
			t.Log("Dialing: ", test.c.DB.Driver)
		}
		_, err := NewConnection(test.ctx, test.c, nil)
		test.Err(t, err, test.c)
	}
}

type ConnectionTestSuite struct {
	suite.Suite
	driver   drivers.Driver
	conn     *Connection
	c        *config.Config
	migrated bool
}

func TestConnection(t *testing.T) {
	t.Parallel()
	dvrs := []drivers.Driver{
		tconf.MySQLTemp,
		tconf.PostgresTemp,
		tconf.SQLiteTemp,
	}
	for _, d := range dvrs {
		ts := &ConnectionTestSuite{
			driver: d,
		}
		t.Run(string(d), func(t *testing.T) {
			t.Parallel()
			suite.Run(t, ts)
		})
	}
}

func (ts *ConnectionTestSuite) SetupSuite() {
	ts.conn, ts.c = testConn(ts.T(), ts.driver)
}

func testConn(t *testing.T, d drivers.Driver) (*Connection, *config.Config) {
	c := tconf.DBConfig(t, d)
	conn, err := Dial(c, nil)
	require.NoError(t, err)
	require.NotNil(t, conn)
	err = conn.DropAll()
	require.NoError(t, err)
	return conn, c
}

func (ts *ConnectionTestSuite) BeforeTest(_, testName string) {
	ts.Require().NotEmpty(ts.conn)
	if ts.migrated {
		return
	}
	if testName == "Test_0_AutoMigration" {
		return
	}
	ts.Test_0_AutoMigration()
}

func (ts *ConnectionTestSuite) TestRunMigration() {
	type ModelX struct {
		gorm.Model
	}
	mx := ModelX{}
	m := migration.NewMigration(time.Now().UTC().String(), mx)
	err := ts.conn.RunMigration(m)
	ts.NoError(err)
	has := ts.conn.DB.Migrator().HasTable(mx)
	ts.True(has)
	err = ts.conn.DB.Migrator().DropTable(mx)
	ts.NoError(err)
	has = ts.conn.DB.Migrator().HasTable(mx)
	ts.False(has)
}

func (ts *ConnectionTestSuite) TestData() {
	type UserWithData struct {
		gorm.Model
		Name       string
		Attributes types.Map
	}
	err := ts.conn.Migrator().AutoMigrate(&UserWithData{})
	ts.NoError(err)
	user1Attrs := types.Map{
		"age":  18,
		"name": "data-1",
		"orgs": types.Map{"orga": "orga"},
		"tags": []string{"tag1", "tag2"},
	}
	user2Attrs := types.Map{
		"name": "data-2",
		"age":  28,
		"tags": []string{"tag1", "tag3"},
		"role": "admin",
		"orgs": types.Map{"orgb": "orgb"},
	}
	users := []UserWithData{
		{
			Name:       "data-1",
			Attributes: user1Attrs,
		}, {
			Name:       "data-2",
			Attributes: user2Attrs,
		}}
	err = ts.conn.Create(&users).Error
	ts.NoError(err)
	var result UserWithData
	q := datatypes.JSONQuery("attributes").HasKey("role")
	err = ts.conn.First(&result, q).Error
	ts.NoError(err)
	ts.Equal(result.Name, users[1].Name)
	var result2 UserWithData
	q = datatypes.JSONQuery("attributes").HasKey("orgs", "orga")
	err = ts.conn.First(&result2, q).Error
	ts.NoError(err)
	ts.Equal(result2.Name, users[0].Name)
	var result3 UserWithData
	q = datatypes.JSONQuery("attributes").Equals("data-1", "name")
	err = ts.conn.First(&result3, q).Error
	ts.NoError(err)
	ts.Equal(result3.Name, users[0].Name)
	var result4 UserWithData
	q = datatypes.JSONQuery("attributes").Equals("orgb", "orgs", "orgb")
	err = ts.conn.First(&result4, q).Error
	ts.NoError(err)
	ts.Equal(result4.Name, users[1].Name)
	user1Attrs["age"] = 19
	b, err := user1Attrs.JSON()
	ts.NoError(err)
	jsonMap := map[string]interface{}{"Attributes": string(b)}
	err = ts.conn.Where(&UserWithData{Name: "data-1"}).Assign(jsonMap).FirstOrCreate(&UserWithData{}).Error
	ts.NoError(err)
	var result5 UserWithData
	q = datatypes.JSONQuery("attributes").Equals(19, "age")
	err = ts.conn.First(&result5, q).Error
	ts.NoError(err)
	err = ts.conn.Migrator().DropTable(UserWithData{})
	ts.NoError(err)
}

func (ts *ConnectionTestSuite) Test_0_AutoMigration() {
	ts.migrated = true
	const tableIdx = "idx_test_model_bs_value"
	var indexes = []string{ModelBIndex}
	ts.T().Cleanup(func() {
		plan.Clear()
	})
	ma := ModelA{}
	mb := ModelB{}
	AddAutoMigration(time.Now().UTC().String(), ma)
	AddAutoMigrationWithIndexes(time.Now().UTC().String(), mb, indexes)
	conn := ts.conn
	err := conn.AutoMigrate()
	ts.Require().NoError(err)
	has := conn.DB.Migrator().HasTable(ma)
	ts.True(has)
	has = conn.DB.Migrator().HasTable(mb)
	ts.True(has)
	has = conn.DB.Migrator().HasIndex(mb, tableIdx)
	ts.True(has)
	err = conn.AutoMigrate()
	ts.Require().NoError(err)
}

func (ts *ConnectionTestSuite) Test_1_Has() {
	const testValue = "test-value"
	conn := ts.conn
	var has bool
	ma := &ModelA{Value: testValue}
	err := conn.Create(ma).Error
	ts.NoError(err)
	has, err = conn.Has(ma)
	ts.NoError(err)
	ts.True(has)
	has, err = conn.Has(ma, "Value = ?", testValue)
	ts.NoError(err)
	ts.True(has)
	err = conn.Delete(ma).Error
	ts.NoError(err)
	has, err = conn.Has(ma)
	ts.NoError(err)
	ts.False(has)
	// false
	ma.ID = math.MaxInt8 - 1
	has, err = conn.Has(ma)
	ts.NoError(err)
	ts.False(has)
	// false
	mb := &ModelB{}
	has, err = conn.Has(mb)
	ts.NoError(err)
	ts.False(has)
	// error
	has, err = conn.Has(nil)
	ts.Error(err)
	ts.False(has)
}

func (ts *ConnectionTestSuite) Test_2_TableNames() {
	conn := ts.conn
	names, err := conn.TableNames()
	ts.NoError(err)
	ts.ElementsMatch([]string{
		"test_migrations",
		"test_model_as",
		"test_model_bs",
	}, names)
}

func (ts *ConnectionTestSuite) Test_3_TruncateAll() {
	conn := ts.conn
	models := []interface{}{
		&ModelA{},
		&ModelB{},
	}
	for _, model := range models {
		err := conn.Create(model).Error
		ts.NoError(err)
		var cma int64
		err = conn.Model(model).Count(&cma).Error
		ts.NoError(err)
		ts.Equal(int64(1), cma)
	}
	err := conn.TruncateAll()
	ts.NoError(err)
	for _, model := range models {
		var cma int64
		err = conn.Model(model).Count(&cma).Error
		ts.NoError(err)
		ts.Equal(int64(0), cma)
	}
}

func (ts *ConnectionTestSuite) Test_4_DropAll() {
	conn := ts.conn
	err := conn.DropAll()
	ts.NoError(err)
	has := conn.Migrator().HasTable(ModelA{})
	ts.Require().False(has)
	has = conn.Migrator().HasTable(ModelB{})
	ts.Require().False(has)
}

/*
func (ts *ConnectionTestSuite) Test_5_DropDatabase() {
	conn := ts.conn
	name := conn.Migrator().CurrentDatabase()
	ts.NotEmpty(name)
	err := conn.DropDatabase()
	ts.NoError(err)
	names, err := conn.TableNames()
	ts.NoError(err)
	ts.Empty(names)
	if conn.Name() != conf.SQLite {
		name = conn.Migrator().CurrentDatabase()
		ts.Empty(name)
	}
}
*/
