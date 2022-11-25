package jaeger_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/jaegertracing/jaeger/model"
	iJaeger "github.com/mrasu/GravityR/infra/jaeger"
	jaegermock "github.com/mrasu/GravityR/mocks/jaeger"
	"github.com/mrasu/GravityR/otel/jaeger"
	"github.com/mrasu/GravityR/otel/omodel"
	"github.com/mrasu/GravityR/thelper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
	"time"
)

func TestTraceFetcher_FetchCompactedTraces(t *testing.T) {
	now := time.Now()
	cases := map[string]struct {
		resp []*model.Trace
		want []*omodel.TraceTree
	}{
		"returns TraceTree": {
			resp: []*model.Trace{
				{
					Spans: []*model.Span{
						{
							TraceID:       model.TraceID{Low: 1},
							SpanID:        model.SpanID(2),
							OperationName: "root-operation",
							StartTime:     now,
							Duration:      10 * time.Millisecond,
							Process: &model.Process{
								ServiceName: "my-service",
							},
							Tags: []model.KeyValue{
								{Key: "foo", VType: model.ValueType_STRING, VStr: "bar"},
								{Key: "hello", VType: model.ValueType_INT64, VInt64: 123},
							},
							References: nil,
						},
					},
				},
			},
			want: []*omodel.TraceTree{
				{
					Root: &omodel.Span{
						Name:              "root-operation",
						TraceId:           thelper.DecodeHex(t, "00000000000000000000000000000001"),
						SpanId:            thelper.DecodeHex(t, "0000000000000002"),
						StartTimeUnixNano: uint64(now.UnixNano()),
						EndTimeUnixNano:   uint64(now.Add(10 * time.Millisecond).UnixNano()),
						SpanAttributes: map[string]*omodel.AnyValue{
							"foo":   {Val: &omodel.AnyValueString{StringValue: "bar"}},
							"hello": {Val: &omodel.AnyValueInt{IntValue: 123}},
						},
						ServiceName: "my-service",
					},
				},
			},
		},
		"removes the same service traces except db query": {
			resp: []*model.Trace{
				{
					Spans: []*model.Span{
						{
							SpanID:        model.SpanID(1),
							OperationName: "root-operation",
							StartTime:     now,
							Process: &model.Process{
								ServiceName: "service1",
							},
							References: nil,
						},
						{
							SpanID:        model.SpanID(2),
							OperationName: "root-operation2",
							StartTime:     now,
							Process: &model.Process{
								ServiceName: "service1",
							},
							References: []model.SpanRef{
								{
									SpanID: model.SpanID(1),
								},
							},
						},
						{
							SpanID:        model.SpanID(3),
							OperationName: "next-service-operation",
							StartTime:     now,
							Process: &model.Process{
								ServiceName: "service2",
							},
							References: []model.SpanRef{
								{
									SpanID: model.SpanID(2),
								},
							},
						},
						{
							SpanID:        model.SpanID(4),
							OperationName: "next-service-operation2",
							StartTime:     now,
							Process: &model.Process{
								ServiceName: "service2",
							},
							Tags: []model.KeyValue{
								{Key: "db.system", VType: model.ValueType_STRING, VStr: "awesome-db"},
							},
							References: []model.SpanRef{
								{
									SpanID: model.SpanID(3),
								},
							},
						},
					},
				},
			},
			want: []*omodel.TraceTree{
				{
					Root: &omodel.Span{
						Name:              "root-operation",
						TraceId:           thelper.DecodeHex(t, "00000000000000000000000000000000"),
						SpanId:            thelper.DecodeHex(t, "0000000000000001"),
						StartTimeUnixNano: uint64(now.UnixNano()),
						EndTimeUnixNano:   uint64(now.UnixNano()),
						SpanAttributes:    map[string]*omodel.AnyValue{},
						ServiceName:       "service1",
						Children: []*omodel.Span{
							{
								Name:              "next-service-operation",
								TraceId:           thelper.DecodeHex(t, "00000000000000000000000000000000"),
								SpanId:            thelper.DecodeHex(t, "0000000000000003"),
								StartTimeUnixNano: uint64(now.UnixNano()),
								EndTimeUnixNano:   uint64(now.UnixNano()),
								SpanAttributes:    map[string]*omodel.AnyValue{},
								ServiceName:       "service2",
								Children: []*omodel.Span{
									{
										Name:              "next-service-operation2",
										TraceId:           thelper.DecodeHex(t, "00000000000000000000000000000000"),
										SpanId:            thelper.DecodeHex(t, "0000000000000004"),
										StartTimeUnixNano: uint64(now.UnixNano()),
										EndTimeUnixNano:   uint64(now.UnixNano()),
										SpanAttributes: map[string]*omodel.AnyValue{
											"db.system": {Val: &omodel.AnyValueString{StringValue: "awesome-db"}},
										},
										ServiceName: "awesome-db (service2)",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			useJaegerClient(t, tt.resp, func(cli *iJaeger.Client) {
				tf := jaeger.NewTraceFetcher(cli)
				traces, err := tf.FetchCompactedTraces(100, time.Now(), time.Now(), "bff", 10*time.Millisecond)
				require.NoError(t, err)

				if diff := cmp.Diff(tt.want, traces); diff != "" {
					t.Errorf(diff)
				}
			})
		})
	}
}

func useJaegerClient(t *testing.T, traces []*model.Trace, fn func(*iJaeger.Client)) {
	t.Helper()

	mockReader := &jaegermock.Reader{}
	mockReader.On("FindTraces", mock.Anything, mock.Anything).Return(traces, nil)

	thelper.UseGRPCServer(t, mockReader, func(addr net.Addr) {
		cli, err := iJaeger.Open(addr.String(), false)
		require.NoError(t, err)
		defer cli.Close()

		fn(cli)
	})
}
