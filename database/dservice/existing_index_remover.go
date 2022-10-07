package dservice

import (
	"github.com/mrasu/GravityR/database/dmodel"
	"github.com/samber/lo"
)

type IndexGetter interface {
	GetIndexes(database string, tables []string) ([]*dmodel.IndexTarget, error)
}

type ExistingIndexRemover struct {
	idxGetter  IndexGetter
	idxTargets []*dmodel.IndexTarget
	dbName     string
}

func NewExistingIndexRemover(idxGetter IndexGetter, dbName string, its []*dmodel.IndexTarget) *ExistingIndexRemover {
	return &ExistingIndexRemover{
		idxGetter:  idxGetter,
		idxTargets: its,
		dbName:     dbName,
	}
}

func (r *ExistingIndexRemover) Remove() ([]*dmodel.IndexTarget, error) {
	idxTargets := lo.Uniq(r.idxTargets)

	tNames := lo.Map(idxTargets, func(it *dmodel.IndexTarget, _ int) string { return it.TableName })
	tNames = lo.Uniq(tNames)

	existingIdxes, err := r.idxGetter.GetIndexes(r.dbName, tNames)
	if err != nil {
		return nil, err
	}

	its := r.removeMatchedIndexTargets(existingIdxes, idxTargets)

	return its, nil
}

func (r *ExistingIndexRemover) removeMatchedIndexTargets(idxInfos []*dmodel.IndexTarget, its []*dmodel.IndexTarget) []*dmodel.IndexTarget {
	idxGroups := lo.GroupBy(idxInfos, func(i *dmodel.IndexTarget) string { return i.TableName })

	return lo.Filter(its, func(it *dmodel.IndexTarget, _ int) bool {
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
