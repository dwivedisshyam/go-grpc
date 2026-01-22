package main

import (
	"context"
	"log"

	pb "github.com/go-grpc/blog/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addr string = "localhost:50052"

type Client struct {
	blogClient pb.BlogServiceClient
}

func main() {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("cannot connect to grpc", err)
	}

	defer conn.Close()

	c := pb.NewBlogServiceClient(conn)

	client := &Client{
		blogClient: c,
	}

	client.readBlog()

}

func (c Client) createBlog() {
	req := &pb.Blog{
		AuthorId: "Shyam",
		Title:    "My Fist blog",
		Content:  "Content of the first blog",
	}

	res, err := c.blogClient.CreageBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling create RPC: %v", err)
	}

	log.Println("New Blog Id:", res.Id)
}

func (c Client) readBlog() {
	req := &pb.BlogId{
		Id: "000000000000000000000000",
	}

	res, err := c.blogClient.ReadBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling create RPC: %v", err)
	}

	log.Printf("Blog: %v\n", res)
}
