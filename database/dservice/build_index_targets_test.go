package dservice_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/mrasu/GravityR/database/dmodel"
	"github.com/mrasu/GravityR/database/dservice"
	"github.com/mrasu/GravityR/thelper/tdata"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBuildIndexTargets_SingleTable(t *testing.T) {
	tables := []*dmodel.TableSchema{
		{
			Name:        "users",
			PrimaryKeys: []string{"pk1", "pk2"},
			Columns: []*dmodel.ColumnSchema{
				{Name: "pk1"},
				{Name: "pk2"},
				{Name: "name"},
				{Name: "is_admin"},
			},
		},
	}

	tests := []struct {
		name                 string
		asTableName          string
		scopeFields          []*dmodel.Field
		expectedTargetFields [][]*dmodel.IndexField
	}{
		{
			name:        "simple",
			asTableName: "",
			scopeFields: []*dmodel.Field{
				{
					Columns: []*dmodel.FieldColumn{
						{Table: "users", Name: "name", Type: dmodel.FieldCondition},
						{Table: "users", Name: "is_admin", Type: dmodel.FieldCondition},
					},
				},
			},
			expectedTargetFields: [][]*dmodel.IndexField{
				{
					{Name: "is_admin", Type: dmodel.FieldCondition},
				},
				{
					{Name: "name", Type: dmodel.FieldCondition},
				},
				{
					{Name: "is_admin", Type: dmodel.FieldCondition},
					{Name: "name", Type: dmodel.FieldCondition},
				},
				{
					{Name: "name", Type: dmodel.FieldCondition},
					{Name: "is_admin", Type: dmodel.FieldCondition},
				},
			},
		},
		{
			name:        "no table name in fields",
			asTableName: "",
			scopeFields: []*dmodel.Field{
				{
					Columns: []*dmodel.FieldColumn{
						{Name: "name", Type: dmodel.FieldCondition},
						{Name: "is_admin", Type: dmodel.FieldCondition},
					},
				},
			},
			expectedTargetFields: [][]*dmodel.IndexField{
				{
					{Name: "is_admin", Type: dmodel.FieldCondition},
				},
				{
					{Name: "name", Type: dmodel.FieldCondition},
				},
				{
					{Name: "is_admin", Type: dmodel.FieldCondition},
					{Name: "name", Type: dmodel.FieldCondition},
				},
				{
					{Name: "name", Type: dmodel.FieldCondition},
					{Name: "is_admin", Type: dmodel.FieldCondition},
				},
			},
		},
		{
			name:        "aliased table",
			asTableName: "u",
			scopeFields: []*dmodel.Field{
				{
					Columns: []*dmodel.FieldColumn{
						{Table: "u", Name: "name", Type: dmodel.FieldCondition},
						{Table: "u", Name: "is_admin", Type: dmodel.FieldCondition},
					},
				},
			},
			expectedTargetFields: [][]*dmodel.IndexField{
				{
					{Name: "is_admin", Type: dmodel.FieldCondition},
				},
				{
					{Name: "name", Type: dmodel.FieldCondition},
				},
				{
					{Name: "is_admin", Type: dmodel.FieldCondition},
					{Name: "name", Type: dmodel.FieldCondition},
				},
				{
					{Name: "name", Type: dmodel.FieldCondition},
					{Name: "is_admin", Type: dmodel.FieldCondition},
				},
			},
		},
		{
			name: "includes reference column",
			scopeFields: []*dmodel.Field{
				{
					Columns: []*dmodel.FieldColumn{
						{Table: "users", Name: "name", Type: dmodel.FieldReference},
						{Table: "users", Name: "is_admin", Type: dmodel.FieldCondition},
					},
				},
			},
			expectedTargetFields: [][]*dmodel.IndexField{
				{
					{Name: "is_admin", Type: dmodel.FieldCondition},
				},
				{
					{Name: "is_admin", Type: dmodel.FieldCondition},
					{Name: "name", Type: dmodel.FieldReference},
				},
			},
		},
		{
			name: "includes primary key column",
			scopeFields: []*dmodel.Field{
				{
					Columns: []*dmodel.FieldColumn{
						{Table: "users", Name: "pk1", Type: dmodel.FieldCondition},
						{Table: "users", Name: "pk2", Type: dmodel.FieldCondition},
						{Table: "users", Name: "is_admin", Type: dmodel.FieldCondition},
					},
				},
			},
			expectedTargetFields: [][]*dmodel.IndexField{
				{
					{Name: "is_admin", Type: dmodel.FieldCondition},
				},
				{
					{Name: "is_admin", Type: dmodel.FieldCondition},
					{Name: "pk1", Type: dmodel.FieldCondition},
				},
				{
					{Name: "is_admin", Type: dmodel.FieldCondition},
					{Name: "pk2", Type: dmodel.FieldCondition},
				},
				{
					{Name: "pk1", Type: dmodel.FieldCondition},
					{Name: "is_admin", Type: dmodel.FieldCondition},
				},
				{
					{Name: "pk2", Type: dmodel.FieldCondition},
					{Name: "is_admin", Type: dmodel.FieldCondition},
				},
				{
					{Name: "is_admin", Type: dmodel.FieldCondition},
					{Name: "pk1", Type: dmodel.FieldCondition},
					{Name: "pk2", Type: dmodel.FieldCondition},
				},
				{
					{Name: "is_admin", Type: dmodel.FieldCondition},
					{Name: "pk2", Type: dmodel.FieldCondition},
					{Name: "pk1", Type: dmodel.FieldCondition},
				},
				{
					{Name: "pk1", Type: dmodel.FieldCondition},
					{Name: "is_admin", Type: dmodel.FieldCondition},
					{Name: "pk2", Type: dmodel.FieldCondition},
				},
				{
					{Name: "pk1", Type: dmodel.FieldCondition},
					{Name: "pk2", Type: dmodel.FieldCondition},
					{Name: "is_admin", Type: dmodel.FieldCondition},
				},
				{
					{Name: "pk2", Type: dmodel.FieldCondition},
					{Name: "is_admin", Type: dmodel.FieldCondition},
					{Name: "pk1", Type: dmodel.FieldCondition},
				},
				{
					{Name: "pk2", Type: dmodel.FieldCondition},
					{Name: "pk1", Type: dmodel.FieldCondition},
					{Name: "is_admin", Type: dmodel.FieldCondition},
				},
			},
		},
		{
			name: "includes primary key column reference",
			scopeFields: []*dmodel.Field{
				{
					Columns: []*dmodel.FieldColumn{
						{Table: "users", Name: "pk1", Type: dmodel.FieldReference},
						{Table: "users", Name: "pk2", Type: dmodel.FieldReference},
						{Table: "users", Name: "is_admin", Type: dmodel.FieldReference},
					},
				},
			},
			expectedTargetFields: [][]*dmodel.IndexField{
				{
					{Name: "is_admin", Type: dmodel.FieldReference},
					{Name: "pk1", Type: dmodel.FieldReference},
					{Name: "pk2", Type: dmodel.FieldReference},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := &dmodel.StmtScope{
				Name: dmodel.RootScopeName,
				Tables: []*dmodel.Table{
					{AsName: tt.asTableName, Name: "users"},
				},
				Fields: tt.scopeFields,
			}
			actualTargets, err := dservice.BuildIndexTargets(tables, []*dmodel.StmtScope{scope})
			require.NoError(t, err)

			var expected []*dmodel.IndexTargetTable
			for _, fields := range tt.expectedTargetFields {
				expected = append(expected, &dmodel.IndexTargetTable{
					TableName:   "users",
					IndexFields: fields,
				})
			}

			if diff := cmp.Diff(expected, actualTargets); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestBuildIndexTargets_MultipleTables(t *testing.T) {
	tables := []*dmodel.TableSchema{
		{
			Name:        "users",
			PrimaryKeys: []string{"id"},
			Columns: []*dmodel.ColumnSchema{
				{Name: "id"},
				{Name: "name"},
				{Name: "is_admin"},
			},
		},
		{
			Name:        "tasks",
			PrimaryKeys: []string{"id"},
			Columns: []*dmodel.ColumnSchema{
				{Name: "id"},
				{Name: "user_id"},
				{Name: "description"},
			},
		},
	}

	tests := []struct {
		name                string
		asTableName         string
		scopeFields         []*dmodel.Field
		expectedIndexTarget []*dmodel.IndexTargetTable
	}{
		{
			name:        "simple",
			asTableName: "u",
			scopeFields: []*dmodel.Field{
				{
					Columns: []*dmodel.FieldColumn{
						{Table: "tasks", Name: "id", Type: dmodel.FieldCondition},
					},
				},
				{
					Columns: []*dmodel.FieldColumn{
						{Table: "u", Name: "id", Type: dmodel.FieldCondition},
						{Table: "tasks", Name: "user_id", Type: dmodel.FieldCondition},
					},
				},
			},
			expectedIndexTarget: []*dmodel.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*dmodel.IndexField{
						{Name: "user_id", Type: dmodel.FieldCondition},
					},
				},
				{
					TableName: "tasks",
					IndexFields: []*dmodel.IndexField{
						{Name: "id", Type: dmodel.FieldCondition},
						{Name: "user_id", Type: dmodel.FieldCondition},
					},
				},
				{
					TableName: "tasks",
					IndexFields: []*dmodel.IndexField{
						{Name: "user_id", Type: dmodel.FieldCondition},
						{Name: "id", Type: dmodel.FieldCondition},
					},
				},
			},
		},
		{
			name:        "includes reference only column",
			asTableName: "u",
			scopeFields: []*dmodel.Field{
				{
					Columns: []*dmodel.FieldColumn{
						{Table: "tasks", Name: "description", Type: dmodel.FieldReference},
					},
				},
				{
					Columns: []*dmodel.FieldColumn{
						{Table: "u", Name: "id", Type: dmodel.FieldCondition},
						{Table: "tasks", Name: "user_id", Type: dmodel.FieldCondition},
					},
				},
			},
			expectedIndexTarget: []*dmodel.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*dmodel.IndexField{
						{Name: "user_id", Type: dmodel.FieldCondition},
					},
				},
				{
					TableName: "tasks",
					IndexFields: []*dmodel.IndexField{
						{Name: "user_id", Type: dmodel.FieldCondition},
						{Name: "description", Type: dmodel.FieldReference},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := &dmodel.StmtScope{
				Name: dmodel.RootScopeName,
				Tables: []*dmodel.Table{
					{AsName: tt.asTableName, Name: "users"},
					{Name: "tasks"},
				},
				Fields: tt.scopeFields,
			}
			actualTargets, err := dservice.BuildIndexTargets(tables, []*dmodel.StmtScope{scope})
			require.NoError(t, err)

			if diff := cmp.Diff(tt.expectedIndexTarget, actualTargets); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestBuildIndexTargets_PostgresSubquery(t *testing.T) {
	tables := []*dmodel.TableSchema{
		{
			Name:        "users",
			PrimaryKeys: []string{"id"},
			Columns: []*dmodel.ColumnSchema{
				{Name: "id"},
				{Name: "name"},
				{Name: "is_admin"},
				{Name: "email"},
			},
		},
		{
			Name:        "tasks",
			PrimaryKeys: []string{"id"},
			Columns: []*dmodel.ColumnSchema{
				{Name: "id"},
				{Name: "user_id"},
				{Name: "title"},
				{Name: "status"},
				{Name: "description"},
			},
		},
	}

	tests := []struct {
		name                string
		scopes              []*dmodel.StmtScope
		expectedIndexTarget []*dmodel.IndexTargetTable
	}{
		{
			name:   "subquery in SELECT",
			scopes: tdata.PostgresSubqueryInSelectData.Scopes,
			expectedIndexTarget: []*dmodel.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*dmodel.IndexField{
						{Name: "user_id", Type: dmodel.FieldReference},
					},
				},
			},
		},
		{
			name:   "subquery in SELECT's aliased field",
			scopes: tdata.PostgresSubqueryInSelectAliasedFieldData.Scopes,
			expectedIndexTarget: []*dmodel.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*dmodel.IndexField{
						{Name: "status", Type: dmodel.FieldReference},
					},
				},
				{
					TableName: "users",
					IndexFields: []*dmodel.IndexField{
						{Name: "is_admin", Type: dmodel.FieldReference},
						{Name: "name", Type: dmodel.FieldReference},
					},
				},
			},
		},
		{
			name:   "subquery with comparison",
			scopes: tdata.PostgresSubqueryWithComparisonData.Scopes,
			expectedIndexTarget: []*dmodel.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*dmodel.IndexField{
						{Name: "description", Type: dmodel.FieldReference},
						{Name: "status", Type: dmodel.FieldReference},
					},
				},
			},
		},
		{
			name:   "subquery in SELECT's function",
			scopes: tdata.PostgresSubqueryInSelectFunctionData.Scopes,
			expectedIndexTarget: []*dmodel.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*dmodel.IndexField{
						{Name: "description", Type: dmodel.FieldReference},
						{Name: "id", Type: dmodel.FieldReference},
					},
				},
			},
		},
		{
			name:   "subquery in FROM",
			scopes: tdata.PostgresSubqueryInFromData.Scopes,
			expectedIndexTarget: []*dmodel.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*dmodel.IndexField{
						{Name: "id", Type: dmodel.FieldReference},
						{Name: "user_id", Type: dmodel.FieldReference},
					},
				},
			},
		},
		{
			name:   "subquery in JOIN",
			scopes: tdata.PostgresSubqueryInJoinData.Scopes,
			expectedIndexTarget: []*dmodel.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*dmodel.IndexField{
						{Name: "user_id", Type: dmodel.FieldCondition},
					},
				},
				{
					TableName: "tasks",
					IndexFields: []*dmodel.IndexField{
						{Name: "user_id", Type: dmodel.FieldCondition},
						{Name: "id", Type: dmodel.FieldReference},
					},
				},
				{
					TableName: "users",
					IndexFields: []*dmodel.IndexField{
						{Name: "id", Type: dmodel.FieldCondition},
						{Name: "name", Type: dmodel.FieldReference},
					},
				},
			},
		},
		{
			name:   "subquery using star",
			scopes: tdata.PostgresSubqueryUsingStarData.Scopes,
			expectedIndexTarget: []*dmodel.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*dmodel.IndexField{
						{Name: "description", Type: dmodel.FieldReference},
					},
				},
			},
		},
		{
			name:   "subquery with LATERAL JOIN",
			scopes: tdata.PostgresSubqueryWithLateralJoinData.Scopes,
			expectedIndexTarget: []*dmodel.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*dmodel.IndexField{
						{Name: "user_id", Type: dmodel.FieldCondition},
					},
				},
				{
					TableName: "tasks",
					IndexFields: []*dmodel.IndexField{
						{Name: "user_id", Type: dmodel.FieldCondition},
						{Name: "description", Type: dmodel.FieldReference},
					},
				},
			},
		},
		{
			name:   "subquery with LATERAL object reference",
			scopes: tdata.PostgresSubqueryWithLateralObjectReferenceData.Scopes,
			expectedIndexTarget: []*dmodel.IndexTargetTable{
				{
					TableName: "users",
					IndexFields: []*dmodel.IndexField{
						{Name: "name", Type: dmodel.FieldReference},
					},
				},
			},
		},
		{
			name:   "subquery with aggregate function",
			scopes: tdata.PostgresSubqueryWithAggregateFunctionData.Scopes,
			expectedIndexTarget: []*dmodel.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*dmodel.IndexField{
						{Name: "description", Type: dmodel.FieldReference},
					},
				},
			},
		},
		{
			name:   "correlated subquery",
			scopes: tdata.PostgresCorrelatedSubqueryData.Scopes,
			expectedIndexTarget: []*dmodel.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*dmodel.IndexField{
						{Name: "user_id", Type: dmodel.FieldCondition},
					},
				},
				{
					TableName: "tasks",
					IndexFields: []*dmodel.IndexField{
						{Name: "user_id", Type: dmodel.FieldCondition},
						{Name: "description", Type: dmodel.FieldReference},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualTargets, err := dservice.BuildIndexTargets(tables, tt.scopes)
			require.NoError(t, err)

			if diff := cmp.Diff(tt.expectedIndexTarget, actualTargets); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestBuildIndexTargets_HasuraSubquery(t *testing.T) {
	tables := []*dmodel.TableSchema{
		{
			Name:        "users",
			PrimaryKeys: []string{"id"},
			Columns: []*dmodel.ColumnSchema{
				{Name: "id"},
				{Name: "name"},
				{Name: "is_admin"},
				{Name: "email"},
			},
		},
		{
			Name:        "tasks",
			PrimaryKeys: []string{"id"},
			Columns: []*dmodel.ColumnSchema{
				{Name: "id"},
				{Name: "user_id"},
				{Name: "title"},
				{Name: "status"},
				{Name: "description"},
			},
		},
	}

	tests := []struct {
		name                string
		scopes              []*dmodel.StmtScope
		expectedIndexTarget []*dmodel.IndexTargetTable
	}{
		{
			name:   "subquery with where",
			scopes: tdata.HasuraSubqueryWithWhereData.Scopes,
			expectedIndexTarget: []*dmodel.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*dmodel.IndexField{
						{Name: "user_id", Type: dmodel.FieldCondition},
					},
				},
				{
					TableName: "tasks",
					IndexFields: []*dmodel.IndexField{
						{Name: "user_id", Type: dmodel.FieldCondition},
						{Name: "status", Type: dmodel.FieldReference},
						{Name: "title", Type: dmodel.FieldReference},
					},
				},
				{
					TableName: "users",
					IndexFields: []*dmodel.IndexField{
						{Name: "email", Type: dmodel.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*dmodel.IndexField{
						{Name: "email", Type: dmodel.FieldCondition},
						{Name: "id", Type: dmodel.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*dmodel.IndexField{
						{Name: "id", Type: dmodel.FieldCondition},
						{Name: "email", Type: dmodel.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*dmodel.IndexField{
						{Name: "email", Type: dmodel.FieldCondition},
						{Name: "id", Type: dmodel.FieldCondition},
						{Name: "name", Type: dmodel.FieldReference},
					},
				},
				{
					TableName: "users",
					IndexFields: []*dmodel.IndexField{
						{Name: "id", Type: dmodel.FieldCondition},
						{Name: "email", Type: dmodel.FieldCondition},
						{Name: "name", Type: dmodel.FieldReference},
					},
				},
			},
		},
		{
			name:   "subquery with variables",
			scopes: tdata.HasuraSubqueryWithVariablesData.Scopes,
			expectedIndexTarget: []*dmodel.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*dmodel.IndexField{
						{Name: "user_id", Type: dmodel.FieldCondition},
					},
				},
				{
					TableName: "tasks",
					IndexFields: []*dmodel.IndexField{
						{Name: "user_id", Type: dmodel.FieldCondition},
						{Name: "description", Type: dmodel.FieldReference},
						{Name: "id", Type: dmodel.FieldReference},
						{Name: "status", Type: dmodel.FieldReference},
					},
				},
				{
					TableName: "users",
					IndexFields: []*dmodel.IndexField{
						{Name: "email", Type: dmodel.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*dmodel.IndexField{
						{Name: "email", Type: dmodel.FieldCondition},
						{Name: "id", Type: dmodel.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*dmodel.IndexField{
						{Name: "id", Type: dmodel.FieldCondition},
						{Name: "email", Type: dmodel.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*dmodel.IndexField{
						{Name: "email", Type: dmodel.FieldCondition},
						{Name: "id", Type: dmodel.FieldCondition},
						{Name: "name", Type: dmodel.FieldReference},
					},
				},
				{
					TableName: "users",
					IndexFields: []*dmodel.IndexField{
						{Name: "id", Type: dmodel.FieldCondition},
						{Name: "email", Type: dmodel.FieldCondition},
						{Name: "name", Type: dmodel.FieldReference},
					},
				},
			},
		},
		{
			name:   "subquery with aggregation",
			scopes: tdata.HasuraSubqueryWithAggregationData.Scopes,
			expectedIndexTarget: []*dmodel.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*dmodel.IndexField{
						{Name: "user_id", Type: dmodel.FieldCondition},
					},
				},
				{
					TableName: "tasks",
					IndexFields: []*dmodel.IndexField{
						{Name: "user_id", Type: dmodel.FieldCondition},
						{Name: "description", Type: dmodel.FieldReference},
						{Name: "status", Type: dmodel.FieldReference},
						{Name: "title", Type: dmodel.FieldReference},
					},
				},
				{
					TableName: "users",
					IndexFields: []*dmodel.IndexField{
						{Name: "email", Type: dmodel.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*dmodel.IndexField{
						{Name: "email", Type: dmodel.FieldCondition},
						{Name: "id", Type: dmodel.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*dmodel.IndexField{
						{Name: "id", Type: dmodel.FieldCondition},
						{Name: "email", Type: dmodel.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*dmodel.IndexField{
						{Name: "email", Type: dmodel.FieldCondition},
						{Name: "id", Type: dmodel.FieldCondition},
						{Name: "name", Type: dmodel.FieldReference},
					},
				},
				{
					TableName: "users",
					IndexFields: []*dmodel.IndexField{
						{Name: "id", Type: dmodel.FieldCondition},
						{Name: "email", Type: dmodel.FieldCondition},
						{Name: "name", Type: dmodel.FieldReference},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualTargets, err := dservice.BuildIndexTargets(tables, tt.scopes)
			require.NoError(t, err)

			if diff := cmp.Diff(tt.expectedIndexTarget, actualTargets); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
