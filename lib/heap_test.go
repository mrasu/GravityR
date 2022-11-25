package lib_test

import (
	"github.com/mrasu/GravityR/lib"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLimitedHeap_Push(t *testing.T) {
	less := func(v1, v2 int) bool {
		return v1 < v2
	}
	tests := []struct {
		name string
		vals []int
		want []int
	}{
		{
			name: "Keep when not reach the limit",
			vals: []int{1, 5, 3, 2, 4},
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "Omit the smallest element when exceeds the limit",
			vals: []int{1, 5, 3, 2, 4, 6},
			want: []int{2, 3, 4, 5, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := lib.NewLimitedHeap(5, less)
			for _, v := range tt.vals {
				h.Push(v)
			}
			assert.Equal(t, tt.want, h.PopAll())
		})
	}
}
