package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	flag.Parse()
	shutdown, err := initTracer()
	if err != nil {
		return err
	}
	defer shutdown()
	server := &Server{
		shouldErr: true,
	}
	go func() {
		if err := server.Start(*port, server); err != nil {
			log.Fatal(err)
		}
	}()
	time.Sleep(time.Second)
	return sayHello(*port)
}

func sayHello(port int) error {
	addr := fmt.Sprintf(":%d", port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "world"})
	if err != nil {
		return fmt.Errorf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
	return nil
}
