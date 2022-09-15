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
			Name:        "todos",
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
						{Table: "todos", Name: "id", Type: common_model.FieldCondition},
					},
				},
				{
					Columns: []*common_model.FieldColumn{
						{Table: "users", Name: "id", Type: common_model.FieldCondition},
						{Table: "todos", Name: "user_id", Type: common_model.FieldCondition},
					},
				},
			},
			expectedIndexTarget: []*common_model.IndexTargetTable{
				{
					TableName: "todos",
					IndexFields: []*common_model.IndexField{
						{Name: "user_id", Type: common_model.FieldCondition},
					},
				},
				{
					TableName: "todos",
					IndexFields: []*common_model.IndexField{
						{Name: "id", Type: common_model.FieldCondition},
						{Name: "user_id", Type: common_model.FieldCondition},
					},
				},
				{
					TableName: "todos",
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
						{Table: "todos", Name: "description", Type: common_model.FieldReference},
					},
				},
				{
					Columns: []*common_model.FieldColumn{
						{Table: "u", Name: "id", Type: common_model.FieldCondition},
						{Table: "todos", Name: "user_id", Type: common_model.FieldCondition},
					},
				},
			},
			expectedIndexTarget: []*common_model.IndexTargetTable{
				{
					TableName: "todos",
					IndexFields: []*common_model.IndexField{
						{Name: "user_id", Type: common_model.FieldCondition},
					},
				},
				{
					TableName: "todos",
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

func TestBuildIndexTargets_Nested(t *testing.T) {
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
			Name:        "todos",
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
			name: "SELECT (SELECT user_id FROM todos limit 1)",
			scopes: []*common_model.StmtScope{
				{
					Name: "<root>",
					Fields: []*common_model.Field{
						{Columns: []*common_model.FieldColumn{{ReferenceName: "<field0>", Type: common_model.FieldSubquery}}},
					},
					FieldScopes: []*common_model.StmtScope{
						{
							Name: "<field0>",
							Fields: []*common_model.Field{
								{Columns: []*common_model.FieldColumn{{Name: "user_id", Type: common_model.FieldReference}}},
							},
							Tables: []*common_model.Table{
								{Name: "todos"},
							},
						},
					},
				},
			},
			expectedIndexTarget: []*common_model.IndexTargetTable{
				{
					TableName: "todos",
					IndexFields: []*common_model.IndexField{
						{Name: "user_id", Type: common_model.FieldReference},
					},
				},
			},
		},
		{
			name: "SELECT name, is_admin, (SELECT COUNT(description) FROM todos) AS description_count FROM users",
			scopes: []*common_model.StmtScope{
				{
					Name: "<root>",
					Fields: []*common_model.Field{
						{Columns: []*common_model.FieldColumn{{Name: "name", Type: common_model.FieldReference}}},
						{Columns: []*common_model.FieldColumn{{Name: "is_admin", Type: common_model.FieldReference}}},
						{AsName: "description_count", Columns: []*common_model.FieldColumn{{ReferenceName: "<field0>", Type: common_model.FieldSubquery}}},
					},
					FieldScopes: []*common_model.StmtScope{
						{
							Name: "<field0>",
							Fields: []*common_model.Field{
								{Columns: []*common_model.FieldColumn{{Name: "description", Type: common_model.FieldReference}}},
							},
							Tables: []*common_model.Table{
								{Name: "todos"},
							},
						},
					},
					Tables: []*common_model.Table{
						{Name: "users"},
					},
				},
			},
			expectedIndexTarget: []*common_model.IndexTargetTable{
				{
					TableName: "todos",
					IndexFields: []*common_model.IndexField{
						{Name: "description", Type: common_model.FieldReference},
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
			name: "subquery in JOIN",
			scopes: []*common_model.StmtScope{
				{
					Name: "<root>",
					Fields: []*common_model.Field{
						{Columns: []*common_model.FieldColumn{{Type: common_model.FieldStar}}},
						{Columns: []*common_model.FieldColumn{
							{Table: "users", Name: "email", Type: common_model.FieldCondition},
							{Table: "t", Name: "user_id", Type: common_model.FieldCondition},
						}},
					},
					Tables: []*common_model.Table{
						{Name: "users"},
						{Name: "<select0>", AsName: "t"},
					},
					Scopes: []*common_model.StmtScope{
						{
							Name: "<select0>",
							Fields: []*common_model.Field{
								{Columns: []*common_model.FieldColumn{{Name: "status", Type: common_model.FieldReference}}},
								{Columns: []*common_model.FieldColumn{{Name: "user_id", Type: common_model.FieldReference}}},
							},
							Tables: []*common_model.Table{
								{Name: "todos"},
							},
						},
					},
				},
			},
			expectedIndexTarget: []*common_model.IndexTargetTable{
				{
					TableName: "todos",
					IndexFields: []*common_model.IndexField{
						{Name: "status", Type: common_model.FieldReference},
						{Name: "user_id", Type: common_model.FieldReference},
					},
				},
				{
					TableName: "users",
					IndexFields: []*common_model.IndexField{
						{Name: "email", Type: common_model.FieldCondition},
					},
				},
			},
		},
		{
			name: "subquery using star",
			scopes: []*common_model.StmtScope{
				{
					Name: "<root>",
					Fields: []*common_model.Field{
						{Columns: []*common_model.FieldColumn{
							{Name: "description", Type: common_model.FieldReference},
						}},
					},
					Tables: []*common_model.Table{
						{Name: "<select0>", AsName: "t"},
					},
					Scopes: []*common_model.StmtScope{
						{
							Name: "<select0>",
							Fields: []*common_model.Field{
								{Columns: []*common_model.FieldColumn{{Type: common_model.FieldStar}}},
							},
							Tables: []*common_model.Table{
								{Name: "todos"},
							},
						},
					},
				},
			},
			expectedIndexTarget: []*common_model.IndexTargetTable{
				{
					TableName: "todos",
					IndexFields: []*common_model.IndexField{
						{Name: "description", Type: common_model.FieldReference},
					},
				},
			},
		},
		{
			name: "subquery with LATERAL JOIN",
			scopes: []*common_model.StmtScope{
				{
					Name: "<root>",
					Fields: []*common_model.Field{
						{Columns: []*common_model.FieldColumn{
							{Name: "description", Type: common_model.FieldReference},
						}},
						{Columns: []*common_model.FieldColumn{
							{Table: "u", Name: "id", Type: common_model.FieldCondition},
							{Table: "t", Name: "user_id", Type: common_model.FieldCondition},
						}},
					},
					Tables: []*common_model.Table{
						{Name: "<select0>", AsName: "u"},
						{Name: "<select1>", AsName: "t", IsLateral: true},
					},
					Scopes: []*common_model.StmtScope{
						{
							Name: "<select0>",
							Fields: []*common_model.Field{
								{Columns: []*common_model.FieldColumn{
									{Name: "email", Type: common_model.FieldReference},
									{Name: "id", Type: common_model.FieldReference},
								}},
							},
							Tables: []*common_model.Table{
								{Name: "users"},
							},
						},
						{
							Name: "<select1>",
							Fields: []*common_model.Field{
								{Columns: []*common_model.FieldColumn{{Name: "description", Type: common_model.FieldReference}}},
								{Columns: []*common_model.FieldColumn{
									{Name: "user_id", Type: common_model.FieldCondition},
									{Table: "u", Name: "id", Type: common_model.FieldCondition},
								}},
							},
							Tables: []*common_model.Table{
								{Name: "todos"},
							},
						},
					},
				},
			},
			expectedIndexTarget: []*common_model.IndexTargetTable{
				{
					TableName: "todos",
					IndexFields: []*common_model.IndexField{
						{Name: "user_id", Type: common_model.FieldCondition},
					},
				},
				{
					TableName: "todos",
					IndexFields: []*common_model.IndexField{
						{Name: "user_id", Type: common_model.FieldCondition},
						{Name: "description", Type: common_model.FieldReference},
					},
				},
				{
					TableName: "users",
					IndexFields: []*common_model.IndexField{
						{Name: "email", Type: common_model.FieldReference},
						{Name: "id", Type: common_model.FieldReference},
					},
				},
			},
		},
		{
			name:   "hasura query",
			scopes: tdata.Hasura1.Scopes,
			expectedIndexTarget: []*common_model.IndexTargetTable{
				{
					TableName: "todos",
					IndexFields: []*common_model.IndexField{
						{Name: "user_id", Type: common_model.FieldCondition},
					},
				},
				{
					TableName: "todos",
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
			name:   "hasura query2",
			scopes: tdata.Hasura2.Scopes,
			expectedIndexTarget: []*common_model.IndexTargetTable{
				{
					TableName: "todos",
					IndexFields: []*common_model.IndexField{
						{Name: "user_id", Type: common_model.FieldCondition},
					},
				},
				{
					TableName: "todos",
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
			name:   "hasura query3",
			scopes: tdata.Hasura3.Scopes,
			expectedIndexTarget: []*common_model.IndexTargetTable{
				{
					TableName: "todos",
					IndexFields: []*common_model.IndexField{
						{Name: "user_id", Type: common_model.FieldCondition},
					},
				},
				{
					TableName: "todos",
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
