package main

import (
	"context"
	"fmt"

	pb "github.com/go-grpc/blog/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreageBlog(ctx context.Context, req *pb.Blog) (*pb.BlogId, error) {
	data := BlogItem{
		AuthorId: req.AuthorId,
		Title:    req.Title,
		Content:  req.Content,
	}

	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, status.Error(
			codes.Internal, fmt.Sprintf("Internal Err %v", err),
		)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Error(
			codes.Internal,
			fmt.Sprintf("Internal Err %v", err),
		)
	}

	return &pb.BlogId{Id: oid.Hex()}, nil
}

func (s *Server) ReadBlog(ctx context.Context, req *pb.BlogId) (*pb.Blog, error) {
	oid, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"Cannot Parse id",
		)
	}

	data := &BlogItem{}
	filter := bson.M{"_id": oid}

	res := collection.FindOne(ctx, filter)

	err = res.Decode(data)
	if err != nil {
		return nil, status.Error(
			codes.Internal, fmt.Sprintf("Internal Err %v", err),
		)
	}

	return documentToBlog(data), nil
}
