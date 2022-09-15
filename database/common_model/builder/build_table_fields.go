package builder

import (
	"fmt"
	"github.com/mrasu/GravityR/database/common_model"
	"github.com/mrasu/GravityR/lib"
	"github.com/samber/lo"
)

type tableFieldColumns map[string][]*common_model.FieldColumn

type outputCol struct {
	name           string
	isFieldCreated bool
	tableName      string
}

type tableFieldsBuilder struct {
	tableSchemaMap map[string]*common_model.TableSchema
}

func buildTableFields(tableSchemaMap map[string]*common_model.TableSchema, scope *common_model.StmtScope) (tableFieldColumns, error) {
	tfb := &tableFieldsBuilder{tableSchemaMap: tableSchemaMap}
	_, tableFields, err := tfb.calculateReferences(scope, map[string][]*outputCol{})
	if err != nil {
		return nil, err
	}

	return tableFields, nil
}

func (tfb *tableFieldsBuilder) calculateReferences(scope *common_model.StmtScope, externalCols map[string][]*outputCol) ([]*outputCol, tableFieldColumns, error) {
	info := tfb.newsScopeInfo(scope, externalCols)

	for _, s := range scope.Scopes {
		visibleCols := externalCols

		sName := s.Name
		if t, ok := info.tableMap[s.Name]; ok {
			if t.IsLateral {
				visibleCols = info.scopeOutputCols
			}
			sName = t.AsOrName()
		}

		cols, newFieldColumns, err := tfb.calculateReferences(s, visibleCols)
		if err != nil {
			return nil, nil, err
		}

		info.setVisitedScopeInfo(s, sName, cols, newFieldColumns)
	}

	for _, s := range scope.FieldScopes {
		cols, newFieldColumns, err := tfb.calculateReferences(s, info.scopeOutputCols)
		if err != nil {
			return nil, nil, err
		}

		info.setVisitedFieldScopeInfo(s, cols, newFieldColumns)
	}

	var outputCols []*outputCol
	for _, f := range scope.Fields {
		for _, c := range f.Columns {
			if err := info.setVisitedFieldColumnInfo(c); err != nil {
				return nil, nil, err
			}
		}

		n := f.Name()
		if n != "" {
			outputCols = append(outputCols, &outputCol{name: n, isFieldCreated: true})
		}
	}

	cols, err := info.listStarReferencingCols(tfb.tableSchemaMap)
	if err != nil {
		return nil, nil, err
	}

	foundCols := lib.NewSetS(lo.Map(outputCols, func(v *outputCol, _ int) string { return v.name }))
	cols = lo.Reject(cols, func(v *outputCol, _ int) bool { return foundCols.Contains(v.name) })

	outputCols = append(outputCols, cols...)
	return outputCols, info.tableFields, nil
}

func (tfb *tableFieldsBuilder) buildColumnTableSchemaMap(scope *common_model.StmtScope) map[string]*common_model.TableSchema {
	res := map[string]*common_model.TableSchema{}

	for _, st := range scope.Tables {
		tName := st.Name
		if tName == "" {
			tName = st.AsName
		}
		if t, ok := tfb.tableSchemaMap[tName]; ok {
			for _, c := range t.Columns {
				res[c.Name] = t
			}
		}
	}

	return res
}

type scopeInfo struct {
	scope *common_model.StmtScope

	tableFields tableFieldColumns
	aliasTs     map[string]string
	tableMap    map[string]*common_model.Table

	scopeOutputCols map[string][]*outputCol

	subScopes      *lib.Set[string]
	subScopeTables *lib.Set[string]

	starTables *lib.Set[string]

	columnTableMap map[string]*common_model.TableSchema
}

func (tfb *tableFieldsBuilder) newsScopeInfo(scope *common_model.StmtScope, externalCols map[string][]*outputCol) *scopeInfo {
	aliasTs := map[string]string{}
	tableMap := map[string]*common_model.Table{}
	for _, t := range scope.Tables {
		aliasTs[t.AsName] = t.Name
		tableMap[t.Name] = t
	}

	scopeOutputCols := map[string][]*outputCol{}
	subScopes := lib.NewSet[string]()
	for tName, cols := range externalCols {
		for _, colName := range cols {
			scopeOutputCols[tName] = append(scopeOutputCols[tName], colName)
		}
		subScopes.Add(tName)
	}

	columnTableMap := tfb.buildColumnTableSchemaMap(scope)

	return &scopeInfo{
		scope: scope,

		tableFields:     tableFieldColumns{},
		aliasTs:         aliasTs,
		tableMap:        tableMap,
		scopeOutputCols: scopeOutputCols,
		subScopes:       subScopes,
		subScopeTables:  lib.NewSet[string](),
		starTables:      lib.NewSet[string](),
		columnTableMap:  columnTableMap,
	}
}

