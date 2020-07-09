package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"strconv"

	pb "local/grpc/proto"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server is used to implement random.RandomServer
type server struct {
	pb.UnimplementedRandomServer
}

// GetNumber implements random.GenerateRandomNumber
func (s *server) GenerateRandomNumber(ctx context.Context, in *pb.RandomNumberRequest) (*pb.RandomNumberReply, error) {
	log.Printf("Received: %v", in.GetMaxLength())

	maxLen, err := strconv.Atoi(in.GetMaxLength())
	if err != nil {
		log.Fatalf("failed: %v", err)
		return nil, err
	}

	number := rand.Intn(maxLen)
	log.Printf("Generated random number: %v", number)
	return &pb.RandomNumberReply{Number: strconv.Itoa(number)}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterRandomServer(s, &server{})
	log.Printf("Server started at: localhost%v", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
