package mysql

import (
	"github.com/mrasu/GravityR/database"
	"github.com/mrasu/GravityR/infra/mysql"
	"github.com/samber/lo"
)

type IndexGetter struct {
	db *mysql.DB
}

func NewIndexGetter(db *mysql.DB) *IndexGetter {
	return &IndexGetter{db: db}
}

func (ig *IndexGetter) GetIndexes(dbName string, tables []string) ([]*database.IndexTarget, error) {
	infos, err := ig.db.GetIndexes(dbName, tables)
	if err != nil {
		return nil, err
	}

	return lo.Map(infos, func(info *mysql.IndexInfo, _ int) *database.IndexTarget {
		return database.NewIndexTarget(info.TableName, info.Columns)
	}), nil
}
