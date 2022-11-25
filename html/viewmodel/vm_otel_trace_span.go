package viewmodel

type VmOtelTraceSpan struct {
	Name            string             `json:"name"`
	SpanId          string             `json:"spanId"`
	StartTimeMillis int                `json:"startTimeMillis"`
	EndTimeMillis   int                `json:"endTimeMillis"`
	ServiceName     string             `json:"serviceName"`
	Children        []*VmOtelTraceSpan `json:"children"`
}
