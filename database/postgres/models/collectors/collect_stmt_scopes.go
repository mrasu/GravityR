package collectors

import (
	"fmt"
	"github.com/auxten/postgresql-parser/pkg/sql/parser"
	"github.com/auxten/postgresql-parser/pkg/sql/sem/tree"
	"github.com/mrasu/GravityR/database/db_models"
	"github.com/mrasu/GravityR/lib"
	"reflect"
)

type scopeCollector struct {
	scopes []*db_models.StmtScope
	errors []error

	scopeStack *lib.Stack[db_models.StmtScope]
	lastError  error
}

func CollectStmtScopes(stmt *parser.Statement) ([]*db_models.StmtScope, []error) {
	sc := &scopeCollector{scopeStack: lib.NewStack[db_models.StmtScope]()}
	if _, ok := stmt.AST.(*tree.Select); !ok {
		return nil, []error{lib.NewUnsupportedError("only SELECT query is supported")}
	}

	errs := walk(sc, stmt)
	if len(errs) > 0 {
		return nil, errs
	}

	if len(sc.errors) > 0 {
		return nil, sc.errors
	}

	return sc.scopes, nil
}

func (sc *scopeCollector) Enter(node interface{}) bool {
	selectNode, ok := node.(*tree.Select)
	if !ok {
		return false
	}
	stmt, ok := selectNode.Select.(*tree.SelectClause)
	if !ok {
		return false
	}

	top := sc.scopeStack.Top()
	var parent *db_models.StmtScope
	if top != nil {
		parent = top
	}
	scope := &db_models.StmtScope{Parent: parent, Name: db_models.RootScopeName}
	var foundFields []*db_models.Field
	for _, expr := range stmt.Exprs {
		cols := sc.collectExprReferences(expr.Expr, db_models.FieldReference)
		if err := sc.getAndClearLastError(); err != nil {
			sc.errors = append(sc.errors, err)
			continue
		}

		foundFields = append(foundFields, &db_models.Field{
			AsName:  string(expr.As),
			Columns: cols,
		})
	}
	for _, t := range stmt.From.Tables {
		tables, fields := sc.collectTables(t)
		if err := sc.getAndClearLastError(); err != nil {
			sc.errors = append(sc.errors, err)
		}

		scope.Tables = tables
		foundFields = append(foundFields, fields...)
	}
	if stmt.Where != nil {
		f, err := sc.createFieldFromExpr(stmt.Where.Expr, db_models.FieldCondition)
		if err == nil {
			foundFields = append(foundFields, f)
		} else {
			sc.errors = append(sc.errors, err)
		}
	}
	if stmt.GroupBy != nil {
		for _, g := range stmt.GroupBy {
			f, err := sc.createFieldFromExpr(g, db_models.FieldReference)
			if err == nil {
				foundFields = append(foundFields, f)
			} else {
				sc.errors = append(sc.errors, err)
			}
		}
	}
	if stmt.Having != nil {
		f, err := sc.createFieldFromExpr(stmt.Having.Expr, db_models.FieldReference)
		if err == nil {
			foundFields = append(foundFields, f)
		} else {
			sc.errors = append(sc.errors, err)
		}
	}
	if stmt.Window != nil {
		sc.errors = append(sc.errors, lib.NewUnsupportedError("not supporting WINDOW"))
	}
	if selectNode.OrderBy != nil {
		for _, o := range selectNode.OrderBy {
			f, err := sc.createFieldFromExpr(o.Expr, db_models.FieldReference)
			if err == nil {
				foundFields = append(foundFields, f)
			} else {
				sc.errors = append(sc.errors, err)
			}
		}
	}
	if selectNode.Limit != nil {
		f, err := sc.createFieldFromExpr(selectNode.Limit.Count, db_models.FieldReference)
		if err == nil {
			foundFields = append(foundFields, f)
		} else {
			sc.errors = append(sc.errors, err)
		}

		f, err = sc.createFieldFromExpr(selectNode.Limit.Offset, db_models.FieldReference)
		if err == nil {
			foundFields = append(foundFields, f)
		} else {
			sc.errors = append(sc.errors, err)
		}
	}
	if selectNode.With != nil {
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

	return false
}

func (sc *scopeCollector) Leave(node interface{}) bool {
	if _, ok := node.(*tree.SelectClause); ok {
		sc.scopeStack.Pop()
	}

	return false
}

func (sc *scopeCollector) getAndClearLastError() error {
	err := sc.lastError
	sc.lastError = nil
	return err
}

func (sc *scopeCollector) createFieldFromExpr(expr tree.Expr, fieldType db_models.FieldType) (*db_models.Field, error) {
	cols := sc.collectExprReferences(expr, fieldType)
	if err := sc.getAndClearLastError(); err != nil {
		return nil, err
	}

	return &db_models.Field{Columns: cols}, nil
}

func (sc *scopeCollector) collectExprReferences(expr tree.Expr, defaultType db_models.FieldType) []*db_models.FieldColumn {
	if expr == nil {
		return nil
	}

	var res []*db_models.FieldColumn

	switch e := expr.(type) {
	case *tree.RangeCond:
		res = append(res, sc.collectExprReferences(e.Left, defaultType)...)
		res = append(res, sc.collectExprReferences(e.From, defaultType)...)
		res = append(res, sc.collectExprReferences(e.To, defaultType)...)
	case *tree.BinaryExpr:
		res = append(res, sc.collectExprReferences(e.Left, defaultType)...)
		res = append(res, sc.collectExprReferences(e.Right, defaultType)...)
	case *tree.CaseExpr:
		res = append(res, sc.collectExprReferences(e.Expr, defaultType)...)
		for _, w := range e.Whens {
			res = append(res, sc.collectExprReferences(w.Cond, defaultType)...)
			res = append(res, sc.collectExprReferences(w.Val, defaultType)...)
		}
		res = append(res, sc.collectExprReferences(e.Else, defaultType)...)
	case *tree.Subquery:
		// ignore the content of subquery as it will be in different scope.Scopes
		return []*db_models.FieldColumn{{Type: db_models.FieldSubquery}}
	case *tree.CommentOnColumn:
		// do nothing as comment doesn't relate to column reference
	case *tree.NumVal:
		// do nothing as number is not reference but constant
	case *tree.StrVal:
		// do nothing as string is not reference but constant
	case *tree.DBool:
		// do nothing as bool is not reference but constant
	case *tree.DBitArray:
		// do nothing as bitarray is not reference but constant
	case *tree.DInt:
		// do nothing as int is not reference but constant
	case *tree.DFloat:
		// do nothing as float is not reference but constant
	case *tree.DDecimal:
		// do nothing as decimal is not reference but constant
	case *tree.DString:
		// do nothing as string is not reference but constant
	case *tree.DCollatedString:
		// do nothing as string with locale is not reference but constant
	case *tree.DBytes:
		// do nothing as bytes is not reference but constant
	case *tree.DUuid:
		// do nothing as UUID is not reference but constant
	case *tree.DIPAddr:
		// do nothing as ip address is not reference but constant
	case *tree.DDate:
		// do nothing as date is not reference but constant
	case *tree.DTime:
		// do nothing as time is not reference but constant
	case *tree.DTimeTZ:
		// do nothing as time with zone is not reference but constant
	case *tree.DTimestamp:
		// do nothing as timestamp is not reference but constant
	case *tree.DTimestampTZ:
		// do nothing as timestamp with zone is not reference but constant
	case *tree.DInterval:
		// do nothing as interval is not reference but constant
	case *tree.DJSON:
		// do nothing as json is not reference but constant
	case *tree.DTuple:
		for _, d := range e.D {
			res = append(res, sc.collectExprReferences(d, defaultType)...)
		}
	case *tree.DArray:
		for _, d := range e.Array {
			res = append(res, sc.collectExprReferences(d, defaultType)...)
		}
	case *tree.DOid:
		// do nothing as oid is not reference but constant
	case *tree.DOidWrapper:
		res = append(res, sc.collectExprReferences(e.Wrapped, defaultType)...)
	case *tree.AndExpr:
		res = append(res, sc.collectExprReferences(e.Left, defaultType)...)
		res = append(res, sc.collectExprReferences(e.Right, defaultType)...)
	case *tree.OrExpr:
		res = append(res, sc.collectExprReferences(e.Left, defaultType)...)
		res = append(res, sc.collectExprReferences(e.Right, defaultType)...)
	case *tree.NotExpr:
		res = append(res, sc.collectExprReferences(e.Expr, defaultType)...)
	case *tree.ParenExpr:
		res = append(res, sc.collectExprReferences(e.Expr, defaultType)...)
	case *tree.ComparisonExpr:
		res = append(res, sc.collectExprReferences(e.Left, defaultType)...)
		res = append(res, sc.collectExprReferences(e.Right, defaultType)...)
	case *tree.IsOfTypeExpr:
		// cockroachdb's special expression?
		res = append(res, sc.collectExprReferences(e.Expr, defaultType)...)
	case *tree.IfErrExpr:
		// cockroachdb's special expression?
		res = append(res, sc.collectExprReferences(e.Cond, defaultType)...)
		res = append(res, sc.collectExprReferences(e.Else, defaultType)...)
		res = append(res, sc.collectExprReferences(e.ErrCode, defaultType)...)
	case *tree.IfExpr:
		// cockroachdb's special expression?
		res = append(res, sc.collectExprReferences(e.Cond, defaultType)...)
		res = append(res, sc.collectExprReferences(e.True, defaultType)...)
		res = append(res, sc.collectExprReferences(e.Else, defaultType)...)
	case *tree.NullIfExpr:
		res = append(res, sc.collectExprReferences(e.Expr1, defaultType)...)
		res = append(res, sc.collectExprReferences(e.Expr2, defaultType)...)
	case *tree.CoalesceExpr:
		for _, ex := range e.Exprs {
			res = append(res, sc.collectExprReferences(ex, defaultType)...)
		}
	case tree.DefaultVal:
		sc.lastError = lib.NewInvalidAstError("DEFAULT is not expected")
	case tree.PartitionMaxVal:
		// MAXVALUE appears in partition definition
		sc.lastError = lib.NewInvalidAstError("MAXVALUE is not expected")
	case tree.PartitionMinVal:
		// MAXVALUE appears in partition definition
		sc.lastError = lib.NewInvalidAstError("MINVALUE is not expected")
	case *tree.Placeholder:
		// do nothing as `$1` for prepared statement doesn't relate to column reference
	case *tree.Tuple:
		for _, ex := range e.Exprs {
			res = append(res, sc.collectExprReferences(ex, defaultType)...)
		}
	case *tree.Array:
		for _, ex := range e.Exprs {
			res = append(res, sc.collectExprReferences(ex, defaultType)...)
		}
	case *tree.ArrayFlatten:
		// cockroachdb's special expression?
		res = append(res, sc.collectExprReferences(e.Subquery, defaultType)...)
	case *tree.TypedDummy:
		// do nothing as dummy doesn't relate to column reference
	case *tree.UnaryExpr:
		res = append(res, sc.collectExprReferences(e.Expr, defaultType)...)
	case *tree.FuncExpr:
		for _, arg := range e.Exprs {
			if _, ok := arg.(tree.UnqualifiedStar); ok {
				// do nothing for star inside a function as it doesn't reference anything
			} else {
				res = append(res, sc.collectExprReferences(arg, defaultType)...)
			}
		}
		res = append(res, sc.collectExprReferences(e.Filter, defaultType)...)

		// WindowDef is for "OVER(...)"
		if e.WindowDef != nil {
			for _, p := range e.WindowDef.Partitions {
				res = append(res, sc.collectExprReferences(p, defaultType)...)
			}
			for _, o := range e.WindowDef.OrderBy {
				res = append(res, sc.collectExprReferences(o.Expr, defaultType)...)
			}
			if frame := e.WindowDef.Frame; frame != nil {
				res = append(res, sc.collectExprReferences(frame.Bounds.StartBound.OffsetExpr, defaultType)...)
				res = append(res, sc.collectExprReferences(frame.Bounds.EndBound.OffsetExpr, defaultType)...)
			}
		}

		// cockroachdb's special?
		for _, o := range e.OrderBy {
			res = append(res, sc.collectExprReferences(o.Expr, defaultType)...)
		}
	case *tree.CastExpr:
		res = append(res, sc.collectExprReferences(e.Expr, defaultType)...)
	case *tree.IndirectionExpr:
		// cockroachdb's special?
		res = append(res, sc.collectExprReferences(e.Expr, defaultType)...)
	case *tree.AnnotateTypeExpr:
		// cockroachdb's special?
		res = append(res, sc.collectExprReferences(e.Expr, defaultType)...)
	case *tree.CollateExpr:
		res = append(res, sc.collectExprReferences(e.Expr, defaultType)...)
	case *tree.TupleStar:
		if n, ok := e.Expr.(*tree.UnresolvedName); ok {
			if n.NumParts > 1 {
				sc.lastError = lib.NewInvalidAstError("schema or catalog/db specification at tuple `(Table).schema.col` is not expected")
			} else {
				return []*db_models.FieldColumn{{Table: n.Parts[0], Type: db_models.FieldStar}}
			}
		} else {
			res = append(res, sc.collectExprReferences(e.Expr, defaultType)...)
		}
	case *tree.ColumnAccessExpr:
		if n, ok := e.Expr.(*tree.UnresolvedName); ok {
			if n.NumParts > 1 {
				sc.lastError = lib.NewInvalidAstError("schema or catalog/db specification at tuple `(Table).schema.col` is not expected")
			} else {
				return []*db_models.FieldColumn{{Table: n.Parts[0], Name: e.ColName, Type: defaultType}}
			}
		} else {
			res = append(res, sc.collectExprReferences(e.Expr, defaultType)...)
		}
	case *tree.IndexedVar:
		// todo irrelevant?
	case *tree.UnresolvedName:
		if e.NumParts > 2 {
			sc.lastError = lib.NewUnsupportedError("schema or catalog/db specification is not supported")
		} else {
			if e.Star {
				return []*db_models.FieldColumn{{Table: e.Parts[1], Name: "", Type: db_models.FieldStar}}
			} else {
				return []*db_models.FieldColumn{{Table: e.Parts[1], Name: e.Parts[0], Type: defaultType}}
			}
		}
	case tree.UnqualifiedStar:
		return []*db_models.FieldColumn{{Table: "", Type: db_models.FieldStar}}
	case *tree.AllColumnsSelector:
		if e.TableName.NumParts > 1 {
			sc.lastError = lib.NewUnsupportedError("schema or catalog/db specification is not supported")
		} else {
			return []*db_models.FieldColumn{{Table: e.TableName.Parts[0], Type: db_models.FieldStar}}
		}
	case *tree.ColumnItem:
		if e.TableName.NumParts > 1 {
			sc.lastError = lib.NewUnsupportedError("schema or catalog/db specification is not supported")
		} else {
			return []*db_models.FieldColumn{{Table: e.TableName.Parts[0], Name: e.ColumnName.Normalize()}}
		}
	default:
		t := reflect.TypeOf(e).Name()
		if t == "dNull" {
			// do nothing as NULL doesn't reference column
		} else {
			sc.lastError = lib.NewUnsupportedError(
				fmt.Sprintf("not supporting query found. %s (type: %s)", e.String(), reflect.TypeOf(e).Name()),
			)
		}
		return nil
	}

	return res
}

func (sc *scopeCollector) collectTables(table tree.TableExpr) ([]*db_models.Table, []*db_models.Field) {
	if table == nil {
		return nil, nil
	}

	var tables []*db_models.Table
	var fields []*db_models.Field
	switch t := table.(type) {
	case *tree.AliasedTableExpr:
		switch tt := t.Expr.(type) {
		case *tree.Subquery:
			// ignore the content of subquery as it will be in different scope.Scopes
			tables = []*db_models.Table{{AsName: t.As.Alias.Normalize()}}
		case *tree.UnresolvedObjectName:
			if tt.NumParts > 1 {
				sc.lastError = lib.NewUnsupportedError("schema or catalog/db specification is not supported")
			} else {
				tables = []*db_models.Table{{AsName: t.As.Alias.Normalize(), Name: tt.Parts[0]}}
			}
		case *tree.StatementSource:
			sc.lastError = lib.NewUnsupportedError("referencing source like [table] is not supported")
		case *tree.AliasedTableExpr:
			sc.lastError = lib.NewInvalidAstError("nested AS is unexpected")
		case *tree.ParenTableExpr:
			return sc.collectTables(tt)
		case *tree.JoinTableExpr:
			tables, fields = sc.collectJoinReferences(tt)
		case *tree.RowsFromExpr:
			for _, item := range tt.Items {
				fcs := sc.collectExprReferences(item, db_models.FieldReference)
				fields = append(fields, &db_models.Field{Columns: fcs})
			}
		case *tree.TableName:
			tables = []*db_models.Table{{AsName: t.As.Alias.Normalize(), Name: tt.Table()}}
		case *tree.TableRef:
			sc.lastError = lib.NewUnsupportedError("numeric table reference is not supported")
		}
	case *tree.ParenTableExpr:
		return sc.collectTables(t)
	case *tree.JoinTableExpr:
		tables, fields = sc.collectJoinReferences(t)
	case *tree.RowsFromExpr:
		for _, item := range t.Items {
			fcs := sc.collectExprReferences(item, db_models.FieldReference)
			fields = append(fields, &db_models.Field{Columns: fcs})
		}
	case *tree.Subquery:
		// ignore the content of subquery as it will be in different scope.Scopes
		return nil, nil
	case *tree.StatementSource:
		sc.lastError = lib.NewUnsupportedError("referencing source like [table] is not supported")
	default:
		sc.lastError = lib.NewUnsupportedError(
			fmt.Sprintf("not supporting query for FROM found. (type: %s)", reflect.TypeOf(t).Name()),
		)
	}

	return tables, fields
}

func (sc *scopeCollector) collectJoinReferences(j *tree.JoinTableExpr) ([]*db_models.Table, []*db_models.Field) {
	if j == nil {
		return nil, nil
	}

	tables, fields := sc.collectTables(j.Left)
	rTables, rFields := sc.collectTables(j.Right)
	tables = append(tables, rTables...)
	fields = append(fields, rFields...)

	if j.Cond != nil {
		switch c := j.Cond.(type) {
		case tree.NaturalJoinCond:
			sc.lastError = lib.NewUnsupportedError(
				fmt.Sprintf("not supporting NATURAL JOIN"),
			)
		case *tree.OnJoinCond:
			fields = append(fields, &db_models.Field{Columns: sc.collectExprReferences(c.Expr, db_models.FieldCondition)})
		case *tree.UsingJoinCond:
			// todo: not supporting the same column name?
			for _, col := range c.Cols {
				fields = append(fields, &db_models.Field{Columns: []*db_models.FieldColumn{{Table: "", Name: col.Normalize(), Type: db_models.FieldCondition}}})
			}
		default:
			sc.lastError = lib.NewUnsupportedError(
				fmt.Sprintf("not supporting query for JOIN found. (type: %s)", reflect.TypeOf(c).Name()),
			)
		}
	}

	return tables, fields
}
