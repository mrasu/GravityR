package lib

import "fmt"

type Set[T comparable] struct {
	vals map[T]bool
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{map[T]bool{}}
}

func NewSetS[T comparable](origins []T) *Set[T] {
	vals := map[T]bool{}
	for _, o := range origins {
		vals[o] = true
	}
	return &Set[T]{vals}
}

func (s *Set[T]) Add(val T) {
	s.vals[val] = true
}

func (s *Set[T]) Merge(os *Set[T]) {
	for _, val := range os.Values() {
		s.Add(val)
	}
}

func (s *Set[T]) Delete(val T) {
	delete(s.vals, val)
}

func (s *Set[T]) Values() []T {
	var vals []T
	for val, _ := range s.vals {
		vals = append(vals, val)
	}

	return vals
}

func (s *Set[T]) Contains(val T) bool {
	if _, ok := s.vals[val]; ok {
		return true
	}
	return false
}

func (s *Set[T]) Count() int {
	return len(s.vals)
}

func (s *Set[T]) String() string {
	return fmt.Sprintf("%v", s.Values())
}
