package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"google.golang.org/grpc"
	"github.com/muhammadakr/go-test-user/proto" 
)

type server struct {
	proto.UnimplementedUserServiceServer
}

// Login implementation for UserService
func (s *server) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	// Dummy user authentication logic (replace with your own)
	if req.GetUsername() == "testuser" && req.GetPassword() == "testpassword" {
		return &proto.LoginResponse{
			Message: "Authentication successful",
			Success: true,
		}, nil
	}
	return &proto.LoginResponse{
		Message: "Invalid credentials",
		Success: false,
	}, nil
}

func main() {
	// Start gRPC server
	listen, err := net.Listen("tcp", ":3001")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	// Register the UserService
	proto.RegisterUserServiceServer(grpcServer, &server{})

	// Start the server
	fmt.Println("User Service is running on port 3001")
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

