package dmodel_test

import (
	"github.com/mrasu/GravityR/database/dmodel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewIndexTargetFromText(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected *dmodel.IndexTarget
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
			it, err := dmodel.NewIndexTargetFromText(tt.text)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, it)
		})
	}
}

func TestIndexTarget_HasSameIdxColumns(t *testing.T) {
	tests := []struct {
		name     string
		s        *dmodel.IndexTarget
		t        *dmodel.IndexTarget
		expected bool
	}{
		{
			name:     "same column",
			s:        dmodel.NewIndexTarget("table", []string{"col1", "col2"}),
			t:        dmodel.NewIndexTarget("table", []string{"col1", "col2"}),
			expected: true,
		},
		{
			name:     "same column but different order",
			s:        dmodel.NewIndexTarget("table", []string{"col1", "col2"}),
			t:        dmodel.NewIndexTarget("table", []string{"col2", "col1"}),
			expected: false,
		},
		{
			name:     "less column",
			s:        dmodel.NewIndexTarget("table", []string{"col1", "col2"}),
			t:        dmodel.NewIndexTarget("table", []string{"col1"}),
			expected: false,
		},
		{
			name:     "more column",
			s:        dmodel.NewIndexTarget("table", []string{"col1", "col2"}),
			t:        dmodel.NewIndexTarget("table", []string{"col1", "col2", "col3"}),
			expected: false,
		},
		{
			name:     "different column",
			s:        dmodel.NewIndexTarget("table", []string{"col1", "col2"}),
			t:        dmodel.NewIndexTarget("table", []string{"col3", "col4"}),
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

func buildIndexTarget(t *testing.T, tName string, cNames []string) *dmodel.IndexTarget {
	t.Helper()

	var cs []*dmodel.IndexColumn
	for _, cName := range cNames {
		c, err := dmodel.NewIndexColumn(cName)
		require.NoError(t, err)
		cs = append(cs, c)
	}

	return &dmodel.IndexTarget{TableName: tName, Columns: cs}
}
