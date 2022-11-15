// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.9
// source: messengerpb/messenger.proto

package messengerpb

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

// MessengerServiceClient is the client API for MessengerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MessengerServiceClient interface {
	SignIn(ctx context.Context, in *SignInData, opts ...grpc.CallOption) (*UserData, error)
	SignUp(ctx context.Context, in *SignUpData, opts ...grpc.CallOption) (*UserData, error)
	CheckName(ctx context.Context, in *CheckNameMessage, opts ...grpc.CallOption) (*CheckNameAck, error)
	CheckLogin(ctx context.Context, in *CheckLoginMessage, opts ...grpc.CallOption) (*CheckLoginAck, error)
	RequestAddToFriendsList(ctx context.Context, in *RequestAddToFriendsListMessage, opts ...grpc.CallOption) (*RequestAddToFriendsListAck, error)
	ListenAddToFriendsList(ctx context.Context, in *UserData, opts ...grpc.CallOption) (MessengerService_ListenAddToFriendsListClient, error)
	SendMessage(ctx context.Context, opts ...grpc.CallOption) (MessengerService_SendMessageClient, error)
	ReceiveMessage(ctx context.Context, in *UserData, opts ...grpc.CallOption) (MessengerService_ReceiveMessageClient, error)
}

type messengerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMessengerServiceClient(cc grpc.ClientConnInterface) MessengerServiceClient {
	return &messengerServiceClient{cc}
}

