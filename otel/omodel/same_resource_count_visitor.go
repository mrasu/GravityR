package omodel

import (
	"github.com/samber/lo"
)

type sameResourceCountVisitor struct {
	visitedCounts map[string]int
}

func (v *sameResourceCountVisitor) Enter(span *Span) {
	service := span.ServiceName
	v.visitedCounts[service] += 1
}

func (v *sameResourceCountVisitor) Leave(*Span) {}

func (v *sameResourceCountVisitor) GetMax() int {
	return lo.Max(lo.Values(v.visitedCounts))
}
