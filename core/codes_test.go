package core

import (
	"errors"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/models/code"
	"github.com/stretchr/testify/assert"
	"gopkg.in/data-dog/go-sqlmock.v2"
)

func TestAPI_CreateCode(t *testing.T) {
	t.Parallel()
	a := createAPI(t)
	ctx := testContext(a)
	c, err := a.CreateSignupCode(ctx, code.SingleUse)
	assert.NoError(t, err)
	assert.NotEmpty(t, c)
	a.conn.Error = errors.New("force error")
	_, err = a.CreateSignupCode(ctx, code.SingleUse)
	assert.Error(t, err)
	a.conn.Error = nil
}

func TestAPI_CreateCode_Error(t *testing.T) {
	t.Parallel()
	a, mock := mockAPI(t)
	mock.ExpectBegin()
	mock.ExpectBegin()
	has := "SELECT * FROM `test_signup_codes` WHERE token = ? AND " +
		"`test_signup_codes`.`deleted_at` IS NULL ORDER BY " +
		"`test_signup_codes`.`id` LIMIT 1"
	mock.ExpectQuery(regexp.QuoteMeta(has)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnError(errors.New("mock error"))
	mock.ExpectRollback()
	mock.ExpectRollback()
	rtx := testContext(a)
	_, err := a.CreateSignupCode(rtx, code.SingleUse)
	assert.Error(t, err)
}

func TestAPI_CreateCodes(t *testing.T) {
	t.Parallel()
	const count = 100
	a := createAPI(t)
	ctx := testContext(a)
	list, err := a.CreateSignupCodes(ctx, code.SingleUse, count)
	assert.Error(t, err)
	ctx.SetAdminID(uuid.New())
	list, err = a.CreateSignupCodes(ctx, code.SingleUse, count)
	assert.NoError(t, err)
	assert.Len(t, list, count)
	assertUnique(t, list)
	a.conn.Error = errors.New("force error")
	_, err = a.CreateSignupCodes(ctx, code.SingleUse, count)
	assert.Error(t, err)
	a.conn.Error = nil
}

func TestAPI_CreateCodes_Error(t *testing.T) {
	t.Parallel()
	a, mock := mockAPI(t)
	mock.ExpectBegin()
	mock.ExpectBegin()
	has := "SELECT * FROM `test_signup_codes` WHERE token = ? AND " +
		"`test_signup_codes`.`deleted_at` IS NULL ORDER BY " +
		"`test_signup_codes`.`id` LIMIT 1"
	mock.ExpectQuery(regexp.QuoteMeta(has)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnError(errors.New("mock error"))
	mock.ExpectRollback()
	mock.ExpectRollback()
	rtx := testContext(a)
	_, err := a.CreateSignupCodes(rtx, code.SingleUse, 1)
	assert.Error(t, err)
}

func TestAPI_CheckSignupCode(t *testing.T) {
	t.Parallel()
	a := createAPI(t)
	ctx := testContext(a)
	test, err := a.CreateSignupCode(ctx, code.SingleUse)
	assert.NoError(t, err)
	assert.NotEmpty(t, test)
	sc, err := a.CheckSignupCode(test)
	assert.NoError(t, err)
	assert.Equal(t, test, sc.Token)
	err = a.DeleteSignupCode(test)
	assert.NoError(t, err)
	_, err = a.CheckSignupCode("")
	assert.Error(t, err)
	a.conn.Error = errors.New("force error")
	_, err = a.CheckSignupCode(test)
	assert.Error(t, err)
	a.conn.Error = nil
}

func TestAPI_CheckSignupCode_Error(t *testing.T) {
	t.Parallel()
	a, mock := mockAPI(t)
	mock.ExpectBegin()
	mock.ExpectBegin()
	has := "SELECT * FROM `test_signup_codes` WHERE token = ? AND " +
		"`test_signup_codes`.`deleted_at` IS NULL ORDER BY " +
		"`test_signup_codes`.`id` LIMIT 1"
	mock.ExpectQuery(regexp.QuoteMeta(has)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnError(errors.New("mock error"))
	mock.ExpectRollback()
	mock.ExpectRollback()
	_, err := a.CheckSignupCode("1234")
	assert.Error(t, err)
}

func assertUnique(t *testing.T, s []string) {
	unique := make(map[string]bool, len(s))
	us := make([]string, len(unique))
	for _, elem := range s {
		if len(elem) != 0 {
			if !unique[elem] {
				us = append(us, elem)
				unique[elem] = true
			}
		}
	}
	assert.Exactly(t, s, us)
}

/*
func TestNewMigrateWithIndexes_Error(t *testing.T) {
	t.Parallel()
	db, mock := tdb.MockDB(t)
	var indexes = []string{ModelBIndex}
	// create := "CREATE TABLE `model_as` (`id` bigint unsigned AUTO_INCREMENT," +
	// 	"`created_at` datetime(3) NULL,`updated_at` datetime(3) NULL," +
	// 	"`deleted_at` datetime(3) NULL,`value` longtext,PRIMARY KEY (`id`)," +
	// 	"INDEX idx_model_as_deleted_at (`deleted_at`))"
	mock.ExpectExec(".*").
		WillReturnResult(sqlmock.NewResult(0,0))
	err := migrateIndexes(db, nil, indexes)
	assert.Error(t, err)
}
*/
