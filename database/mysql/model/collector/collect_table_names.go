package collector

import (
	"fmt"
	"github.com/mrasu/GravityR/lib"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/model"
)

type tableCollector struct {
	tables    []*table
	errors    []error
	cteTables [][]string
}

type table struct {
	name  model.CIStr
	alias model.CIStr
}

func CollectTableNames(rootNode ast.StmtNode) ([]string, []error) {
	t := &tableCollector{}
	if _, ok := rootNode.(*ast.SetOprStmt); ok {
		return nil, []error{lib.NewUnsupportedError("not supporting UNION, EXCEPT etc.")}
	}
	if _, ok := rootNode.(*ast.SelectStmt); !ok {
		return nil, []error{lib.NewUnsupportedError("only SELECT query is supported")}
	}

	rootNode.Accept(t)
	if len(t.errors) > 0 {
		return nil, t.errors
	}

	s := lib.NewSet[string]()
	for _, t := range t.tables {
		s.Add(t.name.L)
	}

	return s.Values(), nil
}

func (tc *tableCollector) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.SelectStmt); ok {
		tc.pushCteTables(stmt)
	} else if tSrc, ok := in.(*ast.TableSource); ok {
		tc.addTableIfNotDerived(tSrc)
	}

	return in, false
}

func (tc *tableCollector) Leave(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.SelectStmt); ok {
		tc.popCteTables(stmt)
	}

	return in, true
}

func (tc *tableCollector) pushCteTables(stmt *ast.SelectStmt) {
	if stmt.With == nil {
		return
	}

	var ts []string
	for _, cte := range stmt.With.CTEs {
		ts = append(ts, cte.Name.L)
	}
	tc.cteTables = append(tc.cteTables, ts)
}

func (tc *tableCollector) popCteTables(stmt *ast.SelectStmt) {
	if stmt.With == nil {
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

func (tc *tableCollector) addTableIfNotDerived(tSrc *ast.TableSource) {
	srcOrig := tSrc.Source
	if src, ok := srcOrig.(*ast.TableName); ok {
		if src.Schema.O != "" {
			tc.errors = append(tc.errors, lib.NewUnsupportedError(
				fmt.Sprintf("not supporting query referencing database explicitly: %s.%s", src.Schema.O, src.Name.O),
			))
			return
		}
		if !tc.isCteTable(src.Name.L) {
			tc.tables = append(tc.tables, &table{
				name:  src.Name,
				alias: tSrc.AsName,
			})
		}
	} else if _, ok := srcOrig.(*ast.SelectStmt); ok {
		// do nothing
	} else if _, ok := srcOrig.(*ast.SetOprStmt); ok {
		tc.errors = append(tc.errors, lib.NewUnsupportedError("not supporting UNION, EXCEPT etc."))
	} else if _, ok := srcOrig.(*ast.Join); ok {
		// Join happens when "WITH" is used but SELECT query doesn't have "WITH" inside it.
		tc.errors = append(tc.errors, NewUnexpectedChildError(tSrc, srcOrig))
	} else {
		tc.errors = append(tc.errors, NewUnexpectedChildError(tSrc, srcOrig))
	}
}
