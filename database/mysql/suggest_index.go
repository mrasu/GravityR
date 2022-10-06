package mysql

import (
	"github.com/mrasu/GravityR/database"
	"github.com/mrasu/GravityR/database/mysql/model"
	"github.com/mrasu/GravityR/database/mysql/model/collector"
	"github.com/mrasu/GravityR/database/service"
	"github.com/mrasu/GravityR/infra/mysql"
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
)

func SuggestIndex(db *mysql.DB, database, query string, aTree *model.ExplainAnalyzeTree) ([]*database.IndexTargetTable, []error) {
	rootNode, err := parse(query)
	if err != nil {
		return nil, []error{err}
	}
	tNames, errs := collector.CollectTableNames(rootNode)
	if len(errs) > 0 {
		return nil, errs
	}

	tables, err := collector.CollectTableSchemas(db, database, tNames)
	if err != nil {
		return nil, []error{err}
	}

	scopes, errs := collector.CollectStmtScopes(rootNode)
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

func parse(sql string) (ast.StmtNode, error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return nil, err
	}

	return stmtNodes[0], nil
}
