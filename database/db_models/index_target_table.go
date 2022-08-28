package db_models

import (
	"fmt"
	"github.com/mrasu/GravityR/lib"
)

type IndexTargetTable struct {
	TableName           string
	AffectingTableNames []string
	IndexFields         []*IndexField
}

func (itt *IndexTargetTable) ToIndexTarget() *IndexTarget {
	it := &IndexTarget{
		TableName: itt.TableName,
	}
	for _, f := range itt.IndexFields {
		it.Columns = append(it.Columns, &IndexColumn{name: f.Name})
	}
	return it
}

func (itt *IndexTargetTable) String() string {
	txt := fmt.Sprintf(
		"IndexTargetTable(table: %s, affectTable: %s, columns: %s)",
		itt.TableName,
		lib.JoinF(itt.AffectingTableNames, ", ", func(t string) string { return t }),
		lib.JoinF(itt.IndexFields, ", ", func(f *IndexField) string { return f.Name }),
	)
	return txt
}
