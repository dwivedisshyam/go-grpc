package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	pb "github.com/go-grpc/calculator/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var addr string = "localhost:50052"

type Server struct {
	pb.CalculatorServiceServer
}

func main() {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)

	}
	log.Printf("Listening on %s\n", addr)
	c := grpc.NewServer()

	pb.RegisterCalculatorServiceServer(c, &Server{})
	reflection.Register(c)

	if err = c.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}

}

func (s *Server) Sum(ctx context.Context, req *pb.SumRequest) (*pb.SumResponse, error) {
	return &pb.SumResponse{
		Result: req.FirstNumber + req.SecondNumber,
	}, nil
}

func (s *Server) PrimeNumberDecomposition(req *pb.PrimeNumberDecompositionRequest, stream pb.CalculatorService_PrimeNumberDecompositionServer) error {

	k := int64(2)

	n := req.Number

	for n > 1 {
		if n%k == 0 {
			stream.Send(&pb.PrimeNumberDecompositionResponse{
				PrimeFactor: k,
			})
			n = n / k
		} else {
			k++
		}
	}

	return nil
}

func (s *Server) Average(stream pb.CalculatorService_AverageServer) error {
	log.Println("Average invoked")

	var sum int64
	var length int64 = 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("Sum:", sum, "Count:", length)
			stream.SendAndClose(&pb.AvgResponse{Avg: float64(sum) / float64(length)})
			return nil
		}

		if err != nil {
			log.Fatal("error calculating average", err)
			return err
		}

		sum += int64(req.Number)
		length++
	}
}

func (s *Server) CurrentMax(stream pb.CalculatorService_CurrentMaxServer) error {
	log.Println("CurrentMax invoked")

	var currentMax int32 = -1
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Println("error while receving", err)
			return nil
		}

		if req.Number > currentMax {
			currentMax = req.Number
		}

		err = stream.Send(&pb.CurrentMaxResponse{Number: currentMax})
		if err != nil {
			log.Println("failed to send the response")
		}
	}
}

func (s *Server) Sqrt(ctx context.Context, req *pb.SqrtRequest) (*pb.SqrtResponse, error) {
	log.Println("sqrt invoked")

	if req.Number < 0 {
		return nil, status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("Got a negative number: %d", req.Number),
		)
	}

	return &pb.SqrtResponse{
			Sqrt: math.Sqrt(float64(req.Number)),
		},
		nil
}
