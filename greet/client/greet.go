package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/go-grpc/greet/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func doGreet(c pb.GreetServiceClient) {
	log.Println("doGreet was invoked")

	req := &pb.GreetRequest{
		FirstName: "John",
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC: %v", err)
	}

	log.Printf("Response from Greet: %v", res.Result)
}

func doGreetWithDeadline(c pb.GreetServiceClient, timeout time.Duration) {
	log.Println("doGreetWithDeadline was invoked")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req := &pb.GreetRequest{
		FirstName: "John",
	}

	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {
		e, ok := status.FromError(err)
		if ok {
			if e.Code() == codes.DeadlineExceeded {
				log.Println("deadline exceeded")
				return
			} else {
				log.Println("non gRPC error", err)
			}
		} else {
			log.Println("non gRPC error", err)
		}

		log.Fatalf("Error while calling Greet RPC: %v", err)
	}

	log.Printf("Response from GreetWithDeadline: %v", res.Result)
}

func dogreetManyTimes(c pb.GreetServiceClient) {
	log.Println("greetManyTimes was invoked")

	req := &pb.GreetRequest{
		FirstName: "John",
	}

	stream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes RPC: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while receiving from GreetManyTimes stream: %v", err)
		}
		log.Printf("Response from GreetManyTimes: %v", res.Result)
	}
}

func doLongGreet(c pb.GreetServiceClient) {
	log.Println("doLongGreet invoked")

	reqs := []*pb.GreetRequest{
		{FirstName: "John"},
		{FirstName: "Undertaker"},
		{FirstName: "Ponting"},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error calling long stream %v\n", err)
	}

	for _, req := range reqs {
		log.Printf("Sending req %v\n", req)
		stream.Send(req)
		time.Sleep(1 * time.Second)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Failed to connect %v\n", err)
	}

	log.Println("LongGreed:", res.Result)
}

func doGreetEveryone(c pb.GreetServiceClient) {
	log.Println("doGreetEveryone invoked")

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatal("error while creating stream", err)
	}

	reqs := []*pb.GreetRequest{
		{FirstName: "Clement"},
		{FirstName: "Shyam"},
		{FirstName: "John"},
	}

	waitc := make(chan struct{})

	go func() {
		for _, req := range reqs {
			log.Println("Sending: ", req.FirstName)
			stream.Send(req)
			time.Sleep(time.Second * 1)
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
				log.Println("error while receving", err)
				break
			}

			log.Println("Received:", res.Result)
		}

		close(waitc)
	}()

	<-waitc
}
