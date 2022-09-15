package database

import (
	"fmt"
	"github.com/mrasu/GravityR/database/common_model"
	"github.com/mrasu/GravityR/lib"
	"github.com/pkg/errors"
	"math/rand"
	"strconv"
	"time"
)

type IndexExaminer interface {
	Execute() (int64, error)
	CreateIndex(string, *common_model.IndexTarget) error
	DropIndex(string, *common_model.IndexTarget) error
}

type IndexEfficiencyExaminer struct {
	ie IndexExaminer
}

func NewIndexEfficiencyExaminer(ie IndexExaminer) *IndexEfficiencyExaminer {
	return &IndexEfficiencyExaminer{
		ie: ie,
	}
}

func (iee *IndexEfficiencyExaminer) Run(its []*common_model.IndexTarget) (*common_model.ExaminationResult, error) {
	origExecTimeMillis, err := iee.execute()
	if err != nil {
		return nil, err
	}

	var irs []*common_model.ExaminationIndexResult
	results := map[string][]*common_model.ExaminationIndexResult{}
	id := strconv.FormatInt(time.Now().Unix(), 10)
	for _, it := range its {
		res, err := iee.examine(id, it)
		if err != nil {
			return nil, err
		}

		irs = append(irs, res)
		if rs, ok := results[it.TableName]; ok {
			results[it.TableName] = append(rs, res)
		} else {
			results[it.TableName] = []*common_model.ExaminationIndexResult{res}
		}
	}
	er := common_model.NewExaminationResult(irs, origExecTimeMillis)

	fmt.Printf("\n\n====RESULT====\n")
	fmt.Printf("Original Time: %dms\n", origExecTimeMillis)

	for tName, rs := range results {
		fmt.Printf("==Table: %s\n", tName)
		lib.Sort(rs, func(r *common_model.ExaminationIndexResult) int64 { return r.ExecutionTimeMillis })
		for _, r := range rs {
			fmt.Printf("Time: %dms (%.1f%% reduced), Columns(%s)\n",
				r.ExecutionTimeMillis,
				(1-float64(r.ExecutionTimeMillis)/float64(origExecTimeMillis))*100,
				lib.Join(r.IndexTarget.Columns, ",", func(i *common_model.IndexColumn) string { return i.SafeName() }),
			)
		}
	}

	return er, nil
}

const maxIdxNameLen = 30

func (iee *IndexEfficiencyExaminer) examine(id string, it *common_model.IndexTarget) (res *common_model.ExaminationIndexResult, err error) {
	idxName := fmt.Sprintf("%s_%s",
		id,
		lib.Join(it.Columns, "_", func(i *common_model.IndexColumn) string { return i.SafeName() }),
	)
	if len(idxName) > maxIdxNameLen {
		idxName = fmt.Sprintf("%s_%d", idxName[:maxIdxNameLen], rand.Intn(999999999))
	}

	fmt.Printf("Creating index(%s)...\n", idxName)
	err = iee.ie.CreateIndex(idxName, it)
	if err != nil {
		return nil, err
	}

	defer func() {
		fmt.Printf("Dropping index(%s)...\n", idxName)
		err2 := iee.ie.DropIndex(idxName, it)
		if err == nil {
			err = errors.Wrap(err2, "failed to drop index")
		}
	}()

	millis, err := iee.execute()
	if err != nil {
		return nil, err
	}

	return common_model.NewExaminationIndexResult(it, millis), nil
}

func (iee *IndexEfficiencyExaminer) execute() (int64, error) {
	fmt.Println("Executing query...")
	return iee.ie.Execute()
}
