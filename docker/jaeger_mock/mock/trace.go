package main

import (
	"github.com/jaegertracing/jaeger/model"
	"time"
)

type Trace struct {
	Spans []*Span
}

type Span struct {
	TraceID          string
	SpanID           string
	OperationName    string
	StartTime        time.Time
	DurationMilliSec int
	Tags             map[string]any
	ServiceName      string
	SpanRef          *SpanRef
}

type SpanRef struct {
	TraceID string
	SpanID  string
}

func (t *Trace) ToJTrace() (*model.Trace, error) {
	var spans []*model.Span
	for _, s := range t.Spans {
		span, err := s.toJSpan()
		if err != nil {
			return nil, err
		}
		spans = append(spans, span)
	}

	return &model.Trace{Spans: spans}, nil
}

func (s *Span) toJSpan() (*model.Span, error) {
	traceID, spanID, err := s.toJTraceAndJSpanID(s.TraceID, s.SpanID)
	if err != nil {
		return nil, err
	}

	var tags []model.KeyValue
	for k, v := range s.Tags {
		tags = append(tags, s.jKeyValue(k, v))
	}

	refs, err := s.jReferences()
	if err != nil {
		return nil, err
	}

	return &model.Span{
		TraceID:       traceID,
		SpanID:        spanID,
		OperationName: s.OperationName,
		StartTime:     s.StartTime,
		Duration:      time.Duration(s.DurationMilliSec) * time.Millisecond,
		Tags:          tags,
		Process: &model.Process{
			ServiceName: s.ServiceName,
		},
		References: refs,
	}, nil
}

func (s *Span) jKeyValue(key string, val any) model.KeyValue {
	switch v := val.(type) {
	case int:
		return model.KeyValue{Key: key, VInt64: int64(v)}
	case int64:
		return model.KeyValue{Key: key, VInt64: v}
	case int32:
		return model.KeyValue{Key: key, VInt64: int64(v)}
	case float64:
		return model.KeyValue{Key: key, VFloat64: v}
	case float32:
		return model.KeyValue{Key: key, VFloat64: float64(v)}
	case bool:
		return model.KeyValue{Key: key, VBool: v}
	case string:
		return model.KeyValue{Key: key, VStr: v}
	default:
		panic("unexpected key value")
	}
}

func (s *Span) jReferences() ([]model.SpanRef, error) {
	if s.SpanRef == nil {
		return nil, nil
	}

	traceID, spanID, err := s.toJTraceAndJSpanID(s.SpanRef.TraceID, s.SpanRef.SpanID)
	if err != nil {
		return nil, err
	}

	return []model.SpanRef{
		{
			TraceID: traceID,
			SpanID:  spanID,
		},
	}, nil

}

func (s *Span) toJTraceAndJSpanID(traceID, spanID string) (model.TraceID, model.SpanID, error) {
	jTraceID, err := model.TraceIDFromString(traceID)
	if err != nil {
		return model.TraceID{}, 0, err
	}
	jSpanID, err := model.SpanIDFromString(spanID)
	if err != nil {
		return model.TraceID{}, 0, err
	}

	return jTraceID, jSpanID, nil
}
