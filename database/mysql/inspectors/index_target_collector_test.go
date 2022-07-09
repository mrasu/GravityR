package inspectors_test

import (
	"github.com/mrasu/GravityR/database/mysql/inspectors"
	"github.com/mrasu/GravityR/database/mysql/models"
	"github.com/mrasu/GravityR/lib"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestCollectIndexTargets_SingleTable(t *testing.T) {
	tables := []*models.TableSchema{
		{
			Name:        "users",
			PrimaryKeys: []string{"pk1", "pk2"},
			Columns: []*models.ColumnSchema{
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
		scopeFields          []*models.Field
		expectedTargetFields [][]*models.IndexField
	}{
		{
			name:        "simple",
			asTableName: "",
			scopeFields: []*models.Field{
				{
					Columns: []*models.FieldColumn{
						{Table: "users", Name: "name", Type: models.FieldCondition},
						{Table: "users", Name: "is_admin", Type: models.FieldCondition},
					},
				},
			},
			expectedTargetFields: [][]*models.IndexField{
				{
					{Name: "is_admin", Type: models.FieldCondition},
				},
				{
					{Name: "is_admin", Type: models.FieldCondition},
					{Name: "name", Type: models.FieldCondition},
				},
				{
					{Name: "name", Type: models.FieldCondition},
				},
				{
					{Name: "name", Type: models.FieldCondition},
					{Name: "is_admin", Type: models.FieldCondition},
				},
			},
		},
		{
			name:        "no table name in fields",
			asTableName: "",
			scopeFields: []*models.Field{
				{
					Columns: []*models.FieldColumn{
						{Name: "name", Type: models.FieldCondition},
						{Name: "is_admin", Type: models.FieldCondition},
					},
				},
			},
			expectedTargetFields: [][]*models.IndexField{
				{
					{Name: "is_admin", Type: models.FieldCondition},
				},
				{
					{Name: "is_admin", Type: models.FieldCondition},
					{Name: "name", Type: models.FieldCondition},
				},
				{
					{Name: "name", Type: models.FieldCondition},
				},
				{
					{Name: "name", Type: models.FieldCondition},
					{Name: "is_admin", Type: models.FieldCondition},
				},
			},
		},
		{
			name:        "aliased table",
			asTableName: "u",
			scopeFields: []*models.Field{
				{
					Columns: []*models.FieldColumn{
						{Table: "u", Name: "name", Type: models.FieldCondition},
						{Table: "u", Name: "is_admin", Type: models.FieldCondition},
					},
				},
			},
			expectedTargetFields: [][]*models.IndexField{
				{
					{Name: "is_admin", Type: models.FieldCondition},
				},
				{
					{Name: "is_admin", Type: models.FieldCondition},
					{Name: "name", Type: models.FieldCondition},
				},
				{
					{Name: "name", Type: models.FieldCondition},
				},
				{
					{Name: "name", Type: models.FieldCondition},
					{Name: "is_admin", Type: models.FieldCondition},
				},
			},
		},
		{
			name: "includes reference column",
			scopeFields: []*models.Field{
				{
					Columns: []*models.FieldColumn{
						{Table: "users", Name: "name", Type: models.FieldReference},
						{Table: "users", Name: "is_admin", Type: models.FieldCondition},
					},
				},
			},
			expectedTargetFields: [][]*models.IndexField{
				{
					{Name: "is_admin", Type: models.FieldCondition},
				},
				{
					{Name: "is_admin", Type: models.FieldCondition},
					{Name: "name", Type: models.FieldReference},
				},
			},
		},
		{
			name: "includes primary key column",
			scopeFields: []*models.Field{
				{
					Columns: []*models.FieldColumn{
						{Table: "users", Name: "pk1", Type: models.FieldCondition},
						{Table: "users", Name: "pk2", Type: models.FieldCondition},
						{Table: "users", Name: "is_admin", Type: models.FieldCondition},
					},
				},
			},
			expectedTargetFields: [][]*models.IndexField{
				{
					{Name: "is_admin", Type: models.FieldCondition},
				},
				{
					{Name: "is_admin", Type: models.FieldCondition},
					{Name: "pk2", Type: models.FieldCondition},
				},
				{
					{Name: "pk1", Type: models.FieldCondition},
					{Name: "is_admin", Type: models.FieldCondition},
				},
				{
					{Name: "pk1", Type: models.FieldCondition},
					{Name: "is_admin", Type: models.FieldCondition},
					{Name: "pk2", Type: models.FieldCondition},
				},
				{
					{Name: "pk1", Type: models.FieldCondition},
					{Name: "pk2", Type: models.FieldCondition},
					{Name: "is_admin", Type: models.FieldCondition},
				},
				{
					{Name: "pk2", Type: models.FieldCondition},
				},
				{
					{Name: "pk2", Type: models.FieldCondition},
					{Name: "is_admin", Type: models.FieldCondition},
				},
				{
					{Name: "pk2", Type: models.FieldCondition},
					{Name: "pk1", Type: models.FieldCondition},
					{Name: "is_admin", Type: models.FieldCondition},
				},
			},
		},
		{
			name: "includes primary key column reference",
			scopeFields: []*models.Field{
				{
					Columns: []*models.FieldColumn{
						{Table: "users", Name: "pk1", Type: models.FieldReference},
						{Table: "users", Name: "pk2", Type: models.FieldReference},
						{Table: "users", Name: "is_admin", Type: models.FieldReference},
					},
				},
			},
			expectedTargetFields: [][]*models.IndexField{
				{
					{Name: "is_admin", Type: models.FieldReference},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := &models.StmtScope{
				Name: models.RootScopeName,
				Tables: []*models.Table{
					{AsName: tt.asTableName, Name: "users"},
				},
				Fields: tt.scopeFields,
			}
			actualTargets, err := inspectors.CollectIndexTargets(tables, []*models.StmtScope{scope})
			assert.NoError(t, err)

			sortTargets(actualTargets)

			var expected []*models.IndexTargetTable
			for _, fields := range tt.expectedTargetFields {
				expected = append(expected, &models.IndexTargetTable{
					TableName:           "users",
					AffectingTableNames: []string{"<root>"},
					IndexFields:         fields,
				})
			}
			assert.Equal(t, expected, actualTargets, toIndexTargetCompareString(expected, actualTargets))
		})
	}
}

func TestCollectIndexTargets_MultipleTables(t *testing.T) {
	tables := []*models.TableSchema{
		{
			Name:        "users",
			PrimaryKeys: []string{"id"},
			Columns: []*models.ColumnSchema{
				{Name: "id"},
				{Name: "name"},
				{Name: "is_admin"},
			},
		},
		{
			Name:        "todos",
			PrimaryKeys: []string{"id"},
			Columns: []*models.ColumnSchema{
				{Name: "id"},
				{Name: "user_id"},
				{Name: "description"},
			},
		},
	}

	tests := []struct {
		name                string
		asTableName         string
		scopeFields         []*models.Field
		expectedIndexTarget []*models.IndexTargetTable
	}{
		{
			name:        "simple",
			asTableName: "u",
			scopeFields: []*models.Field{
				{
					Columns: []*models.FieldColumn{
						{Table: "todos", Name: "id", Type: models.FieldCondition},
					},
				},
				{
					Columns: []*models.FieldColumn{
						{Table: "users", Name: "id", Type: models.FieldCondition},
						{Table: "todos", Name: "user_id", Type: models.FieldCondition},
					},
				},
			},
			expectedIndexTarget: []*models.IndexTargetTable{
				{
					TableName:           "todos",
					AffectingTableNames: []string{"<root>"},
					IndexFields: []*models.IndexField{
						{Name: "id", Type: models.FieldCondition},
						{Name: "user_id", Type: models.FieldCondition},
					},
				},
				{
					TableName:           "todos",
					AffectingTableNames: []string{"<root>"},
					IndexFields: []*models.IndexField{
						{Name: "user_id", Type: models.FieldCondition},
					},
				},
			},
		},
		{
			name:        "includes reference only column",
			asTableName: "u",
			scopeFields: []*models.Field{
				{
					Columns: []*models.FieldColumn{
						{Table: "todos", Name: "description", Type: models.FieldReference},
					},
				},
				{
					Columns: []*models.FieldColumn{
						{Table: "u", Name: "id", Type: models.FieldCondition},
						{Table: "todos", Name: "user_id", Type: models.FieldCondition},
					},
				},
			},
			expectedIndexTarget: []*models.IndexTargetTable{
				{
					TableName:           "todos",
					AffectingTableNames: []string{"<root>"},
					IndexFields: []*models.IndexField{
						{Name: "user_id", Type: models.FieldCondition},
					},
				},
				{
					TableName:           "todos",
					AffectingTableNames: []string{"<root>"},
					IndexFields: []*models.IndexField{
						{Name: "user_id", Type: models.FieldCondition},
						{Name: "description", Type: models.FieldReference},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := &models.StmtScope{
				Name: models.RootScopeName,
				Tables: []*models.Table{
					{AsName: tt.asTableName, Name: "users"},
				},
				Fields: tt.scopeFields,
			}
			actualTargets, err := inspectors.CollectIndexTargets(tables, []*models.StmtScope{scope})
			assert.NoError(t, err)

			sortTargets(actualTargets)
			assert.Equal(t, tt.expectedIndexTarget, actualTargets, toIndexTargetCompareString(tt.expectedIndexTarget, actualTargets))
		})
	}
}

func sortTargets(targets []*models.IndexTargetTable) {
	lib.SortF(targets, func(t *models.IndexTargetTable) string {
		var namesI []string
		for _, f := range t.IndexFields {
			namesI = append(namesI, f.Name)
		}
		return strings.Join(namesI, "")
	})
}

func toIndexTargetCompareString(expected, actual []*models.IndexTargetTable) string {
	txt := "------\nexpected:\n"
	txt += lib.JoinF(expected, "\n", func(it *models.IndexTargetTable) string { return it.String() })
	txt += "\nvs\nactual:\n"
	txt += lib.JoinF(actual, "\n", func(it *models.IndexTargetTable) string { return it.String() })
	txt += "\n------\n"

	return txt
}
