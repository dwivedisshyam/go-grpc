package main

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/go-grpc/greet/proto"
)

var addr string = "localhost:50051"

func main() {
	tls := true

	opts := []grpc.DialOption{}

	if tls {
		certFile := "ssl/ca.crt"
		creds, err := credentials.NewClientTLSFromFile(certFile, "")

		if err != nil {
			log.Fatalf("Error while loading CA trust certificate: %v\n", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		creds := grpc.WithTransportCredentials(insecure.NewCredentials())
		opts = append(opts, creds)
	}

	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		log.Fatalf("Failed to connect %v\n", err)
	}
	defer conn.Close()

	log.Printf("Connected to %s\n", addr)

	c := pb.NewGreetServiceClient(conn)
	doGreet(c)

	// dogreetManyTimes(c)
	// doLongGreet(c)
	// doGreetEveryone(c)

	// doGreetWithDeadline(c, 2*time.Second)
}
