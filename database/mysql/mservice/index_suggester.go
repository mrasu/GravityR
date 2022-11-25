package mservice

import (
	"github.com/mrasu/GravityR/database/dmodel"
	"github.com/mrasu/GravityR/database/dservice"
	"github.com/mrasu/GravityR/database/mysql/mservice/parser"
	"github.com/mrasu/GravityR/infra/mysql"
	tParser "github.com/pingcap/tidb/parser"
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

func (is *IndexSuggester) Suggest(query string) ([]*dmodel.IndexTarget, error) {
	p := tParser.New()
	rootNode, err := Parse(p, query)
	if err != nil {
		return nil, err
	}

	its, errs := parser.ListPossibleIndexes(is.db, is.dbName, rootNode)
	if len(errs) > 0 {
		return nil, errs[0]
	}

	its, err = is.removeExistingIndexTargets(its)
	if err != nil {
		return nil, err
	}

	return its, nil
}

func (is *IndexSuggester) removeExistingIndexTargets(its []*dmodel.IndexTarget) ([]*dmodel.IndexTarget, error) {
	idxGetter := NewIndexGetter(is.db)
	res, err := dservice.NewExistingIndexRemover(idxGetter, is.dbName, its).Remove()
	if err != nil {
		return nil, err
	}

	return res, nil
}
