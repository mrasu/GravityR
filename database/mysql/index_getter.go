package mysql

import (
	"github.com/mrasu/GravityR/database/common_model"
	"github.com/mrasu/GravityR/infra/mysql"
	"github.com/samber/lo"
)

type IndexGetter struct {
	db *mysql.DB
}

func NewIndexGetter(db *mysql.DB) *IndexGetter {
	return &IndexGetter{db: db}
}

func (ig *IndexGetter) GetIndexes(database string, tables []string) ([]*common_model.IndexTarget, error) {
	infos, err := ig.db.GetIndexes(database, tables)
	if err != nil {
		return nil, err
	}

	return lo.Map(infos, func(info *mysql.IndexInfo, _ int) *common_model.IndexTarget {
		return common_model.NewIndexTarget(info.TableName, info.Columns)
	}), nil
}
