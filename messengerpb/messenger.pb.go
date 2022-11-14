// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.9
// source: messengerpb/messenger.proto

package messengerpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type SignUpData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FirstName  string      `protobuf:"bytes,1,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	SecondName string      `protobuf:"bytes,2,opt,name=second_name,json=secondName,proto3" json:"second_name,omitempty"`
	SignInData *SignInData `protobuf:"bytes,3,opt,name=signInData,proto3" json:"signInData,omitempty"`
}

func (x *SignUpData) Reset() {
	*x = SignUpData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messengerpb_messenger_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignUpData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignUpData) ProtoMessage() {}

func (x *SignUpData) ProtoReflect() protoreflect.Message {
	mi := &file_messengerpb_messenger_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignUpData.ProtoReflect.Descriptor instead.
func (*SignUpData) Descriptor() ([]byte, []int) {
	return file_messengerpb_messenger_proto_rawDescGZIP(), []int{0}
}

func (x *SignUpData) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}

func (x *SignUpData) GetSecondName() string {
	if x != nil {
		return x.SecondName
	}
	return ""
}

func (x *SignUpData) GetSignInData() *SignInData {
	if x != nil {
		return x.SignInData
	}
	return nil
}

type SignInData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Login    string `protobuf:"bytes,1,opt,name=login,proto3" json:"login,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *SignInData) Reset() {
	*x = SignInData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messengerpb_messenger_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignInData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignInData) ProtoMessage() {}

func (x *SignInData) ProtoReflect() protoreflect.Message {
	mi := &file_messengerpb_messenger_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignInData.ProtoReflect.Descriptor instead.
func (*SignInData) Descriptor() ([]byte, []int) {
	return file_messengerpb_messenger_proto_rawDescGZIP(), []int{1}
}

func (x *SignInData) GetLogin() string {
	if x != nil {
		return x.Login
	}
	return ""
}

func (x *SignInData) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type UserID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *UserID) Reset() {
	*x = UserID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messengerpb_messenger_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserID) ProtoMessage() {}

