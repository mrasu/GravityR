package rdb

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mrasu/GravityR/database/db_models"
	"github.com/mrasu/GravityR/lib"
	"github.com/pkg/errors"
	"math/rand"
	"strconv"
	"time"
)

type dbType int

const (
	MySQL dbType = iota + 1
	PostgreSQL
)

type IndexExaminer struct {
	DbType dbType
}

func (ie *IndexExaminer) Run(db *sqlx.DB, query string, its []*db_models.IndexTarget) (*db_models.ExaminationResult, error) {
	origExecTimeMillis, err := ie.examineExecTime(db, query)
	if err != nil {
		return nil, err
	}

	var irs []*db_models.ExaminationIndexResult
	results := map[string][]*db_models.ExaminationIndexResult{}
	id := strconv.FormatInt(time.Now().Unix(), 10)
	for _, it := range its {
		res, err := ie.examineIndex(id, db, query, it)
		if err != nil {
			return nil, err
		}

		irs = append(irs, res)
		if rs, ok := results[it.TableName]; ok {
			results[it.TableName] = append(rs, res)
		} else {
			results[it.TableName] = []*db_models.ExaminationIndexResult{res}
		}
	}
	er := db_models.NewExaminationResult(irs, origExecTimeMillis)

	fmt.Printf("\n\n====RESULT====\n")
	fmt.Printf("Original Time: %dms\n", origExecTimeMillis)

	for tName, rs := range results {
		fmt.Printf("==Table: %s\n", tName)
		lib.SortF(rs, func(r *db_models.ExaminationIndexResult) int64 { return r.ExecutionTimeMillis })
		for _, r := range rs {
			fmt.Printf("Time: %dms (%.1f%% reduced), Columns(%s)\n",
				r.ExecutionTimeMillis,
				(1-float64(r.ExecutionTimeMillis)/float64(origExecTimeMillis))*100,
				lib.JoinF(r.IndexTarget.Columns, ",", func(i *db_models.IndexColumn) string { return i.SafeName() }),
			)
		}
	}

	return er, nil
}

const maxIdxNameLen = 30

func (ie *IndexExaminer) examineIndex(id string, db *sqlx.DB, query string, it *db_models.IndexTarget) (res *db_models.ExaminationIndexResult, err error) {
	idxName := fmt.Sprintf("%s_%s",
		id,
		lib.JoinF(it.Columns, "_", func(i *db_models.IndexColumn) string { return i.SafeName() }),
	)
	if len(idxName) > maxIdxNameLen {
		idxName = fmt.Sprintf("%s_%d", idxName[:maxIdxNameLen], rand.Intn(999999999))
	}

	fmt.Printf("Creating index(%s)...\n", idxName)
	addQuery := ie.createAddIndexQuery(idxName, it)
	dropQuery := ie.createDropIndexQuery(idxName, it)

	_, err = db.Exec(addQuery)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add index")
	}
	defer func() {
		fmt.Printf("Dropping index(%s)...\n", idxName)
		_, err2 := db.Exec(dropQuery)
		if err == nil {
			err = errors.Wrap(err2, "failed to drop index")
		}
	}()

	millis, err := ie.examineExecTime(db, query)
	if err != nil {
		return nil, err
	}

	return db_models.NewExaminationIndexResult(it, millis), nil
}

func (ie *IndexExaminer) examineExecTime(db *sqlx.DB, query string) (int64, error) {
	fmt.Println("Executing query...")
	start := time.Now()
	_, err := db.Exec(query)
	if err != nil {
		return 0, err
	}
	elapsed := time.Since(start)

	return elapsed.Milliseconds(), nil
}

func (ie *IndexExaminer) createAddIndexQuery(name string, it *db_models.IndexTarget) string {
	switch ie.DbType {
	case MySQL:
		return fmt.Sprintf(
			"ALTER TABLE `%s` ADD INDEX `%s` (%s)",
			it.TableName, name,
			lib.JoinF(it.Columns, ",", func(i *db_models.IndexColumn) string { return "`" + i.SafeName() + "`" }),
		)
	case PostgreSQL:
		return fmt.Sprintf(
			`CREATE INDEX "%s" ON "%s" (%s)`,
			name, it.TableName,
			lib.JoinF(it.Columns, ",", func(i *db_models.IndexColumn) string { return `"` + i.SafeName() + `"` }),
		)
	default:
		panic("Invalid DbType exists")
	}
}

func (ie *IndexExaminer) createDropIndexQuery(name string, it *db_models.IndexTarget) string {
	switch ie.DbType {
	case MySQL:
		return fmt.Sprintf("ALTER TABLE `%s` DROP INDEX `%s`", it.TableName, name)
	case PostgreSQL:
		return fmt.Sprintf(`DROP INDEX "%s"`, name)
	default:
		panic("Invalid DbType exists")
	}
}
