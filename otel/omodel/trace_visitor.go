package omodel

type traceVisitor interface {
	Enter(*Span)
	Leave(*Span)
}
