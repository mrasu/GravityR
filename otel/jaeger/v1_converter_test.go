package jaeger_test

import (
	"github.com/google/go-cmp/cmp"
	common "github.com/jaegertracing/jaeger/proto-gen/otel/common/v1"
	resource "github.com/jaegertracing/jaeger/proto-gen/otel/resource/v1"
	trace "github.com/jaegertracing/jaeger/proto-gen/otel/trace/v1"
	"github.com/mrasu/GravityR/otel/jaeger"
	"github.com/mrasu/GravityR/otel/omodel"
	"github.com/mrasu/GravityR/thelper"
	"testing"
)

func TestToTraceTrees_convertFields(t *testing.T) {
	resourceSpans := []*trace.ResourceSpans{
		{
			Resource: &resource.Resource{
				Attributes: []*common.KeyValue{
					{Key: "service.name", Value: &common.AnyValue{Value: &common.AnyValue_StringValue{StringValue: "root-service"}}},
				},
			},
			InstrumentationLibrarySpans: []*trace.InstrumentationLibrarySpans{
				{
					Spans: []*trace.Span{
						{
							TraceId:           thelper.DecodeHex(t, "00000000000000000000000000000001"),
							SpanId:            thelper.DecodeHex(t, "0000000000000001"),
							ParentSpanId:      nil,
							Name:              "root-operation",
							StartTimeUnixNano: 1,
							EndTimeUnixNano:   2,
							Attributes: []*common.KeyValue{
								{Key: "string", Value: &common.AnyValue{Value: &common.AnyValue_StringValue{StringValue: "bar"}}},
								{Key: "bool", Value: &common.AnyValue{Value: &common.AnyValue_BoolValue{BoolValue: true}}},
								{Key: "int", Value: &common.AnyValue{Value: &common.AnyValue_IntValue{IntValue: 123}}},
								{Key: "double", Value: &common.AnyValue{Value: &common.AnyValue_DoubleValue{DoubleValue: 1.23}}},
								{Key: "array", Value: &common.AnyValue{Value: &common.AnyValue_ArrayValue{ArrayValue: &common.ArrayValue{
									Values: []*common.AnyValue{
										{Value: &common.AnyValue_StringValue{StringValue: "array1"}},
										{Value: &common.AnyValue_IntValue{IntValue: 2}},
									},
								}}}},
								{Key: "kvlist", Value: &common.AnyValue{Value: &common.AnyValue_KvlistValue{KvlistValue: &common.KeyValueList{
									Values: []*common.KeyValue{
										{Key: "key1", Value: &common.AnyValue{Value: &common.AnyValue_StringValue{StringValue: "val1"}}},
										{Key: "key2", Value: &common.AnyValue{Value: &common.AnyValue_IntValue{IntValue: 2}}},
									},
								}}}},
								{Key: "bytes", Value: &common.AnyValue{Value: &common.AnyValue_BytesValue{BytesValue: []byte{1, 2, 3}}}},
							},
						},
					},
				},
			},
		},
	}
	want := []*omodel.TraceTree{
		{
			Root: &omodel.Span{
				Name:              "root-operation",
				TraceId:           thelper.DecodeHex(t, "00000000000000000000000000000001"),
				SpanId:            thelper.DecodeHex(t, "0000000000000001"),
				StartTimeUnixNano: 1,
				EndTimeUnixNano:   2,
				SpanAttributes: map[string]*omodel.AnyValue{
					"string": {Val: &omodel.AnyValueString{StringValue: "bar"}},
					"bool":   {Val: &omodel.AnyValueBool{BoolValue: true}},
					"int":    {Val: &omodel.AnyValueInt{IntValue: 123}},
					"double": {Val: &omodel.AnyValueDouble{DoubleValue: 1.23}},
					"array": {Val: &omodel.AnyValueArray{
						ArrayValues: []omodel.AnyValueDatum{
							&omodel.AnyValueString{StringValue: "array1"},
							&omodel.AnyValueInt{IntValue: 2},
						},
					}},
					"kvlist": {Val: &omodel.AnyValueKV{
						KVValue: map[string]omodel.AnyValueDatum{
							"key1": &omodel.AnyValueString{StringValue: "val1"},
							"key2": &omodel.AnyValueInt{IntValue: 2},
						},
					}},
					"bytes": {Val: &omodel.AnyValueBytes{BytesValue: []byte{1, 2, 3}}},
				},
				ServiceName: "root-service",
			},
		},
	}

	traces := jaeger.ToTraceTrees(resourceSpans)
	if diff := cmp.Diff(want, traces); diff != "" {
		t.Errorf(diff)
	}
}

