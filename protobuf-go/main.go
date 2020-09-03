package main

import (
	pb "examples/protobuf-go/proto/srv"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

type UserService struct {
	pb.UnimplementedUserServer
}

func main() {
	server := grpc.NewServer()
	pb.RegisterUserServer(server, new(UserService))

	lis, err := net.Listen("tcp", "0.0.0.0:8877")
	if err != nil {
		log.Fatal(fmt.Errorf("failed to listen: %w", err))
	}
	defer lis.Close()

	fmt.Println("server is listening ...")
	if err := server.Serve(lis); err != nil {
		log.Fatal(fmt.Errorf("failed to serve: %w", err))
	}
}
