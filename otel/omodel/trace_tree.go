package omodel

import (
	"fmt"
	"github.com/mrasu/GravityR/html/viewmodel"
)

type TraceTree struct {
	Root *Span
}

func NewTraceTree(root *Span) *TraceTree {
	return &TraceTree{Root: root}
}

func (tt *TraceTree) visit(v traceVisitor) {
	tt.visitRecursive(v, tt.Root)
}

func (tt *TraceTree) visitRecursive(v traceVisitor, current *Span) {
	v.Enter(current)
	for _, c := range current.Children {
		tt.visitRecursive(v, c)
	}
	v.Leave(current)
}

// Compact returns new TraceTree having only the first traces of services and the last traces going to DB which doesn't using otel.
func (tt *TraceTree) Compact() *TraceTree {
	top := &Span{}
	tt.visit(newTraceCompacter(top))

	if len(top.Children) > 0 {
		return &TraceTree{Root: top.Children[0]}
	} else {
		return &TraceTree{Root: nil}
	}
}

func (tt *TraceTree) MaxSameServiceAccessCount() int {
	v := &sameResourceCountVisitor{
		visitedCounts: map[string]int{},
	}
	tt.visit(v)

	return v.GetMax()
}

func (tt *TraceTree) ToCompactViewModel() *viewmodel.VmOtelCompactTrace {
	v := &traceInfoCollectVisitor{
		countVisitor: &sameResourceCountVisitor{
			visitedCounts: map[string]int{},
		},
		serviceConsumedTime: map[string]uint64{},
	}
	tt.visit(v)

	return &viewmodel.VmOtelCompactTrace{
		TraceId:                  fmt.Sprintf("%x", tt.Root.TraceId),
		SameServiceAccessCount:   v.GetMaxResourceAccessCount(),
		TimeConsumingServiceName: v.GetMostTimeConsumedService(),
		Root:                     tt.Root.ToViewModel(),
	}
}
