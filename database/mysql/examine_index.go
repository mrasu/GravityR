package mysql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mrasu/GravityR/database/mysql/models"
	"github.com/mrasu/GravityR/lib"
	"github.com/pkg/errors"
	"math/rand"
	"strconv"
	"time"
)

type indexExaminer struct{}

func ExamineIndex(db *sqlx.DB, query string, its []*models.IndexTarget) (*models.ExaminationResult, error) {
	ie := &indexExaminer{}

	origExecTimeMillis, err := examineExecTime(db, query)
	if err != nil {
		return nil, err
	}

	var irs []*models.ExaminationIndexResult
	results := map[string][]*models.ExaminationIndexResult{}
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
			results[it.TableName] = []*models.ExaminationIndexResult{res}
		}
	}
	er := models.NewExaminationResult(irs, origExecTimeMillis)

	fmt.Printf("\n\n====RESULT====\n")
	fmt.Printf("Original Time: %dms\n", origExecTimeMillis)

	for tName, rs := range results {
		fmt.Printf("==Table: %s\n", tName)
		lib.SortF(rs, func(r *models.ExaminationIndexResult) int64 { return r.ExecutionTimeMillis })
		for _, r := range rs {
			fmt.Printf("Time: %dms (%.1f%% reduced), Columns(%s)\n",
				r.ExecutionTimeMillis,
				(1-float64(r.ExecutionTimeMillis)/float64(origExecTimeMillis))*100,
				lib.JoinF(r.IndexTarget.Columns, ",", func(i *models.IndexColumn) string { return i.Name }),
			)
		}
	}

	return er, nil
}

const maxIdxNameLen = 30

func (ie *indexExaminer) examineIndex(id string, db *sqlx.DB, query string, it *models.IndexTarget) (res *models.ExaminationIndexResult, err error) {
	idxName := fmt.Sprintf("%s_%s",
		id,
		lib.JoinF(it.Columns, "_", func(i *models.IndexColumn) string { return i.Name }),
	)
	if len(idxName) > maxIdxNameLen {
		idxName = fmt.Sprintf("%s_%d", idxName[:maxIdxNameLen], rand.Intn(999999999))
	}

	fmt.Printf("Creating index(%s)...\n", idxName)
	addQuery := fmt.Sprintf(
		"ALTER TABLE %s ADD INDEX %s (%s)",
		it.TableName, idxName,
		lib.JoinF(it.Columns, ",", func(i *models.IndexColumn) string { return i.Name }),
	)
	dropQuery := fmt.Sprintf("ALTER TABLE %s DROP INDEX %s", it.TableName, idxName)

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

	millis, err := examineExecTime(db, query)
	if err != nil {
		return nil, err
	}

	return models.NewExaminationIndexResult(it, millis), nil
}

func examineExecTime(db *sqlx.DB, query string) (int64, error) {
	fmt.Println("Executing query...")
	start := time.Now()
	_, err := db.Exec(query)
	if err != nil {
		return 0, err
	}
	elapsed := time.Since(start)

	return elapsed.Milliseconds(), nil
}
