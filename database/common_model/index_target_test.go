package common_model_test

import (
	"github.com/mrasu/GravityR/database/common_model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseIndexTarget(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected *common_model.IndexTarget
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
			it, err := common_model.NewIndexTarget(tt.text)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, it)
		})
	}
}

func buildIndexTarget(t *testing.T, tName string, cNames []string) *common_model.IndexTarget {
	t.Helper()

	var cs []*common_model.IndexColumn
	for _, cName := range cNames {
		c, err := common_model.NewIndexColumn(cName)
		require.NoError(t, err)
		cs = append(cs, c)
	}

	return &common_model.IndexTarget{TableName: "users", Columns: cs}
}
