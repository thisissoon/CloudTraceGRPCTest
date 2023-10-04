// Package main implements a server for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	cloudtrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/status"
)

// Serve is used to implement helloworld.GreeterServer.
type Server struct {
	pb.UnimplementedGreeterServer
	shouldErr bool
}

// SayHello implements helloworld.GreeterServer
func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	if s.shouldErr {
		return &pb.HelloReply{}, status.Error(codes.Internal, "force error")
	}
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *Server) Start(port int, service pb.GreeterServer) error {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
	)

	pb.RegisterGreeterServer(srv, service)
	log.Printf("server listening at %v", lis.Addr())
	return srv.Serve(lis)
}

func initTracer() (func(), error) {
	projectID := os.Getenv("PROJECT_ID")

	// Create Google Cloud Trace exporter to be able to retrieve
	// the collected spans.
	exporter, err := cloudtrace.New(cloudtrace.WithProjectID(projectID))
	if err != nil {
		return nil, err
	}
	resource, err := resource.New(
		context.Background(),
		resource.WithAttributes(semconv.ServiceName("GreeterService")))
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource))

	otel.SetTracerProvider(tp)
	return func() {
		err := tp.Shutdown(context.Background())
		if err != nil {
			fmt.Printf("error shutting down trace provider: %+v", err)
		}
	}, nil
}
