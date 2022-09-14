package postgres

import (
	"fmt"
	"github.com/mrasu/GravityR/database/db_models"
	"github.com/mrasu/GravityR/infra/postgres"
	"github.com/mrasu/GravityR/lib"
	"time"
)

type IndexExaminer struct {
	db    *postgres.DB
	query string
}

func NewIndexExaminer(db *postgres.DB, query string) *IndexExaminer {
	return &IndexExaminer{
		db:    db,
		query: query,
	}
}

func (ie *IndexExaminer) Execute() (int64, error) {
	start := time.Now()
	_, err := ie.db.Exec(ie.query)
	if err != nil {
		return 0, err
	}
	elapsed := time.Since(start)

	return elapsed.Milliseconds(), nil
}

func (ie *IndexExaminer) CreateIndex(name string, it *db_models.IndexTarget) error {
	sql := fmt.Sprintf(`CREATE INDEX "%s" ON "%s" (%s)`,
		name, it.TableName,
		lib.JoinF(it.Columns, ",", func(i *db_models.IndexColumn) string { return `"` + i.SafeName() + `"` }),
	)
	_, err := ie.db.Exec(sql)
	return err
}

func (ie *IndexExaminer) DropIndex(name string, it *db_models.IndexTarget) error {
	sql := fmt.Sprintf(`DROP INDEX "%s"`, name)
	_, err := ie.db.Exec(sql)
	return err
}
