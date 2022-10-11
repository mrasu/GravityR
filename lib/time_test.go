package lib_test

import (
	"github.com/mrasu/GravityR/lib"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNormalizeTimeByHour(t *testing.T) {
	tests := []struct {
		name     string
		t        time.Time
		expected time.Time
	}{
		{
			name:     "Make zeros under hour",
			t:        time.Date(2001, 2, 3, 4, 5, 6, 7, time.UTC),
			expected: time.Date(2001, 2, 3, 4, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, lib.NormalizeTimeByHour(tt.t))
		})
	}
}