func (si *scopeInfo) setVisitedScopeInfo(s *common_model.StmtScope, scopeName string, cols []*outputCol, newFieldColumns tableFieldColumns) {
	si.scopeOutputCols[scopeName] = cols

	si.addFieldColumns(newFieldColumns)
	si.subScopes.Add(s.Name)

	if t, ok := si.tableMap[s.Name]; ok {
		si.subScopeTables.Add(t.AsName)
	}
}

func (si *scopeInfo) setVisitedFieldScopeInfo(s *common_model.StmtScope, cols []*outputCol, newFieldColumns tableFieldColumns) {
	si.scopeOutputCols[s.Name] = cols
	si.addFieldColumns(newFieldColumns)
}

func (si *scopeInfo) setVisitedFieldColumnInfo(c *common_model.FieldColumn) error {
	if c.Type == common_model.FieldSubquery {
		return nil
	}

	if c.Type == common_model.FieldStar {
		// Ignore "*" field for index candidates as referenced columns can be changed by just adding one column
		si.starTables.Add(c.Table)
		return nil
	}

	if c.Table == "" {
		if t, ok := si.columnTableMap[c.Name]; ok {
			si.tableFields[t.Name] = append(si.tableFields[t.Name], c)
		} else if si.subScopeTables.Contains(c.Name) {
			// ignore field referencing table from subquery
			// e.g.) SELECT "_e" FROM (SELECT 1) AS "_e"
		} else {
			for _, fCols := range si.scopeOutputCols {
				for _, fc := range fCols {
					if c.Name == fc.name {
						if !fc.isFieldCreated && fc.tableName != "" {
							si.tableFields[fc.tableName] = append(si.tableFields[fc.tableName], c)
						}
						return nil
					}
				}
			}
			return lib.NewInvalidAstError(fmt.Sprintf("no table found for %s", c.Name))
		}

		return nil
	}

	if si.subScopes.Contains(c.Table) || si.subScopeTables.Contains(c.Table) {
		if fCols, ok := si.scopeOutputCols[c.Table]; ok {
			for _, fc := range fCols {
				if c.Name == fc.name {
					if !fc.isFieldCreated && fc.tableName != "" {
						si.tableFields[fc.tableName] = append(si.tableFields[fc.tableName], c)
						return nil
					}
				}
			}
		}

		return nil
	}

	if aName, ok := si.aliasTs[c.Table]; ok {
		if si.subScopes.Contains(aName) {
			return nil
		}
		si.tableFields[aName] = append(si.tableFields[aName], c)
	} else {
		si.tableFields[c.Table] = append(si.tableFields[c.Table], c)
	}

	return nil
}

func (si *scopeInfo) listStarReferencingCols(tableSchemaMap map[string]*common_model.TableSchema) ([]*outputCol, error) {
	var schemas []*common_model.TableSchema
	for _, tName := range si.starTables.Values() {
		if tName == "" {
			for _, t := range si.scope.Tables {
				if tSchema, ok := tableSchemaMap[t.Name]; ok {
					schemas = append(schemas, tSchema)
				}
			}
			continue
		}

		if _, ok := si.scopeOutputCols[tName]; ok {
			continue
		}

		if tSchema, ok := tableSchemaMap[tName]; ok {
			schemas = append(schemas, tSchema)
		} else {
			return nil, lib.NewInvalidAstError(fmt.Sprintf("no table found: %s", tName))
		}
	}

	var outputCols []*outputCol
	for _, schema := range schemas {
		for _, c := range schema.Columns {
			outputCols = append(outputCols, &outputCol{name: c.Name, isFieldCreated: false, tableName: schema.Name})
		}
	}

	outputCols = lib.UniqBy(outputCols, func(v *outputCol) string { return v.name })
	return outputCols, nil
}

func (si *scopeInfo) addFieldColumns(newFieldColumns tableFieldColumns) {
	lib.MergeValues(si.tableFields, newFieldColumns)
	for tName, cols := range si.tableFields {
		si.tableFields[tName] = lib.UniqBy(cols, func(c *common_model.FieldColumn) string { return c.Name })
	}
}
