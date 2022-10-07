package pservice

import (
	"github.com/auxten/postgresql-parser/pkg/sql/parser"
	"github.com/mrasu/GravityR/database/dmodel"
	"github.com/mrasu/GravityR/database/dservice"
	"github.com/mrasu/GravityR/database/postgres/pmodel"
	"github.com/mrasu/GravityR/infra/postgres"
	"github.com/pkg/errors"
)

type IndexSuggester struct {
	db     *postgres.DB
	schema string
}

func NewIndexSuggester(db *postgres.DB, schema string) *IndexSuggester {
	return &IndexSuggester{
		db:     db,
		schema: schema,
	}
}

func (is *IndexSuggester) Suggest(query string, aTree *pmodel.ExplainAnalyzeTree) ([]*dmodel.IndexTarget, error) {
	stmt, err := is.parse(query)
	if err != nil {
		return nil, err
	}

	its, errs := is.listPossibleIndexes(stmt, aTree)
	if len(errs) > 0 {
		return nil, errs[0]
	}

	its, err = is.removeExistingIndexTargets(is.db, is.schema, its)
	if err != nil {
		return nil, err
	}

	return its, nil
}

func (is *IndexSuggester) parse(query string) (*parser.Statement, error) {
	stmts, err := parser.Parse(query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse sql")
	}

	if len(stmts) > 1 {
		return nil, errors.New("not supporting query having multiple statements")
	}
	return &stmts[0], nil
}

func (is *IndexSuggester) listPossibleIndexes(stmt *parser.Statement, aTree *pmodel.ExplainAnalyzeTree) ([]*dmodel.IndexTarget, []error) {
	tNames, errs := CollectTableNames(stmt)
	if len(errs) > 0 {
		return nil, errs
	}

	tables, err := CollectTableSchemas(is.db, is.schema, tNames)
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

func (is *IndexSuggester) removeExistingIndexTargets(db *postgres.DB, dbName string, its []*dmodel.IndexTarget) ([]*dmodel.IndexTarget, error) {
	idxGetter := NewIndexGetter(db)
	res, err := dservice.NewExistingIndexRemover(idxGetter, dbName, its).Remove()
	if err != nil {
		return nil, err
	}

	return res, nil
}
