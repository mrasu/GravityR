package postgres

import (
	"github.com/auxten/postgresql-parser/pkg/sql/parser"
	"github.com/mrasu/GravityR/database"
	"github.com/mrasu/GravityR/database/postgres/model"
	"github.com/mrasu/GravityR/database/postgres/model/collector"
	"github.com/mrasu/GravityR/database/service"
	"github.com/mrasu/GravityR/infra/postgres"
	"github.com/pkg/errors"
)

func SuggestIndex(db *postgres.DB, schema, query string, aTree *model.ExplainAnalyzeTree) ([]*database.IndexTargetTable, []error) {
	stmt, err := parse(query)
	if err != nil {
		return nil, []error{err}
	}
	tNames, errs := collector.CollectTableNames(stmt)
	if len(errs) > 0 {
		return nil, errs
	}

	tables, err := collector.CollectTableSchemas(db, schema, tNames)
	if err != nil {
		return nil, []error{err}
	}

	scopes, errs := collector.CollectStmtScopes(stmt, "public")
	if len(errs) > 0 {
		return nil, errs
	}

	idxCandidates, err := service.BuildIndexTargets(tables, scopes)
	if err != nil {
		return nil, []error{err}
	}

	tableResults := aTree.ToSingleTableResults()
	return service.BuildExplainedIndexTargets(idxCandidates, scopes, tableResults)
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
