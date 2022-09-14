package rdb

import "database/sql"

type DB interface {
	Exec(query string) (sql.Result, error)
}
