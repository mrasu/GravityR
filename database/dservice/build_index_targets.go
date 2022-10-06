package dservice

import (
	"github.com/mrasu/GravityR/database/dmodel"
	"github.com/mrasu/GravityR/lib"
	"github.com/pkg/errors"
)

type indexTargetBuilder struct {
	tableSchemaMap map[string]*dmodel.TableSchema
}

type index []*dmodel.IndexField

func newIndexTargetBuilder(tablesSchemas []*dmodel.TableSchema) *indexTargetBuilder {
	tableMap := map[string]*dmodel.TableSchema{}
	for _, t := range tablesSchemas {
		tableMap[t.Name] = t
	}
	return &indexTargetBuilder{tableSchemaMap: tableMap}
}

func BuildIndexTargets(tablesSchemas []*dmodel.TableSchema, scopes []*dmodel.StmtScope) ([]*dmodel.IndexTargetTable, error) {
	itc := newIndexTargetBuilder(tablesSchemas)

	var tables []*dmodel.IndexTargetTable
	for _, scope := range scopes {
		ts, err := itc.buildFromScope(scope)
		if err != nil {
			return nil, err
		}

		tables = append(tables, ts...)
	}

	dmodel.SortIndexTargetTable(tables)
	return tables, nil
}

func (itb *indexTargetBuilder) buildFromScope(scope *dmodel.StmtScope) ([]*dmodel.IndexTargetTable, error) {
	var targets []*dmodel.IndexTargetTable

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
			targets = append(targets, &dmodel.IndexTargetTable{
				TableName:   tName,
				IndexFields: fs,
			})
		}
	}

	return targets, nil
}

func (itb *indexTargetBuilder) buildIndexComposition(fieldColumns *lib.Set[*dmodel.FieldColumn], pkCols []string) []index {
	candidateFields := itb.extractEfficientIndexes(fieldColumns)
	idxFieldList := itb.removePkOnlyFields(pkCols, candidateFields)

	return idxFieldList
}

func (itb *indexTargetBuilder) extractEfficientIndexes(fieldColumns *lib.Set[*dmodel.FieldColumn]) []index {
	compositions := itb.compositeFieldColumns(fieldColumns)

	var efficientIndexFields []index
	for _, fields := range compositions {
		for i := 0; i < len(fields); i++ {
			if fields[i].Type != dmodel.FieldReference {
				continue
			}

			if len(fields) != fieldColumns.Count() {
				goto next
			}

			for j := i + 1; j < len(fields); j++ {
				if fields[j].Type != dmodel.FieldReference {
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

func (itb *indexTargetBuilder) compositeFieldColumns(fColumns *lib.Set[*dmodel.FieldColumn]) []index {
	var compositions []index
	for _, fc := range fColumns.Values() {
		if !(fc.Type == dmodel.FieldReference || fc.Type == dmodel.FieldCondition) {
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
			r2 = append(r2, &dmodel.IndexField{
				Name: fc.Name,
				Type: fc.Type,
			})

			compositions = append(compositions, r2)
		}
	}

	return compositions
}

func (itb *indexTargetBuilder) removePkOnlyFields(pkCols []string, indexFieldList []index) []index {
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
