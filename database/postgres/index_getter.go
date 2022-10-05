package postgres

import (
	"github.com/mrasu/GravityR/database/common_model"
	"github.com/mrasu/GravityR/infra/postgres"
	"github.com/samber/lo"
)

type IndexGetter struct {
	db *postgres.DB
}

func NewIndexGetter(db *postgres.DB) *IndexGetter {
	return &IndexGetter{db: db}
}

func (ig *IndexGetter) GetIndexes(database string, tables []string) ([]*common_model.IndexTarget, error) {
	infos, err := ig.db.GetIndexes(database, tables)
	if err != nil {
		return nil, err
	}

	return lo.Map(infos, func(info *postgres.IndexInfo, _ int) *common_model.IndexTarget {
		return common_model.NewIndexTarget(info.TableName, info.Columns)
	}), nil
}
