package dig

import (
	"github.com/jaegertracing/jaeger/model"
	"github.com/mrasu/GravityR/infra/jaeger"
	jaegermock "github.com/mrasu/GravityR/mocks/jaeger"
	"github.com/mrasu/GravityR/thelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net"
	"os"
	"testing"
	"time"
)

func Test_runJaegerDig(t *testing.T) {
	pbTraces := buildDummyPbTraces()
	mockReader := &jaegermock.Reader{}
	mockReader.On("FindTraces", mock.Anything, mock.Anything).Return(pbTraces, nil)

	thelper.InjectClientDist()
	thelper.UseGRPCServer(t, mockReader, func(addr net.Addr) {
		now := time.Now()
		jr := jaegerRunner{
			serviceName:          "service1",
			grpcAddress:          addr.String(),
			durationMilli:        0,
			sameServiceThreshold: 1,
			count:                1000,
			uiURL:                "http://localhost:16686",
			start:                now.Add(-30 * time.Minute),
			end:                  now,
		}
		cli, err := jaeger.Open(addr.String(), false)
		require.NoError(t, err)

		thelper.CreateTemp(t, "tmp.html", func(tmpfile *os.File) {
			err := jr.dig(tmpfile.Name(), cli)
			require.NoError(t, err)

			html, err := os.ReadFile(tmpfile.Name())
			require.NoError(t, err)
			assert.Contains(t, string(html), `uiPath":"http://localhost:16686"`)
			assert.Contains(t, string(html), `traceId":"00000000000000000000000000000001"`)
		})
	})
}

func buildDummyPbTraces() []*model.Trace {
	now := time.Now()

	return []*model.Trace{
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

}
