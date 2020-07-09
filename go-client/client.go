package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "local/grpc/proto"

	"google.golang.org/grpc"
)

const (
	address          = "localhost:50051"
	defaultMaxLength = "200"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewRandomClient(conn)

	maxLength := defaultMaxLength
	if len(os.Args) > 1 {
		maxLength = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GenerateRandomNumber(ctx, &pb.RandomNumberRequest{MaxLength: maxLength})
	if err != nil {
		log.Fatalf("could not get random number: %v", err)
	}
	log.Printf("Random number: %s", r.GetNumber())
}
