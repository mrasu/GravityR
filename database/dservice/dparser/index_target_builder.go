package dparser

import (
	"github.com/mrasu/GravityR/database/dmodel"
	"github.com/mrasu/GravityR/lib"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type IndexTargetBuilder struct {
	tableSchemaMap map[string]*TableSchema
}

type index []*IndexField

func NewIndexTargetBuilder(tablesSchemas []*TableSchema) *IndexTargetBuilder {
	tableMap := map[string]*TableSchema{}
	for _, t := range tablesSchemas {
		tableMap[t.Name] = t
	}
	return &IndexTargetBuilder{tableSchemaMap: tableMap}
}

func (itb *IndexTargetBuilder) Build(scopes []*StmtScope) ([]*dmodel.IndexTarget, error) {
	var tables []*IndexTargetTable
	for _, scope := range scopes {
		ts, err := itb.buildFromScope(scope)
		if err != nil {
			return nil, err
		}

		tables = append(tables, ts...)
	}

	SortIndexTargetTable(tables)
	idxTargets := lo.Map(tables, func(it *IndexTargetTable, _ int) *dmodel.IndexTarget { return it.ToIndexTarget() })
	return idxTargets, nil
}

func (itb *IndexTargetBuilder) buildFromScope(scope *StmtScope) ([]*IndexTargetTable, error) {
	var targets []*IndexTargetTable

	tableFields, err := buildTableFields(itb.tableSchemaMap, scope)
	if err != nil {
		return nil, err
	}

	for tName, f := range tableFields {
		tSchema := itb.tableSchemaMap[tName]
		if tSchema == nil {
			return nil, errors.Errorf("unknown table exists: %s", tName)
		}
		ifs := itb.buildIndexComposition(lib.NewSetS(f), tSchema.PrimaryKeys)

		for _, fs := range ifs {
			targets = append(targets, &IndexTargetTable{
				TableName:   tName,
				IndexFields: fs,
			})
		}
	}

	return targets, nil
}

func (itb *IndexTargetBuilder) buildIndexComposition(fieldColumns *lib.Set[*FieldColumn], pkCols []string) []index {
	candidateFields := itb.extractEfficientIndexes(fieldColumns)
	idxFieldList := itb.removePkOnlyFields(pkCols, candidateFields)

	return idxFieldList
}

func (itb *IndexTargetBuilder) extractEfficientIndexes(fieldColumns *lib.Set[*FieldColumn]) []index {
	compositions := itb.compositeFieldColumns(fieldColumns)

	var efficientIndexFields []index
	for _, fields := range compositions {
		for i := 0; i < len(fields); i++ {
			if fields[i].Type != FieldReference {
				continue
			}

			if len(fields) != fieldColumns.Count() {
				goto next
			}

			for j := i + 1; j < len(fields); j++ {
				if fields[j].Type != FieldReference {
					goto next
				}

				// pick ordered fields only to avoid duplicating index
				if fields[j-1].Name > fields[j].Name {
					goto next
				}
			}
		}
		efficientIndexFields = append(efficientIndexFields, fields)
	next:
	}

	return efficientIndexFields
}

func (itb *IndexTargetBuilder) compositeFieldColumns(fColumns *lib.Set[*FieldColumn]) []index {
	var compositions []index
	for _, fc := range fColumns.Values() {
		if !(fc.Type == FieldReference || fc.Type == FieldCondition) {
			continue
		}

		fColumns.Delete(fc)
		res2 := itb.compositeFieldColumns(fColumns)
		fColumns.Add(fc)

		compositions = append(compositions, index{{
			Name: fc.Name,
			Type: fc.Type,
		}})
		for _, r := range res2 {
			r2 := r[:]
			r2 = append(r2, &IndexField{
				Name: fc.Name,
				Type: fc.Type,
			})

			compositions = append(compositions, r2)
		}
	}

	return compositions
}

func (itb *IndexTargetBuilder) removePkOnlyFields(pkCols []string, indexFieldList []index) []index {
	pks := lib.NewSetS(pkCols)
	var res []index
	for _, fs := range indexFieldList {
		isAllPk := true
		for _, f := range fs {
			if !pks.Contains(f.Name) {
				isAllPk = false
				break
			}
		}

		if !isAllPk {
			res = append(res, fs)
		}
	}

	return res
}
