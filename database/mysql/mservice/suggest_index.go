package mservice

import (
	"github.com/mrasu/GravityR/database/dmodel"
	"github.com/mrasu/GravityR/database/dservice"
	"github.com/mrasu/GravityR/database/mysql/mmodel"
	"github.com/mrasu/GravityR/infra/mysql"
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
)

func SuggestIndex(db *mysql.DB, database, query string, aTree *mmodel.ExplainAnalyzeTree) ([]*dmodel.IndexTargetTable, []error) {
	rootNode, err := parse(query)
	if err != nil {
		return nil, []error{err}
	}
	tNames, errs := CollectTableNames(rootNode)
	if len(errs) > 0 {
		return nil, errs
	}

	tables, err := CollectTableSchemas(db, database, tNames)
	if err != nil {
		return nil, []error{err}
	}

	scopes, errs := CollectStmtScopes(rootNode)
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

func parse(sql string) (ast.StmtNode, error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return nil, err
	}

	return stmtNodes[0], nil
}
