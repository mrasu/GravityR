package viewmodel

type VmExaminationIndexResult struct {
	ExecutionTimeMillis int64          `json:"executionTimeMillis"`
	IndexTarget         *VmIndexTarget `json:"indexTarget"`
}
