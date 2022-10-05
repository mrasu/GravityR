package database

import (
	"github.com/mrasu/GravityR/database/common_model"
	"github.com/samber/lo"
)

type IndexGetter interface {
	GetIndexes(database string, tables []string) ([]*common_model.IndexTarget, error)
}

type ExistingIndexRemover struct {
	idxGetter       IndexGetter
	idxTargetTables []*common_model.IndexTargetTable
	dbName          string
}

func NewExistingIndexRemover(idxGetter IndexGetter, dbName string, idxTargetTables []*common_model.IndexTargetTable) *ExistingIndexRemover {
	return &ExistingIndexRemover{
		idxGetter:       idxGetter,
		idxTargetTables: idxTargetTables,
		dbName:          dbName,
	}
}

func (r *ExistingIndexRemover) Remove() ([]*common_model.IndexTarget, error) {
	idxTargets := lo.Map(r.idxTargetTables, func(it *common_model.IndexTargetTable, _ int) *common_model.IndexTarget { return it.ToIndexTarget() })
	idxTargets = lo.Uniq(idxTargets)

	tNames := lo.Map(idxTargets, func(it *common_model.IndexTarget, _ int) string { return it.TableName })
	tNames = lo.Uniq(tNames)

	existingIdxes, err := r.idxGetter.GetIndexes(r.dbName, tNames)
	if err != nil {
		return nil, err
	}

	its := r.removeMatchedIndexTargets(existingIdxes, idxTargets)

	return its, nil
}

func (r *ExistingIndexRemover) removeMatchedIndexTargets(idxInfos []*common_model.IndexTarget, its []*common_model.IndexTarget) []*common_model.IndexTarget {
	idxGroups := lo.GroupBy(idxInfos, func(i *common_model.IndexTarget) string { return i.TableName })

	return lo.Filter(its, func(it *common_model.IndexTarget, _ int) bool {
		if idxes, ok := idxGroups[it.TableName]; ok {
			for _, id := range idxes {
				if it.HasSameIdxColumns(id) {
					return false
				}
			}
			return true
		} else {
			return true
		}
	})
}
