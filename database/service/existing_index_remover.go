package service

import (
	"github.com/mrasu/GravityR/database"
	"github.com/samber/lo"
)

type IndexGetter interface {
	GetIndexes(database string, tables []string) ([]*database.IndexTarget, error)
}

type ExistingIndexRemover struct {
	idxGetter       IndexGetter
	idxTargetTables []*database.IndexTargetTable
	dbName          string
}

func NewExistingIndexRemover(idxGetter IndexGetter, dbName string, idxTargetTables []*database.IndexTargetTable) *ExistingIndexRemover {
	return &ExistingIndexRemover{
		idxGetter:       idxGetter,
		idxTargetTables: idxTargetTables,
		dbName:          dbName,
	}
}

func (r *ExistingIndexRemover) Remove() ([]*database.IndexTarget, error) {
	idxTargets := lo.Map(r.idxTargetTables, func(it *database.IndexTargetTable, _ int) *database.IndexTarget { return it.ToIndexTarget() })
	idxTargets = lo.Uniq(idxTargets)

	tNames := lo.Map(idxTargets, func(it *database.IndexTarget, _ int) string { return it.TableName })
	tNames = lo.Uniq(tNames)

	existingIdxes, err := r.idxGetter.GetIndexes(r.dbName, tNames)
	if err != nil {
		return nil, err
	}

	its := r.removeMatchedIndexTargets(existingIdxes, idxTargets)

	return its, nil
}

func (r *ExistingIndexRemover) removeMatchedIndexTargets(idxInfos []*database.IndexTarget, its []*database.IndexTarget) []*database.IndexTarget {
	idxGroups := lo.GroupBy(idxInfos, func(i *database.IndexTarget) string { return i.TableName })

	return lo.Filter(its, func(it *database.IndexTarget, _ int) bool {
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