func (x *UserID) ProtoReflect() protoreflect.Message {
	mi := &file_messengerpb_messenger_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserID.ProtoReflect.Descriptor instead.
func (*UserID) Descriptor() ([]byte, []int) {
	return file_messengerpb_messenger_proto_rawDescGZIP(), []int{2}
}

func (x *UserID) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sender   *UserID `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender,omitempty"`
	Receiver *UserID `protobuf:"bytes,2,opt,name=receiver,proto3" json:"receiver,omitempty"`
	Message  string  `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messengerpb_messenger_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_messengerpb_messenger_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_messengerpb_messenger_proto_rawDescGZIP(), []int{3}
}

func (x *Message) GetSender() *UserID {
	if x != nil {
		return x.Sender
	}
	return nil
}

func (x *Message) GetReceiver() *UserID {
	if x != nil {
		return x.Receiver
	}
	return nil
}

func (x *Message) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type MessageAck struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *MessageAck) Reset() {
	*x = MessageAck{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messengerpb_messenger_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageAck) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageAck) ProtoMessage() {}

func (x *MessageAck) ProtoReflect() protoreflect.Message {
	mi := &file_messengerpb_messenger_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageAck.ProtoReflect.Descriptor instead.
func (*MessageAck) Descriptor() ([]byte, []int) {
	return file_messengerpb_messenger_proto_rawDescGZIP(), []int{4}
}

func (x *MessageAck) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

var File_messengerpb_messenger_proto protoreflect.FileDescriptor

var file_messengerpb_messenger_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x6d, 0x65, 0x73, 0x73, 0x65, 0x6e, 0x67, 0x65, 0x72, 0x70, 0x62, 0x2f, 0x6d, 0x65,
	0x73, 0x73, 0x65, 0x6e, 0x67, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x6d,
	0x65, 0x73, 0x73, 0x65, 0x6e, 0x67, 0x65, 0x72, 0x70, 0x62, 0x22, 0x85, 0x01, 0x0a, 0x0a, 0x53,
	0x69, 0x67, 0x6e, 0x55, 0x70, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1d, 0x0a, 0x0a, 0x66, 0x69, 0x72,
	0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x66,
	0x69, 0x72, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x65, 0x63, 0x6f,
	0x6e, 0x64, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73,
	0x65, 0x63, 0x6f, 0x6e, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x37, 0x0a, 0x0a, 0x73, 0x69, 0x67,
	0x6e, 0x49, 0x6e, 0x44, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e,
	0x6d, 0x65, 0x73, 0x73, 0x65, 0x6e, 0x67, 0x65, 0x72, 0x70, 0x62, 0x2e, 0x53, 0x69, 0x67, 0x6e,
	0x49, 0x6e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x0a, 0x73, 0x69, 0x67, 0x6e, 0x49, 0x6e, 0x44, 0x61,
	0x74, 0x61, 0x22, 0x3e, 0x0a, 0x0a, 0x53, 0x69, 0x67, 0x6e, 0x49, 0x6e, 0x44, 0x61, 0x74, 0x61,
	0x12, 0x14, 0x0a, 0x05, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f,
	0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f,
	0x72, 0x64, 0x22, 0x18, 0x0a, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x22, 0x81, 0x01, 0x0a,
	0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2b, 0x0a, 0x06, 0x73, 0x65, 0x6e, 0x64,
	0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x65,
	0x6e, 0x67, 0x65, 0x72, 0x70, 0x62, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x52, 0x06, 0x73,
	0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x2f, 0x0a, 0x08, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x65, 0x6e,
	0x67, 0x65, 0x72, 0x70, 0x62, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x52, 0x08, 0x72, 0x65,
	0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x22, 0x24, 0x0a, 0x0a, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x41, 0x63, 0x6b, 0x12, 0x16,
	0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x32, 0x89, 0x02, 0x0a, 0x10, 0x4d, 0x65, 0x73, 0x73, 0x65,
	0x6e, 0x67, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x38, 0x0a, 0x06, 0x53,
	0x69, 0x67, 0x6e, 0x55, 0x70, 0x12, 0x17, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x65, 0x6e, 0x67, 0x65,
	0x72, 0x70, 0x62, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x55, 0x70, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x13,
	0x2e, 0x6d, 0x65, 0x73, 0x73, 0x65, 0x6e, 0x67, 0x65, 0x72, 0x70, 0x62, 0x2e, 0x55, 0x73, 0x65,
	0x72, 0x49, 0x44, 0x22, 0x00, 0x12, 0x38, 0x0a, 0x06, 0x53, 0x69, 0x67, 0x6e, 0x49, 0x6e, 0x12,
	0x17, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x65, 0x6e, 0x67, 0x65, 0x72, 0x70, 0x62, 0x2e, 0x53, 0x69,
	0x67, 0x6e, 0x49, 0x6e, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x13, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x65,
	0x6e, 0x67, 0x65, 0x72, 0x70, 0x62, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x22, 0x00, 0x12,
	0x40, 0x0a, 0x0b, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x14,
	0x2e, 0x6d, 0x65, 0x73, 0x73, 0x65, 0x6e, 0x67, 0x65, 0x72, 0x70, 0x62, 0x2e, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x1a, 0x17, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x65, 0x6e, 0x67, 0x65, 0x72,
	0x70, 0x62, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x41, 0x63, 0x6b, 0x22, 0x00, 0x28,
	0x01, 0x12, 0x3f, 0x0a, 0x0e, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x13, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x65, 0x6e, 0x67, 0x65, 0x72, 0x70,
	0x62, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x1a, 0x14, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x65,
	0x6e, 0x67, 0x65, 0x72, 0x70, 0x62, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x00,
	0x30, 0x01, 0x42, 0x0f, 0x5a, 0x0d, 0x2e, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x65, 0x6e, 0x67, 0x65,
	0x72, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_messengerpb_messenger_proto_rawDescOnce sync.Once
	file_messengerpb_messenger_proto_rawDescData = file_messengerpb_messenger_proto_rawDesc
)

func file_messengerpb_messenger_proto_rawDescGZIP() []byte {
	file_messengerpb_messenger_proto_rawDescOnce.Do(func() {
		file_messengerpb_messenger_proto_rawDescData = protoimpl.X.CompressGZIP(file_messengerpb_messenger_proto_rawDescData)
	})
	return file_messengerpb_messenger_proto_rawDescData
}

var file_messengerpb_messenger_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_messengerpb_messenger_proto_goTypes = []interface{}{
	(*SignUpData)(nil), // 0: messengerpb.SignUpData
	(*SignInData)(nil), // 1: messengerpb.SignInData
	(*UserID)(nil),     // 2: messengerpb.UserID
	(*Message)(nil),    // 3: messengerpb.Message
	(*MessageAck)(nil), // 4: messengerpb.MessageAck
}
var file_messengerpb_messenger_proto_depIdxs = []int32{
	1, // 0: messengerpb.SignUpData.signInData:type_name -> messengerpb.SignInData
	2, // 1: messengerpb.Message.sender:type_name -> messengerpb.UserID
	2, // 2: messengerpb.Message.receiver:type_name -> messengerpb.UserID
	0, // 3: messengerpb.MessengerService.SignUp:input_type -> messengerpb.SignUpData
	1, // 4: messengerpb.MessengerService.SignIn:input_type -> messengerpb.SignInData
	3, // 5: messengerpb.MessengerService.SendMessage:input_type -> messengerpb.Message
	2, // 6: messengerpb.MessengerService.ReceiveMessage:input_type -> messengerpb.UserID
	2, // 7: messengerpb.MessengerService.SignUp:output_type -> messengerpb.UserID
	2, // 8: messengerpb.MessengerService.SignIn:output_type -> messengerpb.UserID
	4, // 9: messengerpb.MessengerService.SendMessage:output_type -> messengerpb.MessageAck
	3, // 10: messengerpb.MessengerService.ReceiveMessage:output_type -> messengerpb.Message
	7, // [7:11] is the sub-list for method output_type
	3, // [3:7] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_messengerpb_messenger_proto_init() }
func file_messengerpb_messenger_proto_init() {
	if File_messengerpb_messenger_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_messengerpb_messenger_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignUpData); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_messengerpb_messenger_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignInData); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_messengerpb_messenger_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserID); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_messengerpb_messenger_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_messengerpb_messenger_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageAck); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_messengerpb_messenger_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_messengerpb_messenger_proto_goTypes,
		DependencyIndexes: file_messengerpb_messenger_proto_depIdxs,
		MessageInfos:      file_messengerpb_messenger_proto_msgTypes,
	}.Build()
	File_messengerpb_messenger_proto = out.File
	file_messengerpb_messenger_proto_rawDesc = nil
	file_messengerpb_messenger_proto_goTypes = nil
	file_messengerpb_messenger_proto_depIdxs = nil
}
