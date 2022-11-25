package html

import (
	"github.com/mrasu/GravityR/html/viewmodel"
)

type suggestData struct {
	MySQL    *mysqlData    `json:"mysql"`
	Postgres *postgresData `json:"postgres"`
	Hasura   *hasuraData   `json:"hasura"`
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
	SummaryText        string                                    `json:"summaryText"`
	IndexTargets       []*viewmodel.VmIndexTarget                `json:"indexTargets"`
	CommandOptions     []*viewmodel.VmExaminationCommandOption   `json:"examinationCommandOptions"`
	ExaminationResults *viewmodel.VmExaminationResult            `json:"examinationResult"`
}

type hasuraData struct {
	Postgres *hasuraPostgresData `json:"postgres"`
}

type hasuraPostgresData struct {
	Gql                string                                    `json:"gql"`
	GqlVariables       map[string]interface{}                    `json:"gqlVariables"`
	Query              string                                    `json:"query"`
	AnalyzeNodes       []*viewmodel.VmPostgresExplainAnalyzeNode `json:"analyzeNodes"`
	SummaryText        string                                    `json:"summaryText"`
	IndexTargets       []*viewmodel.VmIndexTarget                `json:"indexTargets"`
	CommandOptions     []*viewmodel.VmExaminationCommandOption   `json:"examinationCommandOptions"`
	ExaminationResults *viewmodel.VmExaminationResult            `json:"examinationResult"`
}

type digData struct {
	PerformanceInsights *performanceInsightsData `json:"performanceInsights"`
	Jaeger              *jaegerData              `json:"jaeger"`
}

type performanceInsightsData struct {
	SqlDbLoads          []*viewmodel.VmTimeDbLoad `json:"sqlDbLoads"`
	TokenizedSqlDbLoads []*viewmodel.VmTimeDbLoad `json:"tokenizedSqlDbLoads"`
}

type jaegerData struct {
	UIPath               string                          `json:"uiPath"`
	SlowThresholdMilli   int                             `json:"slowThresholdMilli"`
	SameServiceThreshold int                             `json:"sameServiceThreshold"`
	SlowTraces           []*viewmodel.VmOtelCompactTrace `json:"slowTraces"`
	SameServiceTraces    []*viewmodel.VmOtelCompactTrace `json:"sameServiceTraces"`
}

type BuildOption struct {
	SuggestData *suggestData
	DigData     *digData
}

func NewSuggestMySQLDataBuildOption(
	query string,
	analyzeNodes []*viewmodel.VmMysqlExplainAnalyzeNode,
	indexTargets []*viewmodel.VmIndexTarget,
	commandOptions []*viewmodel.VmExaminationCommandOption,
	examinationResults *viewmodel.VmExaminationResult,
) *BuildOption {
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

func NewSuggestPostgresDataBuildOption(
	query string,
	analyzeNodes []*viewmodel.VmPostgresExplainAnalyzeNode,
	sumaryText string,
	indexTargets []*viewmodel.VmIndexTarget,
	commandOptions []*viewmodel.VmExaminationCommandOption,
	examinationResults *viewmodel.VmExaminationResult,
) *BuildOption {
	return &BuildOption{SuggestData: &suggestData{
		Postgres: &postgresData{
			Query:              query,
			AnalyzeNodes:       analyzeNodes,
			SummaryText:        sumaryText,
			IndexTargets:       indexTargets,
			CommandOptions:     commandOptions,
			ExaminationResults: examinationResults,
		},
	}}
}

func NewSuggestHasuraDataBuildOption(
	gql string,
	gqlVariables map[string]interface{},
	query string,
	analyzeNodes []*viewmodel.VmPostgresExplainAnalyzeNode,
	summaryText string,
	indexTargets []*viewmodel.VmIndexTarget,
	commandOptions []*viewmodel.VmExaminationCommandOption,
	examinationResults *viewmodel.VmExaminationResult,
) *BuildOption {
	return &BuildOption{SuggestData: &suggestData{
		Hasura: &hasuraData{
			Postgres: &hasuraPostgresData{
				Gql:                gql,
				GqlVariables:       gqlVariables,
				Query:              query,
				AnalyzeNodes:       analyzeNodes,
				SummaryText:        summaryText,
				IndexTargets:       indexTargets,
				CommandOptions:     commandOptions,
				ExaminationResults: examinationResults,
			},
		},
	}}
}

func NewPerformanceInsightsDataBuildOption(sqlCpuUsages []*viewmodel.VmTimeDbLoad, tokenizedSqlCpuUsages []*viewmodel.VmTimeDbLoad) *BuildOption {
	return &BuildOption{DigData: &digData{
		PerformanceInsights: &performanceInsightsData{
			SqlDbLoads:          sqlCpuUsages,
			TokenizedSqlDbLoads: tokenizedSqlCpuUsages,
		},
	}}
}

func NewDigJaegerBuildOption(
	uIPath string,
	slowThresholdMilli int,
	sameServiceThreshold int,
	slowTraces []*viewmodel.VmOtelCompactTrace,
	sameServiceTraces []*viewmodel.VmOtelCompactTrace,
) *BuildOption {
	return &BuildOption{DigData: &digData{
		Jaeger: &jaegerData{
			UIPath:               uIPath,
			SlowThresholdMilli:   slowThresholdMilli,
			SameServiceThreshold: sameServiceThreshold,
			SlowTraces:           slowTraces,
			SameServiceTraces:    sameServiceTraces,
		},
	}}
}

func (bo *BuildOption) createGrMap() map[string]interface{} {
	return map[string]interface{}{
		"suggestData": bo.SuggestData,
		"digData":     bo.DigData,
	}
}
