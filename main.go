package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/muhammadakr/go-test-user/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type server struct {
	proto.UnimplementedUserAuthServiceServer
}

var jwtKey = []byte("your-secret-key")

func (s *server) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	fmt.Printf("User Service received login request: username=%s password=%s\n", req.Username, req.Password)
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
			Token:   tokenString,
		}, nil
	}
	return &proto.LoginResponse{
		Success: false,
		Message: "Invalid username or password",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//Load CA certificate to verify clients
	caCert, err := ioutil.ReadFile("ca.crt")
	if err != nil {
		log.Fatalf("failed to load CA cert: %v", err)
	}
	clientCACertPool := x509.NewCertPool()
	if ok := clientCACertPool.AppendCertsFromPEM(caCert); !ok {
		log.Fatalf("failed to append client CA cert to pool")
	}

	//Load Server certificate and key
	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatalf("failed to load server cert/key: %v", err)
	}

	//Create mTLS config
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    clientCACertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	//Start gRPC server with TLS
	grpcServer := grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig)))
	proto.RegisterUserAuthServiceServer(grpcServer, &server{})

	fmt.Println("User Auth Service is running on mTLS port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
