package builder_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/mrasu/GravityR/database/common_model"
	"github.com/mrasu/GravityR/database/common_model/builder"
	"github.com/mrasu/GravityR/thelper/tdata"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBuildIndexTargets_SingleTable(t *testing.T) {
	tables := []*common_model.TableSchema{
		{
			Name:        "users",
			PrimaryKeys: []string{"pk1", "pk2"},
			Columns: []*common_model.ColumnSchema{
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
		scopeFields          []*common_model.Field
		expectedTargetFields [][]*common_model.IndexField
	}{
		{
			name:        "simple",
			asTableName: "",
			scopeFields: []*common_model.Field{
				{
					Columns: []*common_model.FieldColumn{
						{Table: "users", Name: "name", Type: common_model.FieldCondition},
						{Table: "users", Name: "is_admin", Type: common_model.FieldCondition},
					},
				},
			},
			expectedTargetFields: [][]*common_model.IndexField{
				{
					{Name: "is_admin", Type: common_model.FieldCondition},
				},
				{
					{Name: "name", Type: common_model.FieldCondition},
				},
				{
					{Name: "is_admin", Type: common_model.FieldCondition},
					{Name: "name", Type: common_model.FieldCondition},
				},
				{
					{Name: "name", Type: common_model.FieldCondition},
					{Name: "is_admin", Type: common_model.FieldCondition},
				},
			},
		},
		{
			name:        "no table name in fields",
			asTableName: "",
			scopeFields: []*common_model.Field{
				{
					Columns: []*common_model.FieldColumn{
						{Name: "name", Type: common_model.FieldCondition},
						{Name: "is_admin", Type: common_model.FieldCondition},
					},
				},
			},
			expectedTargetFields: [][]*common_model.IndexField{
				{
					{Name: "is_admin", Type: common_model.FieldCondition},
				},
				{
					{Name: "name", Type: common_model.FieldCondition},
				},
				{
					{Name: "is_admin", Type: common_model.FieldCondition},
					{Name: "name", Type: common_model.FieldCondition},
				},
				{
					{Name: "name", Type: common_model.FieldCondition},
					{Name: "is_admin", Type: common_model.FieldCondition},
				},
			},
		},
		{
			name:        "aliased table",
			asTableName: "u",
			scopeFields: []*common_model.Field{
				{
					Columns: []*common_model.FieldColumn{
						{Table: "u", Name: "name", Type: common_model.FieldCondition},
						{Table: "u", Name: "is_admin", Type: common_model.FieldCondition},
					},
				},
			},
			expectedTargetFields: [][]*common_model.IndexField{
				{
					{Name: "is_admin", Type: common_model.FieldCondition},
				},
				{
					{Name: "name", Type: common_model.FieldCondition},
				},
				{
					{Name: "is_admin", Type: common_model.FieldCondition},
					{Name: "name", Type: common_model.FieldCondition},
				},
				{
					{Name: "name", Type: common_model.FieldCondition},
					{Name: "is_admin", Type: common_model.FieldCondition},
				},
			},
		},
		{
			name: "includes reference column",
			scopeFields: []*common_model.Field{
				{
					Columns: []*common_model.FieldColumn{
						{Table: "users", Name: "name", Type: common_model.FieldReference},
						{Table: "users", Name: "is_admin", Type: common_model.FieldCondition},
					},
				},
			},
			expectedTargetFields: [][]*common_model.IndexField{
				{
					{Name: "is_admin", Type: common_model.FieldCondition},
				},
				{
					{Name: "is_admin", Type: common_model.FieldCondition},
					{Name: "name", Type: common_model.FieldReference},
				},
			},
		},
		{
			name: "includes primary key column",
			scopeFields: []*common_model.Field{
				{
					Columns: []*common_model.FieldColumn{
						{Table: "users", Name: "pk1", Type: common_model.FieldCondition},
						{Table: "users", Name: "pk2", Type: common_model.FieldCondition},
						{Table: "users", Name: "is_admin", Type: common_model.FieldCondition},
					},
				},
			},
			expectedTargetFields: [][]*common_model.IndexField{
				{
					{Name: "is_admin", Type: common_model.FieldCondition},
				},
				{
					{Name: "is_admin", Type: common_model.FieldCondition},
					{Name: "pk1", Type: common_model.FieldCondition},
				},
				{
					{Name: "is_admin", Type: common_model.FieldCondition},
					{Name: "pk2", Type: common_model.FieldCondition},
				},
				{
					{Name: "pk1", Type: common_model.FieldCondition},
					{Name: "is_admin", Type: common_model.FieldCondition},
				},
				{
					{Name: "pk2", Type: common_model.FieldCondition},
					{Name: "is_admin", Type: common_model.FieldCondition},
				},
				{
					{Name: "is_admin", Type: common_model.FieldCondition},
					{Name: "pk1", Type: common_model.FieldCondition},
					{Name: "pk2", Type: common_model.FieldCondition},
				},
				{
					{Name: "is_admin", Type: common_model.FieldCondition},
					{Name: "pk2", Type: common_model.FieldCondition},
					{Name: "pk1", Type: common_model.FieldCondition},
				},
				{
					{Name: "pk1", Type: common_model.FieldCondition},
					{Name: "is_admin", Type: common_model.FieldCondition},
					{Name: "pk2", Type: common_model.FieldCondition},
				},
				{
					{Name: "pk1", Type: common_model.FieldCondition},
					{Name: "pk2", Type: common_model.FieldCondition},
					{Name: "is_admin", Type: common_model.FieldCondition},
				},
				{
					{Name: "pk2", Type: common_model.FieldCondition},
					{Name: "is_admin", Type: common_model.FieldCondition},
					{Name: "pk1", Type: common_model.FieldCondition},
				},
				{
					{Name: "pk2", Type: common_model.FieldCondition},
					{Name: "pk1", Type: common_model.FieldCondition},
					{Name: "is_admin", Type: common_model.FieldCondition},
				},
			},
		},
		{
			name: "includes primary key column reference",
			scopeFields: []*common_model.Field{
				{
					Columns: []*common_model.FieldColumn{
						{Table: "users", Name: "pk1", Type: common_model.FieldReference},
						{Table: "users", Name: "pk2", Type: common_model.FieldReference},
						{Table: "users", Name: "is_admin", Type: common_model.FieldReference},
					},
				},
			},
			expectedTargetFields: [][]*common_model.IndexField{
				{
					{Name: "is_admin", Type: common_model.FieldReference},
					{Name: "pk1", Type: common_model.FieldReference},
					{Name: "pk2", Type: common_model.FieldReference},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := &common_model.StmtScope{
				Name: common_model.RootScopeName,
				Tables: []*common_model.Table{
					{AsName: tt.asTableName, Name: "users"},
				},
				Fields: tt.scopeFields,
			}
			actualTargets, err := builder.BuildIndexTargets(tables, []*common_model.StmtScope{scope})
			require.NoError(t, err)

			var expected []*common_model.IndexTargetTable
			for _, fields := range tt.expectedTargetFields {
				expected = append(expected, &common_model.IndexTargetTable{
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
	tables := []*common_model.TableSchema{
		{
			Name:        "users",
			PrimaryKeys: []string{"id"},
			Columns: []*common_model.ColumnSchema{
				{Name: "id"},
				{Name: "name"},
				{Name: "is_admin"},
			},
		},
		{
			Name:        "tasks",
			PrimaryKeys: []string{"id"},
			Columns: []*common_model.ColumnSchema{
				{Name: "id"},
				{Name: "user_id"},
				{Name: "description"},
			},
		},
	}

	tests := []struct {
		name                string
		asTableName         string
		scopeFields         []*common_model.Field
		expectedIndexTarget []*common_model.IndexTargetTable
	}{
		{
			name:        "simple",
			asTableName: "u",
			scopeFields: []*common_model.Field{
				{
					Columns: []*common_model.FieldColumn{
						{Table: "tasks", Name: "id", Type: common_model.FieldCondition},
					},
				},
				{
					Columns: []*common_model.FieldColumn{
						{Table: "u", Name: "id", Type: common_model.FieldCondition},
						{Table: "tasks", Name: "user_id", Type: common_model.FieldCondition},
					},
				},
			},
			expectedIndexTarget: []*common_model.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*common_model.IndexField{
						{Name: "user_id", Type: common_model.FieldCondition},
					},
				},
				{
					TableName: "tasks",
					IndexFields: []*common_model.IndexField{
						{Name: "id", Type: common_model.FieldCondition},
						{Name: "user_id", Type: common_model.FieldCondition},
					},
				},
				{
					TableName: "tasks",
					IndexFields: []*common_model.IndexField{
						{Name: "user_id", Type: common_model.FieldCondition},
						{Name: "id", Type: common_model.FieldCondition},
					},
				},
			},
		},
		{
			name:        "includes reference only column",
			asTableName: "u",
			scopeFields: []*common_model.Field{
				{
					Columns: []*common_model.FieldColumn{
						{Table: "tasks", Name: "description", Type: common_model.FieldReference},
					},
				},
				{
					Columns: []*common_model.FieldColumn{
						{Table: "u", Name: "id", Type: common_model.FieldCondition},
						{Table: "tasks", Name: "user_id", Type: common_model.FieldCondition},
					},
				},
			},
			expectedIndexTarget: []*common_model.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*common_model.IndexField{
						{Name: "user_id", Type: common_model.FieldCondition},
					},
				},
				{
					TableName: "tasks",
					IndexFields: []*common_model.IndexField{
						{Name: "user_id", Type: common_model.FieldCondition},
						{Name: "description", Type: common_model.FieldReference},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := &common_model.StmtScope{
				Name: common_model.RootScopeName,
				Tables: []*common_model.Table{
					{AsName: tt.asTableName, Name: "users"},
					{Name: "tasks"},
				},
				Fields: tt.scopeFields,
			}
			actualTargets, err := builder.BuildIndexTargets(tables, []*common_model.StmtScope{scope})
			require.NoError(t, err)

			if diff := cmp.Diff(tt.expectedIndexTarget, actualTargets); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestBuildIndexTargets_PostgresSubquery(t *testing.T) {
	tables := []*common_model.TableSchema{
		{
			Name:        "users",
			PrimaryKeys: []string{"id"},
			Columns: []*common_model.ColumnSchema{
				{Name: "id"},
				{Name: "name"},
				{Name: "is_admin"},
				{Name: "email"},
			},
		},
		{
			Name:        "tasks",
			PrimaryKeys: []string{"id"},
			Columns: []*common_model.ColumnSchema{
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
		scopes              []*common_model.StmtScope
		expectedIndexTarget []*common_model.IndexTargetTable
	}{
		{
			name:   "subquery in SELECT",
			scopes: tdata.PostgresSubqueryInSelectData.Scopes,
			expectedIndexTarget: []*common_model.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*common_model.IndexField{
						{Name: "user_id", Type: common_model.FieldReference},
					},
				},
			},
		},
		{
			name:   "subquery in SELECT's aliased field",
			scopes: tdata.PostgresSubqueryInSelectAliasedFieldData.Scopes,
			expectedIndexTarget: []*common_model.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*common_model.IndexField{
						{Name: "status", Type: common_model.FieldReference},
					},
				},
				{
					TableName: "users",
					IndexFields: []*common_model.IndexField{
						{Name: "is_admin", Type: common_model.FieldReference},
						{Name: "name", Type: common_model.FieldReference},
					},
				},
			},
		},
		{
			name:   "subquery with comparison",
			scopes: tdata.PostgresSubqueryWithComparisonData.Scopes,
			expectedIndexTarget: []*common_model.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*common_model.IndexField{
						{Name: "description", Type: common_model.FieldReference},
						{Name: "status", Type: common_model.FieldReference},
					},
				},
			},
		},
		{
			name:   "subquery in SELECT's function",
			scopes: tdata.PostgresSubqueryInSelectFunctionData.Scopes,
			expectedIndexTarget: []*common_model.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*common_model.IndexField{
						{Name: "description", Type: common_model.FieldReference},
						{Name: "id", Type: common_model.FieldReference},
					},
				},
			},
		},
		{
			name:   "subquery in FROM",
			scopes: tdata.PostgresSubqueryInFromData.Scopes,
			expectedIndexTarget: []*common_model.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*common_model.IndexField{
						{Name: "id", Type: common_model.FieldReference},
						{Name: "user_id", Type: common_model.FieldReference},
					},
				},
			},
		},
		{
			name:   "subquery in JOIN",
			scopes: tdata.PostgresSubqueryInJoinData.Scopes,
			expectedIndexTarget: []*common_model.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*common_model.IndexField{
						{Name: "user_id", Type: common_model.FieldCondition},
					},
				},
				{
					TableName: "tasks",
					IndexFields: []*common_model.IndexField{
						{Name: "user_id", Type: common_model.FieldCondition},
						{Name: "id", Type: common_model.FieldReference},
					},
				},
				{
					TableName: "users",
					IndexFields: []*common_model.IndexField{
						{Name: "id", Type: common_model.FieldCondition},
						{Name: "name", Type: common_model.FieldReference},
					},
				},
			},
		},
		{
			name:   "subquery using star",
			scopes: tdata.PostgresSubqueryUsingStarData.Scopes,
			expectedIndexTarget: []*common_model.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*common_model.IndexField{
						{Name: "description", Type: common_model.FieldReference},
					},
				},
			},
		},
		{
			name:   "subquery with LATERAL JOIN",
			scopes: tdata.PostgresSubqueryWithLateralJoinData.Scopes,
			expectedIndexTarget: []*common_model.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*common_model.IndexField{
						{Name: "user_id", Type: common_model.FieldCondition},
					},
				},
				{
					TableName: "tasks",
					IndexFields: []*common_model.IndexField{
						{Name: "user_id", Type: common_model.FieldCondition},
						{Name: "description", Type: common_model.FieldReference},
					},
				},
			},
		},
		{
			name:   "subquery with LATERAL object reference",
			scopes: tdata.PostgresSubqueryWithLateralObjectReferenceData.Scopes,
			expectedIndexTarget: []*common_model.IndexTargetTable{
				{
					TableName: "users",
					IndexFields: []*common_model.IndexField{
						{Name: "name", Type: common_model.FieldReference},
					},
				},
			},
		},
		{
			name:   "subquery with aggregate function",
			scopes: tdata.PostgresSubqueryWithAggregateFunctionData.Scopes,
			expectedIndexTarget: []*common_model.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*common_model.IndexField{
						{Name: "description", Type: common_model.FieldReference},
					},
				},
			},
		},
		{
			name:   "correlated subquery",
			scopes: tdata.PostgresCorrelatedSubqueryData.Scopes,
			expectedIndexTarget: []*common_model.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*common_model.IndexField{
						{Name: "user_id", Type: common_model.FieldCondition},
					},
				},
				{
					TableName: "tasks",
					IndexFields: []*common_model.IndexField{
						{Name: "user_id", Type: common_model.FieldCondition},
						{Name: "description", Type: common_model.FieldReference},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualTargets, err := builder.BuildIndexTargets(tables, tt.scopes)
			require.NoError(t, err)

			if diff := cmp.Diff(tt.expectedIndexTarget, actualTargets); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestBuildIndexTargets_HasuraSubquery(t *testing.T) {
	tables := []*common_model.TableSchema{
		{
			Name:        "users",
			PrimaryKeys: []string{"id"},
			Columns: []*common_model.ColumnSchema{
				{Name: "id"},
				{Name: "name"},
				{Name: "is_admin"},
				{Name: "email"},
			},
		},
		{
			Name:        "tasks",
			PrimaryKeys: []string{"id"},
			Columns: []*common_model.ColumnSchema{
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
		scopes              []*common_model.StmtScope
		expectedIndexTarget []*common_model.IndexTargetTable
	}{
		{
			name:   "subquery with where",
			scopes: tdata.HasuraSubqueryWithWhereData.Scopes,
			expectedIndexTarget: []*common_model.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*common_model.IndexField{
						{Name: "user_id", Type: common_model.FieldCondition},
					},
				},
				{
					TableName: "tasks",
					IndexFields: []*common_model.IndexField{
						{Name: "user_id", Type: common_model.FieldCondition},
						{Name: "status", Type: common_model.FieldReference},
						{Name: "title", Type: common_model.FieldReference},
					},
				},
				{
					TableName: "users",
					IndexFields: []*common_model.IndexField{
						{Name: "email", Type: common_model.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*common_model.IndexField{
						{Name: "email", Type: common_model.FieldCondition},
						{Name: "id", Type: common_model.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*common_model.IndexField{
						{Name: "id", Type: common_model.FieldCondition},
						{Name: "email", Type: common_model.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*common_model.IndexField{
						{Name: "email", Type: common_model.FieldCondition},
						{Name: "id", Type: common_model.FieldCondition},
						{Name: "name", Type: common_model.FieldReference},
					},
				},
				{
					TableName: "users",
					IndexFields: []*common_model.IndexField{
						{Name: "id", Type: common_model.FieldCondition},
						{Name: "email", Type: common_model.FieldCondition},
						{Name: "name", Type: common_model.FieldReference},
					},
				},
			},
		},
		{
			name:   "subquery with variables",
			scopes: tdata.HasuraSubqueryWithVariablesData.Scopes,
			expectedIndexTarget: []*common_model.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*common_model.IndexField{
						{Name: "user_id", Type: common_model.FieldCondition},
					},
				},
				{
					TableName: "tasks",
					IndexFields: []*common_model.IndexField{
						{Name: "user_id", Type: common_model.FieldCondition},
						{Name: "description", Type: common_model.FieldReference},
						{Name: "id", Type: common_model.FieldReference},
						{Name: "status", Type: common_model.FieldReference},
					},
				},
				{
					TableName: "users",
					IndexFields: []*common_model.IndexField{
						{Name: "email", Type: common_model.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*common_model.IndexField{
						{Name: "email", Type: common_model.FieldCondition},
						{Name: "id", Type: common_model.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*common_model.IndexField{
						{Name: "id", Type: common_model.FieldCondition},
						{Name: "email", Type: common_model.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*common_model.IndexField{
						{Name: "email", Type: common_model.FieldCondition},
						{Name: "id", Type: common_model.FieldCondition},
						{Name: "name", Type: common_model.FieldReference},
					},
				},
				{
					TableName: "users",
					IndexFields: []*common_model.IndexField{
						{Name: "id", Type: common_model.FieldCondition},
						{Name: "email", Type: common_model.FieldCondition},
						{Name: "name", Type: common_model.FieldReference},
					},
				},
			},
		},
		{
			name:   "subquery with aggregation",
			scopes: tdata.HasuraSubqueryWithAggregationData.Scopes,
			expectedIndexTarget: []*common_model.IndexTargetTable{
				{
					TableName: "tasks",
					IndexFields: []*common_model.IndexField{
						{Name: "user_id", Type: common_model.FieldCondition},
					},
				},
				{
					TableName: "tasks",
					IndexFields: []*common_model.IndexField{
						{Name: "user_id", Type: common_model.FieldCondition},
						{Name: "description", Type: common_model.FieldReference},
						{Name: "status", Type: common_model.FieldReference},
						{Name: "title", Type: common_model.FieldReference},
					},
				},
				{
					TableName: "users",
					IndexFields: []*common_model.IndexField{
						{Name: "email", Type: common_model.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*common_model.IndexField{
						{Name: "email", Type: common_model.FieldCondition},
						{Name: "id", Type: common_model.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*common_model.IndexField{
						{Name: "id", Type: common_model.FieldCondition},
						{Name: "email", Type: common_model.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*common_model.IndexField{
						{Name: "email", Type: common_model.FieldCondition},
						{Name: "id", Type: common_model.FieldCondition},
						{Name: "name", Type: common_model.FieldReference},
					},
				},
				{
					TableName: "users",
					IndexFields: []*common_model.IndexField{
						{Name: "id", Type: common_model.FieldCondition},
						{Name: "email", Type: common_model.FieldCondition},
						{Name: "name", Type: common_model.FieldReference},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualTargets, err := builder.BuildIndexTargets(tables, tt.scopes)
			require.NoError(t, err)

			if diff := cmp.Diff(tt.expectedIndexTarget, actualTargets); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
