package viewmodel

type VmExaminationResult struct {
	OriginalTimeMillis int64                       `json:"originalTimeMillis"`
	IndexResults       []*VmExaminationIndexResult `json:"indexResults"`
}
