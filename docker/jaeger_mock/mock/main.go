package main

import (
	"context"
	"encoding/json"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/jaegertracing/jaeger/cmd/query/app/apiv3"
	"github.com/jaegertracing/jaeger/cmd/query/app/querysvc"
	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/proto-gen/api_v3"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name: "serve",
				Action: func(*cli.Context) error {
					return serve()
				},
			},
			{
				Name: "show-traces",
				Action: func(*cli.Context) error {
					return showTraces("localhost:16685")
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func serve() error {
	port := 16685

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	s := newGRPCServer()
	spanReader := &mockSpanReader{getTracesFn: readTraces}
	depReader := &mockDependencyReader{}
	queryOpts := querysvc.QueryServiceOptions{}
	qs := querysvc.NewQueryService(
		spanReader, depReader, queryOpts,
	)
	api_v3.RegisterQueryServiceServer(s, &apiv3.Handler{QueryService: qs})

	log.Printf("start Jaeger server at %v\n", port)
	return s.Serve(listener)
}

func newGRPCServer() *grpc.Server {
	return grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
					log.Printf("access to %s\n", info.FullMethod)
					return handler(ctx, req)
				},
			),
		),
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
					log.Printf("access to %s\n", info.FullMethod)
					return handler(srv, ss)
				},
			),
		),
	)
}

const dataDir = "mockdata"

func readTraces() ([]*model.Trace, error) {
	entries, err := os.ReadDir(dataDir)
	if err != nil {
		return nil, err
	}

	var jTraces []*model.Trace
	for _, e := range entries {
		data, err := ioutil.ReadFile(path.Join(dataDir, e.Name()))
		if err != nil {
			return nil, err
		}
		var tt *Trace
		err = json.Unmarshal(data, &tt)
		if err != nil {
			return nil, err
		}

		jt, err := tt.ToJTrace()
		if err != nil {
			return nil, err
		}
		jTraces = append(jTraces, jt)
	}

	return jTraces, nil
}

func showTraces(address string) error {
	d := Dumper{address: address}
	return d.showAsJSON()
}
