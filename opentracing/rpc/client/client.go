package main

import (
	"context"
	"examples/opentracing/rpc/interceptors"
	pb "examples/opentracing/rpc/server/proto"
	"examples/opentracing/utils"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"log"
)

// configurations
const (
	tracerServiceName = "myRPCClient"
	tracerAgentHost   = "localhost"
	tracerAgentPort   = 32772

	serverRPCHost = "localhost"
	serverRPCPort = 8976
)

func main() {
	// register tracer
	closer := utils.InitTracer(tracerServiceName, tracerAgentHost, tracerAgentPort)
	defer closer.Close()

	// start span
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("call sum and product")
	defer span.Finish()

	// make context with span
	ctx := opentracing.ContextWithSpan(context.Background(), span)

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", serverRPCHost, serverRPCPort), grpc.WithInsecure(),
		grpc.WithChainUnaryInterceptor(interceptors.UnaryClientTraceInterceptor))
	if err != nil {
		err = fmt.Errorf("main: grpc dial: %w", err)
		log.Fatal(err)
	}
	defer conn.Close()

	mathClient := pb.NewMathClient(conn)
	if v, err := mathClient.Sum(ctx, &pb.Numbers{
		Numbers: []int32{12, 13},
	}); err != nil {
		err = fmt.Errorf("main: invoke sum: %w", err)
		log.Fatal(err)
	} else {
		fmt.Println("sum:", v.GetValue())
	}

	if v, err := mathClient.Product(ctx, &pb.Numbers{
		Numbers: []int32{6, 3},
	}); err != nil {
		err = fmt.Errorf("main: invoke product: %w", err)
		log.Fatal(err)
	} else {
		fmt.Println("product:", v.GetValue())
	}
}
