package omodel_test

import (
	"github.com/mrasu/GravityR/otel/omodel"
	"github.com/stretchr/testify/assert"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"testing"
)

const dbSystemKey = string(semconv.DBSystemKey)

func TestSpan_GetDBSystem(t *testing.T) {
	cases := map[string]struct {
		attrs map[string]*omodel.AnyValue
		want  string
	}{
		"returns dbname": {
			attrs: map[string]*omodel.AnyValue{
				dbSystemKey: {Val: &omodel.AnyValueString{StringValue: "some-db"}},
			},
			want: "some-db",
		},
		"empty when no db.system": {
			attrs: map[string]*omodel.AnyValue{},
			want:  "",
		},
		"empty when db.system is not string": {
			attrs: map[string]*omodel.AnyValue{
				dbSystemKey: {Val: &omodel.AnyValueInt{IntValue: 123}},
			},
			want: "",
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			s := &omodel.Span{SpanAttributes: tt.attrs}

			assert.Equal(t, tt.want, s.GetDBSystem())
		})
	}
}

func TestSpan_IsLastDbSpan(t *testing.T) {
	cases := map[string]struct {
		attrs    map[string]*omodel.AnyValue
		children []*omodel.Span
		want     bool
	}{
		"true when the span is the last and for db.system": {
			attrs: map[string]*omodel.AnyValue{
				dbSystemKey: {Val: &omodel.AnyValueString{StringValue: "some-db"}},
			},
			children: nil,
			want:     true,
		},
		"false when the span is not the last": {
			attrs: map[string]*omodel.AnyValue{
				dbSystemKey: {Val: &omodel.AnyValueString{StringValue: "some-db"}},
			},
			children: []*omodel.Span{{}},
			want:     false,
		},
		"false when the span doesn't have db.system": {
			attrs: map[string]*omodel.AnyValue{
				"aaa": {Val: &omodel.AnyValueString{StringValue: "some-db"}},
			},
			children: nil,
			want:     false,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			s := &omodel.Span{SpanAttributes: tt.attrs, Children: tt.children}

			assert.Equal(t, tt.want, s.IsLastDbSpan())
		})
	}
}

func TestSpan_ShallowCopyWithoutDependency(t *testing.T) {
	span := &omodel.Span{
		Parent:            &omodel.Span{},
		Name:              "dummy",
		TraceId:           []byte{1},
		SpanId:            []byte{2},
		StartTimeUnixNano: 3,
		EndTimeUnixNano:   4,
		SpanAttributes: map[string]*omodel.AnyValue{
			"foo": {Val: &omodel.AnyValueString{StringValue: "bar"}},
		},
		ServiceName: "test-service",
		Children:    []*omodel.Span{{}},
	}

	cases := map[string]struct {
		span *omodel.Span
		want *omodel.Span
	}{
		"creates new Span without parent and children": {
			span: span,
			want: &omodel.Span{
				Parent:            nil,
				Name:              span.Name,
				TraceId:           span.TraceId,
				SpanId:            span.SpanId,
				StartTimeUnixNano: span.StartTimeUnixNano,
				EndTimeUnixNano:   span.EndTimeUnixNano,
				SpanAttributes:    span.SpanAttributes,
				ServiceName:       span.ServiceName,
				Children:          nil,
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.span.ShallowCopyWithoutDependency())
		})
	}
}

func TestSpan_IsSame(t *testing.T) {
	cases := map[string]struct {
		span      *omodel.Span
		otherSpan *omodel.Span
		want      bool
	}{
		"true when TraceID and SpanID is the same": {
			span: &omodel.Span{
				TraceId: []byte{1},
				SpanId:  []byte{2},
			},
			otherSpan: &omodel.Span{
				TraceId: []byte{1},
				SpanId:  []byte{2},
			},
			want: true,
		},
		"false when TraceID is not the same": {
			span: &omodel.Span{
				TraceId: []byte{1},
				SpanId:  []byte{2},
			},
			otherSpan: &omodel.Span{
				TraceId: []byte{9},
				SpanId:  []byte{2},
			},
			want: false,
		},
		"false when SpanID is not the same": {
			span: &omodel.Span{
				TraceId: []byte{1},
				SpanId:  []byte{2},
			},
			otherSpan: &omodel.Span{
				TraceId: []byte{1},
				SpanId:  []byte{9},
			},
			want: false,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.span.IsSame(tt.otherSpan))
		})
	}
}

func TestSpan_CalcConsumedTime(t *testing.T) {
	cases := map[string]struct {
		span *omodel.Span
		want uint64
	}{
		"returns consumed time when no children": {
			span: &omodel.Span{
				StartTimeUnixNano: 1,
				EndTimeUnixNano:   101,
			},
			want: uint64(100),
		},
		"reduces child's time": {
			span: &omodel.Span{
				StartTimeUnixNano: 1,
				EndTimeUnixNano:   101,
				Children: []*omodel.Span{
					{
						StartTimeUnixNano: 11,
						EndTimeUnixNano:   21,
					},
				},
			},
			want: uint64(90),
		},
		"not reduce both when children runs parallel": {
			span: &omodel.Span{
				StartTimeUnixNano: 1,
				EndTimeUnixNano:   101,
				Children: []*omodel.Span{
					{
						StartTimeUnixNano: 11,
						EndTimeUnixNano:   21,
					},
					{
						StartTimeUnixNano: 15,
						EndTimeUnixNano:   17,
					},
				},
			},
			want: uint64(90),
		},
		"not reduce accumulated when children runs parallel": {
			span: &omodel.Span{
				StartTimeUnixNano: 1,
				EndTimeUnixNano:   101,
				Children: []*omodel.Span{
					{
						StartTimeUnixNano: 11,
						EndTimeUnixNano:   21,
					},
					{
						StartTimeUnixNano: 15,
						EndTimeUnixNano:   31,
					},
				},
			},
			want: uint64(80),
		},
		"reduces both when children runs series": {
			span: &omodel.Span{
				StartTimeUnixNano: 1,
				EndTimeUnixNano:   101,
				Children: []*omodel.Span{
					{
						StartTimeUnixNano: 11,
						EndTimeUnixNano:   21,
					},
					{
						StartTimeUnixNano: 31,
						EndTimeUnixNano:   41,
					},
				},
			},
			want: uint64(80),
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.span.CalcConsumedTime())
		})
	}
}
