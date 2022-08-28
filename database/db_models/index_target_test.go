package db_models_test

import (
	"github.com/mrasu/GravityR/database/db_models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseIndexTarget(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected *db_models.IndexTarget
	}{
		{
			name:     "single column",
			text:     "users:status",
			expected: buildIndexTarget(t, "users", []string{"status"}),
		},
		{
			name:     "multiple columns",
			text:     "users:status+name",
			expected: buildIndexTarget(t, "users", []string{"status", "name"}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			it, err := db_models.NewIndexTarget(tt.text)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, it)
		})
	}
}

func buildIndexTarget(t *testing.T, tName string, cNames []string) *db_models.IndexTarget {
	t.Helper()

	var cs []*db_models.IndexColumn
	for _, cName := range cNames {
		c, err := db_models.NewIndexColumn(cName)
		require.NoError(t, err)
		cs = append(cs, c)
	}

	return &db_models.IndexTarget{TableName: "users", Columns: cs}
}
