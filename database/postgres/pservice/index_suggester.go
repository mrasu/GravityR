package pservice

import (
	"github.com/mrasu/GravityR/database/dmodel"
	"github.com/mrasu/GravityR/database/dservice"
	"github.com/mrasu/GravityR/database/postgres/pservice/parser"
	"github.com/mrasu/GravityR/infra/postgres"
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

func (is *IndexSuggester) Suggest(query string) ([]*dmodel.IndexTarget, error) {
	stmt, err := parser.Parse(query)
	if err != nil {
		return nil, err
	}

	its, errs := parser.ListPossibleIndexes(is.db, is.schema, stmt)
	if len(errs) > 0 {
		return nil, errs[0]
	}

	its, err = is.removeExistingIndexTargets(is.db, is.schema, its)
	if err != nil {
		return nil, err
	}

	return its, nil
}

func (is *IndexSuggester) removeExistingIndexTargets(db *postgres.DB, dbName string, its []*dmodel.IndexTarget) ([]*dmodel.IndexTarget, error) {
	idxGetter := NewIndexGetter(db)
	res, err := dservice.NewExistingIndexRemover(idxGetter, dbName, its).Remove()
	if err != nil {
		return nil, err
	}

	return res, nil
}
