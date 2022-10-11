package lib_test

import (
	"github.com/mrasu/GravityR/lib"
	"github.com/stretchr/testify/assert"
	"testing"
)

type data struct {
	val int
}

func TestStack_Push(t *testing.T) {
	tests := []struct {
		name           string
		vals           []*data
		expectedValues []*data
	}{
		{
			name:           "Push once",
			vals:           []*data{{1}},
			expectedValues: []*data{{1}},
		},
		{
			name:           "Push twice",
			vals:           []*data{{1}, {2}},
			expectedValues: []*data{{2}, {1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := lib.NewStack[data]()
			for _, v := range tt.vals {
				s.Push(v)
			}
			actualValues := popAllValues(s)
			assert.Equal(t, tt.expectedValues, actualValues)
		})
	}
}

func TestStack_Top(t *testing.T) {
	tests := []struct {
		name     string
		vals     []*data
		expected *data
	}{
		{
			name:     "Return the last element",
			vals:     []*data{{1}, {2}, {3}},
			expected: &data{3},
		},
		{
			name:     "Return nil when empty",
			vals:     []*data{},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := lib.NewStack[data]()
			for _, v := range tt.vals {
				s.Push(v)
			}

			assert.Equal(t, tt.expected, s.Top())
		})
	}
}

func popAllValues[T any](s *lib.Stack[T]) []*T {
	var res []*T
	for {
		v := s.Pop()
		if v == nil {
			break
		}
		res = append(res, v)
	}

	return res
}
