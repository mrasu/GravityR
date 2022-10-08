package parser

import (
	"github.com/auxten/postgresql-parser/pkg/sql/parser"
	"github.com/mrasu/GravityR/database/dmodel"
	"github.com/mrasu/GravityR/database/dservice/dparser"
	"github.com/mrasu/GravityR/infra/postgres"
)

func ListPossibleIndexes(db *postgres.DB, schema string, stmt *parser.Statement) ([]*dmodel.IndexTarget, []error) {
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

	its, err := dparser.NewIndexTargetBuilder(tables).Build(scopes)
	if err != nil {
		return nil, []error{err}
	}

	return its, nil
}
