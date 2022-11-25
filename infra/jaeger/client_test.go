package jaeger_test

import (
	"fmt"
	"github.com/gogo/protobuf/types"
	"github.com/google/go-cmp/cmp"
	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/proto-gen/api_v3"
	common "github.com/jaegertracing/jaeger/proto-gen/otel/common/v1"
	trace "github.com/jaegertracing/jaeger/proto-gen/otel/trace/v1"
	"github.com/mrasu/GravityR/infra/jaeger"
	"github.com/mrasu/GravityR/lib"
	jaegermock "github.com/mrasu/GravityR/mocks/jaeger"
	"github.com/mrasu/GravityR/thelper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
	"time"
)

func TestClient_FindTraces(t *testing.T) {
	now := time.Now()
	libInstrumentation := &common.InstrumentationLibrary{
		Name:    "awesome-library",
		Version: "v1.2.3",
	}

	pbTraces := []*model.Trace{
		{
			Spans: []*model.Span{
				{
					TraceID:       model.TraceID{Low: 1},
					SpanID:        model.SpanID(1),
					OperationName: "root-operation",
					StartTime:     now,
					Process: &model.Process{
						ServiceName: "service1",
					},
					Tags: []model.KeyValue{
						{Key: "otel.library.name", VType: model.ValueType_STRING, VStr: "awesome-library"},
						{Key: "otel.library.version", VType: model.ValueType_STRING, VStr: "v1.2.3"},
						{Key: "span.kind", VType: model.ValueType_STRING, VStr: "client"},
					},
					References: nil,
				},
				{
					TraceID:       model.TraceID{Low: 1},
					SpanID:        model.SpanID(2),
					OperationName: "root-operation2",
					StartTime:     now,
					Process: &model.Process{
						ServiceName: "service2",
					},
					Tags: []model.KeyValue{
						{Key: "otel.library.name", VType: model.ValueType_STRING, VStr: "awesome-library"},
						{Key: "otel.library.version", VType: model.ValueType_STRING, VStr: "v1.2.3"},
						{Key: "span.kind", VType: model.ValueType_STRING, VStr: "client"},
					},
					References: []model.SpanRef{
						{
							TraceID: model.TraceID{Low: 1},
							SpanID:  model.SpanID(1),
						},
					},
				},
			},
		},
		{
			Spans: []*model.Span{
				{
					TraceID:       model.TraceID{Low: 2},
					SpanID:        model.SpanID(1),
					OperationName: "root-operation",
					StartTime:     now,
					Process: &model.Process{
						ServiceName: "service1",
					},
					Tags: []model.KeyValue{
						{Key: "otel.library.name", VType: model.ValueType_STRING, VStr: "awesome-library"},
						{Key: "otel.library.version", VType: model.ValueType_STRING, VStr: "v1.2.3"},
						{Key: "span.kind", VType: model.ValueType_STRING, VStr: "client"},
					},
					References: nil,
				},
			},
		},
	}

	want := []*jaeger.Trace{
		{
			Spans: []*trace.ResourceSpans{
				{
					Resource: thelper.BuildJaegerResource("service1"),
					InstrumentationLibrarySpans: []*trace.InstrumentationLibrarySpans{
						{
							InstrumentationLibrary: libInstrumentation,
							Spans: []*trace.Span{
								{
									TraceId:           thelper.DecodeHex(t, "00000000000000000000000000000001"),
									SpanId:            thelper.DecodeHex(t, "0000000000000001"),
									StartTimeUnixNano: uint64(now.UnixNano()),
									EndTimeUnixNano:   uint64(now.UnixNano()),
									ParentSpanId:      thelper.DecodeHex(t, "0000000000000000"),
									Name:              "root-operation",
									Kind:              trace.Span_SPAN_KIND_CLIENT,
								},
							},
						},
					},
				},
				{
					Resource: thelper.BuildJaegerResource("service2"),
					InstrumentationLibrarySpans: []*trace.InstrumentationLibrarySpans{
						{
							InstrumentationLibrary: libInstrumentation,
							Spans: []*trace.Span{
								{
									TraceId:           thelper.DecodeHex(t, "00000000000000000000000000000001"),
									SpanId:            thelper.DecodeHex(t, "0000000000000002"),
									StartTimeUnixNano: uint64(now.UnixNano()),
									EndTimeUnixNano:   uint64(now.UnixNano()),
									ParentSpanId:      thelper.DecodeHex(t, "0000000000000001"),
									Name:              "root-operation2",
									Kind:              trace.Span_SPAN_KIND_CLIENT,
								},
							},
						},
					},
				},
			},
		},
		{
			Spans: []*trace.ResourceSpans{
				{
					Resource: thelper.BuildJaegerResource("service1"),
					InstrumentationLibrarySpans: []*trace.InstrumentationLibrarySpans{
						{
							InstrumentationLibrary: libInstrumentation,
							Spans: []*trace.Span{
								{
									TraceId:           thelper.DecodeHex(t, "00000000000000000000000000000002"),
									SpanId:            thelper.DecodeHex(t, "0000000000000001"),
									StartTimeUnixNano: uint64(now.UnixNano()),
									EndTimeUnixNano:   uint64(now.UnixNano()),
									ParentSpanId:      thelper.DecodeHex(t, "0000000000000000"),
									Name:              "root-operation",
									Kind:              trace.Span_SPAN_KIND_CLIENT,
								},
							},
						},
					},
				},
			},
		},
	}

	params := &api_v3.TraceQueryParameters{
		StartTimeMin: &types.Timestamp{Seconds: int64(now.Second())},
		StartTimeMax: &types.Timestamp{Seconds: int64(now.Second())},
	}

	mockReader := &jaegermock.Reader{}
	mockReader.On("FindTraces", mock.Anything, mock.Anything).Return(pbTraces, nil)

	thelper.UseGRPCServer(t, mockReader, func(addr net.Addr) {
		cli, err := jaeger.Open(addr.String(), false)
		require.NoError(t, err)
		defer cli.Close()

		traces, err := cli.FindTraces(params)
		require.NoError(t, err)
		sortedTraces := sortTraces(traces)

		if diff := cmp.Diff(want, sortedTraces); diff != "" {
			t.Errorf(diff)
		}
	})
}

func sortTraces(traces []*jaeger.Trace) []*jaeger.Trace {
	for _, tr := range traces {
		lib.Sort(tr.Spans, func(s *trace.ResourceSpans) string {
			return fmt.Sprintf("%x", s.InstrumentationLibrarySpans[0].Spans[0].SpanId)
		})
	}

	lib.Sort(traces, func(t *jaeger.Trace) string {
		return fmt.Sprintf("%x", t.Spans[0].InstrumentationLibrarySpans[0].Spans[0].TraceId)
	})

	res := make([]*jaeger.Trace, len(traces))
	copy(res, traces)
	return res
}
