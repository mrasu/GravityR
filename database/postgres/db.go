package postgres

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func OpenPostgresDB(cfg *Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", cfg.GetKVConnectionString())
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect")
	}

	return db, nil
}