func TestToTraceTrees_convertChildren(t *testing.T) {
	resourceSpans := []*trace.ResourceSpans{
		{
			InstrumentationLibrarySpans: []*trace.InstrumentationLibrarySpans{
				{
					Spans: []*trace.Span{
						{
							TraceId:           thelper.DecodeHex(t, "00000000000000000000000000000001"),
							SpanId:            thelper.DecodeHex(t, "0000000000000001"),
							ParentSpanId:      nil,
							StartTimeUnixNano: 1,
							Name:              "operation1",
						},
					},
				},
			},
		},
		{
			InstrumentationLibrarySpans: []*trace.InstrumentationLibrarySpans{
				{
					Spans: []*trace.Span{
						{
							TraceId:           thelper.DecodeHex(t, "00000000000000000000000000000001"),
							SpanId:            thelper.DecodeHex(t, "0000000000000002"),
							ParentSpanId:      thelper.DecodeHex(t, "0000000000000001"),
							StartTimeUnixNano: 2,
							Name:              "operation1-2",
						},
					},
				},
			},
		},
		{
			InstrumentationLibrarySpans: []*trace.InstrumentationLibrarySpans{
				{
					Spans: []*trace.Span{
						{
							TraceId:           thelper.DecodeHex(t, "00000000000000000000000000000001"),
							SpanId:            thelper.DecodeHex(t, "0000000000000003"),
							ParentSpanId:      thelper.DecodeHex(t, "0000000000000002"),
							StartTimeUnixNano: 3,
							Name:              "operation1-2-3",
						},
					},
				},
			},
		},
		{
			InstrumentationLibrarySpans: []*trace.InstrumentationLibrarySpans{
				{
					Spans: []*trace.Span{
						{
							TraceId:           thelper.DecodeHex(t, "00000000000000000000000000000001"),
							SpanId:            thelper.DecodeHex(t, "0000000000000004"),
							ParentSpanId:      thelper.DecodeHex(t, "0000000000000003"),
							StartTimeUnixNano: 4,
							Name:              "operation1-2-3-4",
						},
					},
				},
			},
		},
		{
			InstrumentationLibrarySpans: []*trace.InstrumentationLibrarySpans{
				{
					Spans: []*trace.Span{
						{
							TraceId:           thelper.DecodeHex(t, "00000000000000000000000000000001"),
							SpanId:            thelper.DecodeHex(t, "0000000000000005"),
							ParentSpanId:      thelper.DecodeHex(t, "0000000000000001"),
							StartTimeUnixNano: 5,
							Name:              "operation1-5",
						},
					},
				},
			},
		},
		{
			InstrumentationLibrarySpans: []*trace.InstrumentationLibrarySpans{
				{
					Spans: []*trace.Span{
						{
							TraceId:           thelper.DecodeHex(t, "00000000000000000000000000000001"),
							SpanId:            thelper.DecodeHex(t, "0000000000000006"),
							ParentSpanId:      thelper.DecodeHex(t, "0000000000000002"),
							StartTimeUnixNano: 6,
							Name:              "operation1-2-6",
						},
					},
				},
			},
		},
	}

	op := &omodel.Span{
		Name:              "operation1",
		TraceId:           thelper.DecodeHex(t, "00000000000000000000000000000001"),
		SpanId:            thelper.DecodeHex(t, "0000000000000001"),
		StartTimeUnixNano: 1,
		SpanAttributes:    map[string]*omodel.AnyValue{},
		Children: []*omodel.Span{
			{
				Name:              "operation1-2",
				TraceId:           thelper.DecodeHex(t, "00000000000000000000000000000001"),
				SpanId:            thelper.DecodeHex(t, "0000000000000002"),
				SpanAttributes:    map[string]*omodel.AnyValue{},
				StartTimeUnixNano: 2,
				Children: []*omodel.Span{
					{
						Name:              "operation1-2-3",
						TraceId:           thelper.DecodeHex(t, "00000000000000000000000000000001"),
						SpanId:            thelper.DecodeHex(t, "0000000000000003"),
						SpanAttributes:    map[string]*omodel.AnyValue{},
						StartTimeUnixNano: 3,
						Children: []*omodel.Span{
							{

								Name:              "operation1-2-3-4",
								TraceId:           thelper.DecodeHex(t, "00000000000000000000000000000001"),
								SpanId:            thelper.DecodeHex(t, "0000000000000004"),
								SpanAttributes:    map[string]*omodel.AnyValue{},
								StartTimeUnixNano: 4,
							},
						},
					},
					{
						Name:              "operation1-2-6",
						TraceId:           thelper.DecodeHex(t, "00000000000000000000000000000001"),
						SpanId:            thelper.DecodeHex(t, "0000000000000006"),
						SpanAttributes:    map[string]*omodel.AnyValue{},
						StartTimeUnixNano: 6,
					},
				},
			},
			{
				Name:              "operation1-5",
				TraceId:           thelper.DecodeHex(t, "00000000000000000000000000000001"),
				SpanId:            thelper.DecodeHex(t, "0000000000000005"),
				SpanAttributes:    map[string]*omodel.AnyValue{},
				StartTimeUnixNano: 5,
			},
		},
	}
	op.Children[0].Parent = op
	op.Children[1].Parent = op
	op.Children[0].Children[0].Parent = op.Children[0]
	op.Children[0].Children[1].Parent = op.Children[0]
	op.Children[0].Children[0].Children[0].Parent = op.Children[0].Children[0]

	want := []*omodel.TraceTree{{Root: op}}
	traces := jaeger.ToTraceTrees(resourceSpans)
	if diff := cmp.Diff(want, traces); diff != "" {
		t.Errorf(diff)
	}
}

