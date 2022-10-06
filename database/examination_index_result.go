package database

import (
	"github.com/mrasu/GravityR/html/viewmodel"
)

type ExaminationIndexResult struct {
	IndexTarget         *IndexTarget
	ExecutionTimeMillis int64
}

func NewExaminationIndexResult(it *IndexTarget, time int64) *ExaminationIndexResult {
	return &ExaminationIndexResult{
		IndexTarget:         it,
		ExecutionTimeMillis: time,
	}
}

func (eir *ExaminationIndexResult) ToViewModel() *viewmodel.VmExaminationIndexResult {
	return &viewmodel.VmExaminationIndexResult{
		ExecutionTimeMillis: eir.ExecutionTimeMillis,
		IndexTarget:         eir.IndexTarget.ToViewModel(),
	}
}
