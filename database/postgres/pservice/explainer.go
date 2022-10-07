package pservice

import (
	"github.com/mrasu/GravityR/database/postgres/pmodel"
	"github.com/mrasu/GravityR/infra/postgres"
)

type Explainer struct {
	db *postgres.DB
}

func NewExplainer(db *postgres.DB) *Explainer {
	return &Explainer{
		db: db,
	}
}

func (e *Explainer) ExplainWithAnalyze(query string) (*pmodel.ExplainAnalyzeTree, error) {
	explainLines, err := e.db.ExplainWithAnalyze(query)
	if err != nil {
		return nil, err
	}

	return NewExplainAnalyzeTreeBuilder().Build(explainLines)
}
