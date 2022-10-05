package hasura

import (
	"github.com/mrasu/GravityR/database/common_model"
	"github.com/mrasu/GravityR/infra/hasura"
	"github.com/samber/lo"
)

type IndexGetter struct {
	cli *hasura.Client
}

func NewIndexGetter(cli *hasura.Client) *IndexGetter {
	return &IndexGetter{cli: cli}
}

func (ig *IndexGetter) GetIndexes(database string, tables []string) ([]*common_model.IndexTarget, error) {
	infos, err := ig.cli.GetIndexes(database, tables)
	if err != nil {
		return nil, err
	}

	return lo.Map(infos, func(info *hasura.IndexInfo, _ int) *common_model.IndexTarget {
		return common_model.NewIndexTarget(info.TableName, info.Columns)
	}), nil
}
