package models_test

import (
	"github.com/mrasu/GravityR/database/mysql/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseIndexTarget(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected *models.IndexTarget
	}{
		{
			name: "single column",
			text: "users:status",
			expected: &models.IndexTarget{TableName: "users", Columns: []*models.IndexColumn{
				{Name: "status"},
			}},
		},
		{
			name: "multiple columns",
			text: "users:status+name",
			expected: &models.IndexTarget{TableName: "users", Columns: []*models.IndexColumn{
				{Name: "status"},
				{Name: "name"},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			it, err := models.ParseIndexTarget(tt.text)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, it)
		})
	}
}

func TestParseIndexTargetWithBacktick(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected *models.IndexTarget
	}{
		{
			name: "single column",
			text: "users`:status",
			expected: &models.IndexTarget{TableName: "users", Columns: []*models.IndexColumn{
				{Name: "status"},
			}},
		},
		{
			name: "multiple columns",
			text: "users`:status`+name",
			expected: &models.IndexTarget{TableName: "users", Columns: []*models.IndexColumn{
				{Name: "status"},
				{Name: "name"},
			}},
		},
		{
			name: "special characters exist",
			text: "users,+users`:status,status`+name,name",
			expected: &models.IndexTarget{TableName: "users,+users", Columns: []*models.IndexColumn{
				{Name: "status,status"},
				{Name: "name,name"},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			it, err := models.ParseIndexTargetWithBacktick(tt.text)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, it)
		})
	}
}
