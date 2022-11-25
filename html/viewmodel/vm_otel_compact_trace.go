package viewmodel

type VmOtelCompactTrace struct {
	TraceId                  string           `json:"traceId"`
	SameServiceAccessCount   int              `json:"sameServiceAccessCount"`
	TimeConsumingServiceName string           `json:"timeConsumingServiceName"`
	Root                     *VmOtelTraceSpan `json:"root"`
}
