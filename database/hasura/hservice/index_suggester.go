package hservice

import (
	"github.com/mrasu/GravityR/database/dmodel"
	"github.com/mrasu/GravityR/database/dservice"
	"github.com/mrasu/GravityR/database/hasura/hservice/parser"
	"github.com/mrasu/GravityR/database/postgres/pservice"
	"github.com/mrasu/GravityR/infra/hasura"
)

type IndexSuggester struct {
	cli *hasura.Client
}

func NewIndexSuggester(cli *hasura.Client) *IndexSuggester {
	return &IndexSuggester{
		cli: cli,
	}
}

func (is *IndexSuggester) Suggest(query string) ([]*dmodel.IndexTarget, error) {
	stmt, err := pservice.Parse(query)
	if err != nil {
		return nil, err
	}

	its, errs := parser.ListPossibleIndexes(is.cli, stmt)
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
	idxGetter := NewIndexGetter(is.cli)
	res, err := dservice.NewExistingIndexRemover(idxGetter, "public", its).Remove()
	if err != nil {
		return nil, err
	}

	return res, nil
}
