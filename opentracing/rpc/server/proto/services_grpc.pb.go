// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

import (
	context "context"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// MathClient is the client API for Math service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MathClient interface {
	Sum(ctx context.Context, in *Numbers, opts ...grpc.CallOption) (*wrappers.Int32Value, error)
	Product(ctx context.Context, in *Numbers, opts ...grpc.CallOption) (*wrappers.Int32Value, error)
}

type mathClient struct {
	cc grpc.ClientConnInterface
}

func NewMathClient(cc grpc.ClientConnInterface) MathClient {
	return &mathClient{cc}
}

func (c *mathClient) Sum(ctx context.Context, in *Numbers, opts ...grpc.CallOption) (*wrappers.Int32Value, error) {
	out := new(wrappers.Int32Value)
	err := c.cc.Invoke(ctx, "/proto.Math/Sum", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mathClient) Product(ctx context.Context, in *Numbers, opts ...grpc.CallOption) (*wrappers.Int32Value, error) {
	out := new(wrappers.Int32Value)
	err := c.cc.Invoke(ctx, "/proto.Math/Product", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MathServer is the server API for Math service.
// All implementations must embed UnimplementedMathServer
// for forward compatibility
type MathServer interface {
	Sum(context.Context, *Numbers) (*wrappers.Int32Value, error)
	Product(context.Context, *Numbers) (*wrappers.Int32Value, error)
	mustEmbedUnimplementedMathServer()
}

// UnimplementedMathServer must be embedded to have forward compatible implementations.
type UnimplementedMathServer struct {
}

func (*UnimplementedMathServer) Sum(context.Context, *Numbers) (*wrappers.Int32Value, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Sum not implemented")
}
func (*UnimplementedMathServer) Product(context.Context, *Numbers) (*wrappers.Int32Value, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Product not implemented")
}
func (*UnimplementedMathServer) mustEmbedUnimplementedMathServer() {}

func RegisterMathServer(s *grpc.Server, srv MathServer) {
	s.RegisterService(&_Math_serviceDesc, srv)
}

func _Math_Sum_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Numbers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MathServer).Sum(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Math/Sum",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MathServer).Sum(ctx, req.(*Numbers))
	}
	return interceptor(ctx, in, info, handler)
}

func _Math_Product_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Numbers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MathServer).Product(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Math/Product",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MathServer).Product(ctx, req.(*Numbers))
	}
	return interceptor(ctx, in, info, handler)
}

var _Math_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Math",
	HandlerType: (*MathServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Sum",
			Handler:    _Math_Sum_Handler,
		},
		{
			MethodName: "Product",
			Handler:    _Math_Product_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/services.proto",
}
