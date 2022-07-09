package lib

type Stack[T any] struct {
	vals []*T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

func (s *Stack[T]) Push(val *T) {
	s.vals = append(s.vals, val)
}

func (s *Stack[T]) Pop() *T {
	v := s.vals[len(s.vals)-1]
	s.vals = s.vals[:len(s.vals)-1]

	return v
}

func (s *Stack[T]) Top() *T {
	if len(s.vals) == 0 {
		return nil
	}

	return s.vals[len(s.vals)-1]
}

func (s *Stack[T]) Size() int {
	return len(s.vals)
}