func TestToTraceTrees_sortChildrenByStartTime(t *testing.T) {
	resourceSpans := []*trace.ResourceSpans{
		{
			InstrumentationLibrarySpans: []*trace.InstrumentationLibrarySpans{
				{
					Spans: []*trace.Span{
						{
							TraceId:           thelper.DecodeHex(t, "00000000000000000000000000000001"),
							SpanId:            thelper.DecodeHex(t, "0000000000000002"),
							ParentSpanId:      thelper.DecodeHex(t, "0000000000000001"),
							StartTimeUnixNano: 4,
							Name:              "operation1-2",
						},
					},
				},
			},
		},
		{
			InstrumentationLibrarySpans: []*trace.InstrumentationLibrarySpans{
				{
					Spans: []*trace.Span{
						{
							TraceId:           thelper.DecodeHex(t, "00000000000000000000000000000001"),
							SpanId:            thelper.DecodeHex(t, "0000000000000003"),
							ParentSpanId:      thelper.DecodeHex(t, "0000000000000001"),
							StartTimeUnixNano: 3,
							Name:              "operation1-3",
						},
					},
				},
			},
		},
		{
			InstrumentationLibrarySpans: []*trace.InstrumentationLibrarySpans{
				{
					Spans: []*trace.Span{
						{
							TraceId:           thelper.DecodeHex(t, "00000000000000000000000000000001"),
							SpanId:            thelper.DecodeHex(t, "0000000000000004"),
							ParentSpanId:      thelper.DecodeHex(t, "0000000000000001"),
							StartTimeUnixNano: 10,
							Name:              "operation1-4",
						},
					},
				},
			},
		},
		{
			InstrumentationLibrarySpans: []*trace.InstrumentationLibrarySpans{
				{
					Spans: []*trace.Span{
						{
							TraceId:           thelper.DecodeHex(t, "00000000000000000000000000000001"),
							SpanId:            thelper.DecodeHex(t, "0000000000000001"),
							ParentSpanId:      nil,
							StartTimeUnixNano: 1,
							Name:              "operation1",
						},
					},
				},
			},
		},
	}

	op := &omodel.Span{
		Name:              "operation1",
		TraceId:           thelper.DecodeHex(t, "00000000000000000000000000000001"),
		SpanId:            thelper.DecodeHex(t, "0000000000000001"),
		StartTimeUnixNano: 1,
		SpanAttributes:    map[string]*omodel.AnyValue{},
		Children: []*omodel.Span{
			{
				Name:              "operation1-3",
				TraceId:           thelper.DecodeHex(t, "00000000000000000000000000000001"),
				SpanId:            thelper.DecodeHex(t, "0000000000000003"),
				SpanAttributes:    map[string]*omodel.AnyValue{},
				StartTimeUnixNano: 3,
			},
			{
				Name:              "operation1-2",
				TraceId:           thelper.DecodeHex(t, "00000000000000000000000000000001"),
				SpanId:            thelper.DecodeHex(t, "0000000000000002"),
				SpanAttributes:    map[string]*omodel.AnyValue{},
				StartTimeUnixNano: 4,
			},
			{
				Name:              "operation1-4",
				TraceId:           thelper.DecodeHex(t, "00000000000000000000000000000001"),
				SpanId:            thelper.DecodeHex(t, "0000000000000004"),
				SpanAttributes:    map[string]*omodel.AnyValue{},
				StartTimeUnixNano: 10,
			},
		},
	}
	op.Children[0].Parent = op
	op.Children[1].Parent = op
	op.Children[2].Parent = op

	want := []*omodel.TraceTree{{Root: op}}
	traces := jaeger.ToTraceTrees(resourceSpans)
	if diff := cmp.Diff(want, traces); diff != "" {
		t.Errorf(diff)
	}
}
