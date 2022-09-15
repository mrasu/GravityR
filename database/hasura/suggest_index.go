package hasura

import (
	"fmt"
	"github.com/auxten/postgresql-parser/pkg/sql/parser"
	"github.com/mrasu/GravityR/database/db_models"
	"github.com/mrasu/GravityR/database/db_models/builders"
	"github.com/mrasu/GravityR/database/hasura/models/collectors"
	"github.com/mrasu/GravityR/database/postgres/models"
	pCollectors "github.com/mrasu/GravityR/database/postgres/models/collectors"
	"github.com/mrasu/GravityR/infra/hasura"
	"github.com/mrasu/GravityR/lib"
	"github.com/pkg/errors"
)

// TODO: 既存のインデックスと被るものは除外する
// func SuggestIndex(db *postgres.DB, schema, query string, aTree *models.ExplainAnalyzeTree) ([]*db_models.IndexTargetTable, []error) {
func SuggestIndex(cli *hasura.Client, query string, aTree *models.ExplainAnalyzeTree) ([]*db_models.IndexTargetTable, []error) {
	stmt, err := parse(query)
	if err != nil {
		return nil, []error{err}
	}
	tNames, errs := pCollectors.CollectTableNames(stmt)
	if len(errs) > 0 {
		return nil, errs
	}
	tables, err := collectors.CollectTableSchemas(cli, "public", tNames)
	if err != nil {
		return nil, []error{err}
	}
	for _, t := range tables {
		fmt.Println(t)
	}
	scopes, errs := pCollectors.CollectStmtScopes(stmt)
	if len(errs) > 0 {
		return nil, errs
	}

	idxCandidates, err := builders.BuildIndexTargets(tables, scopes)
	if err != nil {
		return nil, []error{err}
	}
	fmt.Println(lib.Join(idxCandidates, "\n", func(f *db_models.IndexTargetTable) string { return f.String() }))
	fmt.Println()

	tableResults := aTree.ToSingleTableResults()
	fmt.Println(lib.Join(tableResults, "\n", func(st *db_models.SingleTableExplainResult) string { return st.String() }))

	return builders.BuildExplainedIndexTargets(idxCandidates, scopes, tableResults)
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
