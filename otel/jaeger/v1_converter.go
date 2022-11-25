package jaeger

import (
	v1 "github.com/jaegertracing/jaeger/proto-gen/otel/common/v1"
	v1Trace "github.com/jaegertracing/jaeger/proto-gen/otel/trace/v1"
	"github.com/mrasu/GravityR/otel/omodel"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"reflect"
)

type v1Converter struct{}

type spanHolder struct {
	orig *v1Trace.Span
	span *omodel.Span
}

func ToTraceTrees(resourceSpans []*v1Trace.ResourceSpans) []*omodel.TraceTree {
	c := &v1Converter{}
	return c.toTraceTrees(resourceSpans)
}

func (c *v1Converter) toTraceTrees(resourceSpans []*v1Trace.ResourceSpans) []*omodel.TraceTree {
	spanIdMap := map[string]*spanHolder{}
	for _, rSpan := range resourceSpans {
		for _, lSpan := range rSpan.InstrumentationLibrarySpans {
			for _, span := range lSpan.Spans {
				spanIdMap[string(span.SpanId)] = &spanHolder{orig: span, span: c.toOmodelSpan(span, rSpan)}
			}
		}
	}

	var roots []*omodel.Span
	for _, h := range spanIdMap {
		if c.isRootSpan(h.orig) {
			roots = append(roots, h.span)
			continue
		}

		if parent, ok := spanIdMap[string(h.orig.ParentSpanId)]; ok {
			parent.span.Children = append(parent.span.Children, h.span)
			h.span.Parent = parent.span
		} else {
			roots = append(roots, h.span)
		}
	}

	return lo.Map(roots, func(root *omodel.Span, _ int) *omodel.TraceTree {
		root.SortChildren()
		return omodel.NewTraceTree(root)
	})
}

func (c *v1Converter) toOmodelSpan(s *v1Trace.Span, rs *v1Trace.ResourceSpans) *omodel.Span {
	serviceName := ""
	if rs.Resource != nil {
		for _, r := range rs.Resource.Attributes {
			if r.Key == string(semconv.ServiceNameKey) {
				serviceName = r.Value.GetStringValue()
				break
			}
		}
	}

	attrs := map[string]*omodel.AnyValue{}
	for _, attr := range s.Attributes {
		attrs[attr.Key] = c.toAnyValue(attr.Value)
	}

	return &omodel.Span{
		Name:              s.Name,
		TraceId:           s.TraceId,
		SpanId:            s.SpanId,
		StartTimeUnixNano: s.StartTimeUnixNano,
		EndTimeUnixNano:   s.EndTimeUnixNano,
		SpanAttributes:    attrs,
		ServiceName:       serviceName,
	}
}

func (c *v1Converter) toAnyValue(v *v1.AnyValue) *omodel.AnyValue {
	return &omodel.AnyValue{Val: c.toAnyValueData(v)}
}

func (c *v1Converter) toAnyValueData(v *v1.AnyValue) omodel.AnyValueDatum {
	switch vv := v.Value.(type) {
	case *v1.AnyValue_StringValue:
		return &omodel.AnyValueString{StringValue: vv.StringValue}
	case *v1.AnyValue_BoolValue:
		return &omodel.AnyValueBool{BoolValue: vv.BoolValue}
	case *v1.AnyValue_IntValue:
		return &omodel.AnyValueInt{IntValue: vv.IntValue}
	case *v1.AnyValue_DoubleValue:
		return &omodel.AnyValueDouble{DoubleValue: vv.DoubleValue}
	case *v1.AnyValue_ArrayValue:
		var vals []omodel.AnyValueDatum
		for _, v := range vv.ArrayValue.Values {
			vals = append(vals, c.toAnyValueData(v))
		}
		return &omodel.AnyValueArray{ArrayValues: vals}
	case *v1.AnyValue_KvlistValue:
		vals := map[string]omodel.AnyValueDatum{}
		for _, v := range vv.KvlistValue.Values {
			vals[v.Key] = c.toAnyValueData(v.Value)
		}
		return &omodel.AnyValueKV{KVValue: vals}
	case *v1.AnyValue_BytesValue:
		return &omodel.AnyValueBytes{BytesValue: vv.BytesValue}
	default:
		// This must not be happened
		panic(errors.Errorf("not implemented any value: %s", reflect.TypeOf(vv).Name()))
	}
}

func (c *v1Converter) isRootSpan(span *v1Trace.Span) bool {
	for _, b := range span.ParentSpanId {
		if b != 0 {
			return false
		}
	}
	return true
}
