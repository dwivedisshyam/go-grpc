package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/go-grpc/calculator/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var addr string = "localhost:50052"

func main() {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("cannot connect to grpc", err)
	}

	defer conn.Close()

	c := pb.NewCalculatorServiceClient(conn)

	// getSum(c)
	// primeNumberDecomposition(c)
	// average(c)
	// doCurrentMax(c)
	doSqrt(c, 10)
}

func getSum(c pb.CalculatorServiceClient) {
	req := &pb.SumRequest{
		FirstNumber:  10,
		SecondNumber: 25,
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v", err)
	}

	log.Println("Result:", res.Result)
}

func primeNumberDecomposition(c pb.CalculatorServiceClient) {
	req := &pb.PrimeNumberDecompositionRequest{
		Number: 120,
	}

	stream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling PrimeNumberDecomposition RPC: %v", err)
	}
	var factors []int64
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while receiving from stream: %v", err)
		}
		factors = append(factors, msg.PrimeFactor)
	}
	log.Println("Prime factors:", factors)
}

func average(c pb.CalculatorServiceClient) {
	requests := []*pb.AvgRequest{
		{Number: 1},
		{Number: 2},
		{Number: 3},
		{Number: 4},
	}
	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatal("error which calling Average", err)
	}

	for _, req := range requests {
		err = stream.Send(req)
		if err != nil {
			log.Fatal("error while sending avg request", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("failed to receive avg response")
	}

	log.Println("Result:", res.Avg)
}

func doCurrentMax(c pb.CalculatorServiceClient) {
	stream, err := c.CurrentMax(context.Background())
	if err != nil {
		log.Fatal("Failed to create the stream", err)
	}

	numbers := []int32{1, 5, 3, 6, 2, 20}

	waitc := make(chan struct{})

	go func() {
		for _, num := range numbers {
			log.Println("Sending:", num)
			stream.Send(&pb.CurrentMaxRequest{Number: num})
			time.Sleep(1 * time.Second)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println("Failed to get the result", err)
				break
			}

			log.Println("CurrentMax:", res.Number)
		}
		close(waitc)
	}()

	<-waitc
}

func doSqrt(c pb.CalculatorServiceClient, number int32) {
	log.Println("doSqrt invoked")

	res, err := c.Sqrt(context.Background(), &pb.SqrtRequest{Number: number})
	if err != nil {
		e, ok := status.FromError(err)
		if ok {
			if e.Code() == codes.InvalidArgument {
				log.Println("invalid argument passed")
				return
			}
		} else {
			log.Println("Unexpected error")
			return
		}
	}

	log.Println("Sqrt:", res.Sqrt)
}
