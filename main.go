package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/muhammadakr/go-test-user/proto" // Update this import path to match your project
	"google.golang.org/grpc"
)

// Implement the UserAuthServiceServer interface
type server struct {
	proto.UnimplementedUserAuthServiceServer
}

var jwtKey = []byte("your-secret-key")

func (s *server) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	if req.Username == "test" && req.Password == "password" {
		// Create JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": req.Username,
			"exp":      time.Now().Add(time.Hour * 1).Unix(),
		})

		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			return nil, err
		}

		return &proto.LoginResponse{
			Success: true,
			Message: "Login successful",
			Token:   tokenString, // Add this field in your proto
		}, nil
	}
	return &proto.LoginResponse{
		Success: false,
		Message: "Invalid username or password",
	}, nil
}

func main() {
	//start gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//create a new gRPC server and register the UserAuthService
	grpcServer := grpc.NewServer()
	proto.RegisterUserAuthServiceServer(grpcServer, &server{})

	fmt.Println("User Authentication Service is running on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
