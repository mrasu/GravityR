package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gogo/protobuf/types"
	"github.com/jaegertracing/jaeger/proto-gen/api_v3"
	v1 "github.com/jaegertracing/jaeger/proto-gen/otel/common/v1"
	v1Trace "github.com/jaegertracing/jaeger/proto-gen/otel/trace/v1"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"reflect"
	"time"
)

type Dumper struct {
	address string
}

func (d *Dumper) showAsJSON() error {
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial(d.address, opts...)
	if err != nil {
		return err
	}
	defer conn.Close()

	return d.findTraces(conn, func(spans []*v1Trace.ResourceSpans) error {
		a := d.toTrace(spans)
		data, err := json.Marshal(a)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		fmt.Println("===========-")
		return nil
	})
}

func (d *Dumper) findTraces(conn *grpc.ClientConn, fn func([]*v1Trace.ResourceSpans) error) error {
	client := api_v3.NewQueryServiceClient(conn)
	ctx := context.Background()
	now := time.Now()
	req := &api_v3.FindTracesRequest{
		Query: &api_v3.TraceQueryParameters{
			ServiceName:  "bff",
			StartTimeMin: &types.Timestamp{Seconds: now.Add(-30 * time.Minute).Unix()},
			StartTimeMax: &types.Timestamp{Seconds: now.Unix()},
			NumTraces:    10,
		},
	}
	resp, err := client.FindTraces(ctx, req)
	if err != nil {
		return err
	}

	for {
		chunk, err := resp.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		spans := chunk.GetResourceSpans()

		if err := fn(spans); err != nil {
			return err
		}
	}
	return nil
}

func (d *Dumper) toTrace(spans []*v1Trace.ResourceSpans) *Trace {
	var ss []*Span
	for _, s := range spans {
		ss = append(ss, d.toSpan(s))
	}

	return &Trace{Spans: ss}
}

func (d *Dumper) toSpan(span *v1Trace.ResourceSpans) *Span {
	serviceName := ""
	if span.Resource != nil {
		for _, r := range span.Resource.Attributes {
			if r.Key == string(semconv.ServiceNameKey) {
				serviceName = r.Value.GetStringValue()
				break
			}
		}
	}

	s := span.InstrumentationLibrarySpans[0].Spans[0]

	attrs := map[string]any{}
	for _, attr := range s.Attributes {
		attrs[attr.Key] = d.toAnyValue(attr.Value)
	}

	return &Span{
		TraceID:          fmt.Sprintf("%x", s.TraceId),
		SpanID:           fmt.Sprintf("%x", s.SpanId),
		OperationName:    s.Name,
		StartTime:        time.UnixMilli(int64(s.StartTimeUnixNano) / 1000 / 1000),
		DurationMilliSec: int((s.EndTimeUnixNano - s.StartTimeUnixNano) / 1000 / 1000),
		Tags:             attrs,
		ServiceName:      serviceName,
		SpanRef: &SpanRef{
			TraceID: fmt.Sprintf("%x", s.TraceId),
			SpanID:  fmt.Sprintf("%x", s.ParentSpanId),
		},
	}
}

func (d *Dumper) toAnyValue(v *v1.AnyValue) any {
	switch vv := v.Value.(type) {
	case *v1.AnyValue_StringValue:
		return vv.StringValue
	case *v1.AnyValue_BoolValue:
		return vv.BoolValue
	case *v1.AnyValue_IntValue:
		return vv.IntValue
	case *v1.AnyValue_DoubleValue:
		return vv.DoubleValue
	default:
		// This must not be happened
		panic(fmt.Sprintf("not implemented any value: %s", reflect.TypeOf(vv).Name()))
	}
}
