package jaeger

import (
	v1Trace "github.com/jaegertracing/jaeger/proto-gen/otel/trace/v1"
)

type Trace struct {
	Spans []*v1Trace.ResourceSpans
}
