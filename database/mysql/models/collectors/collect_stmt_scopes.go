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

	scopeStack *lib.Stack[scopeHolder]
}

type scopeHolder struct {
	scope     *models.StmtScope
	lastError error
}

func CollectStmtScopes(rootNode ast.StmtNode) ([]*models.StmtScope, []error) {
	sc := &scopeCollector{scopeStack: lib.NewStack[scopeHolder]()}
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
	if stmt, ok := in.(*ast.SelectStmt); ok {
		top := sc.scopeStack.Top()
		var parent *models.StmtScope
		if top != nil {
			parent = top.scope
		}
		scope := &models.StmtScope{Parent: parent, Name: models.RootScopeName}
		holder := &scopeHolder{scope: scope}
		for _, field := range stmt.Fields.Fields {
			var cols []*models.FieldColumn
			if field.Expr != nil {
				cols = holder.collectExprReferences(field.Expr, models.FieldReference)
				if holder.lastError != nil {
					sc.errors = append(sc.errors, holder.lastError)
					holder.lastError = nil
					continue
				}
			} else {
				cols = []*models.FieldColumn{{
					Table: field.WildCard.Table.L,
					Type:  models.FieldStar,
				}}
			}

			scope.Fields = append(scope.Fields, &models.Field{
				AsName:  field.AsName.L,
				Columns: cols,
			})
		}
		if stmt.From != nil {
			tables, fields := holder.collectTables(stmt.From)
			if holder.lastError != nil {
				sc.errors = append(sc.errors, holder.lastError)
				holder.lastError = nil
			}

			scope.Tables = tables
			scope.Fields = append(scope.Fields, fields...)
		}
		// todo: where, group, having, order, window, cte, limit...

		sc.scopeStack.Push(holder)
		if len(sc.scopes) == 0 {
			sc.scopes = append(sc.scopes, scope)
		}

	}

	return in, false
}

func (sc *scopeCollector) Leave(in ast.Node) (ast.Node, bool) {
	if _, ok := in.(*ast.SelectStmt); ok {
		sc.scopeStack.Pop()
	}

	return in, true
}

