package parser

import (
	"fmt"
	"github.com/auxten/postgresql-parser/pkg/sql/parser"
	"github.com/auxten/postgresql-parser/pkg/sql/sem/tree"
	"github.com/mrasu/GravityR/lib"
)

type tableCollector struct {
	tables    []*table
	errors    []error
	cteTables [][]string
}

type table struct {
	name  tree.Name
	alias tree.Name
}

func CollectTableNames(stmt *parser.Statement) ([]string, []error) {
	tc := &tableCollector{}
	if _, ok := stmt.AST.(*tree.Select); !ok {
		return nil, []error{lib.NewUnsupportedError("only SELECT query is supported")}
	}

	errs := walk(tc, stmt)
	if len(errs) > 0 {
		return nil, errs
	}

	if len(tc.errors) > 0 {
		return nil, tc.errors
	}

	return tc.getTableNames(), nil
}

func (tc *tableCollector) getTableNames() []string {
	s := lib.NewSet[string]()
	for _, t := range tc.tables {
		s.Add(t.name.Normalize())
	}

	return s.Values()
}

func (tc *tableCollector) Enter(node interface{}) bool {
	if s, ok := node.(*tree.Select); ok {
		tc.pushCteTables(s)
	}
	if t, ok := node.(*tree.AliasedTableExpr); ok {
		if n, ok := t.Expr.(*tree.TableName); ok {
			tc.addTableIfNotDerived(t.As.Alias, n)
		}
	}

	return false
}

func (tc *tableCollector) Leave(node interface{}) bool {
	if s, ok := node.(*tree.Select); ok {
		tc.popCteTables(s)
	}

	return false
}

func (tc *tableCollector) pushCteTables(s *tree.Select) {
	if s.With == nil {
		return
	}

	var ts []string
	for _, cte := range s.With.CTEList {
		ts = append(ts, cte.Name.Alias.Normalize())
	}
	tc.cteTables = append(tc.cteTables, ts)
}

func (tc *tableCollector) popCteTables(s *tree.Select) {
	if s.With == nil {
		return
	}

	tc.cteTables = tc.cteTables[:len(tc.cteTables)-1]
}

func (tc *tableCollector) isCteTable(name string) bool {
	for _, ctes := range tc.cteTables {
		for _, cte := range ctes {
			if cte == name {
				return true
			}
		}
	}
	return false
}

func (tc *tableCollector) addTableIfNotDerived(asName tree.Name, tName *tree.TableName) {
	if (tName.SchemaName != "public" && tName.SchemaName != "") || tName.CatalogName != "" {
		tc.errors = append(tc.errors, lib.NewUnsupportedError(
			fmt.Sprintf("not supporting query referencing schema or catalog explicitly: %s.%s.%s", tName.CatalogName, tName.SchemaName, tName.TableName),
		))
		return
	}
	if !tc.isCteTable(tName.TableName.Normalize()) {
		tc.tables = append(tc.tables, &table{
			name:  tName.TableName,
			alias: asName,
		})
	}
}
