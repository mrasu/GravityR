package mysql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

func OpenMySQLDB(cfg *Config) (*sqlx.DB, error) {
	dsn := cfg.ToDSN()
	fmt.Printf("DSN: %s\n", dsn)

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect")
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db, nil
}
