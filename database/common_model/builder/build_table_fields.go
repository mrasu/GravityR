package builder

import (
	"fmt"
	"github.com/mrasu/GravityR/database/common_model"
	"github.com/mrasu/GravityR/lib"
	"github.com/samber/lo"
)

type tableOutputCols map[string][]*outputCol
type tableFieldColumns map[string][]*common_model.FieldColumn

type outputCol struct {
	name               string
	referencingColumns []*columnInfo
}

type columnInfo struct {
	tableName string
	name      string
}

type tableFieldsBuilder struct{}

func buildTableFields(dbTables map[string]*common_model.TableSchema, scope *common_model.StmtScope) (tableFieldColumns, error) {
	tfb := &tableFieldsBuilder{}
	info := newScopeInfo(dbTables, scope, tableOutputCols{})
	_, tableFields, err := tfb.calculateReferences(info)
	if err != nil {
		return nil, err
	}

	return tableFields, nil
}

func (tfb *tableFieldsBuilder) calculateReferences(info *scopeInfo) ([]*outputCol, tableFieldColumns, error) {
	scope := info.scope

	tableMap := lo.SliceToMap(scope.Tables, func(t *common_model.Table) (string, *common_model.Table) {
		return t.Name, t
	})

	for _, s := range scope.SubScopes {
		subScopeVisibleCols := info.originalExternalCols

		t := tableMap[s.Name]
		if t.IsLateral {
			subScopeVisibleCols = lo.Assign(info.originalExternalCols, info.getVisibleColumns())
		}

		childInfo := info.createChild(s, subScopeVisibleCols)
		cols, newFieldCols, err := tfb.calculateReferences(childInfo)
		if err != nil {
			return nil, nil, err
		}

		info.addSubqueryInfo(s.Name, t.AsOrName(), cols, newFieldCols)
	}

	visibleCols := info.getVisibleColumns()
	for _, s := range scope.FieldScopes {
		childInfo := info.createChild(s, visibleCols)
		cols, newFieldColumns, err := tfb.calculateReferences(childInfo)
		if err != nil {
			return nil, nil, err
		}

		info.addSubqueryInfo(s.Name, s.Name, cols, newFieldColumns)
	}

	var outputCols []*outputCol
	for _, f := range scope.Fields {
		var allColInfos []*columnInfo
		for _, c := range f.Columns {
			colInfos, err := info.getColumnReferenceColumnInfos(c)
			if err != nil {
				return nil, nil, err
			}
			for _, ci := range colInfos {
				fc := &common_model.FieldColumn{
					Table: ci.tableName,
					Name:  ci.name,
					Type:  c.Type,
				}
				info.addReferencingField(ci.tableName, fc)
			}
			allColInfos = append(allColInfos, colInfos...)

			if c.Type == common_model.FieldStar {
				info.starTables.Add(c.Table)
			}
		}

		n := f.Name()
		if n != "" {
			outputCols = append(outputCols, &outputCol{name: n, referencingColumns: allColInfos})
		}
	}

	colsFromStar, err := info.listStarReferencingCols(outputCols)
	if err != nil {
		return nil, nil, err
	}

	outputCols = append(outputCols, colsFromStar...)
	return outputCols, info.referencingFields, nil
}

type scopeInfo struct {
	dbTables map[string]*common_model.TableSchema

	scope                *common_model.StmtScope
	originalExternalCols tableOutputCols

	referencingFields tableFieldColumns

	refTables *refTables

	starTables *lib.Set[string]
}

func newScopeInfo(dbTables map[string]*common_model.TableSchema, scope *common_model.StmtScope, externalCols tableOutputCols) *scopeInfo {
	referencingColumns := buildReferencingColumns(dbTables, scope.Tables)
	refTables := newRefTables(referencingColumns, scope.Tables, externalCols)

	return &scopeInfo{
		dbTables:             dbTables,
		scope:                scope,
		originalExternalCols: externalCols,

		referencingFields: tableFieldColumns{},
		refTables:         refTables,
		starTables:        lib.NewSet[string](),
	}
}

func buildReferencingColumns(dbTables map[string]*common_model.TableSchema, tables []*common_model.Table) tableOutputCols {
	res := tableOutputCols{}
	for _, table := range tables {
		tName := table.Name
		if t, ok := dbTables[tName]; ok {
			for _, col := range t.Columns {
				res[tName] = append(res[tName], &outputCol{
					name: col.Name,
					referencingColumns: []*columnInfo{{
						tableName: tName,
						name:      col.Name,
					}},
				})
			}
		}
	}

	return res
}

func (si *scopeInfo) createChild(childScope *common_model.StmtScope, externalCols tableOutputCols) *scopeInfo {
	return newScopeInfo(si.dbTables, childScope, externalCols)
}

func (si *scopeInfo) addSubqueryInfo(name, asName string, cols []*outputCol, newFieldColumns tableFieldColumns) {
	si.refTables.addSubquery(name, asName, cols)
	si.addReferencingFields(newFieldColumns)
}