func (c *messengerServiceClient) SignIn(ctx context.Context, in *SignInData, opts ...grpc.CallOption) (*UserData, error) {
	out := new(UserData)
	err := c.cc.Invoke(ctx, "/messengerpb.MessengerService/SignIn", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messengerServiceClient) SignUp(ctx context.Context, in *SignUpData, opts ...grpc.CallOption) (*UserData, error) {
	out := new(UserData)
	err := c.cc.Invoke(ctx, "/messengerpb.MessengerService/SignUp", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messengerServiceClient) CheckName(ctx context.Context, in *CheckNameMessage, opts ...grpc.CallOption) (*CheckNameAck, error) {
	out := new(CheckNameAck)
	err := c.cc.Invoke(ctx, "/messengerpb.MessengerService/CheckName", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messengerServiceClient) CheckLogin(ctx context.Context, in *CheckLoginMessage, opts ...grpc.CallOption) (*CheckLoginAck, error) {
	out := new(CheckLoginAck)
	err := c.cc.Invoke(ctx, "/messengerpb.MessengerService/CheckLogin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messengerServiceClient) RequestAddToFriendsList(ctx context.Context, in *RequestAddToFriendsListMessage, opts ...grpc.CallOption) (*RequestAddToFriendsListAck, error) {
	out := new(RequestAddToFriendsListAck)
	err := c.cc.Invoke(ctx, "/messengerpb.MessengerService/RequestAddToFriendsList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messengerServiceClient) ListenAddToFriendsList(ctx context.Context, in *UserData, opts ...grpc.CallOption) (MessengerService_ListenAddToFriendsListClient, error) {
	stream, err := c.cc.NewStream(ctx, &MessengerService_ServiceDesc.Streams[0], "/messengerpb.MessengerService/ListenAddToFriendsList", opts...)
	if err != nil {
		return nil, err
	}
	x := &messengerServiceListenAddToFriendsListClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MessengerService_ListenAddToFriendsListClient interface {
	Recv() (*UserData, error)
	grpc.ClientStream
}

type messengerServiceListenAddToFriendsListClient struct {
	grpc.ClientStream
}

func (x *messengerServiceListenAddToFriendsListClient) Recv() (*UserData, error) {
	m := new(UserData)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *messengerServiceClient) SendMessage(ctx context.Context, opts ...grpc.CallOption) (MessengerService_SendMessageClient, error) {
	stream, err := c.cc.NewStream(ctx, &MessengerService_ServiceDesc.Streams[1], "/messengerpb.MessengerService/SendMessage", opts...)
	if err != nil {
		return nil, err
	}
	x := &messengerServiceSendMessageClient{stream}
	return x, nil
}

type MessengerService_SendMessageClient interface {
	Send(*Message) error
	CloseAndRecv() (*MessageAck, error)
	grpc.ClientStream
}

type messengerServiceSendMessageClient struct {
	grpc.ClientStream
}

func (x *messengerServiceSendMessageClient) Send(m *Message) error {
	return x.ClientStream.SendMsg(m)
}

func (x *messengerServiceSendMessageClient) CloseAndRecv() (*MessageAck, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(MessageAck)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *messengerServiceClient) ReceiveMessage(ctx context.Context, in *UserData, opts ...grpc.CallOption) (MessengerService_ReceiveMessageClient, error) {
	stream, err := c.cc.NewStream(ctx, &MessengerService_ServiceDesc.Streams[2], "/messengerpb.MessengerService/ReceiveMessage", opts...)
	if err != nil {
		return nil, err
	}
	x := &messengerServiceReceiveMessageClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MessengerService_ReceiveMessageClient interface {
	Recv() (*Message, error)
	grpc.ClientStream
}

type messengerServiceReceiveMessageClient struct {
	grpc.ClientStream
}

func (x *messengerServiceReceiveMessageClient) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MessengerServiceServer is the server API for MessengerService service.
// All implementations must embed UnimplementedMessengerServiceServer
// for forward compatibility
type MessengerServiceServer interface {
	SignIn(context.Context, *SignInData) (*UserData, error)
	SignUp(context.Context, *SignUpData) (*UserData, error)
	CheckName(context.Context, *CheckNameMessage) (*CheckNameAck, error)
	CheckLogin(context.Context, *CheckLoginMessage) (*CheckLoginAck, error)
	RequestAddToFriendsList(context.Context, *RequestAddToFriendsListMessage) (*RequestAddToFriendsListAck, error)
	ListenAddToFriendsList(*UserData, MessengerService_ListenAddToFriendsListServer) error
	SendMessage(MessengerService_SendMessageServer) error
	ReceiveMessage(*UserData, MessengerService_ReceiveMessageServer) error
	mustEmbedUnimplementedMessengerServiceServer()
}

// UnimplementedMessengerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMessengerServiceServer struct {
}

func (UnimplementedMessengerServiceServer) SignIn(context.Context, *SignInData) (*UserData, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignIn not implemented")
}
func (UnimplementedMessengerServiceServer) SignUp(context.Context, *SignUpData) (*UserData, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignUp not implemented")
}
func (UnimplementedMessengerServiceServer) CheckName(context.Context, *CheckNameMessage) (*CheckNameAck, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckName not implemented")
}
func (UnimplementedMessengerServiceServer) CheckLogin(context.Context, *CheckLoginMessage) (*CheckLoginAck, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckLogin not implemented")
}
func (UnimplementedMessengerServiceServer) RequestAddToFriendsList(context.Context, *RequestAddToFriendsListMessage) (*RequestAddToFriendsListAck, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestAddToFriendsList not implemented")
}
func (UnimplementedMessengerServiceServer) ListenAddToFriendsList(*UserData, MessengerService_ListenAddToFriendsListServer) error {
	return status.Errorf(codes.Unimplemented, "method ListenAddToFriendsList not implemented")
}
func (UnimplementedMessengerServiceServer) SendMessage(MessengerService_SendMessageServer) error {
	return status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}
func (UnimplementedMessengerServiceServer) ReceiveMessage(*UserData, MessengerService_ReceiveMessageServer) error {
	return status.Errorf(codes.Unimplemented, "method ReceiveMessage not implemented")
}
func (UnimplementedMessengerServiceServer) mustEmbedUnimplementedMessengerServiceServer() {}

// UnsafeMessengerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MessengerServiceServer will
// result in compilation errors.
type UnsafeMessengerServiceServer interface {
	mustEmbedUnimplementedMessengerServiceServer()
}

func RegisterMessengerServiceServer(s grpc.ServiceRegistrar, srv MessengerServiceServer) {
	s.RegisterService(&MessengerService_ServiceDesc, srv)
}

func _MessengerService_SignIn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignInData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessengerServiceServer).SignIn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messengerpb.MessengerService/SignIn",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessengerServiceServer).SignIn(ctx, req.(*SignInData))
	}
	return interceptor(ctx, in, info, handler)
}

func _MessengerService_SignUp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignUpData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessengerServiceServer).SignUp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messengerpb.MessengerService/SignUp",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessengerServiceServer).SignUp(ctx, req.(*SignUpData))
	}
	return interceptor(ctx, in, info, handler)
}

