package builders_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/mrasu/GravityR/database/db_models"
	"github.com/mrasu/GravityR/database/db_models/builders"
	"github.com/mrasu/GravityR/thelper/tdata"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBuildIndexTargets_SingleTable(t *testing.T) {
	tables := []*db_models.TableSchema{
		{
			Name:        "users",
			PrimaryKeys: []string{"pk1", "pk2"},
			Columns: []*db_models.ColumnSchema{
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
		scopeFields          []*db_models.Field
		expectedTargetFields [][]*db_models.IndexField
	}{
		{
			name:        "simple",
			asTableName: "",
			scopeFields: []*db_models.Field{
				{
					Columns: []*db_models.FieldColumn{
						{Table: "users", Name: "name", Type: db_models.FieldCondition},
						{Table: "users", Name: "is_admin", Type: db_models.FieldCondition},
					},
				},
			},
			expectedTargetFields: [][]*db_models.IndexField{
				{
					{Name: "is_admin", Type: db_models.FieldCondition},
				},
				{
					{Name: "name", Type: db_models.FieldCondition},
				},
				{
					{Name: "is_admin", Type: db_models.FieldCondition},
					{Name: "name", Type: db_models.FieldCondition},
				},
				{
					{Name: "name", Type: db_models.FieldCondition},
					{Name: "is_admin", Type: db_models.FieldCondition},
				},
			},
		},
		{
			name:        "no table name in fields",
			asTableName: "",
			scopeFields: []*db_models.Field{
				{
					Columns: []*db_models.FieldColumn{
						{Name: "name", Type: db_models.FieldCondition},
						{Name: "is_admin", Type: db_models.FieldCondition},
					},
				},
			},
			expectedTargetFields: [][]*db_models.IndexField{
				{
					{Name: "is_admin", Type: db_models.FieldCondition},
				},
				{
					{Name: "name", Type: db_models.FieldCondition},
				},
				{
					{Name: "is_admin", Type: db_models.FieldCondition},
					{Name: "name", Type: db_models.FieldCondition},
				},
				{
					{Name: "name", Type: db_models.FieldCondition},
					{Name: "is_admin", Type: db_models.FieldCondition},
				},
			},
		},
		{
			name:        "aliased table",
			asTableName: "u",
			scopeFields: []*db_models.Field{
				{
					Columns: []*db_models.FieldColumn{
						{Table: "u", Name: "name", Type: db_models.FieldCondition},
						{Table: "u", Name: "is_admin", Type: db_models.FieldCondition},
					},
				},
			},
			expectedTargetFields: [][]*db_models.IndexField{
				{
					{Name: "is_admin", Type: db_models.FieldCondition},
				},
				{
					{Name: "name", Type: db_models.FieldCondition},
				},
				{
					{Name: "is_admin", Type: db_models.FieldCondition},
					{Name: "name", Type: db_models.FieldCondition},
				},
				{
					{Name: "name", Type: db_models.FieldCondition},
					{Name: "is_admin", Type: db_models.FieldCondition},
				},
			},
		},
		{
			name: "includes reference column",
			scopeFields: []*db_models.Field{
				{
					Columns: []*db_models.FieldColumn{
						{Table: "users", Name: "name", Type: db_models.FieldReference},
						{Table: "users", Name: "is_admin", Type: db_models.FieldCondition},
					},
				},
			},
			expectedTargetFields: [][]*db_models.IndexField{
				{
					{Name: "is_admin", Type: db_models.FieldCondition},
				},
				{
					{Name: "is_admin", Type: db_models.FieldCondition},
					{Name: "name", Type: db_models.FieldReference},
				},
			},
		},
		{
			name: "includes primary key column",
			scopeFields: []*db_models.Field{
				{
					Columns: []*db_models.FieldColumn{
						{Table: "users", Name: "pk1", Type: db_models.FieldCondition},
						{Table: "users", Name: "pk2", Type: db_models.FieldCondition},
						{Table: "users", Name: "is_admin", Type: db_models.FieldCondition},
					},
				},
			},
			expectedTargetFields: [][]*db_models.IndexField{
				{
					{Name: "is_admin", Type: db_models.FieldCondition},
				},
				{
					{Name: "is_admin", Type: db_models.FieldCondition},
					{Name: "pk1", Type: db_models.FieldCondition},
				},
				{
					{Name: "is_admin", Type: db_models.FieldCondition},
					{Name: "pk2", Type: db_models.FieldCondition},
				},
				{
					{Name: "pk1", Type: db_models.FieldCondition},
					{Name: "is_admin", Type: db_models.FieldCondition},
				},
				{
					{Name: "pk2", Type: db_models.FieldCondition},
					{Name: "is_admin", Type: db_models.FieldCondition},
				},
				{
					{Name: "is_admin", Type: db_models.FieldCondition},
					{Name: "pk1", Type: db_models.FieldCondition},
					{Name: "pk2", Type: db_models.FieldCondition},
				},
				{
					{Name: "is_admin", Type: db_models.FieldCondition},
					{Name: "pk2", Type: db_models.FieldCondition},
					{Name: "pk1", Type: db_models.FieldCondition},
				},
				{
					{Name: "pk1", Type: db_models.FieldCondition},
					{Name: "is_admin", Type: db_models.FieldCondition},
					{Name: "pk2", Type: db_models.FieldCondition},
				},
				{
					{Name: "pk1", Type: db_models.FieldCondition},
					{Name: "pk2", Type: db_models.FieldCondition},
					{Name: "is_admin", Type: db_models.FieldCondition},
				},
				{
					{Name: "pk2", Type: db_models.FieldCondition},
					{Name: "is_admin", Type: db_models.FieldCondition},
					{Name: "pk1", Type: db_models.FieldCondition},
				},
				{
					{Name: "pk2", Type: db_models.FieldCondition},
					{Name: "pk1", Type: db_models.FieldCondition},
					{Name: "is_admin", Type: db_models.FieldCondition},
				},
			},
		},
		{
			name: "includes primary key column reference",
			scopeFields: []*db_models.Field{
				{
					Columns: []*db_models.FieldColumn{
						{Table: "users", Name: "pk1", Type: db_models.FieldReference},
						{Table: "users", Name: "pk2", Type: db_models.FieldReference},
						{Table: "users", Name: "is_admin", Type: db_models.FieldReference},
					},
				},
			},
			expectedTargetFields: [][]*db_models.IndexField{
				{
					{Name: "is_admin", Type: db_models.FieldReference},
					{Name: "pk1", Type: db_models.FieldReference},
					{Name: "pk2", Type: db_models.FieldReference},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := &db_models.StmtScope{
				Name: db_models.RootScopeName,
				Tables: []*db_models.Table{
					{AsName: tt.asTableName, Name: "users"},
				},
				Fields: tt.scopeFields,
			}
			actualTargets, err := builders.BuildIndexTargets(tables, []*db_models.StmtScope{scope})
			require.NoError(t, err)

			var expected []*db_models.IndexTargetTable
			for _, fields := range tt.expectedTargetFields {
				expected = append(expected, &db_models.IndexTargetTable{
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
	tables := []*db_models.TableSchema{
		{
			Name:        "users",
			PrimaryKeys: []string{"id"},
			Columns: []*db_models.ColumnSchema{
				{Name: "id"},
				{Name: "name"},
				{Name: "is_admin"},
			},
		},
		{
			Name:        "todos",
			PrimaryKeys: []string{"id"},
			Columns: []*db_models.ColumnSchema{
				{Name: "id"},
				{Name: "user_id"},
				{Name: "description"},
			},
		},
	}

	tests := []struct {
		name                string
		asTableName         string
		scopeFields         []*db_models.Field
		expectedIndexTarget []*db_models.IndexTargetTable
	}{
		{
			name:        "simple",
			asTableName: "u",
			scopeFields: []*db_models.Field{
				{
					Columns: []*db_models.FieldColumn{
						{Table: "todos", Name: "id", Type: db_models.FieldCondition},
					},
				},
				{
					Columns: []*db_models.FieldColumn{
						{Table: "users", Name: "id", Type: db_models.FieldCondition},
						{Table: "todos", Name: "user_id", Type: db_models.FieldCondition},
					},
				},
			},
			expectedIndexTarget: []*db_models.IndexTargetTable{
				{
					TableName: "todos",
					IndexFields: []*db_models.IndexField{
						{Name: "user_id", Type: db_models.FieldCondition},
					},
				},
				{
					TableName: "todos",
					IndexFields: []*db_models.IndexField{
						{Name: "id", Type: db_models.FieldCondition},
						{Name: "user_id", Type: db_models.FieldCondition},
					},
				},
				{
					TableName: "todos",
					IndexFields: []*db_models.IndexField{
						{Name: "user_id", Type: db_models.FieldCondition},
						{Name: "id", Type: db_models.FieldCondition},
					},
				},
			},
		},
		{
			name:        "includes reference only column",
			asTableName: "u",
			scopeFields: []*db_models.Field{
				{
					Columns: []*db_models.FieldColumn{
						{Table: "todos", Name: "description", Type: db_models.FieldReference},
					},
				},
				{
					Columns: []*db_models.FieldColumn{
						{Table: "u", Name: "id", Type: db_models.FieldCondition},
						{Table: "todos", Name: "user_id", Type: db_models.FieldCondition},
					},
				},
			},
			expectedIndexTarget: []*db_models.IndexTargetTable{
				{
					TableName: "todos",
					IndexFields: []*db_models.IndexField{
						{Name: "user_id", Type: db_models.FieldCondition},
					},
				},
				{
					TableName: "todos",
					IndexFields: []*db_models.IndexField{
						{Name: "user_id", Type: db_models.FieldCondition},
						{Name: "description", Type: db_models.FieldReference},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := &db_models.StmtScope{
				Name: db_models.RootScopeName,
				Tables: []*db_models.Table{
					{AsName: tt.asTableName, Name: "users"},
				},
				Fields: tt.scopeFields,
			}
			actualTargets, err := builders.BuildIndexTargets(tables, []*db_models.StmtScope{scope})
			require.NoError(t, err)

			if diff := cmp.Diff(tt.expectedIndexTarget, actualTargets); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestBuildIndexTargets_Nested(t *testing.T) {
	tables := []*db_models.TableSchema{
		{
			Name:        "users",
			PrimaryKeys: []string{"id"},
			Columns: []*db_models.ColumnSchema{
				{Name: "id"},
				{Name: "name"},
				{Name: "is_admin"},
				{Name: "email"},
			},
		},
		{
			Name:        "todos",
			PrimaryKeys: []string{"id"},
			Columns: []*db_models.ColumnSchema{
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
		scopes              []*db_models.StmtScope
		expectedIndexTarget []*db_models.IndexTargetTable
	}{
		{
			name: "SELECT (SELECT user_id FROM todos limit 1)",
			scopes: []*db_models.StmtScope{
				{
					Name: "<root>",
					Fields: []*db_models.Field{
						{Columns: []*db_models.FieldColumn{{ReferenceName: "<field0>", Type: db_models.FieldSubquery}}},
					},
					FieldScopes: []*db_models.StmtScope{
						{
							Name: "<field0>",
							Fields: []*db_models.Field{
								{Columns: []*db_models.FieldColumn{{Name: "user_id", Type: db_models.FieldReference}}},
							},
							Tables: []*db_models.Table{
								{Name: "todos"},
							},
						},
					},
				},
			},
			expectedIndexTarget: []*db_models.IndexTargetTable{
				{
					TableName: "todos",
					IndexFields: []*db_models.IndexField{
						{Name: "user_id", Type: db_models.FieldReference},
					},
				},
			},
		},
		{
			name: "SELECT name, is_admin, (SELECT COUNT(description) FROM todos) AS description_count FROM users",
			scopes: []*db_models.StmtScope{
				{
					Name: "<root>",
					Fields: []*db_models.Field{
						{Columns: []*db_models.FieldColumn{{Name: "name", Type: db_models.FieldReference}}},
						{Columns: []*db_models.FieldColumn{{Name: "is_admin", Type: db_models.FieldReference}}},
						{AsName: "description_count", Columns: []*db_models.FieldColumn{{ReferenceName: "<field0>", Type: db_models.FieldSubquery}}},
					},
					FieldScopes: []*db_models.StmtScope{
						{
							Name: "<field0>",
							Fields: []*db_models.Field{
								{Columns: []*db_models.FieldColumn{{Name: "description", Type: db_models.FieldReference}}},
							},
							Tables: []*db_models.Table{
								{Name: "todos"},
							},
						},
					},
					Tables: []*db_models.Table{
						{Name: "users"},
					},
				},
			},
			expectedIndexTarget: []*db_models.IndexTargetTable{
				{
					TableName: "todos",
					IndexFields: []*db_models.IndexField{
						{Name: "description", Type: db_models.FieldReference},
					},
				},
				{
					TableName: "users",
					IndexFields: []*db_models.IndexField{
						{Name: "is_admin", Type: db_models.FieldReference},
						{Name: "name", Type: db_models.FieldReference},
					},
				},
			},
		},
		{
			name: "subquery in JOIN",
			scopes: []*db_models.StmtScope{
				{
					Name: "<root>",
					Fields: []*db_models.Field{
						{Columns: []*db_models.FieldColumn{{Type: db_models.FieldStar}}},
						{Columns: []*db_models.FieldColumn{
							{Table: "users", Name: "email", Type: db_models.FieldCondition},
							{Table: "t", Name: "user_id", Type: db_models.FieldCondition},
						}},
					},
					Tables: []*db_models.Table{
						{Name: "users"},
						{Name: "<select0>", AsName: "t"},
					},
					Scopes: []*db_models.StmtScope{
						{
							Name: "<select0>",
							Fields: []*db_models.Field{
								{Columns: []*db_models.FieldColumn{{Name: "status", Type: db_models.FieldReference}}},
								{Columns: []*db_models.FieldColumn{{Name: "user_id", Type: db_models.FieldReference}}},
							},
							Tables: []*db_models.Table{
								{Name: "todos"},
							},
						},
					},
				},
			},
			expectedIndexTarget: []*db_models.IndexTargetTable{
				{
					TableName: "todos",
					IndexFields: []*db_models.IndexField{
						{Name: "status", Type: db_models.FieldReference},
						{Name: "user_id", Type: db_models.FieldReference},
					},
				},
				{
					TableName: "users",
					IndexFields: []*db_models.IndexField{
						{Name: "email", Type: db_models.FieldCondition},
					},
				},
			},
		},
		{
			name: "subquery using star",
			scopes: []*db_models.StmtScope{
				{
					Name: "<root>",
					Fields: []*db_models.Field{
						{Columns: []*db_models.FieldColumn{
							{Name: "description", Type: db_models.FieldReference},
						}},
					},
					Tables: []*db_models.Table{
						{Name: "<select0>", AsName: "t"},
					},
					Scopes: []*db_models.StmtScope{
						{
							Name: "<select0>",
							Fields: []*db_models.Field{
								{Columns: []*db_models.FieldColumn{{Type: db_models.FieldStar}}},
							},
							Tables: []*db_models.Table{
								{Name: "todos"},
							},
						},
					},
				},
			},
			expectedIndexTarget: []*db_models.IndexTargetTable{
				{
					TableName: "todos",
					IndexFields: []*db_models.IndexField{
						{Name: "description", Type: db_models.FieldReference},
					},
				},
			},
		},
		{
			name: "subquery with LATERAL JOIN",
			scopes: []*db_models.StmtScope{
				{
					Name: "<root>",
					Fields: []*db_models.Field{
						{Columns: []*db_models.FieldColumn{
							{Name: "description", Type: db_models.FieldReference},
						}},
						{Columns: []*db_models.FieldColumn{
							{Table: "u", Name: "id", Type: db_models.FieldCondition},
							{Table: "t", Name: "user_id", Type: db_models.FieldCondition},
						}},
					},
					Tables: []*db_models.Table{
						{Name: "<select0>", AsName: "u"},
						{Name: "<select1>", AsName: "t", IsLateral: true},
					},
					Scopes: []*db_models.StmtScope{
						{
							Name: "<select0>",
							Fields: []*db_models.Field{
								{Columns: []*db_models.FieldColumn{
									{Name: "email", Type: db_models.FieldReference},
									{Name: "id", Type: db_models.FieldReference},
								}},
							},
							Tables: []*db_models.Table{
								{Name: "users"},
							},
						},
						{
							Name: "<select1>",
							Fields: []*db_models.Field{
								{Columns: []*db_models.FieldColumn{{Name: "description", Type: db_models.FieldReference}}},
								{Columns: []*db_models.FieldColumn{
									{Name: "user_id", Type: db_models.FieldCondition},
									{Table: "u", Name: "id", Type: db_models.FieldCondition},
								}},
							},
							Tables: []*db_models.Table{
								{Name: "todos"},
							},
						},
					},
				},
			},
			expectedIndexTarget: []*db_models.IndexTargetTable{
				{
					TableName: "todos",
					IndexFields: []*db_models.IndexField{
						{Name: "user_id", Type: db_models.FieldCondition},
					},
				},
				{
					TableName: "todos",
					IndexFields: []*db_models.IndexField{
						{Name: "user_id", Type: db_models.FieldCondition},
						{Name: "description", Type: db_models.FieldReference},
					},
				},
				{
					TableName: "users",
					IndexFields: []*db_models.IndexField{
						{Name: "email", Type: db_models.FieldReference},
						{Name: "id", Type: db_models.FieldReference},
					},
				},
			},
		},
		{
			name:   "hasura query",
			scopes: tdata.Hasura1.Scopes,
			expectedIndexTarget: []*db_models.IndexTargetTable{
				{
					TableName: "todos",
					IndexFields: []*db_models.IndexField{
						{Name: "user_id", Type: db_models.FieldCondition},
					},
				},
				{
					TableName: "todos",
					IndexFields: []*db_models.IndexField{
						{Name: "user_id", Type: db_models.FieldCondition},
						{Name: "status", Type: db_models.FieldReference},
						{Name: "title", Type: db_models.FieldReference},
					},
				},
				{
					TableName: "users",
					IndexFields: []*db_models.IndexField{
						{Name: "email", Type: db_models.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*db_models.IndexField{
						{Name: "email", Type: db_models.FieldCondition},
						{Name: "id", Type: db_models.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*db_models.IndexField{
						{Name: "id", Type: db_models.FieldCondition},
						{Name: "email", Type: db_models.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*db_models.IndexField{
						{Name: "email", Type: db_models.FieldCondition},
						{Name: "id", Type: db_models.FieldCondition},
						{Name: "name", Type: db_models.FieldReference},
					},
				},
				{
					TableName: "users",
					IndexFields: []*db_models.IndexField{
						{Name: "id", Type: db_models.FieldCondition},
						{Name: "email", Type: db_models.FieldCondition},
						{Name: "name", Type: db_models.FieldReference},
					},
				},
			},
		},
		{
			name:   "hasura query2",
			scopes: tdata.Hasura2.Scopes,
			expectedIndexTarget: []*db_models.IndexTargetTable{
				{
					TableName: "todos",
					IndexFields: []*db_models.IndexField{
						{Name: "user_id", Type: db_models.FieldCondition},
					},
				},
				{
					TableName: "todos",
					IndexFields: []*db_models.IndexField{
						{Name: "user_id", Type: db_models.FieldCondition},
						{Name: "description", Type: db_models.FieldReference},
						{Name: "id", Type: db_models.FieldReference},
						{Name: "status", Type: db_models.FieldReference},
					},
				},
				{
					TableName: "users",
					IndexFields: []*db_models.IndexField{
						{Name: "email", Type: db_models.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*db_models.IndexField{
						{Name: "email", Type: db_models.FieldCondition},
						{Name: "id", Type: db_models.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*db_models.IndexField{
						{Name: "id", Type: db_models.FieldCondition},
						{Name: "email", Type: db_models.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*db_models.IndexField{
						{Name: "email", Type: db_models.FieldCondition},
						{Name: "id", Type: db_models.FieldCondition},
						{Name: "name", Type: db_models.FieldReference},
					},
				},
				{
					TableName: "users",
					IndexFields: []*db_models.IndexField{
						{Name: "id", Type: db_models.FieldCondition},
						{Name: "email", Type: db_models.FieldCondition},
						{Name: "name", Type: db_models.FieldReference},
					},
				},
			},
		},
		{
			name:   "hasura query3",
			scopes: tdata.Hasura3.Scopes,
			expectedIndexTarget: []*db_models.IndexTargetTable{
				{
					TableName: "todos",
					IndexFields: []*db_models.IndexField{
						{Name: "user_id", Type: db_models.FieldCondition},
					},
				},
				{
					TableName: "todos",
					IndexFields: []*db_models.IndexField{
						{Name: "user_id", Type: db_models.FieldCondition},
						{Name: "description", Type: db_models.FieldReference},
						{Name: "status", Type: db_models.FieldReference},
						{Name: "title", Type: db_models.FieldReference},
					},
				},
				{
					TableName: "users",
					IndexFields: []*db_models.IndexField{
						{Name: "email", Type: db_models.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*db_models.IndexField{
						{Name: "email", Type: db_models.FieldCondition},
						{Name: "id", Type: db_models.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*db_models.IndexField{
						{Name: "id", Type: db_models.FieldCondition},
						{Name: "email", Type: db_models.FieldCondition},
					},
				},
				{
					TableName: "users",
					IndexFields: []*db_models.IndexField{
						{Name: "email", Type: db_models.FieldCondition},
						{Name: "id", Type: db_models.FieldCondition},
						{Name: "name", Type: db_models.FieldReference},
					},
				},
				{
					TableName: "users",
					IndexFields: []*db_models.IndexField{
						{Name: "id", Type: db_models.FieldCondition},
						{Name: "email", Type: db_models.FieldCondition},
						{Name: "name", Type: db_models.FieldReference},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualTargets, err := builders.BuildIndexTargets(tables, tt.scopes)
			require.NoError(t, err)

			if diff := cmp.Diff(tt.expectedIndexTarget, actualTargets); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
