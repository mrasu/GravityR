package omodel

type traceInfoCollectVisitor struct {
	countVisitor        *sameResourceCountVisitor
	serviceConsumedTime map[string]uint64
}

func (v *traceInfoCollectVisitor) Enter(span *Span) {
	v.countVisitor.Enter(span)
	v.serviceConsumedTime[span.ServiceName] += span.CalcConsumedTime()
}

func (v *traceInfoCollectVisitor) Leave(span *Span) {
	v.countVisitor.Leave(span)
}

func (v *traceInfoCollectVisitor) GetMaxResourceAccessCount() int {
	return v.countVisitor.GetMax()
}

func (v *traceInfoCollectVisitor) GetMostTimeConsumedService() string {
	maxService := ""
	maxTime := uint64(0)

	for s, t := range v.serviceConsumedTime {
		if maxTime < t {
			maxService = s
			maxTime = t
		}
	}

	return maxService
}
