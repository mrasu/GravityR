package postgres

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type DB struct {
	db *sqlx.DB
}

func OpenPostgresDB(cfg *Config) (*DB, error) {
	db, err := sqlx.Open("postgres", cfg.GetKVConnectionString())
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect")
	}

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

func (db *DB) ExplainWithAnalyze(query string) ([]string, error) {
	rows, err := db.db.Query("EXPLAIN (ANALYZE, BUFFERS) " + query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select")
	}
	defer rows.Close()

	var res []string
	for rows.Next() {
		var txt string
		if err := rows.Scan(&txt); err != nil {
			return nil, errors.Wrap(err, "failed to Scan")
		}
		res = append(res, txt)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to execute EXPLAIN ANALYZE")
	}

	return res, nil
}

func (db *DB) GetTableColumns(schema string, tables []string) ([]*ColumnInfo, error) {
	query, args, err := sqlx.In(columnFetchQuery, schema, tables)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build query to get table schema")
	}
	query = db.db.Rebind(query)

	var cols []*ColumnInfo
	err = db.db.Select(&cols, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute query to get table column")
	}

	return cols, nil
}

func (db *DB) GetIndexes(database string, tables []string) ([]*IndexInfo, error) {
	query, args, err := sqlx.In(indexFetchQuery, database, tables)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build query to get table schema")
	}
	query = db.db.Rebind(query)

	var qInfos []*queryIndexInfo
	err = db.db.Select(&qInfos, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute query to get indexes")
	}

	return toIndexInfo(qInfos), nil
}
