package dparser_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/mrasu/GravityR/database/dmodel"
	"github.com/mrasu/GravityR/database/dservice/dparser"
	"github.com/mrasu/GravityR/thelper/tdata"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBuildIndexTargets_SingleTable(t *testing.T) {
	tables := []*dparser.TableSchema{
		{
			Name:        "users",
			PrimaryKeys: []string{"pk1", "pk2"},
			Columns: []*dparser.ColumnSchema{
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
		scopeFields          []*dparser.Field
		expectedTargetFields [][]*dmodel.IndexColumn
	}{
		{
			name:        "simple",
			asTableName: "",
			scopeFields: []*dparser.Field{
				{
					Columns: []*dparser.FieldColumn{
						{Table: "users", Name: "name", Type: dparser.FieldCondition},
						{Table: "users", Name: "is_admin", Type: dparser.FieldCondition},
					},
				},
			},
			expectedTargetFields: [][]*dmodel.IndexColumn{
				{
					{Name: "is_admin"},
				},
				{
					{Name: "name"},
				},
				{
					{Name: "is_admin"},
					{Name: "name"},
				},
				{
					{Name: "name"},
					{Name: "is_admin"},
				},
			},
		},
		{
			name:        "no table name in fields",
			asTableName: "",
			scopeFields: []*dparser.Field{
				{
					Columns: []*dparser.FieldColumn{
						{Name: "name", Type: dparser.FieldCondition},
						{Name: "is_admin", Type: dparser.FieldCondition},
					},
				},
			},
			expectedTargetFields: [][]*dmodel.IndexColumn{
				{
					{Name: "is_admin"},
				},
				{
					{Name: "name"},
				},
				{
					{Name: "is_admin"},
					{Name: "name"},
				},
				{
					{Name: "name"},
					{Name: "is_admin"},
				},
			},
		},
		{
			name:        "aliased table",
			asTableName: "u",
			scopeFields: []*dparser.Field{
				{
					Columns: []*dparser.FieldColumn{
						{Table: "u", Name: "name", Type: dparser.FieldCondition},
						{Table: "u", Name: "is_admin", Type: dparser.FieldCondition},
					},
				},
			},
			expectedTargetFields: [][]*dmodel.IndexColumn{
				{
					{Name: "is_admin"},
				},
				{
					{Name: "name"},
				},
				{
					{Name: "is_admin"},
					{Name: "name"},
				},
				{
					{Name: "name"},
					{Name: "is_admin"},
				},
			},
		},
		{
			name: "includes reference column",
			scopeFields: []*dparser.Field{
				{
					Columns: []*dparser.FieldColumn{
						{Table: "users", Name: "name", Type: dparser.FieldReference},
						{Table: "users", Name: "is_admin", Type: dparser.FieldCondition},
					},
				},
			},
			expectedTargetFields: [][]*dmodel.IndexColumn{
				{
					{Name: "is_admin"},
				},
				{
					{Name: "is_admin"},
					{Name: "name"},
				},
			},
		},
		{
			name: "includes primary key column",
			scopeFields: []*dparser.Field{
				{
					Columns: []*dparser.FieldColumn{
						{Table: "users", Name: "pk1", Type: dparser.FieldCondition},
						{Table: "users", Name: "pk2", Type: dparser.FieldCondition},
						{Table: "users", Name: "is_admin", Type: dparser.FieldCondition},
					},
				},
			},
			expectedTargetFields: [][]*dmodel.IndexColumn{
				{
					{Name: "is_admin"},
				},
				{
					{Name: "is_admin"},
					{Name: "pk1"},
				},
				{
					{Name: "is_admin"},
					{Name: "pk2"},
				},
				{
					{Name: "pk1"},
					{Name: "is_admin"},
				},
				{
					{Name: "pk2"},
					{Name: "is_admin"},
				},
				{
					{Name: "is_admin"},
					{Name: "pk1"},
					{Name: "pk2"},
				},
				{
					{Name: "is_admin"},
					{Name: "pk2"},
					{Name: "pk1"},
				},
				{
					{Name: "pk1"},
					{Name: "is_admin"},
					{Name: "pk2"},
				},
				{
					{Name: "pk1"},
					{Name: "pk2"},
					{Name: "is_admin"},
				},
				{
					{Name: "pk2"},
					{Name: "is_admin"},
					{Name: "pk1"},
				},
				{
					{Name: "pk2"},
					{Name: "pk1"},
					{Name: "is_admin"},
				},
			},
		},
		{
			name: "includes primary key column reference",
			scopeFields: []*dparser.Field{
				{
					Columns: []*dparser.FieldColumn{
						{Table: "users", Name: "pk1", Type: dparser.FieldReference},
						{Table: "users", Name: "pk2", Type: dparser.FieldReference},
						{Table: "users", Name: "is_admin", Type: dparser.FieldReference},
					},
				},
			},
			expectedTargetFields: [][]*dmodel.IndexColumn{
				{
					{Name: "is_admin"},
					{Name: "pk1"},
					{Name: "pk2"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := &dparser.StmtScope{
				Name: dparser.RootScopeName,
				Tables: []*dparser.Table{
					{AsName: tt.asTableName, Name: "users"},
				},
				Fields: tt.scopeFields,
			}
			actualTargets, err := dparser.NewIndexTargetBuilder(tables).Build([]*dparser.StmtScope{scope})
			require.NoError(t, err)

			var expected []*dmodel.IndexTarget
			for _, fields := range tt.expectedTargetFields {
				expected = append(expected, &dmodel.IndexTarget{
					TableName: "users",
					Columns:   fields,
				})
			}

			if diff := cmp.Diff(expected, actualTargets); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestBuildIndexTargets_MultipleTables(t *testing.T) {
	tables := []*dparser.TableSchema{
		{
			Name:        "users",
			PrimaryKeys: []string{"id"},
			Columns: []*dparser.ColumnSchema{
				{Name: "id"},
				{Name: "name"},
				{Name: "is_admin"},
			},
		},
		{
			Name:        "tasks",
			PrimaryKeys: []string{"id"},
			Columns: []*dparser.ColumnSchema{
				{Name: "id"},
				{Name: "user_id"},
				{Name: "description"},
			},
		},
	}

	tests := []struct {
		name                string
		asTableName         string
		scopeFields         []*dparser.Field
		expectedIndexTarget []*dmodel.IndexTarget
	}{
		{
			name:        "simple",
			asTableName: "u",
			scopeFields: []*dparser.Field{
				{
					Columns: []*dparser.FieldColumn{
						{Table: "tasks", Name: "id", Type: dparser.FieldCondition},
					},
				},
				{
					Columns: []*dparser.FieldColumn{
						{Table: "u", Name: "id"},
						{Table: "tasks", Name: "user_id", Type: dparser.FieldCondition},
					},
				},
			},
			expectedIndexTarget: []*dmodel.IndexTarget{
				{
					TableName: "tasks",
					Columns: []*dmodel.IndexColumn{
						{Name: "user_id"},
					},
				},
				{
					TableName: "tasks",
					Columns: []*dmodel.IndexColumn{
						{Name: "id"},
						{Name: "user_id"},
					},
				},
				{
					TableName: "tasks",
					Columns: []*dmodel.IndexColumn{
						{Name: "user_id"},
						{Name: "id"},
					},
				},
			},
		},
		{
			name:        "includes reference only column",
			asTableName: "u",
			scopeFields: []*dparser.Field{
				{
					Columns: []*dparser.FieldColumn{
						{Table: "tasks", Name: "description", Type: dparser.FieldReference},
					},
				},
				{
					Columns: []*dparser.FieldColumn{
						{Table: "u", Name: "id", Type: dparser.FieldCondition},
						{Table: "tasks", Name: "user_id", Type: dparser.FieldCondition},
					},
				},
			},
			expectedIndexTarget: []*dmodel.IndexTarget{
				{
					TableName: "tasks",
					Columns: []*dmodel.IndexColumn{
						{Name: "user_id"},
					},
				},
				{
					TableName: "tasks",
					Columns: []*dmodel.IndexColumn{
						{Name: "user_id"},
						{Name: "description"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := &dparser.StmtScope{
				Name: dparser.RootScopeName,
				Tables: []*dparser.Table{
					{AsName: tt.asTableName, Name: "users"},
					{Name: "tasks"},
				},
				Fields: tt.scopeFields,
			}
			actualTargets, err := dparser.NewIndexTargetBuilder(tables).Build([]*dparser.StmtScope{scope})
			require.NoError(t, err)

			if diff := cmp.Diff(tt.expectedIndexTarget, actualTargets); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestBuildIndexTargets_PostgresSubquery(t *testing.T) {
	tables := []*dparser.TableSchema{
		{
			Name:        "users",
			PrimaryKeys: []string{"id"},
			Columns: []*dparser.ColumnSchema{
				{Name: "id"},
				{Name: "name"},
				{Name: "is_admin"},
				{Name: "email"},
			},
		},
		{
			Name:        "tasks",
			PrimaryKeys: []string{"id"},
			Columns: []*dparser.ColumnSchema{
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
		scopes              []*dparser.StmtScope
		expectedIndexTarget []*dmodel.IndexTarget
	}{
		{
			name:   "subquery in SELECT",
			scopes: tdata.PostgresSubqueryInSelectData.Scopes,
			expectedIndexTarget: []*dmodel.IndexTarget{
				{
					TableName: "tasks",
					Columns: []*dmodel.IndexColumn{
						{Name: "user_id"},
					},
				},
			},
		},
		{
			name:   "subquery in SELECT's aliased field",
			scopes: tdata.PostgresSubqueryInSelectAliasedFieldData.Scopes,
			expectedIndexTarget: []*dmodel.IndexTarget{
				{
					TableName: "tasks",
					Columns: []*dmodel.IndexColumn{
						{Name: "status"},
					},
				},
				{
					TableName: "users",
					Columns: []*dmodel.IndexColumn{
						{Name: "is_admin"},
						{Name: "name"},
					},
				},
			},
		},
		{
			name:   "subquery with comparison",
			scopes: tdata.PostgresSubqueryWithComparisonData.Scopes,
			expectedIndexTarget: []*dmodel.IndexTarget{
				{
					TableName: "tasks",
					Columns: []*dmodel.IndexColumn{
						{Name: "description"},
						{Name: "status"},
					},
				},
			},
		},
		{
			name:   "subquery in SELECT's function",
			scopes: tdata.PostgresSubqueryInSelectFunctionData.Scopes,
			expectedIndexTarget: []*dmodel.IndexTarget{
				{
					TableName: "tasks",
					Columns: []*dmodel.IndexColumn{
						{Name: "description"},
						{Name: "id"},
					},
				},
			},
		},
		{
			name:   "subquery in FROM",
			scopes: tdata.PostgresSubqueryInFromData.Scopes,
			expectedIndexTarget: []*dmodel.IndexTarget{
				{
					TableName: "tasks",
					Columns: []*dmodel.IndexColumn{
						{Name: "id"},
						{Name: "user_id"},
					},
				},
			},
		},
		{
			name:   "subquery in JOIN",
			scopes: tdata.PostgresSubqueryInJoinData.Scopes,
			expectedIndexTarget: []*dmodel.IndexTarget{
				{
					TableName: "tasks",
					Columns: []*dmodel.IndexColumn{
						{Name: "user_id"},
					},
				},
				{
					TableName: "tasks",
					Columns: []*dmodel.IndexColumn{
						{Name: "user_id"},
						{Name: "id"},
					},
				},
				{
					TableName: "users",
					Columns: []*dmodel.IndexColumn{
						{Name: "id"},
						{Name: "name"},
					},
				},
			},
		},
		{
			name:   "subquery using star",
			scopes: tdata.PostgresSubqueryUsingStarData.Scopes,
			expectedIndexTarget: []*dmodel.IndexTarget{
				{
					TableName: "tasks",
					Columns: []*dmodel.IndexColumn{
						{Name: "description"},
					},
				},
			},
		},
		{
			name:   "subquery with LATERAL JOIN",
			scopes: tdata.PostgresSubqueryWithLateralJoinData.Scopes,
			expectedIndexTarget: []*dmodel.IndexTarget{
				{
					TableName: "tasks",
					Columns: []*dmodel.IndexColumn{
						{Name: "user_id"},
					},
				},
				{
					TableName: "tasks",
					Columns: []*dmodel.IndexColumn{
						{Name: "user_id"},
						{Name: "description"},
					},
				},
			},
		},
		{
			name:   "subquery with LATERAL object reference",
			scopes: tdata.PostgresSubqueryWithLateralObjectReferenceData.Scopes,
			expectedIndexTarget: []*dmodel.IndexTarget{
				{
					TableName: "users",
					Columns: []*dmodel.IndexColumn{
						{Name: "name"},
					},
				},
			},
		},
		{
			name:   "subquery with aggregate function",
			scopes: tdata.PostgresSubqueryWithAggregateFunctionData.Scopes,
			expectedIndexTarget: []*dmodel.IndexTarget{
				{
					TableName: "tasks",
					Columns: []*dmodel.IndexColumn{
						{Name: "description"},
					},
				},
			},
		},
		{
			name:   "correlated subquery",
			scopes: tdata.PostgresCorrelatedSubqueryData.Scopes,
			expectedIndexTarget: []*dmodel.IndexTarget{
				{
					TableName: "tasks",
					Columns: []*dmodel.IndexColumn{
						{Name: "user_id"},
					},
				},
				{
					TableName: "tasks",
					Columns: []*dmodel.IndexColumn{
						{Name: "user_id"},
						{Name: "description"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualTargets, err := dparser.NewIndexTargetBuilder(tables).Build(tt.scopes)
			require.NoError(t, err)

			if diff := cmp.Diff(tt.expectedIndexTarget, actualTargets); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestBuildIndexTargets_HasuraSubquery(t *testing.T) {
	tables := []*dparser.TableSchema{
		{
			Name:        "users",
			PrimaryKeys: []string{"id"},
			Columns: []*dparser.ColumnSchema{
				{Name: "id"},
				{Name: "name"},
				{Name: "is_admin"},
				{Name: "email"},
			},
		},
		{
			Name:        "tasks",
			PrimaryKeys: []string{"id"},
			Columns: []*dparser.ColumnSchema{
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
		scopes              []*dparser.StmtScope
		expectedIndexTarget []*dmodel.IndexTarget
	}{
		{
			name:   "subquery with where",
			scopes: tdata.HasuraSubqueryWithWhereData.Scopes,
			expectedIndexTarget: []*dmodel.IndexTarget{
				{
					TableName: "tasks",
					Columns: []*dmodel.IndexColumn{
						{Name: "user_id"},
					},
				},
				{
					TableName: "tasks",
					Columns: []*dmodel.IndexColumn{
						{Name: "user_id"},
						{Name: "status"},
						{Name: "title"},
					},
				},
				{
					TableName: "users",
					Columns: []*dmodel.IndexColumn{
						{Name: "email"},
					},
				},
				{
					TableName: "users",
					Columns: []*dmodel.IndexColumn{
						{Name: "email"},
						{Name: "id"},
					},
				},
				{
					TableName: "users",
					Columns: []*dmodel.IndexColumn{
						{Name: "id"},
						{Name: "email"},
					},
				},
				{
					TableName: "users",
					Columns: []*dmodel.IndexColumn{
						{Name: "email"},
						{Name: "id"},
						{Name: "name"},
					},
				},
				{
					TableName: "users",
					Columns: []*dmodel.IndexColumn{
						{Name: "id"},
						{Name: "email"},
						{Name: "name"},
					},
				},
			},
		},
		{
			name:   "subquery with variables",
			scopes: tdata.HasuraSubqueryWithVariablesData.Scopes,
			expectedIndexTarget: []*dmodel.IndexTarget{
				{
					TableName: "tasks",
					Columns: []*dmodel.IndexColumn{
						{Name: "user_id"},
					},
				},
				{
					TableName: "tasks",
					Columns: []*dmodel.IndexColumn{
						{Name: "user_id"},
						{Name: "description"},
						{Name: "id"},
						{Name: "status"},
					},
				},
				{
					TableName: "users",
					Columns: []*dmodel.IndexColumn{
						{Name: "email"},
					},
				},
				{
					TableName: "users",
					Columns: []*dmodel.IndexColumn{
						{Name: "email"},
						{Name: "id"},
					},
				},
				{
					TableName: "users",
					Columns: []*dmodel.IndexColumn{
						{Name: "id"},
						{Name: "email"},
					},
				},
				{
					TableName: "users",
					Columns: []*dmodel.IndexColumn{
						{Name: "email"},
						{Name: "id"},
						{Name: "name"},
					},
				},
				{
					TableName: "users",
					Columns: []*dmodel.IndexColumn{
						{Name: "id"},
						{Name: "email"},
						{Name: "name"},
					},
				},
			},
		},
		{
			name:   "subquery with aggregation",
			scopes: tdata.HasuraSubqueryWithAggregationData.Scopes,
			expectedIndexTarget: []*dmodel.IndexTarget{
				{
					TableName: "tasks",
					Columns: []*dmodel.IndexColumn{
						{Name: "user_id"},
					},
				},
				{
					TableName: "tasks",
					Columns: []*dmodel.IndexColumn{
						{Name: "user_id"},
						{Name: "description"},
						{Name: "status"},
						{Name: "title"},
					},
				},
				{
					TableName: "users",
					Columns: []*dmodel.IndexColumn{
						{Name: "email"},
					},
				},
				{
					TableName: "users",
					Columns: []*dmodel.IndexColumn{
						{Name: "email"},
						{Name: "id"},
					},
				},
				{
					TableName: "users",
					Columns: []*dmodel.IndexColumn{
						{Name: "id"},
						{Name: "email"},
					},
				},
				{
					TableName: "users",
					Columns: []*dmodel.IndexColumn{
						{Name: "email"},
						{Name: "id"},
						{Name: "name"},
					},
				},
				{
					TableName: "users",
					Columns: []*dmodel.IndexColumn{
						{Name: "id"},
						{Name: "email"},
						{Name: "name"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualTargets, err := dparser.NewIndexTargetBuilder(tables).Build(tt.scopes)
			require.NoError(t, err)

			if diff := cmp.Diff(tt.expectedIndexTarget, actualTargets); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
