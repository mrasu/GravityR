package collectors

import (
	"fmt"
	"github.com/mrasu/GravityR/database/mysql/models"
	"github.com/mrasu/GravityR/lib"
	"github.com/pingcap/tidb/parser/ast"
	driver "github.com/pingcap/tidb/types/parser_driver"
	"reflect"
)

type scopeCollector struct {
	scopes []*models.StmtScope
	errors []error

	scopeStack *lib.Stack[models.StmtScope]
	lastError  error
}

func CollectStmtScopes(rootNode ast.StmtNode) ([]*models.StmtScope, []error) {
	sc := &scopeCollector{scopeStack: lib.NewStack[models.StmtScope]()}
	if _, ok := rootNode.(*ast.SetOprStmt); ok {
		return nil, []error{lib.NewUnsupportedError("not supporting UNION, EXCEPT etc.")}
	}
	if _, ok := rootNode.(*ast.SelectStmt); !ok {
		return nil, []error{lib.NewUnsupportedError("only SELECT query is supported")}
	}

	rootNode.Accept(sc)
	if len(sc.errors) > 0 {
		return nil, sc.errors
	}

	return sc.scopes, nil
}

func (sc *scopeCollector) Enter(in ast.Node) (ast.Node, bool) {
	stmt, ok := in.(*ast.SelectStmt)
	if !ok {
		return in, false
	}

	top := sc.scopeStack.Top()
	var parent *models.StmtScope
	if top != nil {
		parent = top
	}
	scope := &models.StmtScope{Parent: parent, Name: models.RootScopeName}
	var foundFields []*models.Field
	for _, field := range stmt.Fields.Fields {
		var cols []*models.FieldColumn
		if field.Expr != nil {
			cols = sc.collectExprReferences(field.Expr, models.FieldReference)
			if err := sc.getAndClearLastError(); err != nil {
				sc.errors = append(sc.errors, err)
				continue
			}
		} else {
			cols = []*models.FieldColumn{{
				Table: field.WildCard.Table.L,
				Type:  models.FieldStar,
			}}
		}

		foundFields = append(foundFields, &models.Field{
			AsName:  field.AsName.L,
			Columns: cols,
		})
	}
	if stmt.From != nil {
		tables, fields := sc.collectTables(stmt.From)
		if err := sc.getAndClearLastError(); err != nil {
			sc.errors = append(sc.errors, err)
		}

		scope.Tables = tables
		foundFields = append(foundFields, fields...)
	}
	if stmt.Where != nil {
		f, err := sc.createFieldFromExpr(stmt.Where, models.FieldCondition)
		if err == nil {
			foundFields = append(foundFields, f)
		} else {
			sc.errors = append(sc.errors, err)
		}
	}
	if stmt.GroupBy != nil {
		for _, item := range stmt.GroupBy.Items {
			f, err := sc.createFieldFromExpr(item.Expr, models.FieldReference)
			if err == nil {
				foundFields = append(foundFields, f)
			} else {
				sc.errors = append(sc.errors, err)
			}
		}
	}
	if stmt.Having != nil {
		f, err := sc.createFieldFromExpr(stmt.Having.Expr, models.FieldReference)
		if err == nil {
			foundFields = append(foundFields, f)
		} else {
			sc.errors = append(sc.errors, err)
		}
	}
	if stmt.WindowSpecs != nil {
		sc.errors = append(sc.errors, lib.NewUnsupportedError("not supporting WINDOW"))
	}
	if stmt.OrderBy != nil {
		for _, item := range stmt.OrderBy.Items {
			f, err := sc.createFieldFromExpr(item.Expr, models.FieldReference)
			if err == nil {
				foundFields = append(foundFields, f)
			} else {
				sc.errors = append(sc.errors, err)
			}
		}
	}
	if stmt.Limit != nil {
		f, err := sc.createFieldFromExpr(stmt.Limit.Count, models.FieldReference)
		if err == nil {
			foundFields = append(foundFields, f)
		} else {
			sc.errors = append(sc.errors, err)
		}

		f, err = sc.createFieldFromExpr(stmt.Limit.Offset, models.FieldReference)
		if err == nil {
			foundFields = append(foundFields, f)
		} else {
			sc.errors = append(sc.errors, err)
		}
	}
	if stmt.Kind == ast.SelectStmtKindValues {
		sc.lastError = lib.NewInvalidAstError("SELECT for VALUES is not expected")
	}
	if stmt.With != nil {
		sc.lastError = lib.NewInvalidAstError("WITH is not expected")
	}

	for _, f := range foundFields {
		if len(f.Columns) > 0 {
			scope.Fields = append(scope.Fields, f)
		}
	}

	sc.scopeStack.Push(scope)
	if len(sc.scopes) == 0 {
		sc.scopes = append(sc.scopes, scope)
	}

	return in, false
}

func (sc *scopeCollector) Leave(in ast.Node) (ast.Node, bool) {
	if _, ok := in.(*ast.SelectStmt); ok {
		sc.scopeStack.Pop()
	}

	return in, true
}

