package hasura

import (
	"github.com/mrasu/GravityR/database"
	"github.com/mrasu/GravityR/infra/hasura"
	"github.com/samber/lo"
)

type IndexGetter struct {
	cli *hasura.Client
}

func NewIndexGetter(cli *hasura.Client) *IndexGetter {
	return &IndexGetter{cli: cli}
}

func (ig *IndexGetter) GetIndexes(dbSchema string, tables []string) ([]*database.IndexTarget, error) {
	infos, err := ig.cli.GetIndexes(dbSchema, tables)
	if err != nil {
		return nil, err
	}

	return lo.Map(infos, func(info *hasura.IndexInfo, _ int) *database.IndexTarget {
		return database.NewIndexTarget(info.TableName, info.Columns)
	}), nil
}
