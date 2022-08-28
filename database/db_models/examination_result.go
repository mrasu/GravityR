package db_models

import "github.com/mrasu/GravityR/html/viewmodel"

type ExaminationResult struct {
	OriginalTimeMillis int64
	IndexResults       []*ExaminationIndexResult
}

func NewExaminationResult(irs []*ExaminationIndexResult, time int64) *ExaminationResult {
	return &ExaminationResult{
		OriginalTimeMillis: time,
		IndexResults:       irs,
	}
}

func (er *ExaminationResult) ToViewModel() *viewmodel.VmExaminationResult {
	var virs []*viewmodel.VmExaminationIndexResult
	for _, ir := range er.IndexResults {
		virs = append(virs, ir.ToViewModel())
	}

	return &viewmodel.VmExaminationResult{
		OriginalTimeMillis: er.OriginalTimeMillis,
		IndexResults:       virs,
	}
}
