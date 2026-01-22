package main

import (
	"context"
	"log"
	"net"

	pb "github.com/go-grpc/blog/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var collection *mongo.Collection
var addr string = "localhost:50052"

type Server struct {
	pb.BlogServiceServer
}

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:root@localhost:27017/"))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("blogdb").Collection("blog")

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)

	}
	log.Printf("Listening on %s\n", addr)
	c := grpc.NewServer()

	pb.RegisterBlogServiceServer(c, &Server{})
	reflection.Register(c)

	if err = c.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
