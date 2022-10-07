package mservice

import (
	"github.com/mrasu/GravityR/database/dmodel"
	"github.com/mrasu/GravityR/database/dservice"
	"github.com/mrasu/GravityR/database/mysql/mmodel"
	"github.com/mrasu/GravityR/infra/mysql"
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
)

type IndexSuggester struct {
	db     *mysql.DB
	dbName string
}

func NewIndexSuggester(db *mysql.DB, dbName string) *IndexSuggester {
	return &IndexSuggester{
		db:     db,
		dbName: dbName,
	}
}

func (is *IndexSuggester) Suggest(query string, aTree *mmodel.ExplainAnalyzeTree) ([]*dmodel.IndexTarget, error) {
	rootNode, err := is.parse(query)
	if err != nil {
		return nil, err
	}

	its, errs := is.listPossibleIndexes(rootNode, aTree)
	if len(errs) > 0 {
		return nil, errs[0]
	}

	its, err = is.removeExistingIndexTargets(its)
	if err != nil {
		return nil, err
	}

	return its, nil
}

func (is *IndexSuggester) parse(sql string) (ast.StmtNode, error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return nil, err
	}

	return stmtNodes[0], nil
}

func (is *IndexSuggester) listPossibleIndexes(rootNode ast.StmtNode, aTree *mmodel.ExplainAnalyzeTree) ([]*dmodel.IndexTarget, []error) {
	tNames, errs := CollectTableNames(rootNode)
	if len(errs) > 0 {
		return nil, errs
	}

	tables, err := CollectTableSchemas(is.db, is.dbName, tNames)
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

func (is *IndexSuggester) removeExistingIndexTargets(its []*dmodel.IndexTarget) ([]*dmodel.IndexTarget, error) {
	idxGetter := NewIndexGetter(is.db)
	res, err := dservice.NewExistingIndexRemover(idxGetter, is.dbName, its).Remove()
	if err != nil {
		return nil, err
	}

	return res, nil
}
