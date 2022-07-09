package html

import "github.com/mrasu/GravityR/html/viewmodel"

type BuildOption struct {
	Query              string
	AnalyzeNodes       []*viewmodel.VmExplainAnalyzeNode
	IndexTargets       []*viewmodel.VmIndexTarget
	CommandOptions     []*viewmodel.VmExaminationCommandOption
	ExaminationResults *viewmodel.VmExaminationResult
}
