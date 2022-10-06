package pservice

import (
	"github.com/auxten/postgresql-parser/pkg/sql/parser"
	"github.com/mrasu/GravityR/database/dmodel"
	"github.com/mrasu/GravityR/database/dservice"
	"github.com/mrasu/GravityR/database/postgres/pmodel"
	iPostgres "github.com/mrasu/GravityR/infra/postgres"
	"github.com/pkg/errors"
)

func SuggestIndex(db *iPostgres.DB, schema, query string, aTree *pmodel.ExplainAnalyzeTree) ([]*dmodel.IndexTargetTable, []error) {
	stmt, err := parse(query)
	if err != nil {
		return nil, []error{err}
	}
	tNames, errs := CollectTableNames(stmt)
	if len(errs) > 0 {
		return nil, errs
	}

	tables, err := CollectTableSchemas(db, schema, tNames)
	if err != nil {
		return nil, []error{err}
	}

	scopes, errs := CollectStmtScopes(stmt, "public")
	if len(errs) > 0 {
		return nil, errs
	}

	idxCandidates, err := dservice.BuildIndexTargets(tables, scopes)
	if err != nil {
		return nil, []error{err}
	}

	tableResults := aTree.ToSingleTableResults()
	return dservice.BuildExplainedIndexTargets(idxCandidates, scopes, tableResults)
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
