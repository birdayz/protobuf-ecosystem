// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        (unknown)
// source: protoconfig/v1/options.proto

package protoconfigv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type FieldOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Env *string `protobuf:"bytes,1,opt,name=env,proto3,oneof" json:"env,omitempty"`
}

func (x *FieldOptions) Reset() {
	*x = FieldOptions{}
	mi := &file_protoconfig_v1_options_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FieldOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FieldOptions) ProtoMessage() {}

func (x *FieldOptions) ProtoReflect() protoreflect.Message {
	mi := &file_protoconfig_v1_options_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FieldOptions.ProtoReflect.Descriptor instead.
func (*FieldOptions) Descriptor() ([]byte, []int) {
	return file_protoconfig_v1_options_proto_rawDescGZIP(), []int{0}
}

func (x *FieldOptions) GetEnv() string {
	if x != nil && x.Env != nil {
		return *x.Env
	}
	return ""
}

var file_protoconfig_v1_options_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*FieldOptions)(nil),
		Field:         50116,
		Name:          "protoconfig.v1.options",
		Tag:           "bytes,50116,opt,name=options",
		Filename:      "protoconfig/v1/options.proto",
	},
}

// Extension fields to descriptorpb.FieldOptions.
var (
	// optional protoconfig.v1.FieldOptions options = 50116;
	E_Options = &file_protoconfig_v1_options_proto_extTypes[0]
)

var File_protoconfig_v1_options_proto protoreflect.FileDescriptor

var file_protoconfig_v1_options_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x76, 0x31,
	0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x1a, 0x20,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x2d, 0x0a, 0x0c, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x12, 0x15, 0x0a, 0x03, 0x65, 0x6e, 0x76, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52,
	0x03, 0x65, 0x6e, 0x76, 0x88, 0x01, 0x01, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x65, 0x6e, 0x76, 0x3a,
	0x5a, 0x0a, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x1d, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65,
	0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xc4, 0x87, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e,
	0x76, 0x31, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52,
	0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x88, 0x01, 0x01, 0x42, 0xd8, 0x01, 0x0a, 0x12,
	0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e,
	0x76, 0x31, 0x42, 0x0c, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x50, 0x01, 0x5a, 0x5b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62,
	0x69, 0x72, 0x64, 0x61, 0x79, 0x7a, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2d,
	0x65, 0x63, 0x6f, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x65, 0x6e, 0x2f,
	0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x76,
	0x31, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x76, 0x31, 0xa2,
	0x02, 0x03, 0x50, 0x58, 0x58, 0xaa, 0x02, 0x0e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x1a, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0f, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protoconfig_v1_options_proto_rawDescOnce sync.Once
	file_protoconfig_v1_options_proto_rawDescData = file_protoconfig_v1_options_proto_rawDesc
)

func file_protoconfig_v1_options_proto_rawDescGZIP() []byte {
	file_protoconfig_v1_options_proto_rawDescOnce.Do(func() {
		file_protoconfig_v1_options_proto_rawDescData = protoimpl.X.CompressGZIP(file_protoconfig_v1_options_proto_rawDescData)
	})
	return file_protoconfig_v1_options_proto_rawDescData
}

var file_protoconfig_v1_options_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_protoconfig_v1_options_proto_goTypes = []any{
	(*FieldOptions)(nil),              // 0: protoconfig.v1.FieldOptions
	(*descriptorpb.FieldOptions)(nil), // 1: google.protobuf.FieldOptions
}
var file_protoconfig_v1_options_proto_depIdxs = []int32{
	1, // 0: protoconfig.v1.options:extendee -> google.protobuf.FieldOptions
	0, // 1: protoconfig.v1.options:type_name -> protoconfig.v1.FieldOptions
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	1, // [1:2] is the sub-list for extension type_name
	0, // [0:1] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_protoconfig_v1_options_proto_init() }
func file_protoconfig_v1_options_proto_init() {
	if File_protoconfig_v1_options_proto != nil {
		return
	}
	file_protoconfig_v1_options_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_protoconfig_v1_options_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 1,
			NumServices:   0,
		},
		GoTypes:           file_protoconfig_v1_options_proto_goTypes,
		DependencyIndexes: file_protoconfig_v1_options_proto_depIdxs,
		MessageInfos:      file_protoconfig_v1_options_proto_msgTypes,
		ExtensionInfos:    file_protoconfig_v1_options_proto_extTypes,
	}.Build()
	File_protoconfig_v1_options_proto = out.File
	file_protoconfig_v1_options_proto_rawDesc = nil
	file_protoconfig_v1_options_proto_goTypes = nil
	file_protoconfig_v1_options_proto_depIdxs = nil
}