func (sh *scopeHolder) collectExprReferences(expr ast.ExprNode, defaultType models.FieldType) []*models.FieldColumn {
	if expr == nil {
		return nil
	}

	var res []*models.FieldColumn
	switch e := expr.(type) {
	case *ast.BetweenExpr:
		res = append(res, sh.collectExprReferences(e.Expr, defaultType)...)
		res = append(res, sh.collectExprReferences(e.Left, defaultType)...)
		res = append(res, sh.collectExprReferences(e.Right, defaultType)...)
	case *ast.BinaryOperationExpr:
		res = append(res, sh.collectExprReferences(e.L, defaultType)...)
		res = append(res, sh.collectExprReferences(e.R, defaultType)...)
	case *ast.CaseExpr:
		res = append(res, sh.collectExprReferences(e.Value, defaultType)...)
		for _, w := range e.WhenClauses {
			res = append(res, sh.collectExprReferences(w.Expr, defaultType)...)
			res = append(res, sh.collectExprReferences(w.Result, defaultType)...)
		}
		res = append(res, sh.collectExprReferences(e.ElseClause, defaultType)...)
	case *ast.SubqueryExpr:
		// ignore the content of subquery as it will be in different scope.Scopes
		return []*models.FieldColumn{{Type: models.FieldSubquery}}
	case *ast.CompareSubqueryExpr:
		res = append(res, sh.collectExprReferences(e.L, defaultType)...)
		res = append(res, sh.collectExprReferences(e.R, defaultType)...)
	case *ast.TableNameExpr:
		sh.lastError = lib.NewInvalidAstError("Table definition is not expected")
	case *ast.ColumnNameExpr:
		return []*models.FieldColumn{{Table: e.Name.Table.L, Name: e.Name.Name.L, Type: defaultType}}
	case *ast.DefaultExpr:
		sh.lastError = lib.NewInvalidAstError("DEFAULT is not expected")
	case *ast.ExistsSubqueryExpr:
		res = append(res, sh.collectExprReferences(e.Sel, defaultType)...)
	case *ast.PatternInExpr:
		res = append(res, sh.collectExprReferences(e.Expr, defaultType)...)
		for _, v := range e.List {
			res = append(res, sh.collectExprReferences(v, defaultType)...)
		}
		res = append(res, sh.collectExprReferences(e.Sel, defaultType)...)
	case *ast.IsNullExpr:
		res = append(res, sh.collectExprReferences(e.Expr, defaultType)...)
	case *ast.IsTruthExpr:
		res = append(res, sh.collectExprReferences(e.Expr, defaultType)...)
	case *ast.PatternLikeExpr:
		res = append(res, sh.collectExprReferences(e.Expr, defaultType)...)
		res = append(res, sh.collectExprReferences(e.Pattern, defaultType)...)
	case *ast.ParenthesesExpr:
		res = append(res, sh.collectExprReferences(e.Expr, defaultType)...)
	case *ast.PositionExpr:
		sh.lastError = lib.NewInvalidAstError("ORDER or GROUP is not expected")
	case *ast.PatternRegexpExpr:
		res = append(res, sh.collectExprReferences(e.Expr, defaultType)...)
		res = append(res, sh.collectExprReferences(e.Pattern, defaultType)...)
	case *ast.RowExpr:
		for _, v := range e.Values {
			res = append(res, sh.collectExprReferences(v, defaultType)...)
		}
	case *ast.UnaryOperationExpr:
		res = append(res, sh.collectExprReferences(e.V, defaultType)...)
	case *ast.ValuesExpr:
		sh.lastError = lib.NewInvalidAstError("VALUES is not expected")
	case *ast.VariableExpr:
		// do nothing as variable doesn't relate to table's column
	case *ast.MaxValueExpr:
		sh.lastError = lib.NewInvalidAstError("MAXVALUE is not expected")
	case *ast.MatchAgainst:
		res = append(res, sh.collectExprReferences(e.Against, defaultType)...)
		for _, c := range e.ColumnNames {
			res = append(res, &models.FieldColumn{Table: c.Table.L, Name: c.Name.L, Type: defaultType})
		}
	case *ast.SetCollationExpr:
		// do nothing as COLLATE doesn't relate to table's column
	case *ast.FuncCallExpr:
		for _, arg := range e.Args {
			res = append(res, sh.collectExprReferences(arg, defaultType)...)
		}
	case *ast.FuncCastExpr:
		res = append(res, sh.collectExprReferences(e.Expr, defaultType)...)
	case *ast.TrimDirectionExpr:
		// do nothing as Direction of trim doesn't relate to table's column
		// Direction is "LEADING" at `TRIM(LEADING ' ' FROM '  hello  ')` for example.
	case *ast.AggregateFuncExpr:
		for _, arg := range e.Args {
			res = append(res, sh.collectExprReferences(arg, defaultType)...)
		}
	case *ast.WindowFuncExpr:
		for _, arg := range e.Args {
			res = append(res, sh.collectExprReferences(arg, defaultType)...)
		}
		if e.Spec.PartitionBy != nil {
			for _, i := range e.Spec.PartitionBy.Items {
				res = append(res, sh.collectExprReferences(i.Expr, defaultType)...)
			}
		}
		if e.Spec.OrderBy != nil {
			for _, i := range e.Spec.OrderBy.Items {
				res = append(res, sh.collectExprReferences(i.Expr, defaultType)...)
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
		sh.lastError = lib.NewUnsupportedError(
			fmt.Sprintf("not supporting query exists from %d. type: %s", e.OriginTextPosition(), reflect.TypeOf(e).Name()),
		)
		return nil
	}

	return res
}

func (sh *scopeHolder) collectTables(ref *ast.TableRefsClause) ([]*models.Table, []*models.Field) {
	if ref == nil {
		return nil, nil
	}

	tables, fields := sh.collectJoinReferences(ref.TableRefs)

	return tables, fields
}

func (sh *scopeHolder) collectReferencingTables(resultSet ast.ResultSetNode) ([]*models.Table, []*models.Field) {
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
				sh.lastError = lib.NewUnsupportedError(
					fmt.Sprintf("not supporting query referencing database explicitly: %s.%s", src.Schema.O, src.Name.O),
				)
				return nil, nil
			}
			tables = append(tables, &models.Table{AsName: n.AsName.L, Name: src.Name.L})
		case *ast.SelectStmt:
			tables = append(tables, &models.Table{AsName: n.AsName.L})
		case *ast.SetOprStmt:
			sh.lastError = lib.NewUnsupportedError("not supporting UNION, EXCEPT etc.")
		case *ast.Join:
			tables, fields = sh.collectJoinReferences(src)
		default:
			sh.lastError = lib.NewUnsupportedError(
				fmt.Sprintf("not supporting table source exists from %d. type: %s", n.OriginTextPosition(), reflect.TypeOf(n).Name()),
			)
		}
	case *ast.Join:
		tables, fields = sh.collectJoinReferences(n)
	default:
		sh.lastError = lib.NewUnsupportedError(
			fmt.Sprintf("not supporting table reference exists from %d. type: %s", n.OriginTextPosition(), reflect.TypeOf(n).Name()),
		)
	}

	return tables, fields
}

func (sh *scopeHolder) collectJoinReferences(j *ast.Join) ([]*models.Table, []*models.Field) {
	if j == nil {
		return nil, nil
	}

	tables, fields := sh.collectReferencingTables(j.Left)
	rTables, rFields := sh.collectReferencingTables(j.Right)
	tables = append(tables, rTables...)
	fields = append(fields, rFields...)

	if j.On != nil {
		fields = append(fields, &models.Field{Columns: sh.collectExprReferences(j.On.Expr, models.FieldCondition)})
	}

	return tables, fields
}
