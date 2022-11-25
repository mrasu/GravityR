package omodel

import (
	"fmt"
	"github.com/mrasu/GravityR/lib"
)

type traceCompacter struct {
	parentStack *lib.Stack[Span]
}

func newTraceCompacter(top *Span) *traceCompacter {
	return &traceCompacter{
		parentStack: lib.NewStack[Span](top),
	}
}

func (c *traceCompacter) Enter(span *Span) {
	if !c.remains(span) {
		return
	}

	newParent := span.ShallowCopyWithoutDependency()
	if span.IsLastDbSpan() {
		newParent.ServiceName = fmt.Sprintf("%s (%s)", span.GetDBSystem(), newParent.ServiceName)
	}

	parent := c.parentStack.Top()
	if parent != nil {
		parent.Children = append(parent.Children, newParent)
	}

	c.parentStack.Push(newParent)
}

func (c *traceCompacter) Leave(span *Span) {
	if !c.parentStack.Top().IsSame(span) {
		return
	}

	c.parentStack.Pop()
}

func (c *traceCompacter) remains(span *Span) bool {
	if span.Parent == nil {
		return true
	}
	if span.Parent.ServiceName != span.ServiceName {
		return true
	}
	if span.IsLastDbSpan() {
		return true
	}

	return false
}
