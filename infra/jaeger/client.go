package jaeger

import (
	"context"
	"github.com/jaegertracing/jaeger/proto-gen/api_v3"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"time"
)

type Client struct {
	conn *grpc.ClientConn
}

func Open(address string, isSecure bool) (*Client, error) {
	opts := []grpc.DialOption{
		grpc.WithBlock(),
	}
	if !isSecure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, address, opts...)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot connect to %s", address)
	}

	return NewClient(conn), nil
}

func NewClient(conn *grpc.ClientConn) *Client {
	return &Client{conn: conn}
}

func (c *Client) Close() error {
	err := c.conn.Close()
	if err != nil {
		return errors.Wrap(err, "failed to close grpc connection")
	}
	return nil
}

func (c *Client) FindTraces(param *api_v3.TraceQueryParameters) ([]*Trace, error) {
	client := api_v3.NewQueryServiceClient(c.conn)
	ctx := context.Background()
	req := &api_v3.FindTracesRequest{
		Query: param,
	}
	resp, err := client.FindTraces(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find traces")
	}

	var traces []*Trace
	for {
		chunk, err := resp.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, errors.Wrap(err, "failed to receive result to find traces")
		}
		spans := chunk.GetResourceSpans()

		traces = append(traces, &Trace{Spans: spans})
	}

	return traces, nil
}