func (sc *scopeCollector) getAndClearLastError() error {
	err := sc.lastError
	sc.lastError = nil
	return err
}

func (sc *scopeCollector) createFieldFromExpr(expr ast.ExprNode, fieldType models.FieldType) (*models.Field, error) {
	cols := sc.collectExprReferences(expr, fieldType)
	if err := sc.getAndClearLastError(); err != nil {
		return nil, err
	}

	return &models.Field{Columns: cols}, nil
}

func (sc *scopeCollector) collectExprReferences(expr ast.ExprNode, defaultType models.FieldType) []*models.FieldColumn {
	if expr == nil {
		return nil
	}

	var res []*models.FieldColumn
	switch e := expr.(type) {
	case *ast.BetweenExpr:
		res = append(res, sc.collectExprReferences(e.Expr, defaultType)...)
		res = append(res, sc.collectExprReferences(e.Left, defaultType)...)
		res = append(res, sc.collectExprReferences(e.Right, defaultType)...)
	case *ast.BinaryOperationExpr:
		res = append(res, sc.collectExprReferences(e.L, defaultType)...)
		res = append(res, sc.collectExprReferences(e.R, defaultType)...)
	case *ast.CaseExpr:
		res = append(res, sc.collectExprReferences(e.Value, defaultType)...)
		for _, w := range e.WhenClauses {
			res = append(res, sc.collectExprReferences(w.Expr, defaultType)...)
			res = append(res, sc.collectExprReferences(w.Result, defaultType)...)
		}
		res = append(res, sc.collectExprReferences(e.ElseClause, defaultType)...)
	case *ast.SubqueryExpr:
		// ignore the content of subquery as it will be in different scope.Scopes
		return []*models.FieldColumn{{Type: models.FieldSubquery}}
	case *ast.CompareSubqueryExpr:
		res = append(res, sc.collectExprReferences(e.L, defaultType)...)
		res = append(res, sc.collectExprReferences(e.R, defaultType)...)
	case *ast.TableNameExpr:
		sc.lastError = lib.NewInvalidAstError("Table definition is not expected")
	case *ast.ColumnNameExpr:
		return []*models.FieldColumn{{Table: e.Name.Table.L, Name: e.Name.Name.L, Type: defaultType}}
	case *ast.DefaultExpr:
		sc.lastError = lib.NewInvalidAstError("DEFAULT is not expected")
	case *ast.ExistsSubqueryExpr:
		res = append(res, sc.collectExprReferences(e.Sel, defaultType)...)
	case *ast.PatternInExpr:
		res = append(res, sc.collectExprReferences(e.Expr, defaultType)...)
		for _, v := range e.List {
			res = append(res, sc.collectExprReferences(v, defaultType)...)
		}
		res = append(res, sc.collectExprReferences(e.Sel, defaultType)...)
	case *ast.IsNullExpr:
		res = append(res, sc.collectExprReferences(e.Expr, defaultType)...)
	case *ast.IsTruthExpr:
		res = append(res, sc.collectExprReferences(e.Expr, defaultType)...)
	case *ast.PatternLikeExpr:
		res = append(res, sc.collectExprReferences(e.Expr, defaultType)...)
		res = append(res, sc.collectExprReferences(e.Pattern, defaultType)...)
	case *ast.ParenthesesExpr:
		res = append(res, sc.collectExprReferences(e.Expr, defaultType)...)
	case *ast.PositionExpr:
		// do nothing as position references a column which is defined in SELECT clause
		// e.g.) SELECT COUNT(*), name FROM users GROUP BY 2; <- `2` is `name`
	case *ast.PatternRegexpExpr:
		res = append(res, sc.collectExprReferences(e.Expr, defaultType)...)
		res = append(res, sc.collectExprReferences(e.Pattern, defaultType)...)
	case *ast.RowExpr:
		for _, v := range e.Values {
			res = append(res, sc.collectExprReferences(v, defaultType)...)
		}
	case *ast.UnaryOperationExpr:
		res = append(res, sc.collectExprReferences(e.V, defaultType)...)
	case *ast.ValuesExpr:
		sc.lastError = lib.NewInvalidAstError("VALUES is not expected")
	case *ast.VariableExpr:
		// do nothing as variable doesn't relate to table's column
	case *ast.MaxValueExpr:
		// MAXVALUE appears in table definition for partition
		sc.lastError = lib.NewInvalidAstError("MAXVALUE is not expected")
	case *ast.MatchAgainst:
		res = append(res, sc.collectExprReferences(e.Against, defaultType)...)
		for _, c := range e.ColumnNames {
			res = append(res, &models.FieldColumn{Table: c.Table.L, Name: c.Name.L, Type: defaultType})
		}
	case *ast.SetCollationExpr:
		// do nothing as COLLATE doesn't relate to table's column
	case *ast.FuncCallExpr:
		for _, arg := range e.Args {
			res = append(res, sc.collectExprReferences(arg, defaultType)...)
		}
	case *ast.FuncCastExpr:
		res = append(res, sc.collectExprReferences(e.Expr, defaultType)...)
	case *ast.TrimDirectionExpr:
		// do nothing as Direction of trim doesn't relate to table's column
		// Direction is "LEADING" at `TRIM(LEADING ' ' FROM '  hello  ')` for example.
	case *ast.AggregateFuncExpr:
		for _, arg := range e.Args {
			res = append(res, sc.collectExprReferences(arg, defaultType)...)
		}
	case *ast.WindowFuncExpr:
		for _, arg := range e.Args {
			res = append(res, sc.collectExprReferences(arg, defaultType)...)
		}
		if e.Spec.PartitionBy != nil {
			for _, i := range e.Spec.PartitionBy.Items {
				res = append(res, sc.collectExprReferences(i.Expr, defaultType)...)
			}
		}
		if e.Spec.OrderBy != nil {
			for _, i := range e.Spec.OrderBy.Items {
				res = append(res, sc.collectExprReferences(i.Expr, defaultType)...)
			}
		}

		// `Frame.Extent.Start` or End has Expr field, but not dig Frame because they don't reference column.
		// c.f.) https://dev.mysql.com/doc/refman/8.0/en/window-functions-frames.html
		// > expr can be a
		// > * ? parameter marker (for use in a prepared statement)
		// > * nonnegative numeric literal
		// > * temporal interval of the form INTERVAL val unit. For INTERVAL expressions, val specifies nonnegative interval value, and unit is a keyword indicating the units in which the value should be interpreted. (For details about the permitted units specifiers, see the description of the DATE_ADD() function in Section 12.7, “Date and Time Functions”.)
		//
		// if e.Spec.Frame != nil {}
	case *ast.TimeUnitExpr:
		// do nothing as time unit doesn't relate to table's column
	case *ast.GetFormatSelectorExpr:
		// do nothing as format like DATE in `GET_FORMAT(DATE, 'USA')` is literal
	case *driver.ValueExpr:
		// do nothing as simple value like 1 or 'hello' doesn't relate to table's column
	case *driver.ParamMarkerExpr:
		// do nothing as `?` for prepared statement doesn't relate to table's column
	default:
		sc.lastError = lib.NewUnsupportedError(
			fmt.Sprintf("not supporting query exists from %d. type: %s", e.OriginTextPosition(), reflect.TypeOf(e).Name()),
		)
		return nil
	}

	return res
}

