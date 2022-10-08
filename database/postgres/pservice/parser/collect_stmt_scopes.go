package parser

import (
	"fmt"
	"github.com/auxten/postgresql-parser/pkg/sql/parser"
	"github.com/auxten/postgresql-parser/pkg/sql/sem/tree"
	"github.com/mrasu/GravityR/database/dservice/dparser"
	"github.com/mrasu/GravityR/lib"
	"reflect"
)

type scopeCollector struct {
	targetSchema string

	scopes []*dparser.StmtScope
	errors []error

	scopeStack *lib.Stack[scopeInfo]
	lastError  error

	foundSelectInfos         map[*tree.Select]*selectInfo
	foundSelectPrefixCounter map[string]int
}

type scopeInfo struct {
	node  interface{}
	scope *dparser.StmtScope
}

type selectInfo struct {
	name    string
	fieldOf *scopeInfo
}

func newScopeCollector(targetSchema string) *scopeCollector {
	return &scopeCollector{
		targetSchema:             targetSchema,
		scopeStack:               lib.NewStack[scopeInfo](),
		foundSelectInfos:         map[*tree.Select]*selectInfo{},
		foundSelectPrefixCounter: map[string]int{},
	}
}

func CollectStmtScopes(stmt *parser.Statement, targetSchema string) ([]*dparser.StmtScope, []error) {
	sc := newScopeCollector(targetSchema)
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
	scope, selectNode := sc.newStmtScopeIfEnterRequired(node)
	if selectNode == nil {
		return false
	}

	stmt, ok := selectNode.Select.(*tree.SelectClause)
	if !ok {
		return false
	}

	sc.scopeStack.Push(&scopeInfo{node: node, scope: scope})

	var foundFields []*dparser.Field
	for _, expr := range stmt.Exprs {
		cols := sc.collectExprReferences(expr.Expr, dparser.FieldReference)
		if err := sc.getAndClearLastError(); err != nil {
			sc.errors = append(sc.errors, err)
			continue
		}

		foundFields = append(foundFields, &dparser.Field{
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
		f, err := sc.createFieldFromExpr(stmt.Where.Expr, dparser.FieldCondition)
		if err == nil {
			foundFields = append(foundFields, f)
		} else {
			sc.errors = append(sc.errors, err)
		}
	}
	if stmt.GroupBy != nil {
		for _, g := range stmt.GroupBy {
			f, err := sc.createFieldFromExpr(g, dparser.FieldReference)
			if err == nil {
				foundFields = append(foundFields, f)
			} else {
				sc.errors = append(sc.errors, err)
			}
		}
	}
	if stmt.Having != nil {
		f, err := sc.createFieldFromExpr(stmt.Having.Expr, dparser.FieldReference)
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
			f, err := sc.createFieldFromExpr(o.Expr, dparser.FieldReference)
			if err == nil {
				foundFields = append(foundFields, f)
			} else {
				sc.errors = append(sc.errors, err)
			}
		}
	}
	if selectNode.Limit != nil {
		f, err := sc.createFieldFromExpr(selectNode.Limit.Count, dparser.FieldReference)
		if err == nil {
			foundFields = append(foundFields, f)
		} else {
			sc.errors = append(sc.errors, err)
		}

		f, err = sc.createFieldFromExpr(selectNode.Limit.Offset, dparser.FieldReference)
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

	if len(sc.scopes) == 0 {
		sc.scopes = append(sc.scopes, scope)
	}

	return false
}

func (sc *scopeCollector) Leave(node interface{}) bool {
	if sc.scopeStack.Top().node == node {
		sc.scopeStack.Pop()
	}

	return false
}

func (sc *scopeCollector) newStmtScopeIfEnterRequired(node interface{}) (*dparser.StmtScope, *tree.Select) {
	if s, ok := node.(*tree.Select); ok {
		if s, ok := sc.foundSelectInfos[s]; ok {
			if s.fieldOf != nil {
				// already visited when node is *tree.Subquery
				return nil, nil
			}
		}

		parent := sc.scopeStack.Top()
		scope := sc.newScope(parent, s)
		return scope, s
	}

	if sub, ok := node.(*tree.Subquery); ok {
		ps, ok := sub.Select.(*tree.ParenSelect)
		if !ok {
			return nil, nil
		}
		sInfo, ok := sc.foundSelectInfos[ps.Select]
		if !ok {
			return nil, nil
		}
		if sInfo.fieldOf == nil {
			// do nothing when subquery is not in fields as it will be visited when node is *tree.Select
			return nil, nil
		}

		parent := sInfo.fieldOf
		scope := &dparser.StmtScope{Name: sInfo.name}
		parent.scope.FieldScopes = append(parent.scope.FieldScopes, scope)

		return scope, ps.Select
	}

	return nil, nil
}

func (sc *scopeCollector) getAndClearLastError() error {
	err := sc.lastError
	sc.lastError = nil
	return err
}

func (sc *scopeCollector) createFieldFromExpr(expr tree.Expr, fieldType dparser.FieldType) (*dparser.Field, error) {
	cols := sc.collectExprReferences(expr, fieldType)
	if err := sc.getAndClearLastError(); err != nil {
		return nil, err
	}

	return &dparser.Field{Columns: cols}, nil
}

func (sc *scopeCollector) collectExprReferences(expr tree.Expr, defaultType dparser.FieldType) []*dparser.FieldColumn {
	if expr == nil {
		return nil
	}

	var res []*dparser.FieldColumn

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
		info, err := sc.newSelectInfo(e, "field")
		if err != nil {
			sc.lastError = err
			break
		}

		info.fieldOf = sc.scopeStack.Top()

		// ignore the content of subquery as it will be in different StmtScopes
		return []*dparser.FieldColumn{{
			Type:          dparser.FieldSubquery,
			ReferenceName: info.name,
		}}
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
			if sc.isSchemaSpecified(n, 1) {
				sc.lastError = lib.NewInvalidAstError("schema or catalog/db specification at tuple `(Table).schema.col` is not expected")
			} else {
				return []*dparser.FieldColumn{{Table: n.Parts[0], Type: dparser.FieldStar}}
			}
		} else {
			res = append(res, sc.collectExprReferences(e.Expr, defaultType)...)
		}
	case *tree.ColumnAccessExpr:
		if n, ok := e.Expr.(*tree.UnresolvedName); ok {
			if sc.isSchemaSpecified(n, 1) {
				sc.lastError = lib.NewInvalidAstError("schema or catalog/db specification at tuple `(Table).schema.col` is not expected")
			} else {
				return []*dparser.FieldColumn{{Table: n.Parts[0], Name: e.ColName, Type: defaultType}}
			}
		} else {
			res = append(res, sc.collectExprReferences(e.Expr, defaultType)...)
		}
	case *tree.IndexedVar:
		// todo irrelevant?
	case *tree.UnresolvedName:
		if sc.isSchemaSpecified(e, 2) {
			sc.lastError = lib.NewUnsupportedError("schema or catalog/db specification is not supported")
		} else {
			if e.Star {
				return []*dparser.FieldColumn{{Table: e.Parts[1], Name: "", Type: dparser.FieldStar}}
			} else {
				return []*dparser.FieldColumn{{Table: e.Parts[1], Name: e.Parts[0], Type: defaultType}}
			}
		}
	case tree.UnqualifiedStar:
		return []*dparser.FieldColumn{{Table: "", Type: dparser.FieldStar}}
	case *tree.AllColumnsSelector:
		if sc.isSchemaSpecified2(e.TableName, 1) {
			sc.lastError = lib.NewUnsupportedError("schema or catalog/db specification is not supported")
		} else {
			return []*dparser.FieldColumn{{Table: e.TableName.Parts[0], Type: dparser.FieldStar}}
		}
	case *tree.ColumnItem:
		if sc.isSchemaSpecified2(e.TableName, 1) {
			sc.lastError = lib.NewUnsupportedError("schema or catalog/db specification is not supported")
		} else {
			return []*dparser.FieldColumn{{Table: e.TableName.Parts[0], Name: e.ColumnName.Normalize()}}
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

func (sc *scopeCollector) collectTables(table tree.TableExpr) ([]*dparser.Table, []*dparser.Field) {
	if table == nil {
		return nil, nil
	}

	var tables []*dparser.Table
	var fields []*dparser.Field
	switch t := table.(type) {
	case *tree.AliasedTableExpr:
		switch tt := t.Expr.(type) {
		case *tree.Subquery:
			info, err := sc.newSelectInfo(tt, "select")
			if err != nil {
				sc.lastError = err
				break
			}
			name := ""
			if info != nil {
				name = info.name
			}

			// ignore the content of subquery as it will be in different ScopeStmt
			tables = []*dparser.Table{{
				AsName:    t.As.Alias.Normalize(),
				Name:      name,
				IsLateral: t.Lateral,
			}}
		case *tree.UnresolvedObjectName:
			if sc.isSchemaSpecified2(tt, 1) {
				sc.lastError = lib.NewUnsupportedError("schema or catalog/db specification is not supported")
			} else {
				tables = []*dparser.Table{{AsName: t.As.Alias.Normalize(), Name: tt.Parts[0]}}
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
				fcs := sc.collectExprReferences(item, dparser.FieldReference)
				fields = append(fields, &dparser.Field{Columns: fcs})
			}
		case *tree.TableName:
			if sc.isSchemaSpecifiedTableName(tt) {
				sc.lastError = lib.NewUnsupportedError("schema or catalog/db specification is not supported")
			} else {
				tables = []*dparser.Table{{AsName: t.As.Alias.Normalize(), Name: tt.Table()}}
			}
		case *tree.TableRef:
			sc.lastError = lib.NewUnsupportedError("numeric table reference is not supported")
		}
	case *tree.ParenTableExpr:
		return sc.collectTables(t)
	case *tree.JoinTableExpr:
		tables, fields = sc.collectJoinReferences(t)
	case *tree.RowsFromExpr:
		for _, item := range t.Items {
			fcs := sc.collectExprReferences(item, dparser.FieldReference)
			fields = append(fields, &dparser.Field{Columns: fcs})
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

func (sc *scopeCollector) collectJoinReferences(j *tree.JoinTableExpr) ([]*dparser.Table, []*dparser.Field) {
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
			sc.lastError = lib.NewUnsupportedError("not supporting NATURAL JOIN")
		case *tree.OnJoinCond:
			fields = append(fields, &dparser.Field{Columns: sc.collectExprReferences(c.Expr, dparser.FieldCondition)})
		case *tree.UsingJoinCond:
			// todo: not supporting the same column name?
			for _, col := range c.Cols {
				fields = append(fields, &dparser.Field{Columns: []*dparser.FieldColumn{{Table: "", Name: col.Normalize(), Type: dparser.FieldCondition}}})
			}
		default:
			sc.lastError = lib.NewUnsupportedError(
				fmt.Sprintf("not supporting query for JOIN found. (type: %s)", reflect.TypeOf(c).Name()),
			)
		}
	}

	return tables, fields
}

func (sc *scopeCollector) newSelectInfo(t *tree.Subquery, prefix string) (*selectInfo, error) {
	switch s := t.Select.(type) {
	case *tree.ParenSelect:
		cnt := sc.incrementFoundPrefixCounter(prefix)
		info := &selectInfo{name: fmt.Sprintf("<%s%d>", prefix, cnt)}
		sc.foundSelectInfos[s.Select] = info
		return info, nil
	case *tree.SelectClause:
		return nil, lib.NewInvalidAstError("a subquery having select clause is unexpected")
	case *tree.UnionClause:
		return nil, lib.NewInvalidAstError("a subquery having union clause is unexpected")
	case *tree.ValuesClause:
		return nil, lib.NewInvalidAstError("a subquery having values clause is unexpected")
	case *tree.ValuesClauseWithNames:
		return nil, lib.NewInvalidAstError("a subquery having values clause is unexpected")
	default:
		return nil, lib.NewInvalidAstError("unknown subquery type")
	}
}

func (sc *scopeCollector) incrementFoundPrefixCounter(prefix string) int {
	if _, ok := sc.foundSelectPrefixCounter[prefix]; !ok {
		sc.foundSelectPrefixCounter[prefix] = 0
	}

	orig := sc.foundSelectPrefixCounter[prefix]
	sc.foundSelectPrefixCounter[prefix] = orig + 1

	return orig
}

func (sc *scopeCollector) newScope(parent *scopeInfo, s *tree.Select) *dparser.StmtScope {
	name := dparser.RootScopeName
	if n, ok := sc.foundSelectInfos[s]; ok {
		name = n.name
	}
	scope := &dparser.StmtScope{Name: name}

	if parent != nil {
		parent.scope.SubScopes = append(parent.scope.SubScopes, scope)
	}

	return scope
}

func (sc *scopeCollector) isSchemaSpecified(un *tree.UnresolvedName, schemaPos int) bool {
	if (un.NumParts == schemaPos+1 && un.Parts[schemaPos] != sc.targetSchema) || un.NumParts > schemaPos+1 {
		return true
	}

	return false
}

func (sc *scopeCollector) isSchemaSpecified2(uon *tree.UnresolvedObjectName, schemaPos int) bool {
	if (uon.NumParts == schemaPos+1 && uon.Parts[schemaPos] != sc.targetSchema) || uon.NumParts > schemaPos+1 {
		return true
	}

	return false
}

func (sc *scopeCollector) isSchemaSpecifiedTableName(tt *tree.TableName) bool {
	if tt.CatalogName != "" {
		return true
	}

	sName := tt.SchemaName.Normalize()
	if sName == "" || sName == sc.targetSchema {
		return false
	}

	return true
}
