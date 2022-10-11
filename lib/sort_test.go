package lib_test

import (
	"github.com/mrasu/GravityR/lib"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSort(t *testing.T) {
	tests := []struct {
		name     string
		vals     []int
		expected []int
	}{
		{
			name:     "Sort values",
			vals:     []int{3, 1, 2},
			expected: []int{1, 2, 3},
		},
		{
			name:     "As is when empty",
			vals:     []int{},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lib.Sort(tt.vals, func(v int) int { return v })
			assert.Equal(t, tt.expected, tt.vals)
		})
	}
}
