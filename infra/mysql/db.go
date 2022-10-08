package mysql

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

type DB struct {
	db *sqlx.DB
}

func OpenMySQLDB(cfg *Config) (*DB, error) {
	dsn := cfg.ToDSN()

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect")
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return NewDB(db), nil
}

func NewDB(db *sqlx.DB) *DB {
	return &DB{db: db}
}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) Exec(query string) (sql.Result, error) {
	res, err := db.db.Exec(query)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to execute query: %s", query))
	}

	return res, nil
}

func (db *DB) ExplainWithAnalyze(query string) (string, error) {
	var res []string
	err := db.db.Select(&res, "EXPLAIN ANALYZE FORMAT=TREE "+query)
	if err != nil {
		return "", errors.Wrap(err, "failed to select")
	}

	if len(res) == 0 {
		return "", errors.New("no result for EXPLAIN ANALYZE")
	}

	return res[0], nil
}

func (db *DB) GetTableColumns(dbName string, tables []string) ([]*ColumnInfo, error) {
	query, args, err := sqlx.In(columnFetchQuery, dbName, tables)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build query to get table schema")
	}

	var cols []*ColumnInfo
	err = db.db.Select(&cols, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute query to get table schema")
	}

	return cols, nil
}

func (db *DB) GetIndexes(database string, tables []string) ([]*IndexInfo, error) {
	query, args, err := sqlx.In(indexFetchQuery, database, tables)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build query to get table schema")
	}

	var flatInfos []*flatIndexInfo
	err = db.db.Select(&flatInfos, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute query to get indexes")
	}

	return toIndexInfo(flatInfos), nil
}
