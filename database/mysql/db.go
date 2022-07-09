package mysql

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

func OpenMySQLDB() (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", "root@/gravityr")
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect")
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db, nil
}
