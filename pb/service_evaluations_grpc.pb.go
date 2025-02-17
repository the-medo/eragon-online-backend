// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.22.0
// source: services/evaluations/service_evaluations.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Evaluations_GetEvaluationById_FullMethodName    = "/pb.Evaluations/GetEvaluationById"
	Evaluations_GetEvaluationsByType_FullMethodName = "/pb.Evaluations/GetEvaluationsByType"
)

// EvaluationsClient is the client API for Evaluations service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EvaluationsClient interface {
	// ============= USER EVALUATION =================
	GetEvaluationById(ctx context.Context, in *GetEvaluationByIdRequest, opts ...grpc.CallOption) (*GetEvaluationByIdResponse, error)
	GetEvaluationsByType(ctx context.Context, in *GetEvaluationsByTypeRequest, opts ...grpc.CallOption) (*GetEvaluationsByTypeResponse, error)
}

type evaluationsClient struct {
	cc grpc.ClientConnInterface
}

func NewEvaluationsClient(cc grpc.ClientConnInterface) EvaluationsClient {
	return &evaluationsClient{cc}
}

func (c *evaluationsClient) GetEvaluationById(ctx context.Context, in *GetEvaluationByIdRequest, opts ...grpc.CallOption) (*GetEvaluationByIdResponse, error) {
	out := new(GetEvaluationByIdResponse)
	err := c.cc.Invoke(ctx, Evaluations_GetEvaluationById_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *evaluationsClient) GetEvaluationsByType(ctx context.Context, in *GetEvaluationsByTypeRequest, opts ...grpc.CallOption) (*GetEvaluationsByTypeResponse, error) {
	out := new(GetEvaluationsByTypeResponse)
	err := c.cc.Invoke(ctx, Evaluations_GetEvaluationsByType_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EvaluationsServer is the server API for Evaluations service.
// All implementations must embed UnimplementedEvaluationsServer
// for forward compatibility
type EvaluationsServer interface {
	// ============= USER EVALUATION =================
	GetEvaluationById(context.Context, *GetEvaluationByIdRequest) (*GetEvaluationByIdResponse, error)
	GetEvaluationsByType(context.Context, *GetEvaluationsByTypeRequest) (*GetEvaluationsByTypeResponse, error)
	mustEmbedUnimplementedEvaluationsServer()
}

// UnimplementedEvaluationsServer must be embedded to have forward compatible implementations.
type UnimplementedEvaluationsServer struct {
}

func (UnimplementedEvaluationsServer) GetEvaluationById(context.Context, *GetEvaluationByIdRequest) (*GetEvaluationByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEvaluationById not implemented")
}
func (UnimplementedEvaluationsServer) GetEvaluationsByType(context.Context, *GetEvaluationsByTypeRequest) (*GetEvaluationsByTypeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEvaluationsByType not implemented")
}
func (UnimplementedEvaluationsServer) mustEmbedUnimplementedEvaluationsServer() {}

// UnsafeEvaluationsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EvaluationsServer will
// result in compilation errors.
type UnsafeEvaluationsServer interface {
	mustEmbedUnimplementedEvaluationsServer()
}

func RegisterEvaluationsServer(s grpc.ServiceRegistrar, srv EvaluationsServer) {
	s.RegisterService(&Evaluations_ServiceDesc, srv)
}

func _Evaluations_GetEvaluationById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEvaluationByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EvaluationsServer).GetEvaluationById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Evaluations_GetEvaluationById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EvaluationsServer).GetEvaluationById(ctx, req.(*GetEvaluationByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Evaluations_GetEvaluationsByType_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEvaluationsByTypeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EvaluationsServer).GetEvaluationsByType(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Evaluations_GetEvaluationsByType_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EvaluationsServer).GetEvaluationsByType(ctx, req.(*GetEvaluationsByTypeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Evaluations_ServiceDesc is the grpc.ServiceDesc for Evaluations service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Evaluations_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Evaluations",
	HandlerType: (*EvaluationsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetEvaluationById",
			Handler:    _Evaluations_GetEvaluationById_Handler,
		},
		{
			MethodName: "GetEvaluationsByType",
			Handler:    _Evaluations_GetEvaluationsByType_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services/evaluations/service_evaluations.proto",
}
