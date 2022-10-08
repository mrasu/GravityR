package dparser

import (
	"fmt"
	"github.com/mrasu/GravityR/database/dmodel"
	"github.com/mrasu/GravityR/lib"
	"sort"
)

type IndexTargetTable struct {
	TableName   string
	IndexFields []*IndexField
}

func (itt *IndexTargetTable) ToIndexTarget() *dmodel.IndexTarget {
	it := &dmodel.IndexTarget{
		TableName: itt.TableName,
	}
	for _, f := range itt.IndexFields {
		it.Columns = append(it.Columns, &dmodel.IndexColumn{Name: f.Name})
	}
	return it
}

func (itt *IndexTargetTable) String() string {
	txt := fmt.Sprintf(
		"IndexTargetTable(table: %s columns: %s)",
		itt.TableName,
		lib.Join(itt.IndexFields, ", ", func(f *IndexField) string { return f.Name }),
	)
	return txt
}

func SortIndexTargetTable(tables []*IndexTargetTable) {
	sort.Slice(tables, func(i, j int) bool {
		if tables[i].TableName != tables[j].TableName {
			return tables[i].TableName < tables[j].TableName
		}

		iIndexes := tables[i].IndexFields
		jIndexes := tables[j].IndexFields
		if len(iIndexes) != len(jIndexes) {
			return len(iIndexes) < len(jIndexes)
		}

		for idx := 0; idx < len(iIndexes); idx++ {
			if iIndexes[idx].Name == jIndexes[idx].Name {
				continue
			}
			return iIndexes[idx].Name < jIndexes[idx].Name
		}

		return false
	})
}
