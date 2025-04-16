package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/muhammadakr/go-test-user/proto" // Update this import path to match your project
	"google.golang.org/grpc"
)

// Implement the UserAuthServiceServer interface
type server struct {
	proto.UnimplementedUserAuthServiceServer
}

func (s *server) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	// For testing, we'll mock a user check
	if req.Username == "test" && req.Password == "password" {
		return &proto.LoginResponse{
			Success: true,
			Message: "Login successful",
		}, nil
	}
	return &proto.LoginResponse{
		Success: false,
		Message: "Invalid username or password",
	}, nil
}

func main() {
	// Start gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a new gRPC server and register the UserAuthService
	grpcServer := grpc.NewServer()
	proto.RegisterUserAuthServiceServer(grpcServer, &server{})

	fmt.Println("User Authentication Service is running on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
