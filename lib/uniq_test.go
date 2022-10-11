package lib_test

import (
	"github.com/mrasu/GravityR/lib"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUniqBy(t *testing.T) {
	type data struct{ val int }

	tests := []struct {
		name     string
		vals     []data
		fn       func(data) int
		expected []data
	}{
		{
			name:     "Remove duplications",
			vals:     []data{{1}, {2}, {3}, {2}, {1}},
			fn:       func(d data) int { return d.val },
			expected: []data{{1}, {2}, {3}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := lib.UniqBy(tt.vals, tt.fn)
			assert.ElementsMatch(t, tt.expected, res)
		})
	}
}
