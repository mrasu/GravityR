package lib

import (
	"container/heap"
)

var _ heap.Interface = &anyHeap[any]{}

type LimitedHeap[T any] struct {
	h     *anyHeap[T]
	limit int
}

func NewLimitedHeap[T any](limit int, less func(a, b T) bool) *LimitedHeap[T] {
	h := &anyHeap[T]{less: less}
	heap.Init(h)
	return &LimitedHeap[T]{
		h:     h,
		limit: limit,
	}
}

func (h LimitedHeap[T]) Len() int {
	return h.h.Len()
}

func (h LimitedHeap[T]) Push(x T) {
	heap.Push(h.h, x)
	if h.h.Len() > h.limit {
		heap.Pop(h.h)
	}
}

func (h LimitedHeap[T]) PopAll() []T {
	var res []T
	for h.h.Len() > 0 {
		res = append(res, heap.Pop(h.h).(T))
	}

	return res
}

type anyHeap[T any] struct {
	vals []T
	less func(a, b T) bool
}

func (h *anyHeap[T]) Len() int {
	return len(h.vals)
}

func (h *anyHeap[T]) Less(i, j int) bool {
	return h.less(h.vals[i], h.vals[j])
}

func (h *anyHeap[T]) Swap(i, j int) {
	h.vals[i], h.vals[j] = h.vals[j], h.vals[i]
}

func (h *anyHeap[T]) Push(x any) {
	v := x.(T)
	h.vals = append(h.vals, v)
}

func (h *anyHeap[T]) Pop() any {
	l := len(h.vals) - 1
	v := h.vals[l]
	h.vals = h.vals[0:l]

	return v
}
