package thelper

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/mrasu/GravityR/infra/mysql"
	"github.com/mrasu/GravityR/infra/postgres"
	"github.com/stretchr/testify/assert"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func MockMysqlDB(t *testing.T, f func(*mysql.DB, sqlmock.Sqlmock)) {
	mockDB(t, "mysql", func(db *sqlx.DB, mock sqlmock.Sqlmock) {
		f(mysql.NewDB(db), mock)
	})
}

func MockPostgresDB(t *testing.T, f func(*postgres.DB, sqlmock.Sqlmock)) {
	mockDB(t, "postgres", func(db *sqlx.DB, mock sqlmock.Sqlmock) {
		f(postgres.NewDB(db), mock)
	})
}

func mockDB(t *testing.T, driver string, f func(*sqlx.DB, sqlmock.Sqlmock)) {
	t.Helper()
	mockdb, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockdb.Close()
	db := sqlx.NewDb(mockdb, driver)
	f(db, mock)
	assert.NoError(t, mock.ExpectationsWereMet())
}
