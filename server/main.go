package main

import (
	"context"
	"log"
	"net"
	"os"

	pb "server/kvstore"
	in_memory_store "server/stores"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedKeyValueStoreServer
	store *in_memory_store.Store
}

func (s *server) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	s.store.Set(req.GetKey(), req.GetValue())
	return &pb.SetResponse{Message: "Value set successfully"}, nil
}

func (s *server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	value, ok := s.store.Get(req.GetKey())
	if !ok {
		return &pb.GetResponse{Message: "Key not found"}, nil
	}
	return &pb.GetResponse{Value: value, Message: "Value retrieved successfully"}, nil
}

func (s *server) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	s.store.Delete(req.GetKey())
	return &pb.DeleteResponse{Message: "Key deleted successfully"}, nil
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}

	address := ":" + port
	log.Printf("Starting server on port %s", port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterKeyValueStoreServer(grpcServer, &server{store: in_memory_store.NewStore()})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
