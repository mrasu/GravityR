package hasura

import (
	"fmt"
	"github.com/auxten/postgresql-parser/pkg/sql/parser"
	"github.com/mrasu/GravityR/database/common_model"
	"github.com/mrasu/GravityR/database/common_model/builder"
	"github.com/mrasu/GravityR/database/hasura/model/collector"
	"github.com/mrasu/GravityR/database/postgres/model"
	pCollector "github.com/mrasu/GravityR/database/postgres/model/collector"
	"github.com/mrasu/GravityR/infra/hasura"
	"github.com/mrasu/GravityR/lib"
	"github.com/pkg/errors"
)

// TODO: 既存のインデックスと被るものは除外する
// func SuggestIndex(db *postgres.DB, schema, query string, aTree *model.ExplainAnalyzeTree) ([]*common_model.IndexTargetTable, []error) {
func SuggestIndex(cli *hasura.Client, query string, aTree *model.ExplainAnalyzeTree) ([]*common_model.IndexTargetTable, []error) {
	stmt, err := parse(query)
	if err != nil {
		return nil, []error{err}
	}
	tNames, errs := pCollector.CollectTableNames(stmt)
	if len(errs) > 0 {
		return nil, errs
	}
	tables, err := collector.CollectTableSchemas(cli, "public", tNames)
	if err != nil {
		return nil, []error{err}
	}
	for _, t := range tables {
		fmt.Println(t)
	}
	scopes, errs := pCollector.CollectStmtScopes(stmt, "public")
	if len(errs) > 0 {
		return nil, errs
	}

	idxCandidates, err := builder.BuildIndexTargets(tables, scopes)
	if err != nil {
		return nil, []error{err}
	}
	fmt.Println(lib.Join(idxCandidates, "\n", func(f *common_model.IndexTargetTable) string { return f.String() }))
	fmt.Println()

	tableResults := aTree.ToSingleTableResults()
	fmt.Println(lib.Join(tableResults, "\n", func(st *common_model.SingleTableExplainResult) string { return st.String() }))

	return builder.BuildExplainedIndexTargets(idxCandidates, scopes, tableResults)
}

func parse(sql string) (*parser.Statement, error) {
	stmts, err := parser.Parse(sql)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse sql")
	}

	if len(stmts) > 1 {
		return nil, errors.New("not supporting query having multiple statements")
	}
	return &stmts[0], nil
}
