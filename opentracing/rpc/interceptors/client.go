package interceptors

import (
	"context"
	"examples/opentracing/utils"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
)

func UnaryClientTraceInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	if !opentracing.IsGlobalTracerRegistered() {
		return invoker(ctx, method, req, reply, cc, opts...)
	}

	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return invoker(ctx, method, req, reply, cc, opts...)
	}

	tracer := opentracing.GlobalTracer()
	md := metadata.New(nil)
	outCtx := metadata.NewOutgoingContext(ctx, md)

	if err := tracer.Inject(span.Context(), opentracing.TextMap, utils.MetadataCarrier(md)); err != nil {
		err = fmt.Errorf("UnaryClientTraceInterceptor: cannot inject span to outgoing context: %w", err)
		grpclog.Error(err)
	}

	return invoker(outCtx, method, req, reply, cc, opts...)
}
