// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: ctr/v1/ctr.proto

package ctr

import (
	context "context"
	any "github.com/golang/protobuf/ptypes/any"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_ctr_v1_ctr_proto protoreflect.FileDescriptor

var file_ctr_v1_ctr_proto_rawDesc = []byte{
	0x0a, 0x10, 0x63, 0x74, 0x72, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x74, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x06, 0x63, 0x74, 0x72, 0x2e, 0x76, 0x31, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x32, 0x40, 0x0a, 0x03, 0x43, 0x74, 0x72, 0x12, 0x39, 0x0a, 0x07, 0x50, 0x75, 0x62,
	0x6c, 0x69, 0x73, 0x68, 0x12, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x22, 0x00, 0x42, 0x0c, 0x5a, 0x0a, 0x63, 0x74, 0x72, 0x2f, 0x76, 0x31, 0x3b, 0x63,
	0x74, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_ctr_v1_ctr_proto_goTypes = []interface{}{
	(*any.Any)(nil),     // 0: google.protobuf.Any
	(*empty.Empty)(nil), // 1: google.protobuf.Empty
}
var file_ctr_v1_ctr_proto_depIdxs = []int32{
	0, // 0: ctr.v1.Ctr.Publish:input_type -> google.protobuf.Any
	1, // 1: ctr.v1.Ctr.Publish:output_type -> google.protobuf.Empty
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_ctr_v1_ctr_proto_init() }
func file_ctr_v1_ctr_proto_init() {
	if File_ctr_v1_ctr_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_ctr_v1_ctr_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_ctr_v1_ctr_proto_goTypes,
		DependencyIndexes: file_ctr_v1_ctr_proto_depIdxs,
	}.Build()
	File_ctr_v1_ctr_proto = out.File
	file_ctr_v1_ctr_proto_rawDesc = nil
	file_ctr_v1_ctr_proto_goTypes = nil
	file_ctr_v1_ctr_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// CtrClient is the client API for Ctr service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CtrClient interface {
	Publish(ctx context.Context, in *any.Any, opts ...grpc.CallOption) (*empty.Empty, error)
}

type ctrClient struct {
	cc grpc.ClientConnInterface
}

func NewCtrClient(cc grpc.ClientConnInterface) CtrClient {
	return &ctrClient{cc}
}

func (c *ctrClient) Publish(ctx context.Context, in *any.Any, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/ctr.v1.Ctr/Publish", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CtrServer is the server API for Ctr service.
type CtrServer interface {
	Publish(context.Context, *any.Any) (*empty.Empty, error)
}

// UnimplementedCtrServer can be embedded to have forward compatible implementations.
type UnimplementedCtrServer struct {
}

func (*UnimplementedCtrServer) Publish(context.Context, *any.Any) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Publish not implemented")
}

func RegisterCtrServer(s *grpc.Server, srv CtrServer) {
	s.RegisterService(&_Ctr_serviceDesc, srv)
}

func _Ctr_Publish_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(any.Any)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CtrServer).Publish(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ctr.v1.Ctr/Publish",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CtrServer).Publish(ctx, req.(*any.Any))
	}
	return interceptor(ctx, in, info, handler)
}

var _Ctr_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ctr.v1.Ctr",
	HandlerType: (*CtrServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Publish",
			Handler:    _Ctr_Publish_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ctr/v1/ctr.proto",
}