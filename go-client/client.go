package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "local/grpc/proto"

	"google.golang.org/grpc"
)

const (
	defaultAddress   = "localhost:50051"
	defaultMaxLength = "200"
)

func getRandomNumber(address string, maxLength string) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewRandomClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GenerateRandomNumber(ctx, &pb.RandomNumberRequest{MaxLength: maxLength})
	if err != nil {
		log.Fatalf("could not get random number: %v", err)
	}
	log.Printf("Random number: %s", r.GetNumber())
}

func main() {
	maxLength := defaultMaxLength
	if len(os.Args) > 1 {
		maxLength = os.Args[1]

	}

	address := defaultAddress
	if len(os.Args) > 2 && os.Args[2] != "" {
		address = fmt.Sprintf("localhost:%s", os.Args[2])
	}

	log.Printf("Connecting to server: %s", address)

	for {
		getRandomNumber(address, maxLength)
		time.Sleep(time.Second)
	}
}