package lib_test

import (
	"github.com/mrasu/GravityR/lib"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSet_Add(t *testing.T) {
	tests := []struct {
		name           string
		origins        []int
		val            int
		expectedValues []int
	}{
		{
			name:           "Add a value when no value exist",
			origins:        []int{},
			val:            1,
			expectedValues: []int{1},
		},
		{
			name:           "Add a value when it doesn't exist",
			origins:        []int{1, 2},
			val:            3,
			expectedValues: []int{1, 2, 3},
		},
		{
			name:           "Not add a value when it exists",
			origins:        []int{1, 2},
			val:            1,
			expectedValues: []int{1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := lib.NewSetS(tt.origins)
			s.Add(tt.val)
			assert.ElementsMatch(t, tt.expectedValues, s.Values())
		})
	}
}

func TestSet_Merge(t *testing.T) {
	tests := []struct {
		name           string
		set1           *lib.Set[int]
		set2           *lib.Set[int]
		expectedValues []int
	}{
		{
			name:           "Add all values from another set",
			set1:           lib.NewSetS([]int{1, 2}),
			set2:           lib.NewSetS([]int{3, 4}),
			expectedValues: []int{1, 2, 3, 4},
		},
		{
			name:           "Not add same values",
			set1:           lib.NewSetS([]int{1, 2}),
			set2:           lib.NewSetS([]int{2}),
			expectedValues: []int{1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.set1.Merge(tt.set2)
			assert.ElementsMatch(t, tt.expectedValues, tt.set1.Values())
		})
	}
}

func TestSet_Delete(t *testing.T) {
	tests := []struct {
		name           string
		set            *lib.Set[int]
		val            int
		expectedValues []int
	}{
		{
			name:           "Remove elements from set",
			set:            lib.NewSetS([]int{1, 2, 3, 4}),
			val:            1,
			expectedValues: []int{2, 3, 4},
		},
		{
			name:           "Do nothing when the value doesn't exist",
			set:            lib.NewSetS([]int{1, 2, 3, 4}),
			val:            9999,
			expectedValues: []int{1, 2, 3, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.set.Delete(tt.val)
			assert.ElementsMatch(t, tt.expectedValues, tt.set.Values())
		})
	}
}

func TestSet_Contains(t *testing.T) {
	tests := []struct {
		name     string
		set      *lib.Set[int]
		val      int
		expected bool
	}{
		{
			name:     "Return true when exists",
			set:      lib.NewSetS([]int{1, 2, 3, 4}),
			val:      1,
			expected: true,
		},
		{
			name:     "Return false when not exist",
			set:      lib.NewSetS([]int{1, 2, 3, 4}),
			val:      1111,
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.set.Contains(tt.val)
			assert.Equal(t, tt.expected, res)
		})
	}
}

func TestSet_Count(t *testing.T) {
	tests := []struct {
		name     string
		set      *lib.Set[int]
		expected int
	}{
		{
			name:     "Return the number of elements",
			set:      lib.NewSetS([]int{1, 2, 3, 4}),
			expected: 4,
		},
		{
			name:     "Return zero when no value exists",
			set:      lib.NewSet[int](),
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.set.Count()
			assert.Equal(t, tt.expected, res)
		})
	}
}
