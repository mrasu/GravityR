package builder

import (
	"github.com/mrasu/GravityR/database/common_model"
	"github.com/mrasu/GravityR/lib"
	"github.com/pkg/errors"
)

type indexTargetBuilder struct {
	tableSchemaMap map[string]*common_model.TableSchema
}

type index []*common_model.IndexField

func newIndexTargetBuilder(tablesSchemas []*common_model.TableSchema) *indexTargetBuilder {
	tableMap := map[string]*common_model.TableSchema{}
	for _, t := range tablesSchemas {
		tableMap[t.Name] = t
	}
	return &indexTargetBuilder{tableSchemaMap: tableMap}
}

func BuildIndexTargets(tablesSchemas []*common_model.TableSchema, scopes []*common_model.StmtScope) ([]*common_model.IndexTargetTable, error) {
	itc := newIndexTargetBuilder(tablesSchemas)

	var tables []*common_model.IndexTargetTable
	for _, scope := range scopes {
		ts, err := itc.buildFromScope(scope)
		if err != nil {
			return nil, err
		}

		tables = append(tables, ts...)
	}

	common_model.SortIndexTargetTable(tables)
	return tables, nil
}

func (itb *indexTargetBuilder) buildFromScope(scope *common_model.StmtScope) ([]*common_model.IndexTargetTable, error) {
	var targets []*common_model.IndexTargetTable

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
			targets = append(targets, &common_model.IndexTargetTable{
				TableName:   tName,
				IndexFields: fs,
			})
		}
	}

	return targets, nil
}

func (itb *indexTargetBuilder) buildIndexComposition(fieldColumns *lib.Set[*common_model.FieldColumn], pkCols []string) []index {
	candidateFields := itb.extractEfficientIndexes(fieldColumns)
	idxFieldList := itb.removePkOnlyFields(pkCols, candidateFields)

	return idxFieldList
}

func (itb *indexTargetBuilder) extractEfficientIndexes(fieldColumns *lib.Set[*common_model.FieldColumn]) []index {
	compositions := itb.compositeFieldColumns(fieldColumns)

	var efficientIndexFields []index
	for _, fields := range compositions {
		for i := 0; i < len(fields); i++ {
			if fields[i].Type != common_model.FieldReference {
				continue
			}

			if len(fields) != fieldColumns.Count() {
				goto next
			}

			for j := i + 1; j < len(fields); j++ {
				if fields[j].Type != common_model.FieldReference {
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

func (itb *indexTargetBuilder) compositeFieldColumns(fColumns *lib.Set[*common_model.FieldColumn]) []index {
	var compositions []index
	for _, fc := range fColumns.Values() {
		if !(fc.Type == common_model.FieldReference || fc.Type == common_model.FieldCondition) {
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
			r2 = append(r2, &common_model.IndexField{
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
