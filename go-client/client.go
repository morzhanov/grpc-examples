package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	pb "local/grpc/proto"

	"google.golang.org/grpc"
)

const (
	defaultAddress   = "localhost:50051"
	defaultMaxLength = 200
)

func connectToServer(address string) (pb.RandomClient, context.Context, *grpc.ClientConn, context.CancelFunc) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := pb.NewRandomClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	return c, ctx, conn, cancel
}

func getRandomNumber(address string, maxLength int32) {
	c, ctx, conn, cancel := connectToServer(address)
	defer cancel()
	defer conn.Close()

	r, err := c.GenerateRandomNumber(ctx, &pb.RandomNumberRequest{MaxLength: maxLength})
	if err != nil {
		log.Fatalf("could not get random number: %v", err)
	}
	log.Printf("Random number: %v", r.GetNumber())
}

func getStreamRandomNumbers(address string, maxLength int32) {
	c, ctx, conn, cancel := connectToServer(address)
	defer cancel()
	defer conn.Close()

	stream, err := c.StreamNumbers(ctx, &pb.RandomNumberRequest{MaxLength: maxLength})
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	for {
		number, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.StreamNumbers(_) = _, %v", c, err)
		}
		log.Printf("Server Stream result: %v", number)
	}
}

func sendStreamOfRandomNumbers(address string, maxLength int32) {
	c, ctx, conn, cancel := connectToServer(address)
	defer cancel()
	defer conn.Close()

	stream, err := c.LogStreamOfRandomNumbers(ctx)
	if err != nil {
		log.Fatalf("%v.LogStreamOfRandomNumbers(_) = _, %v", c, err)
	}
	for i := 0; i < 2; i++ {
		number := rand.Intn(int(maxLength))
		if err := stream.Send(&pb.LogRandomNumberRequest{Number: int32(number)}); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("%v.Send(%v) = %v", stream, number, err)
		}
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}
	log.Printf("Client Stream result: %v", reply)
}

func bidirectionalCommunication(address string, maxLength int32) {
	c, ctx, conn, cancel := connectToServer(address)
	defer cancel()
	defer conn.Close()

	stream, err := c.BidirectionalStream(ctx)
	if err != nil {
		log.Fatalf("%v.BidirectionalStream(_) = _, %v", c, err)
	}

	waitc := make(chan struct{})

	go func() {
		for i := 0; i < 2; i++ {
			in, err := stream.Recv()
			if err != nil {
				log.Fatalf("Failed to receive a number : %v", err)
			}
			log.Printf("BidirectionalStream got number %d", in.Number)
		}
		log.Printf("BidirectionalStream done")
		close(waitc)
	}()

	for i := 0; i < 2; i++ {
		number := rand.Intn(int(maxLength))
		if err := stream.Send(&pb.BidirectionalMessage{Number: int32(number)}); err != nil {
			log.Fatalf("Failed to send a number: %v", err)
		}
	}

	stream.CloseSend()
	<-waitc
}

func main() {
	maxLength := int32(defaultMaxLength)
	if len(os.Args) > 1 {
		param, err := strconv.Atoi(os.Args[1])
		if err != nil {
			panic(err)
		}
		maxLength = int32(param)
	}

	address := defaultAddress
	if len(os.Args) > 2 && os.Args[2] != "" {
		address = fmt.Sprintf("localhost:%s", os.Args[2])
	}

	log.Printf("Connecting to server: %s", address)

	for {
		log.Println("SIMPLE RPC")
		getRandomNumber(address, maxLength)
		log.Println("SERVER SIDE STREAMING RPC")
		getStreamRandomNumbers(address, maxLength)
		log.Println("CLIENT SIDE STREAMING RPC")
		sendStreamOfRandomNumbers(address, maxLength)
		log.Println("BIDIRECTIONAL STREAMING RPC")
		bidirectionalCommunication(address, maxLength)
		time.Sleep(time.Second)
	}
}