func (si *scopeInfo) getColumnReferenceColumnInfos(c *common_model.FieldColumn) ([]*columnInfo, error) {
	if c.Type == common_model.FieldSubquery {
		return nil, nil
	}

	if c.Type == common_model.FieldStar {
		// Ignore "*" field for index candidates as referenced columns can be changed by just adding one column
		return nil, nil
	}

	if c.Table == "" {
		cols, err := si.refTables.getColumnInfos(c.Name)
		if err != nil {
			return nil, err
		}

		return cols, nil
	}

	return si.refTables.getReferencingColumnInfos(c)
}

func (si *scopeInfo) addReferencingFields(newFieldColumns tableFieldColumns) {
	for tName, cols := range newFieldColumns {
		for _, col := range cols {
			si.addReferencingField(tName, col)
		}
	}
}

func (si *scopeInfo) addReferencingField(tableName string, col *common_model.FieldColumn) {
	cols, ok := si.referencingFields[tableName]
	if !ok {
		si.referencingFields[tableName] = []*common_model.FieldColumn{col}
		return
	}

	for i, c := range cols {
		if c.Name == col.Name {
			if c.Type == common_model.FieldReference && col.Type == common_model.FieldCondition {
				si.referencingFields[tableName][i] = col
			}
			return
		}
	}
	si.referencingFields[tableName] = append(si.referencingFields[tableName], col)
}

func (si *scopeInfo) listStarReferencingCols(foundCols []*outputCol) ([]*outputCol, error) {
	var res []*outputCol
	for _, tName := range si.starTables.Values() {
		if tName == "" {
			for _, cols := range si.refTables.internalCols {
				res = append(res, cols...)
			}
			continue
		}

		cols, err := si.refTables.getTableColsFromAsName(tName)
		if err != nil {
			return nil, err
		}
		res = append(res, cols...)
	}

	founds := lib.NewSetS(lo.Map(foundCols, func(v *outputCol, _ int) string { return v.name }))
	res = lo.Reject(res, func(v *outputCol, _ int) bool { return founds.Contains(v.name) })
	res = lib.UniqBy(res, func(v *outputCol) string { return v.name })

	return res, nil
}

func (si *scopeInfo) getVisibleColumns() tableOutputCols {
	res := tableOutputCols{}
	for asName, name := range si.refTables.asToName {
		if cols, ok := si.refTables.internalCols[name]; ok {
			res[asName] = cols
		}
	}

	return res
}

type refTables struct {
	internalCols tableOutputCols
	externalCols tableOutputCols

	asToName        map[string]string
	subqueryAsNames *lib.Set[string]
}

func newRefTables(referencingColumns tableOutputCols, tables []*common_model.Table, externalCols tableOutputCols) *refTables {
	asToName := map[string]string{}
	subqueryAsNames := lib.NewSet[string]()
	for tName := range externalCols {
		asToName[tName] = tName
		subqueryAsNames.Add(tName)
	}
	for _, t := range tables {
		asToName[t.AsOrName()] = t.Name
	}

	return &refTables{
		internalCols:    referencingColumns,
		externalCols:    externalCols,
		asToName:        asToName,
		subqueryAsNames: subqueryAsNames,
	}
}

func (t *refTables) addSubquery(name, asName string, cols []*outputCol) {
	t.internalCols[name] = cols
	t.asToName[asName] = name
	t.subqueryAsNames.Add(asName)
}

func (t *refTables) getTableColsFromName(name string) ([]*outputCol, error) {
	if cols, ok := t.internalCols[name]; ok {
		return cols, nil
	}
	res, ok := t.externalCols[name]
	if !ok {
		return nil, lib.NewInvalidAstError(fmt.Sprintf("no table found: %s", name))
	}

	return res, nil
}

func (t *refTables) getTableColsFromAsName(asName string) ([]*outputCol, error) {
	name, ok := t.asToName[asName]
	if !ok {
		return nil, lib.NewInvalidAstError(fmt.Sprintf("no table found: %s", asName))
	}

	return t.getTableColsFromName(name)
}

func (t *refTables) getReferencingColumnInfos(c *common_model.FieldColumn) ([]*columnInfo, error) {
	fCols, err := t.getTableColsFromAsName(c.Table)
	if err != nil {
		return nil, err
	}

	for _, fc := range fCols {
		if c.Name == fc.name {
			return fc.referencingColumns, nil
		}
	}

	return nil, lib.NewInvalidAstError(fmt.Sprintf("no column found for %s", c.Name))
}

func (t *refTables) getColumnInfos(colName string) ([]*columnInfo, error) {
	if cols, ok := t.getColumnInfosFromTableOutputCols(t.internalCols, colName); ok {
		return cols, nil
	}
	if cols, ok := t.getColumnInfosFromTableOutputCols(t.externalCols, colName); ok {
		return cols, nil
	}

	if t.subqueryAsNames.Contains(colName) {
		// ignore field referencing table from subquery
		// e.g.) SELECT "_e" FROM (SELECT 1) AS "_e"
		return nil, nil
	}

	return nil, lib.NewInvalidAstError(fmt.Sprintf("no column found for %s", colName))
}

func (t *refTables) getColumnInfosFromTableOutputCols(tableOutputs tableOutputCols, colName string) ([]*columnInfo, bool) {
	for _, fCols := range tableOutputs {
		for _, fc := range fCols {
			if colName == fc.name {
				return fc.referencingColumns, true
			}
		}
	}

	return nil, false
}
