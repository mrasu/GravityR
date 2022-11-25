package thelper

import (
	"github.com/jaegertracing/jaeger/cmd/query/app/apiv3"
	"github.com/jaegertracing/jaeger/cmd/query/app/querysvc"
	"github.com/jaegertracing/jaeger/proto-gen/api_v3"
	common "github.com/jaegertracing/jaeger/proto-gen/otel/common/v1"
	v1 "github.com/jaegertracing/jaeger/proto-gen/otel/resource/v1"
	"github.com/jaegertracing/jaeger/storage/spanstore"
	"github.com/stretchr/testify/require"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"google.golang.org/grpc"
	"net"
	"testing"
)

func UseGRPCServer(t *testing.T, reader spanstore.Reader, fn func(net.Addr)) {
	t.Helper()

	s := grpc.NewServer()
	qs := querysvc.NewQueryService(
		reader, nil, querysvc.QueryServiceOptions{},
	)
	api_v3.RegisterQueryServiceServer(s, &apiv3.Handler{QueryService: qs})
	lis, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	go func() {
		err := s.Serve(lis)
		require.NoError(t, err)
	}()
	defer s.Stop()

	fn(lis.Addr())
}

func BuildJaegerResource(serviceName string) *v1.Resource {
	return &v1.Resource{
		Attributes: []*common.KeyValue{
			{
				Key: string(semconv.ServiceNameKey),
				Value: &common.AnyValue{
					Value: &common.AnyValue_StringValue{StringValue: serviceName},
				},
			},
		},
	}
}
