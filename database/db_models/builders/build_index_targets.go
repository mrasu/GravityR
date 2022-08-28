package builders

import (
	"fmt"
	"github.com/mrasu/GravityR/database/db_models"
	"github.com/mrasu/GravityR/lib"
	"github.com/pkg/errors"
	"strings"
)

type indexTargetBuilder struct {
	tableSchemas []*db_models.TableSchema

	tableSchemaMap map[string]*db_models.TableSchema
}

func newIndexTargetBuilder(tablesSchemas []*db_models.TableSchema) *indexTargetBuilder {
	tableMap := map[string]*db_models.TableSchema{}
	for _, t := range tablesSchemas {
		tableMap[t.Name] = t
	}
	return &indexTargetBuilder{tableSchemas: tablesSchemas, tableSchemaMap: tableMap}
}

func BuildIndexTargets(tablesSchemas []*db_models.TableSchema, scopes []*db_models.StmtScope) ([]*db_models.IndexTargetTable, error) {
	itc := newIndexTargetBuilder(tablesSchemas)

	var targets []*db_models.IndexTargetTable
	for _, scope := range scopes {
		ts, err := itc.buildFromScope(scope)
		if err != nil {
			return nil, err
		}

		targets = append(targets, ts...)
	}

	return targets, nil
}

func (itb *indexTargetBuilder) buildFromScope(scope *db_models.StmtScope) ([]*db_models.IndexTargetTable, error) {
	var targets []*db_models.IndexTargetTable

	groupedFields, err := itb.groupFieldColumnsByTableName(scope)
	if err != nil {
		return nil, err
	}

	for tName, f := range groupedFields {
		tSchema := itb.tableSchemaMap[tName]
		if tSchema == nil {
			return nil, errors.Errorf("unknown table exists: %s", tName)
		}
		ifs := itb.buildPossibleIndexFields(lib.NewSetS(f), tSchema.PrimaryKeys)

		for _, fs := range ifs {
			if itb.hasRedundantPk(fs, tSchema.PrimaryKeys) {
				continue
			}
			targets = append(targets, &db_models.IndexTargetTable{
				TableName:           tName,
				AffectingTableNames: []string{scope.Name},
				IndexFields:         fs,
			})
		}
	}

	return targets, nil
}

func (itb *indexTargetBuilder) groupFieldColumnsByTableName(scope *db_models.StmtScope) (map[string][]*db_models.FieldColumn, error) {
	columnTableMap := itb.buildColumnTableSchemaMap(scope)

	for _, s := range scope.Scopes {
		// TODO: 内側のselect結果のカラムも候補に加える
		fmt.Printf("=TODO: recursive: %v", s)
	}

	aliasTs := map[string]string{}
	for _, t := range scope.Tables {
		aliasTs[t.AsName] = t.Name
	}
	// TODO: 内側のselect結果のカラムもカラムの情報に加える

	groupedFields := map[string][]*db_models.FieldColumn{}
	for _, f := range scope.Fields {
		for _, c := range f.Columns {
			if c.Table == "" {
				if t, ok := columnTableMap[c.Name]; ok {
					groupedFields[t.Name] = append(groupedFields[t.Name], c)
				} else {
					return nil, lib.NewInvalidAstError(fmt.Sprintf("no table found for %s", c.Name))
				}
			} else {
				if n, ok := aliasTs[c.Table]; ok {
					groupedFields[n] = append(groupedFields[n], c)
				} else {
					groupedFields[c.Table] = append(groupedFields[c.Table], c)
				}
			}
		}
	}

	return groupedFields, nil
}

func (itb *indexTargetBuilder) buildColumnTableSchemaMap(scope *db_models.StmtScope) map[string]*db_models.TableSchema {
	res := map[string]*db_models.TableSchema{}

	for _, st := range scope.Tables {
		tName := st.Name
		if tName == "" {
			tName = st.AsName
		}
		t := itb.tableSchemaMap[tName]
		for _, c := range t.Columns {
			res[c.Name] = t
		}
	}

	return res
}

func (itb *indexTargetBuilder) buildPossibleIndexFields(fColumns *lib.Set[*db_models.FieldColumn], pkCols []string) [][]*db_models.IndexField {
	candidateFields := itb.buildPossibleIndexFieldsRecursive(fColumns)

	pks := lib.NewSetS(pkCols)
	var res [][]*db_models.IndexField
	var refFields [][]*db_models.IndexField
	for _, fs := range candidateFields {
		hasRef := false
		for _, f := range fs {
			if f.Type == db_models.FieldReference {
				hasRef = true
			}
		}
		if hasRef {
			if fColumns.Count() <= len(fs)+len(pkCols) {
				if itb.isAllMissingColsPk(fs, fColumns, pks) {
					refFields = append(refFields, fs)
				}
			}
			continue
		}
		res = append(res, fs)
	}

	if len(refFields) > 0 {
		lib.SortF(refFields, func(fs []*db_models.IndexField) string {
			var names []string
			for _, f := range fs {
				names = append(names, f.Name)
			}
			return strings.Join(names, ",")
		})
		res = append(res, refFields[0])
	}

	return res
}

func (itb *indexTargetBuilder) buildPossibleIndexFieldsRecursive(fColumns *lib.Set[*db_models.FieldColumn]) [][]*db_models.IndexField {
	var res [][]*db_models.IndexField
	for _, fc := range fColumns.Values() {
		if !(fc.Type == db_models.FieldReference || fc.Type == db_models.FieldCondition) {
			continue
		}

		fColumns.Delete(fc)
		res2 := itb.buildPossibleIndexFieldsRecursive(fColumns)
		fColumns.Add(fc)

		res = append(res, []*db_models.IndexField{{
			Name: fc.Name,
			Type: fc.Type,
		}})
		for _, r := range res2 {
			if fc.Type == db_models.FieldCondition && r[len(r)-1].Type == db_models.FieldReference {
				continue
			}

			r2 := r[:]
			r2 = append(r2, &db_models.IndexField{
				Name: fc.Name,
				Type: fc.Type,
			})

			res = append(res, r2)
		}
	}

	return res
}

func (itb *indexTargetBuilder) isAllMissingColsPk(targetFs []*db_models.IndexField, fColumns *lib.Set[*db_models.FieldColumn], pks *lib.Set[string]) bool {
	missingCols := lib.SelectF(fColumns.Values(), func(f *db_models.FieldColumn) bool {
		for _, field := range targetFs {
			if field.Name == f.Name {
				return false
			}
		}
		return true
	})
	isAllMissingPk := true
	for _, mc := range missingCols {
		if !pks.Contains(mc.Name) {
			isAllMissingPk = false
			break
		}
	}

	return isAllMissingPk
}

func (itb *indexTargetBuilder) hasRedundantPk(target []*db_models.IndexField, pkCols []string) bool {
	pks := lib.NewSetS(pkCols)
	for _, t := range target {
		if t.Type == db_models.FieldReference && pks.Contains(t.Name) {
			return true
		}
	}

	for i := len(pkCols); i > 0; i-- {
		if len(target) < i {
			continue
		}

		matched := true
		for j := 0; j < i; j++ {
			if target[len(target)-i+j].Name != pkCols[j] {
				matched = false
				break
			}
		}
		if matched {
			return true
		}
	}

	return false
}