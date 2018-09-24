package main

import (
	"log"
	"time"
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/fibonacci/fibonacci"
)

const (
	address = "172.17.0.6:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	var v int32
	v = 38

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	for i := 0; i < 1000; i++ {
		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
		t1 := time.Now()
		c.Fibo(ctx, &pb.FiboRequest{N: v})
		t2 := time.Now()
		x := float64(t2.Sub(t1).Nanoseconds()) / 1000000
		fmt.Printf("%F \n", x)
		//time.Sleep(parameters.REQUEST_TIME * time.Millisecond)
	}
}
