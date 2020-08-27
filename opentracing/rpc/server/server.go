package main

import (
	"context"
	"examples/opentracing/rpc/interceptors"
	pb "examples/opentracing/rpc/server/proto"
	"examples/opentracing/utils"
	"fmt"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"log"
	"net"
)

type mathServer struct {
	pb.UnimplementedMathServer
}

func (m *mathServer) Sum(ctx context.Context, numbers *pb.Numbers) (*wrappers.Int32Value, error) {
	res := int32(0)
	for _, num := range numbers.GetNumbers() {
		res += num
	}
	return wrapperspb.Int32(res), nil
}

func (m *mathServer) Product(ctx context.Context, numbers *pb.Numbers) (*wrappers.Int32Value, error) {
	res := int32(1)
	for _, num := range numbers.GetNumbers() {
		res *= num
	}
	return wrapperspb.Int32(res), nil
}

// configurations
const (
	tracerServiceName = "myRPCService"
	tracerAgentHost   = "localhost"
	tracerAgentPort   = 32772

	serverPort = 8976
)

func main() {
	closer := utils.InitTracer(tracerServiceName, tracerAgentHost, tracerAgentPort)
	defer closer.Close()

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptors.UnaryServerTraceInterceptor))
	pb.RegisterMathServer(server, new(mathServer))

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", serverPort))
	if err != nil {
		err = fmt.Errorf("main: can't listen: %w", err)
		log.Fatal(err)
	}
	if err = server.Serve(listener); err != nil {
		err = fmt.Errorf("main: grpc server can't serve: %w", err)
		log.Fatal(err)
	}
}
