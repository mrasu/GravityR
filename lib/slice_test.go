package lib_test

import (
	"github.com/mrasu/GravityR/lib"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSliceOrEmpty(t *testing.T) {
	tests := []struct {
		name     string
		vals     []int
		expected []int
	}{
		{
			name:     "Return as is when not empty",
			vals:     []int{1, 2, 3},
			expected: []int{1, 2, 3},
		},
		{
			name:     "Return empty slice when empty",
			vals:     []int{},
			expected: []int{},
		},
		{
			name:     "Return empty slice when nil",
			vals:     nil,
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, lib.SliceOrEmpty(tt.vals))
		})
	}
}
