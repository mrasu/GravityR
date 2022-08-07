package html

import (
	"github.com/mrasu/GravityR/html/viewmodel"
	"github.com/mrasu/GravityR/lib"
)

type suggestData struct {
	Query              string
	AnalyzeNodes       []*viewmodel.VmExplainAnalyzeNode
	IndexTargets       []*viewmodel.VmIndexTarget
	CommandOptions     []*viewmodel.VmExaminationCommandOption
	ExaminationResults *viewmodel.VmExaminationResult
}

type digData struct {
	sqlDbLoads          []*viewmodel.VmTimeDbLoad
	tokenizedSqlDbLoads []*viewmodel.VmTimeDbLoad
}

type BuildOption struct {
	SuggestData *suggestData
	DigData     *digData
}

func NewSuggestDataBuildOption(query string,
	analyzeNodes []*viewmodel.VmExplainAnalyzeNode,
	indexTargets []*viewmodel.VmIndexTarget,
	commandOptions []*viewmodel.VmExaminationCommandOption,
	examinationResults *viewmodel.VmExaminationResult) *BuildOption {
	return &BuildOption{SuggestData: &suggestData{
		Query:              query,
		AnalyzeNodes:       analyzeNodes,
		IndexTargets:       indexTargets,
		CommandOptions:     commandOptions,
		ExaminationResults: examinationResults,
	}}
}

func NewDigDataBuildOption(sqlCpuUsages []*viewmodel.VmTimeDbLoad, tokenizedSqlCpuUsages []*viewmodel.VmTimeDbLoad) *BuildOption {
	return &BuildOption{DigData: &digData{
		sqlDbLoads:          sqlCpuUsages,
		tokenizedSqlDbLoads: tokenizedSqlCpuUsages,
	}}
}

func (bo *BuildOption) createGrMap() map[string]interface{} {
	var suggestData map[string]interface{}
	if bo.SuggestData != nil {
		suggestData = map[string]interface{}{
			"query":                     bo.SuggestData.Query,
			"analyzeNodes":              bo.SuggestData.AnalyzeNodes,
			"indexTargets":              bo.SuggestData.IndexTargets,
			"examinationCommandOptions": bo.SuggestData.CommandOptions,
			"examinationResult":         bo.SuggestData.ExaminationResults,
		}
	}

	var digData map[string]interface{}
	if bo.DigData != nil {
		digData = map[string]interface{}{
			"sqlDbLoads":          lib.SliceOrEmpty(bo.DigData.sqlDbLoads),
			"tokenizedSqlDbLoads": lib.SliceOrEmpty(bo.DigData.tokenizedSqlDbLoads),
		}
	}

	return map[string]interface{}{
		"suggestData": suggestData,
		"digData":     digData,
	}
}
