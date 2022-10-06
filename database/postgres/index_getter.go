package postgres

import (
	"github.com/mrasu/GravityR/database"
	"github.com/mrasu/GravityR/infra/postgres"
	"github.com/samber/lo"
)

type IndexGetter struct {
	db *postgres.DB
}

func NewIndexGetter(db *postgres.DB) *IndexGetter {
	return &IndexGetter{db: db}
}

func (ig *IndexGetter) GetIndexes(dbSchema string, tables []string) ([]*database.IndexTarget, error) {
	infos, err := ig.db.GetIndexes(dbSchema, tables)
	if err != nil {
		return nil, err
	}

	return lo.Map(infos, func(info *postgres.IndexInfo, _ int) *database.IndexTarget {
		return database.NewIndexTarget(info.TableName, info.Columns)
	}), nil
}
