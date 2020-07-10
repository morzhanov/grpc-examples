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

func getParams(in *pb.RandomNumberRequest) (int, error) {
	maxLen, err := strconv.Atoi(in.GetMaxLength())
	if err != nil {
		log.Fatalf("failed: %v", err)
		return 0, err
	}
	return maxLen, err
}

// GetNumber implements random.GenerateRandomNumber (get example)
func (s *server) GenerateRandomNumber(ctx context.Context, in *pb.RandomNumberRequest) (*pb.RandomNumberReply, error) {
	log.Printf("GenerateRandomNumber received: %v", in.GetMaxLength())
	maxLen, err := getParams(in)
	if err != nil {
		return nil, err
	}

	number := rand.Intn(maxLen)

	log.Printf("GenerateRandomNumber: random number: %v", number)
	return &pb.RandomNumberReply{Number: strconv.Itoa(number)}, nil
}

// StreamNumbers implements random.StreamNumbers (server side stream)
func (s *server) StreamNumbers(in *pb.RandomNumberRequest, stream pb.Random_StreamNumbersServer) error {
	log.Printf("StreamNumbers received: %v", in.GetMaxLength())
	maxLen, err := getParams(in)
	if err != nil {
		return err
	}

	log.Printf("StreamNumbers sending stream of random numbers...")
	for i := 0; i < 2; i++ {
		number := rand.Intn(maxLen)
		if err := stream.Send(&pb.RandomNumberReply{Number: strconv.Itoa(number)}); err != nil {
			return err
		}
	}

	return nil
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
