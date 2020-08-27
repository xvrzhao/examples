package interceptors

import (
	"context"
	"examples/opentracing/utils"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
)

// UnaryServerTraceInterceptor is a unary server-side interceptor implemented tracing.
func UnaryServerTraceInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
	resp interface{}, err error) {

	if opentracing.IsGlobalTracerRegistered() {
		tracer := opentracing.GlobalTracer()
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		clientSpanContext, err := tracer.Extract(opentracing.TextMap, utils.MetadataCarrier(md))
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			err = fmt.Errorf("unaryTraceInterceptor: extract spanContext from metadata: %w", err)
			grpclog.Error(err)
		} else {
			span := opentracing.StartSpan(
				info.FullMethod,
				opentracing.ChildOf(clientSpanContext),
				opentracing.Tag{Key: string(ext.Component), Value: "gRPC"})
			defer span.Finish()

			// put this span into ctx
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}

	return handler(ctx, req)
}

// TODO: func StreamServerTraceInterceptor
