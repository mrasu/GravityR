package collectors

import (
	"fmt"
	"github.com/mrasu/GravityR/database/mysql/models"
	"github.com/mrasu/GravityR/lib"
	"github.com/pkg/errors"
	"strings"
)

type indexTargetCollector struct {
	tableSchemas []*models.TableSchema

	tableSchemaMap map[string]*models.TableSchema
}

func newIndexTargetCollector(tablesSchemas []*models.TableSchema) *indexTargetCollector {
	tableMap := map[string]*models.TableSchema{}
	for _, t := range tablesSchemas {
		tableMap[t.Name] = t
	}
	return &indexTargetCollector{tableSchemas: tablesSchemas, tableSchemaMap: tableMap}
}

func CollectIndexTargets(tablesSchemas []*models.TableSchema, scopes []*models.StmtScope) ([]*models.IndexTargetTable, error) {
	itc := newIndexTargetCollector(tablesSchemas)

	var targets []*models.IndexTargetTable
	for _, scope := range scopes {
		ts, err := itc.collectFromScope(scope)
		if err != nil {
			return nil, err
		}

		targets = append(targets, ts...)
	}

	return targets, nil
}

func (itc *indexTargetCollector) collectFromScope(scope *models.StmtScope) ([]*models.IndexTargetTable, error) {
	var targets []*models.IndexTargetTable

	groupedFields, err := itc.groupFieldColumnsByTableName(scope)
	if err != nil {
		return nil, err
	}

	for tName, f := range groupedFields {
		tSchema := itc.tableSchemaMap[tName]
		if tSchema == nil {
			return nil, errors.Errorf("unknown table exists: %s", tName)
		}
		ifs := itc.collectPossibleIndexFields(lib.NewSetS(f), tSchema.PrimaryKeys)

		for _, fs := range ifs {
			if itc.hasRedundantPk(fs, tSchema.PrimaryKeys) {
				continue
			}
			targets = append(targets, &models.IndexTargetTable{
				TableName:           tName,
				AffectingTableNames: []string{scope.Name},
				IndexFields:         fs,
			})
		}
	}

	return targets, nil
}

func (itc *indexTargetCollector) groupFieldColumnsByTableName(scope *models.StmtScope) (map[string][]*models.FieldColumn, error) {
	columnTableMap := itc.buildColumnTableSchemaMap(scope)

	for _, s := range scope.Scopes {
		// TODO: 内側のselect結果のカラムも候補に加える
		fmt.Printf("=TODO: recursive: %v", s)
	}

	aliasTs := map[string]string{}
	for _, t := range scope.Tables {
		aliasTs[t.AsName] = t.Name
	}
	// TODO: 内側のselect結果のカラムもカラムの情報に加える

	groupedFields := map[string][]*models.FieldColumn{}
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

func (itc *indexTargetCollector) buildColumnTableSchemaMap(scope *models.StmtScope) map[string]*models.TableSchema {
	res := map[string]*models.TableSchema{}

	for _, st := range scope.Tables {
		tName := st.Name
		if tName == "" {
			tName = st.AsName
		}
		t := itc.tableSchemaMap[tName]
		for _, c := range t.Columns {
			res[c.Name] = t
		}
	}

	return res
}

func (itc *indexTargetCollector) collectPossibleIndexFields(fColumns *lib.Set[*models.FieldColumn], pkCols []string) [][]*models.IndexField {
	candidateFields := itc.collectPossibleIndexFieldsRecursive(fColumns)

	pks := lib.NewSetS(pkCols)
	var res [][]*models.IndexField
	var refFields [][]*models.IndexField
	for _, fs := range candidateFields {
		hasRef := false
		for _, f := range fs {
			if f.Type == models.FieldReference {
				hasRef = true
			}
		}
		if hasRef {
			if fColumns.Count() <= len(fs)+len(pkCols) {
				if itc.isAllMissingColsPk(fs, fColumns, pks) {
					refFields = append(refFields, fs)
				}
			}
			continue
		}
		res = append(res, fs)
	}

	if len(refFields) > 0 {
		lib.SortF(refFields, func(fs []*models.IndexField) string {
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

func (itc *indexTargetCollector) collectPossibleIndexFieldsRecursive(fColumns *lib.Set[*models.FieldColumn]) [][]*models.IndexField {
	var res [][]*models.IndexField
	for _, fc := range fColumns.Values() {
		if !(fc.Type == models.FieldReference || fc.Type == models.FieldCondition) {
			continue
		}

		fColumns.Delete(fc)
		res2 := itc.collectPossibleIndexFieldsRecursive(fColumns)
		fColumns.Add(fc)

		res = append(res, []*models.IndexField{{
			Name: fc.Name,
			Type: fc.Type,
		}})
		for _, r := range res2 {
			if fc.Type == models.FieldCondition && r[len(r)-1].Type == models.FieldReference {
				continue
			}

			r2 := r[:]
			r2 = append(r2, &models.IndexField{
				Name: fc.Name,
				Type: fc.Type,
			})

			res = append(res, r2)
		}
	}

	return res
}

func (itc *indexTargetCollector) isAllMissingColsPk(targetFs []*models.IndexField, fColumns *lib.Set[*models.FieldColumn], pks *lib.Set[string]) bool {
	missingCols := lib.SelectF(fColumns.Values(), func(f *models.FieldColumn) bool {
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

func (itc *indexTargetCollector) hasRedundantPk(target []*models.IndexField, pkCols []string) bool {
	pks := lib.NewSetS(pkCols)
	for _, t := range target {
		if t.Type == models.FieldReference && pks.Contains(t.Name) {
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