func (sc *scopeCollector) collectTables(ref *ast.TableRefsClause) ([]*models.Table, []*models.Field) {
	if ref == nil {
		return nil, nil
	}

	tables, fields := sc.collectJoinReferences(ref.TableRefs)

	return tables, fields
}

func (sc *scopeCollector) collectReferencingTables(resultSet ast.ResultSetNode) ([]*models.Table, []*models.Field) {
	if resultSet == nil {
		return nil, nil
	}

	var tables []*models.Table
	var fields []*models.Field
	switch n := resultSet.(type) {
	case *ast.TableSource:
		switch src := n.Source.(type) {
		case *ast.TableName:
			if src.Schema.O != "" {
				sc.lastError = lib.NewUnsupportedError(
					fmt.Sprintf("not supporting query referencing database explicitly: %s.%s", src.Schema.O, src.Name.O),
				)
				return nil, nil
			}
			tables = append(tables, &models.Table{AsName: n.AsName.L, Name: src.Name.L})
		case *ast.SelectStmt:
			tables = append(tables, &models.Table{AsName: n.AsName.L})
		case *ast.SetOprStmt:
			sc.lastError = lib.NewUnsupportedError("not supporting UNION, EXCEPT etc.")
		case *ast.Join:
			tables, fields = sc.collectJoinReferences(src)
		default:
			sc.lastError = lib.NewUnsupportedError(
				fmt.Sprintf("not supporting table source exists from %d. type: %s", n.OriginTextPosition(), reflect.TypeOf(n).Name()),
			)
		}
	case *ast.Join:
		tables, fields = sc.collectJoinReferences(n)
	default:
		sc.lastError = lib.NewUnsupportedError(
			fmt.Sprintf("not supporting table reference exists from %d. type: %s", n.OriginTextPosition(), reflect.TypeOf(n).Name()),
		)
	}

	return tables, fields
}

func (sc *scopeCollector) collectJoinReferences(j *ast.Join) ([]*models.Table, []*models.Field) {
	if j == nil {
		return nil, nil
	}

	tables, fields := sc.collectReferencingTables(j.Left)
	rTables, rFields := sc.collectReferencingTables(j.Right)
	tables = append(tables, rTables...)
	fields = append(fields, rFields...)

	if j.On != nil {
		fields = append(fields, &models.Field{Columns: sc.collectExprReferences(j.On.Expr, models.FieldCondition)})
	}

	return tables, fields
}
