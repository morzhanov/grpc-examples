package main

import (
	"context"
	"io"
	"log"
	"math/rand"
	"net"

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

// GetNumber implements random.GenerateRandomNumber (get example)
func (s *server) GenerateRandomNumber(ctx context.Context, in *pb.RandomNumberRequest) (*pb.RandomNumberReply, error) {
	log.Printf("GenerateRandomNumber received: %v", in.GetMaxLength())
	maxLen := int(in.GetMaxLength())
	number := rand.Intn(maxLen)

	log.Printf("GenerateRandomNumber: random number: %v", number)
	return &pb.RandomNumberReply{Number: int32(number)}, nil
}

// StreamNumbers implements random.StreamNumbers (server side stream)
func (s *server) StreamNumbers(in *pb.RandomNumberRequest, stream pb.Random_StreamNumbersServer) error {
	log.Printf("StreamNumbers received: %v", in.GetMaxLength())
	maxLen := int(in.GetMaxLength())

	log.Printf("StreamNumbers sending stream of random numbers...")
	for i := 0; i < 2; i++ {
		number := rand.Intn(maxLen)
		if err := stream.Send(&pb.RandomNumberReply{Number: int32(number)}); err != nil {
			return err
		}
	}

	return nil
}

// LogStreamOfRandomNumbers implements random.LogStreamOfRandomNumbers (client side stream)
func (s *server) LogStreamOfRandomNumbers(stream pb.Random_LogStreamOfRandomNumbersServer) error {
	for {
		number, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.LogRandomNumberReply{
				Message: "Numbers successfully logged out",
			})
		}
		if err != nil {
			return err
		}
		log.Printf("LogStreamOfRandomNumbers: received random number: %v", number)
	}
}

// BidirectionalStream implements random.BidirectionalStream (bidirectional stream)
func (s *server) BidirectionalStream(stream pb.Random_BidirectionalStreamServer) error {
	for i := 100; i < 201; i += 100 {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		log.Printf("BidirectionalStream got number: %v", in.Number)

		number := rand.Intn(i)
		if err := stream.Send(&pb.BidirectionalMessage{Number: int32(number)}); err != nil {
			return err
		}
	}

	log.Printf("BidirectionalStream done")
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