func _MessengerService_CheckName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckNameMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessengerServiceServer).CheckName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messengerpb.MessengerService/CheckName",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessengerServiceServer).CheckName(ctx, req.(*CheckNameMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _MessengerService_CheckLogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckLoginMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessengerServiceServer).CheckLogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messengerpb.MessengerService/CheckLogin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessengerServiceServer).CheckLogin(ctx, req.(*CheckLoginMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _MessengerService_RequestAddToFriendsList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestAddToFriendsListMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessengerServiceServer).RequestAddToFriendsList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messengerpb.MessengerService/RequestAddToFriendsList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessengerServiceServer).RequestAddToFriendsList(ctx, req.(*RequestAddToFriendsListMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _MessengerService_ListenAddToFriendsList_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(UserData)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MessengerServiceServer).ListenAddToFriendsList(m, &messengerServiceListenAddToFriendsListServer{stream})
}

type MessengerService_ListenAddToFriendsListServer interface {
	Send(*UserData) error
	grpc.ServerStream
}

type messengerServiceListenAddToFriendsListServer struct {
	grpc.ServerStream
}

func (x *messengerServiceListenAddToFriendsListServer) Send(m *UserData) error {
	return x.ServerStream.SendMsg(m)
}

func _MessengerService_SendMessage_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MessengerServiceServer).SendMessage(&messengerServiceSendMessageServer{stream})
}

type MessengerService_SendMessageServer interface {
	SendAndClose(*MessageAck) error
	Recv() (*Message, error)
	grpc.ServerStream
}

type messengerServiceSendMessageServer struct {
	grpc.ServerStream
}

func (x *messengerServiceSendMessageServer) SendAndClose(m *MessageAck) error {
	return x.ServerStream.SendMsg(m)
}

func (x *messengerServiceSendMessageServer) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _MessengerService_ReceiveMessage_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(UserData)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MessengerServiceServer).ReceiveMessage(m, &messengerServiceReceiveMessageServer{stream})
}

type MessengerService_ReceiveMessageServer interface {
	Send(*Message) error
	grpc.ServerStream
}

type messengerServiceReceiveMessageServer struct {
	grpc.ServerStream
}

func (x *messengerServiceReceiveMessageServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

// MessengerService_ServiceDesc is the grpc.ServiceDesc for MessengerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MessengerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "messengerpb.MessengerService",
	HandlerType: (*MessengerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SignIn",
			Handler:    _MessengerService_SignIn_Handler,
		},
		{
			MethodName: "SignUp",
			Handler:    _MessengerService_SignUp_Handler,
		},
		{
			MethodName: "CheckName",
			Handler:    _MessengerService_CheckName_Handler,
		},
		{
			MethodName: "CheckLogin",
			Handler:    _MessengerService_CheckLogin_Handler,
		},
		{
			MethodName: "RequestAddToFriendsList",
			Handler:    _MessengerService_RequestAddToFriendsList_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListenAddToFriendsList",
			Handler:       _MessengerService_ListenAddToFriendsList_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SendMessage",
			Handler:       _MessengerService_SendMessage_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "ReceiveMessage",
			Handler:       _MessengerService_ReceiveMessage_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "messengerpb/messenger.proto",
}
