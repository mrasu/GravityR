package jaeger

import (
	"github.com/mrasu/GravityR/lib"
	"github.com/mrasu/GravityR/otel/omodel"
	"github.com/samber/lo"
)

type TracePicker struct {
	traces []*omodel.TraceTree
}

func NewTracePicker(traces []*omodel.TraceTree) *TracePicker {
	return &TracePicker{traces: traces}
}

func (tp *TracePicker) PickSameServiceAccessTraces(threshold int) []*omodel.TraceTree {
	return lo.Filter(tp.traces, func(tree *omodel.TraceTree, _ int) bool {
		return tree.MaxSameServiceAccessCount() > threshold
	})
}

func (tp *TracePicker) PickSlowTraces(num int) []*omodel.TraceTree {
	slowTraces := lib.NewLimitedHeap[*omodel.TraceTree](num, func(a1, a2 *omodel.TraceTree) bool {
		d1 := a1.Root.EndTimeUnixNano - a1.Root.StartTimeUnixNano
		d2 := a2.Root.EndTimeUnixNano - a2.Root.StartTimeUnixNano
		return d1 < d2
	})

	for _, tree := range tp.traces {
		slowTraces.Push(tree)
	}

	return lo.Reverse(slowTraces.PopAll())
}
