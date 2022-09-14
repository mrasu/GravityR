package builders

import (
	"github.com/mrasu/GravityR/database/db_models"
	"github.com/mrasu/GravityR/lib"
	"github.com/pkg/errors"
)

type indexTargetBuilder struct {
	tableSchemaMap map[string]*db_models.TableSchema
}

type index []*db_models.IndexField

func newIndexTargetBuilder(tablesSchemas []*db_models.TableSchema) *indexTargetBuilder {
	tableMap := map[string]*db_models.TableSchema{}
	for _, t := range tablesSchemas {
		tableMap[t.Name] = t
	}
	return &indexTargetBuilder{tableSchemaMap: tableMap}
}

func BuildIndexTargets(tablesSchemas []*db_models.TableSchema, scopes []*db_models.StmtScope) ([]*db_models.IndexTargetTable, error) {
	itc := newIndexTargetBuilder(tablesSchemas)

	var tables []*db_models.IndexTargetTable
	for _, scope := range scopes {
		ts, err := itc.buildFromScope(scope)
		if err != nil {
			return nil, err
		}

		tables = append(tables, ts...)
	}

	db_models.SortIndexTargetTable(tables)
	return tables, nil
}

func (itb *indexTargetBuilder) buildFromScope(scope *db_models.StmtScope) ([]*db_models.IndexTargetTable, error) {
	var targets []*db_models.IndexTargetTable

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
			targets = append(targets, &db_models.IndexTargetTable{
				TableName:   tName,
				IndexFields: fs,
			})
		}
	}

	return targets, nil
}

func (itb *indexTargetBuilder) buildIndexComposition(fieldColumns *lib.Set[*db_models.FieldColumn], pkCols []string) []index {
	candidateFields := itb.extractEfficientIndexes(fieldColumns)
	idxFieldList := itb.removePkOnlyFields(pkCols, candidateFields)

	return idxFieldList
}

func (itb *indexTargetBuilder) extractEfficientIndexes(fieldColumns *lib.Set[*db_models.FieldColumn]) []index {
	compositions := itb.compositeFieldColumns(fieldColumns)

	var efficientIndexFields []index
	for _, fields := range compositions {
		for i := 0; i < len(fields); i++ {
			if fields[i].Type != db_models.FieldReference {
				continue
			}

			if len(fields) != fieldColumns.Count() {
				goto next
			}

			for j := i + 1; j < len(fields); j++ {
				if fields[j].Type != db_models.FieldReference {
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

func (itb *indexTargetBuilder) compositeFieldColumns(fColumns *lib.Set[*db_models.FieldColumn]) []index {
	var compositions []index
	for _, fc := range fColumns.Values() {
		if !(fc.Type == db_models.FieldReference || fc.Type == db_models.FieldCondition) {
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
			r2 = append(r2, &db_models.IndexField{
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
