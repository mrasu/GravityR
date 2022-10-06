package service_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/mrasu/GravityR/database"
	"github.com/mrasu/GravityR/database/service"
	"github.com/mrasu/GravityR/thelper/tdata"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBuildIndexTargets_SingleTable(t *testing.T) {
	tables := []*database.TableSchema{
		{
			Name:        "users",
			PrimaryKeys: []string{"pk1", "pk2"},
			Columns: []*database.ColumnSchema{
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
		scopeFields          []*database.Field
		expectedTargetFields [][]*database.IndexField
	}{
		{
			name:        "simple",
			asTableName: "",
			scopeFields: []*database.Field{
				{
					Columns: []*database.FieldColumn{
						{Table: "users", Name: "name", Type: database.FieldCondition},
						{Table: "users", Name: "is_admin", Type: database.FieldCondition},
					},
				},
			},
			expectedTargetFields: [][]*database.IndexField{
				{
					{Name: "is_admin", Type: database.FieldCondition},
				},
				{
					{Name: "name", Type: database.FieldCondition},
				},
				{
					{Name: "is_admin", Type: database.FieldCondition},
					{Name: "name", Type: database.FieldCondition},
				},
				{
					{Name: "name", Type: database.FieldCondition},
					{Name: "is_admin", Type: database.FieldCondition},
				},
			},
		},
		{
			name:        "no table name in fields",
			asTableName: "",
			scopeFields: []*database.Field{
				{
					Columns: []*database.FieldColumn{
						{Name: "name", Type: database.FieldCondition},
						{Name: "is_admin", Type: database.FieldCondition},
					},
				},
			},
			expectedTargetFields: [][]*database.IndexField{
				{
					{Name: "is_admin", Type: database.FieldCondition},
				},
				{
					{Name: "name", Type: database.FieldCondition},
				},
				{
					{Name: "is_admin", Type: database.FieldCondition},
					{Name: "name", Type: database.FieldCondition},
				},
				{
					{Name: "name", Type: database.FieldCondition},
					{Name: "is_admin", Type: database.FieldCondition},
				},
			},
		},
		{
			name:        "aliased table",
			asTableName: "u",
			scopeFields: []*database.Field{
				{
					Columns: []*database.FieldColumn{
						{Table: "u", Name: "name", Type: database.FieldCondition},
						{Table: "u", Name: "is_admin", Type: database.FieldCondition},
					},
				},
			},
			expectedTargetFields: [][]*database.IndexField{
				{
					{Name: "is_admin", Type: database.FieldCondition},
				},
				{
					{Name: "name", Type: database.FieldCondition},
				},
				{
					{Name: "is_admin", Type: database.FieldCondition},
					{Name: "name", Type: database.FieldCondition},
				},
				{
					{Name: "name", Type: database.FieldCondition},
					{Name: "is_admin", Type: database.FieldCondition},
				},
			},
		},
		{
			name: "includes reference column",
			scopeFields: []*database.Field{
				{
					Columns: []*database.FieldColumn{
						{Table: "users", Name: "name", Type: database.FieldReference},
						{Table: "users", Name: "is_admin", Type: database.FieldCondition},
					},
				},
			},
			expectedTargetFields: [][]*database.IndexField{
				{
					{Name: "is_admin", Type: database.FieldCondition},
				},
				{
					{Name: "is_admin", Type: database.FieldCondition},
					{Name: "name", Type: database.FieldReference},
				},
			},
		},
		{
			name: "includes primary key column",
			scopeFields: []*database.Field{
				{
					Columns: []*database.FieldColumn{
						{Table: "users", Name: "pk1", Type: database.FieldCondition},
						{Table: "users", Name: "pk2", Type: database.FieldCondition},
						{Table: "users", Name: "is_admin", Type: database.FieldCondition},
					},
				},
			},
			expectedTargetFields: [][]*database.IndexField{
				{
					{Name: "is_admin", Type: database.FieldCondition},
				},
				{
					{Name: "is_admin", Type: database.FieldCondition},
					{Name: "pk1", Type: database.FieldCondition},
				},
				{
					{Name: "is_admin", Type: database.FieldCondition},
					{Name: "pk2", Type: database.FieldCondition},
				},
				{
					{Name: "pk1", Type: database.FieldCondition},
					{Name: "is_admin", Type: database.FieldCondition},
				},
				{
					{Name: "pk2", Type: database.FieldCondition},
					{Name: "is_admin", Type: database.FieldCondition},
				},
				{
					{Name: "is_admin", Type: database.FieldCondition},
					{Name: "pk1", Type: database.FieldCondition},
					{Name: "pk2", Type: database.FieldCondition},
				},
				{
					{Name: "is_admin", Type: database.FieldCondition},
					{Name: "pk2", Type: database.FieldCondition},
					{Name: "pk1", Type: database.FieldCondition},
				},
				{
					{Name: "pk1", Type: database.FieldCondition},
					{Name: "is_admin", Type: database.FieldCondition},
					{Name: "pk2", Type: database.FieldCondition},
				},
				{
					{Name: "pk1", Type: database.FieldCondition},
					{Name: "pk2", Type: database.FieldCondition},
					{Name: "is_admin", Type: database.FieldCondition},
				},
				{
					{Name: "pk2", Type: database.FieldCondition},
					{Name: "is_admin", Type: database.FieldCondition},
					{Name: "pk1", Type: database.FieldCondition},
				},
				{
					{Name: "pk2", Type: database.FieldCondition},
					{Name: "pk1", Type: database.FieldCondition},
					{Name: "is_admin", Type: database.FieldCondition},
				},
			},
		},
		{
			name: "includes primary key column reference",
			scopeFields: []*database.Field{
				{
					Columns: []*database.FieldColumn{
						{Table: "users", Name: "pk1", Type: database.FieldReference},
						{Table: "users", Name: "pk2", Type: database.FieldReference},
						{Table: "users", Name: "is_admin", Type: database.FieldReference},
					},
				},
			},
			expectedTargetFields: [][]*database.IndexField{
				{
					{Name: "is_admin", Type: database.FieldReference},
					{Name: "pk1", Type: database.FieldReference},
					{Name: "pk2", Type: database.FieldReference},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := &database.StmtScope{
				Name: database.RootScopeName,
				Tables: []*database.Table{
					{AsName: tt.asTableName, Name: "users"},
				},
				Fields: tt.scopeFields,
			}
			actualTargets, err := service.BuildIndexTargets(tables, []*database.StmtScope{scope})
			require.NoError(t, err)

			var expected []*database.IndexTargetTable
			for _, fields := range tt.expectedTargetFields {
				expected = append(expected, &database.IndexTargetTable{
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
	tables := []*database.TableSchema{
		{
			Name:        "users",
			PrimaryKeys: []string{"id"},
			Columns: []*database.ColumnSchema{
				{Name: "id"},
				{Name: "name"},
				{Name: "is_admin"},
			},
		},
		{
			Name:        "tasks",
			PrimaryKeys: []string{"id"},
			Columns: []*database.ColumnSchema{
				{Name: "id"},
				{Name: "user_id"},
				{Name: "description"},
			},
		},
	}

	tests := []struct {
		name                string
		asTableName         string
		scopeFields         []*database.Field
		expectedIndexTarget []*database.IndexTargetTable
	}{
		{
			name:        "simple",
			asTableName: "u",
			scopeFields: []*database.Field{
				{
					Columns: []*database.FieldColumn{
						{Table: "tasks", Name: "id", Type: database.FieldCondition},
					},
				},
				{
					Columns: []*database.FieldColumn{
						{Table: "u", Name: "id", Type: database.FieldCondition},
						{Table: "tasks", Name: "user_id", Type: database.FieldCondition},
					},
				},
			},
			expectedIndexTarget: []*database.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*database.IndexField{
						{Name: "user_id", Type: database.FieldCondition},
					},
				},
				{
					TableName: "tasks",
					IndexFields: []*database.IndexField{
						{Name: "id", Type: database.FieldCondition},
						{Name: "user_id", Type: database.FieldCondition},
					},
				},
				{
					TableName: "tasks",
					IndexFields: []*database.IndexField{
						{Name: "user_id", Type: database.FieldCondition},
						{Name: "id", Type: database.FieldCondition},
					},
				},
			},
		},
		{
			name:        "includes reference only column",
			asTableName: "u",
			scopeFields: []*database.Field{
				{
					Columns: []*database.FieldColumn{
						{Table: "tasks", Name: "description", Type: database.FieldReference},
					},
				},
				{
					Columns: []*database.FieldColumn{
						{Table: "u", Name: "id", Type: database.FieldCondition},
						{Table: "tasks", Name: "user_id", Type: database.FieldCondition},
					},
				},
			},
			expectedIndexTarget: []*database.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*database.IndexField{
						{Name: "user_id", Type: database.FieldCondition},
					},
				},
				{
					TableName: "tasks",
					IndexFields: []*database.IndexField{
						{Name: "user_id", Type: database.FieldCondition},
						{Name: "description", Type: database.FieldReference},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := &database.StmtScope{
				Name: database.RootScopeName,
				Tables: []*database.Table{
					{AsName: tt.asTableName, Name: "users"},
					{Name: "tasks"},
				},
				Fields: tt.scopeFields,
			}
			actualTargets, err := service.BuildIndexTargets(tables, []*database.StmtScope{scope})
			require.NoError(t, err)

			if diff := cmp.Diff(tt.expectedIndexTarget, actualTargets); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestBuildIndexTargets_PostgresSubquery(t *testing.T) {
	tables := []*database.TableSchema{
		{
			Name:        "users",
			PrimaryKeys: []string{"id"},
			Columns: []*database.ColumnSchema{
				{Name: "id"},
				{Name: "name"},
				{Name: "is_admin"},
				{Name: "email"},
			},
		},
		{
			Name:        "tasks",
			PrimaryKeys: []string{"id"},
			Columns: []*database.ColumnSchema{
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
		scopes              []*database.StmtScope
		expectedIndexTarget []*database.IndexTargetTable
	}{
		{
			name:   "subquery in SELECT",
			scopes: tdata.PostgresSubqueryInSelectData.Scopes,
			expectedIndexTarget: []*database.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*database.IndexField{
						{Name: "user_id", Type: database.FieldReference},
					},
				},
			},
		},
		{
			name:   "subquery in SELECT's aliased field",
			scopes: tdata.PostgresSubqueryInSelectAliasedFieldData.Scopes,
			expectedIndexTarget: []*database.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*database.IndexField{
						{Name: "status", Type: database.FieldReference},
					},
				},
				{
					TableName: "users",
					IndexFields: []*database.IndexField{
						{Name: "is_admin", Type: database.FieldReference},
						{Name: "name", Type: database.FieldReference},
					},
				},
			},
		},
		{
			name:   "subquery with comparison",
			scopes: tdata.PostgresSubqueryWithComparisonData.Scopes,
			expectedIndexTarget: []*database.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*database.IndexField{
						{Name: "description", Type: database.FieldReference},
						{Name: "status", Type: database.FieldReference},
					},
				},
			},
		},
		{
			name:   "subquery in SELECT's function",
			scopes: tdata.PostgresSubqueryInSelectFunctionData.Scopes,
			expectedIndexTarget: []*database.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*database.IndexField{
						{Name: "description", Type: database.FieldReference},
						{Name: "id", Type: database.FieldReference},
					},
				},
			},
		},
		{
			name:   "subquery in FROM",
			scopes: tdata.PostgresSubqueryInFromData.Scopes,
			expectedIndexTarget: []*database.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*database.IndexField{
						{Name: "id", Type: database.FieldReference},
						{Name: "user_id", Type: database.FieldReference},
					},
				},
			},
		},
		{
			name:   "subquery in JOIN",
			scopes: tdata.PostgresSubqueryInJoinData.Scopes,
			expectedIndexTarget: []*database.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*database.IndexField{
						{Name: "user_id", Type: database.FieldCondition},
					},
				},
				{
					TableName: "tasks",
					IndexFields: []*database.IndexField{
						{Name: "user_id", Type: database.FieldCondition},
						{Name: "id", Type: database.FieldReference},
					},
				},
				{
					TableName: "users",
					IndexFields: []*database.IndexField{
						{Name: "id", Type: database.FieldCondition},
						{Name: "name", Type: database.FieldReference},
					},
				},
			},
		},
		{
			name:   "subquery using star",
			scopes: tdata.PostgresSubqueryUsingStarData.Scopes,
			expectedIndexTarget: []*database.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*database.IndexField{
						{Name: "description", Type: database.FieldReference},
					},
				},
			},
		},
		{
			name:   "subquery with LATERAL JOIN",
			scopes: tdata.PostgresSubqueryWithLateralJoinData.Scopes,
			expectedIndexTarget: []*database.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*database.IndexField{
						{Name: "user_id", Type: database.FieldCondition},
					},
				},
				{
					TableName: "tasks",
					IndexFields: []*database.IndexField{
						{Name: "user_id", Type: database.FieldCondition},
						{Name: "description", Type: database.FieldReference},
					},
				},
			},
		},
		{
			name:   "subquery with LATERAL object reference",
			scopes: tdata.PostgresSubqueryWithLateralObjectReferenceData.Scopes,
			expectedIndexTarget: []*database.IndexTargetTable{
				{
					TableName: "users",
					IndexFields: []*database.IndexField{
						{Name: "name", Type: database.FieldReference},
					},
				},
			},
		},
		{
			name:   "subquery with aggregate function",
			scopes: tdata.PostgresSubqueryWithAggregateFunctionData.Scopes,
			expectedIndexTarget: []*database.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*database.IndexField{
						{Name: "description", Type: database.FieldReference},
					},
				},
			},
		},
		{
			name:   "correlated subquery",
			scopes: tdata.PostgresCorrelatedSubqueryData.Scopes,
			expectedIndexTarget: []*database.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*database.IndexField{
						{Name: "user_id", Type: database.FieldCondition},
					},
				},
				{
					TableName: "tasks",
					IndexFields: []*database.IndexField{
						{Name: "user_id", Type: database.FieldCondition},
						{Name: "description", Type: database.FieldReference},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualTargets, err := service.BuildIndexTargets(tables, tt.scopes)
			require.NoError(t, err)

			if diff := cmp.Diff(tt.expectedIndexTarget, actualTargets); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestBuildIndexTargets_HasuraSubquery(t *testing.T) {
	tables := []*database.TableSchema{
		{
			Name:        "users",
			PrimaryKeys: []string{"id"},
			Columns: []*database.ColumnSchema{
				{Name: "id"},
				{Name: "name"},
				{Name: "is_admin"},
				{Name: "email"},
			},
		},
		{
			Name:        "tasks",
			PrimaryKeys: []string{"id"},
			Columns: []*database.ColumnSchema{
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
		scopes              []*database.StmtScope
		expectedIndexTarget []*database.IndexTargetTable
	}{
		{
			name:   "subquery with where",
			scopes: tdata.HasuraSubqueryWithWhereData.Scopes,
			expectedIndexTarget: []*database.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*database.IndexField{
						{Name: "user_id", Type: database.FieldCondition},
					},
				},
				{
					TableName: "tasks",
					IndexFields: []*database.IndexField{
						{Name: "user_id", Type: database.FieldCondition},
						{Name: "status", Type: database.FieldReference},
						{Name: "title", Type: database.FieldReference},
					},
				},
				{
					TableName: "users",
					IndexFields: []*database.IndexField{
						{Name: "email", Type: database.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*database.IndexField{
						{Name: "email", Type: database.FieldCondition},
						{Name: "id", Type: database.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*database.IndexField{
						{Name: "id", Type: database.FieldCondition},
						{Name: "email", Type: database.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*database.IndexField{
						{Name: "email", Type: database.FieldCondition},
						{Name: "id", Type: database.FieldCondition},
						{Name: "name", Type: database.FieldReference},
					},
				},
				{
					TableName: "users",
					IndexFields: []*database.IndexField{
						{Name: "id", Type: database.FieldCondition},
						{Name: "email", Type: database.FieldCondition},
						{Name: "name", Type: database.FieldReference},
					},
				},
			},
		},
		{
			name:   "subquery with variables",
			scopes: tdata.HasuraSubqueryWithVariablesData.Scopes,
			expectedIndexTarget: []*database.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*database.IndexField{
						{Name: "user_id", Type: database.FieldCondition},
					},
				},
				{
					TableName: "tasks",
					IndexFields: []*database.IndexField{
						{Name: "user_id", Type: database.FieldCondition},
						{Name: "description", Type: database.FieldReference},
						{Name: "id", Type: database.FieldReference},
						{Name: "status", Type: database.FieldReference},
					},
				},
				{
					TableName: "users",
					IndexFields: []*database.IndexField{
						{Name: "email", Type: database.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*database.IndexField{
						{Name: "email", Type: database.FieldCondition},
						{Name: "id", Type: database.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*database.IndexField{
						{Name: "id", Type: database.FieldCondition},
						{Name: "email", Type: database.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*database.IndexField{
						{Name: "email", Type: database.FieldCondition},
						{Name: "id", Type: database.FieldCondition},
						{Name: "name", Type: database.FieldReference},
					},
				},
				{
					TableName: "users",
					IndexFields: []*database.IndexField{
						{Name: "id", Type: database.FieldCondition},
						{Name: "email", Type: database.FieldCondition},
						{Name: "name", Type: database.FieldReference},
					},
				},
			},
		},
		{
			name:   "subquery with aggregation",
			scopes: tdata.HasuraSubqueryWithAggregationData.Scopes,
			expectedIndexTarget: []*database.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*database.IndexField{
						{Name: "user_id", Type: database.FieldCondition},
					},
				},
				{
					TableName: "tasks",
					IndexFields: []*database.IndexField{
						{Name: "user_id", Type: database.FieldCondition},
						{Name: "description", Type: database.FieldReference},
						{Name: "status", Type: database.FieldReference},
						{Name: "title", Type: database.FieldReference},
					},
				},
				{
					TableName: "users",
					IndexFields: []*database.IndexField{
						{Name: "email", Type: database.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*database.IndexField{
						{Name: "email", Type: database.FieldCondition},
						{Name: "id", Type: database.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*database.IndexField{
						{Name: "id", Type: database.FieldCondition},
						{Name: "email", Type: database.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*database.IndexField{
						{Name: "email", Type: database.FieldCondition},
						{Name: "id", Type: database.FieldCondition},
						{Name: "name", Type: database.FieldReference},
					},
				},
				{
					TableName: "users",
					IndexFields: []*database.IndexField{
						{Name: "id", Type: database.FieldCondition},
						{Name: "email", Type: database.FieldCondition},
						{Name: "name", Type: database.FieldReference},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualTargets, err := service.BuildIndexTargets(tables, tt.scopes)
			require.NoError(t, err)

			if diff := cmp.Diff(tt.expectedIndexTarget, actualTargets); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
