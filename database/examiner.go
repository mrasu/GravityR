package database

import (
	"fmt"
	"github.com/mrasu/GravityR/database/common_model"
	"github.com/mrasu/GravityR/lib"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
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

	log.Debug().Msg("====RESULT====")
	log.Printf("Original Time: %dms", origExecTimeMillis)

	for tName, rs := range results {
		log.Debug().Msg("==Table: " + tName)
		lib.Sort(rs, func(r *common_model.ExaminationIndexResult) int64 { return r.ExecutionTimeMillis })
		for _, r := range rs {
			log.Printf("Time: %dms (%.1f%% reduced), Columns(%s)",
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

	log.Info().Msg(fmt.Sprintf("Creating index(%s)...", idxName))
	err = iee.ie.CreateIndex(idxName, it)
	if err != nil {
		return nil, err
	}

	defer func() {
		log.Info().Msg(fmt.Sprintf("Dropping index(%s)...", idxName))
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
	log.Info().Msg("Executing query...")
	return iee.ie.Execute()
}
