package database_test

import (
	"github.com/mrasu/GravityR/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewIndexTargetFromText(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected *database.IndexTarget
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
			it, err := database.NewIndexTargetFromText(tt.text)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, it)
		})
	}
}

func TestIndexTarget_HasSameIdxColumns(t *testing.T) {
	tests := []struct {
		name     string
		s        *database.IndexTarget
		t        *database.IndexTarget
		expected bool
	}{
		{
			name:     "same column",
			s:        database.NewIndexTarget("table", []string{"col1", "col2"}),
			t:        database.NewIndexTarget("table", []string{"col1", "col2"}),
			expected: true,
		},
		{
			name:     "same column but different order",
			s:        database.NewIndexTarget("table", []string{"col1", "col2"}),
			t:        database.NewIndexTarget("table", []string{"col2", "col1"}),
			expected: false,
		},
		{
			name:     "less column",
			s:        database.NewIndexTarget("table", []string{"col1", "col2"}),
			t:        database.NewIndexTarget("table", []string{"col1"}),
			expected: false,
		},
		{
			name:     "more column",
			s:        database.NewIndexTarget("table", []string{"col1", "col2"}),
			t:        database.NewIndexTarget("table", []string{"col1", "col2", "col3"}),
			expected: false,
		},
		{
			name:     "different column",
			s:        database.NewIndexTarget("table", []string{"col1", "col2"}),
			t:        database.NewIndexTarget("table", []string{"col3", "col4"}),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.s.HasSameIdxColumns(tt.t)
			assert.Equal(t, tt.expected, res)
		})
	}
}

func buildIndexTarget(t *testing.T, tName string, cNames []string) *database.IndexTarget {
	t.Helper()

	var cs []*database.IndexColumn
	for _, cName := range cNames {
		c, err := database.NewIndexColumn(cName)
		require.NoError(t, err)
		cs = append(cs, c)
	}

	return &database.IndexTarget{TableName: tName, Columns: cs}
}
