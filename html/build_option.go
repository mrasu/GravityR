package html

import (
	"github.com/mrasu/GravityR/html/viewmodel"
	"github.com/mrasu/GravityR/lib"
)

type suggestData struct {
	MySQL    *mysqlData    `json:"mysql"`
	Postgres *postgresData `json:"postgres"`
}

type mysqlData struct {
	Query              string                                  `json:"query"`
	AnalyzeNodes       []*viewmodel.VmMysqlExplainAnalyzeNode  `json:"analyzeNodes"`
	IndexTargets       []*viewmodel.VmIndexTarget              `json:"indexTargets"`
	CommandOptions     []*viewmodel.VmExaminationCommandOption `json:"examinationCommandOptions"`
	ExaminationResults *viewmodel.VmExaminationResult          `json:"examinationResult"`
}

type postgresData struct {
	Query              string                                    `json:"query"`
	AnalyzeNodes       []*viewmodel.VmPostgresExplainAnalyzeNode `json:"analyzeNodes"`
	PlanningText       string                                    `json:"planningText"`
	IndexTargets       []*viewmodel.VmIndexTarget                `json:"indexTargets"`
	CommandOptions     []*viewmodel.VmExaminationCommandOption   `json:"examinationCommandOptions"`
	ExaminationResults *viewmodel.VmExaminationResult            `json:"examinationResult"`
}

type digData struct {
	sqlDbLoads          []*viewmodel.VmTimeDbLoad
	tokenizedSqlDbLoads []*viewmodel.VmTimeDbLoad
}

type BuildOption struct {
	SuggestData *suggestData
	DigData     *digData
}

func NewSuggestMySQLDataBuildOption(query string,
	analyzeNodes []*viewmodel.VmMysqlExplainAnalyzeNode,
	indexTargets []*viewmodel.VmIndexTarget,
	commandOptions []*viewmodel.VmExaminationCommandOption,
	examinationResults *viewmodel.VmExaminationResult) *BuildOption {
	return &BuildOption{SuggestData: &suggestData{
		MySQL: &mysqlData{
			Query:              query,
			AnalyzeNodes:       analyzeNodes,
			IndexTargets:       indexTargets,
			CommandOptions:     commandOptions,
			ExaminationResults: examinationResults,
		},
	}}
}

func NewSuggestPostgresDataBuildOption(query string,
	analyzeNodes []*viewmodel.VmPostgresExplainAnalyzeNode,
	planningText string,
	indexTargets []*viewmodel.VmIndexTarget,
	commandOptions []*viewmodel.VmExaminationCommandOption,
	examinationResults *viewmodel.VmExaminationResult) *BuildOption {
	return &BuildOption{SuggestData: &suggestData{
		Postgres: &postgresData{
			Query:              query,
			AnalyzeNodes:       analyzeNodes,
			PlanningText:       planningText,
			IndexTargets:       indexTargets,
			CommandOptions:     commandOptions,
			ExaminationResults: examinationResults,
		},
	}}
}

func NewDigDataBuildOption(sqlCpuUsages []*viewmodel.VmTimeDbLoad, tokenizedSqlCpuUsages []*viewmodel.VmTimeDbLoad) *BuildOption {
	return &BuildOption{DigData: &digData{
		sqlDbLoads:          sqlCpuUsages,
		tokenizedSqlDbLoads: tokenizedSqlCpuUsages,
	}}
}

func (bo *BuildOption) createGrMap() map[string]interface{} {
	var digData map[string]interface{}
	if bo.DigData != nil {
		digData = map[string]interface{}{
			"sqlDbLoads":          lib.SliceOrEmpty(bo.DigData.sqlDbLoads),
			"tokenizedSqlDbLoads": lib.SliceOrEmpty(bo.DigData.tokenizedSqlDbLoads),
		}
	}

	return map[string]interface{}{
		"suggestData": bo.SuggestData,
		"digData":     digData,
	}
}
