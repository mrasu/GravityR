package lib_test

import (
	"github.com/mrasu/GravityR/lib"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestJoin(t *testing.T) {
	tests := []struct {
		name     string
		vals     []int
		sep      string
		expected string
	}{
		{
			name:     "Concat text from the result of function",
			vals:     []int{1, 2, 3},
			sep:      "*",
			expected: "1*2*3",
		},
		{
			name:     "Empty when arg is empty",
			vals:     []int{},
			sep:      "*",
			expected: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := lib.Join(tt.vals, tt.sep, func(i int) string { return strconv.Itoa(i) })
			assert.Equal(t, tt.expected, res)
		})
	}
}
