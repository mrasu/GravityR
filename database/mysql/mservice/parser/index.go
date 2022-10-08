package parser

import (
	"github.com/mrasu/GravityR/database/dmodel"
	"github.com/mrasu/GravityR/database/dservice/dparser"
	"github.com/mrasu/GravityR/infra/mysql"
	"github.com/pingcap/tidb/parser/ast"
)

func ListPossibleIndexes(db *mysql.DB, dbName string, rootNode ast.StmtNode) ([]*dmodel.IndexTarget, []error) {
	tNames, errs := CollectTableNames(rootNode)
	if len(errs) > 0 {
		return nil, errs
	}

	tables, err := CollectTableSchemas(db, dbName, tNames)
	if err != nil {
		return nil, []error{err}
	}

	scopes, errs := CollectStmtScopes(rootNode)
	if len(errs) > 0 {
		return nil, errs
	}

	its, err := dparser.NewIndexTargetBuilder(tables).Build(scopes)
	if err != nil {
		return nil, []error{err}
	}

	return its, nil
}
