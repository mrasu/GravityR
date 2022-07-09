package thelper

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func MockDB(t *testing.T, f func(*sqlx.DB, sqlmock.Sqlmock)) {
	mockdb, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockdb.Close()
	db := sqlx.NewDb(mockdb, "mysql")
	f(db, mock)
	assert.NoError(t, mock.ExpectationsWereMet())
}
