package builders_test

import (
	"github.com/mrasu/GravityR/database/db_models"
	"github.com/mrasu/GravityR/database/db_models/builders"
	"github.com/mrasu/GravityR/lib"
	"github.com/stretchr/testify/assert"
	"strings"
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
					{Name: "is_admin", Type: db_models.FieldCondition},
					{Name: "name", Type: db_models.FieldCondition},
				},
				{
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
					{Name: "is_admin", Type: db_models.FieldCondition},
					{Name: "name", Type: db_models.FieldCondition},
				},
				{
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
					{Name: "is_admin", Type: db_models.FieldCondition},
					{Name: "name", Type: db_models.FieldCondition},
				},
				{
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
					{Name: "pk2", Type: db_models.FieldCondition},
				},
				{
					{Name: "pk1", Type: db_models.FieldCondition},
					{Name: "is_admin", Type: db_models.FieldCondition},
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
				},
				{
					{Name: "pk2", Type: db_models.FieldCondition},
					{Name: "is_admin", Type: db_models.FieldCondition},
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
			assert.NoError(t, err)

			sortTargets(actualTargets)

			var expected []*db_models.IndexTargetTable
			for _, fields := range tt.expectedTargetFields {
				expected = append(expected, &db_models.IndexTargetTable{
					TableName:           "users",
					AffectingTableNames: []string{"<root>"},
					IndexFields:         fields,
				})
			}
			assert.Equal(t, expected, actualTargets, toIndexTargetCompareString(expected, actualTargets))
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
					TableName:           "todos",
					AffectingTableNames: []string{"<root>"},
					IndexFields: []*db_models.IndexField{
						{Name: "id", Type: db_models.FieldCondition},
						{Name: "user_id", Type: db_models.FieldCondition},
					},
				},
				{
					TableName:           "todos",
					AffectingTableNames: []string{"<root>"},
					IndexFields: []*db_models.IndexField{
						{Name: "user_id", Type: db_models.FieldCondition},
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
					TableName:           "todos",
					AffectingTableNames: []string{"<root>"},
					IndexFields: []*db_models.IndexField{
						{Name: "user_id", Type: db_models.FieldCondition},
					},
				},
				{
					TableName:           "todos",
					AffectingTableNames: []string{"<root>"},
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
			assert.NoError(t, err)

			sortTargets(actualTargets)
			assert.Equal(t, tt.expectedIndexTarget, actualTargets, toIndexTargetCompareString(tt.expectedIndexTarget, actualTargets))
		})
	}
}

func sortTargets(targets []*db_models.IndexTargetTable) {
	lib.SortF(targets, func(t *db_models.IndexTargetTable) string {
		var namesI []string
		for _, f := range t.IndexFields {
			namesI = append(namesI, f.Name)
		}
		return strings.Join(namesI, "")
	})
}

func toIndexTargetCompareString(expected, actual []*db_models.IndexTargetTable) string {
	txt := "------\nexpected:\n"
	txt += lib.JoinF(expected, "\n", func(it *db_models.IndexTargetTable) string { return it.String() })
	txt += "\nvs\nactual:\n"
	txt += lib.JoinF(actual, "\n", func(it *db_models.IndexTargetTable) string { return it.String() })
	txt += "\n------\n"

	return txt
}
