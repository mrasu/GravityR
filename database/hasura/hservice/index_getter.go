package hservice

import (
	"github.com/mrasu/GravityR/database/dmodel"
	"github.com/mrasu/GravityR/infra/hasura"
	"github.com/samber/lo"
)

type IndexGetter struct {
	cli *hasura.Client
}

func NewIndexGetter(cli *hasura.Client) *IndexGetter {
	return &IndexGetter{cli: cli}
}

func (ig *IndexGetter) GetIndexes(dbSchema string, tables []string) ([]*dmodel.IndexTarget, error) {
	infos, err := ig.cli.GetIndexes(dbSchema, tables)
	if err != nil {
		return nil, err
	}

	return lo.Map(infos, func(info *hasura.IndexInfo, _ int) *dmodel.IndexTarget {
		return dmodel.NewIndexTarget(info.TableName, info.Columns)
	}), nil
}
