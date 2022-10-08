package parser

import (
	pParser "github.com/auxten/postgresql-parser/pkg/sql/parser"
	"github.com/mrasu/GravityR/database/dmodel"
	"github.com/mrasu/GravityR/database/dservice/dparser"
	"github.com/mrasu/GravityR/database/postgres/pservice/parser"
	"github.com/mrasu/GravityR/infra/hasura"
)

func ListPossibleIndexes(cli *hasura.Client, stmt *pParser.Statement) ([]*dmodel.IndexTarget, []error) {
	tNames, errs := parser.CollectTableNames(stmt)
	if len(errs) > 0 {
		return nil, errs
	}
	tables, err := CollectTableSchemas(cli, "public", tNames)
	if err != nil {
		return nil, []error{err}
	}

	scopes, errs := parser.CollectStmtScopes(stmt, "public")
	if len(errs) > 0 {
		return nil, errs
	}

	its, err := dparser.NewIndexTargetBuilder(tables).Build(scopes)
	if err != nil {
		return nil, []error{err}
	}

	return its, nil
}
