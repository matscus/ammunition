package database

import (
	"database/sql"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

var (
	poolScheme = PoolScheme{
		Project: "testProject",
		Script:  "testScript",
	}
	testJSON = `{"test1": "test1val","test2": "test2val"}`
)

func newMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func Test_GetAllName(t *testing.T) {
	mockDB, mock := newMock()
	defer mockDB.Close()
	DB = sqlx.NewDb(mockDB, "sqlmock")
	rows := sqlmock.NewRows([]string{"project", "script"}).
		AddRow(poolScheme.Project, poolScheme.Script)
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	u, err := GetAllName()
	assert.NotNil(t, u)
	assert.NoError(t, err)
}

func Test_CreateScheme(t *testing.T) {
	mockDB, mock := newMock()
	defer mockDB.Close()
	DB = sqlx.NewDb(mockDB, "sqlmock")
	mock.ExpectExec("CREATE SCHEMA").WillReturnResult(sqlmock.NewResult(1, 1))
	err := poolScheme.CreateScheme()
	assert.NoError(t, err)
}

func Test_CreateTable(t *testing.T) {
	mockDB, mock := newMock()
	defer mockDB.Close()
	DB = sqlx.NewDb(mockDB, "sqlmock")
	mock.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(1, 1))
	err := poolScheme.CreateTable()
	assert.NoError(t, err)
}

func Test_DropTable(t *testing.T) {
	mockDB, mock := newMock()
	defer mockDB.Close()
	DB = sqlx.NewDb(mockDB, "sqlmock")
	mock.ExpectExec("DROP TABLE").WillReturnResult(sqlmock.NewResult(1, 1))
	err := poolScheme.DropTable()
	assert.NoError(t, err)
}

func Test_AddRelationsSchemeScript(t *testing.T) {
	mockDB, mock := newMock()
	defer mockDB.Close()
	DB = sqlx.NewDb(mockDB, "sqlmock")
	mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))
	err := poolScheme.AddRelationsSchemeScript()
	assert.NoError(t, err)
}

func Test_DeleteRelationsSchemeScript(t *testing.T) {
	mockDB, mock := newMock()
	defer mockDB.Close()
	DB = sqlx.NewDb(mockDB, "sqlmock")
	mock.ExpectExec("DELETE").WillReturnResult(sqlmock.NewResult(1, 1))
	err := poolScheme.DeleteRelationsSchemeScript()
	assert.NoError(t, err)
}

func Test_InsertSingleValuePool(t *testing.T) {
	mockDB, mock := newMock()
	defer mockDB.Close()
	DB = sqlx.NewDb(mockDB, "sqlmock")
	mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))
	result := poolScheme.InsertSingleValuePool(testJSON)
	assert.NotNil(t, result)
}

func Test_InsertMultiValues(t *testing.T) {
	mockDB, mock := newMock()
	defer mockDB.Close()
	DB = sqlx.NewDb(mockDB, "sqlmock")
	mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))
	err := poolScheme.InsertMultiValues([]string{testJSON, testJSON})
	assert.NoError(t, err)
}

func Test_GetPool(t *testing.T) {
	mockDB, mock := newMock()
	defer mockDB.Close()
	DB = sqlx.NewDb(mockDB, "sqlmock")
	rows := sqlmock.NewRows([]string{"pool"}).
		AddRow(testJSON)
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	u, err := poolScheme.GetPool()
	assert.NotNil(t, u)
	assert.NoError(t, err)
}
