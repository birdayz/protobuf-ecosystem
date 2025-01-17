// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        (unknown)
// source: protoconfig/v1/example.proto

package protoconfigv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Test_TestEnum int32

const (
	Test_TEST_ENUM_UNSPECIFIED Test_TestEnum = 0
	Test_TEST_ENUM_SOME_VAL    Test_TestEnum = 1
)

// Enum value maps for Test_TestEnum.
var (
	Test_TestEnum_name = map[int32]string{
		0: "TEST_ENUM_UNSPECIFIED",
		1: "TEST_ENUM_SOME_VAL",
	}
	Test_TestEnum_value = map[string]int32{
		"TEST_ENUM_UNSPECIFIED": 0,
		"TEST_ENUM_SOME_VAL":    1,
	}
)

func (x Test_TestEnum) Enum() *Test_TestEnum {
	p := new(Test_TestEnum)
	*p = x
	return p
}

func (x Test_TestEnum) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Test_TestEnum) Descriptor() protoreflect.EnumDescriptor {
	return file_protoconfig_v1_example_proto_enumTypes[0].Descriptor()
}

func (Test_TestEnum) Type() protoreflect.EnumType {
	return &file_protoconfig_v1_example_proto_enumTypes[0]
}

func (x Test_TestEnum) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Test_TestEnum.Descriptor instead.
func (Test_TestEnum) EnumDescriptor() ([]byte, []int) {
	return file_protoconfig_v1_example_proto_rawDescGZIP(), []int{3, 0}
}

type Test_ExampleEnum int32

const (
	Test_EXAMPLE_ENUM_UNSPECIFIED Test_ExampleEnum = 0
	Test_EXAMPLE_ENUM_EXAMPLE_VAL Test_ExampleEnum = 1
)

// Enum value maps for Test_ExampleEnum.
var (
	Test_ExampleEnum_name = map[int32]string{
		0: "EXAMPLE_ENUM_UNSPECIFIED",
		1: "EXAMPLE_ENUM_EXAMPLE_VAL",
	}
	Test_ExampleEnum_value = map[string]int32{
		"EXAMPLE_ENUM_UNSPECIFIED": 0,
		"EXAMPLE_ENUM_EXAMPLE_VAL": 1,
	}
)

func (x Test_ExampleEnum) Enum() *Test_ExampleEnum {
	p := new(Test_ExampleEnum)
	*p = x
	return p
}

func (x Test_ExampleEnum) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Test_ExampleEnum) Descriptor() protoreflect.EnumDescriptor {
	return file_protoconfig_v1_example_proto_enumTypes[1].Descriptor()
}

func (Test_ExampleEnum) Type() protoreflect.EnumType {
	return &file_protoconfig_v1_example_proto_enumTypes[1]
}

func (x Test_ExampleEnum) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Test_ExampleEnum.Descriptor instead.
func (Test_ExampleEnum) EnumDescriptor() ([]byte, []int) {
	return file_protoconfig_v1_example_proto_rawDescGZIP(), []int{3, 1}
}

type Nested struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StringField      string `protobuf:"bytes,1,opt,name=string_field,json=stringField,proto3" json:"string_field,omitempty"`
	NotUpdatedViaEnv string `protobuf:"bytes,2,opt,name=not_updated_via_env,json=notUpdatedViaEnv,proto3" json:"not_updated_via_env,omitempty"`
}

func (x *Nested) Reset() {
	*x = Nested{}
	mi := &file_protoconfig_v1_example_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Nested) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Nested) ProtoMessage() {}

func (x *Nested) ProtoReflect() protoreflect.Message {
	mi := &file_protoconfig_v1_example_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Nested.ProtoReflect.Descriptor instead.
func (*Nested) Descriptor() ([]byte, []int) {
	return file_protoconfig_v1_example_proto_rawDescGZIP(), []int{0}
}

func (x *Nested) GetStringField() string {
	if x != nil {
		return x.StringField
	}
	return ""
}

func (x *Nested) GetNotUpdatedViaEnv() string {
	if x != nil {
		return x.NotUpdatedViaEnv
	}
	return ""
}

type Nested2 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StringField      string `protobuf:"bytes,1,opt,name=string_field,json=stringField,proto3" json:"string_field,omitempty"`
	NotUpdatedViaEnv string `protobuf:"bytes,2,opt,name=not_updated_via_env,json=notUpdatedViaEnv,proto3" json:"not_updated_via_env,omitempty"`
}

func (x *Nested2) Reset() {
	*x = Nested2{}
	mi := &file_protoconfig_v1_example_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Nested2) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Nested2) ProtoMessage() {}

func (x *Nested2) ProtoReflect() protoreflect.Message {
	mi := &file_protoconfig_v1_example_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Nested2.ProtoReflect.Descriptor instead.
func (*Nested2) Descriptor() ([]byte, []int) {
	return file_protoconfig_v1_example_proto_rawDescGZIP(), []int{1}
}

func (x *Nested2) GetStringField() string {
	if x != nil {
		return x.StringField
	}
	return ""
}

func (x *Nested2) GetNotUpdatedViaEnv() string {
	if x != nil {
		return x.NotUpdatedViaEnv
	}
	return ""
}

type NestedWithNested struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NestedNested *NestedWithNested_NestedNested `protobuf:"bytes,1,opt,name=nested_nested,json=nestedNested,proto3" json:"nested_nested,omitempty"`
}

func (x *NestedWithNested) Reset() {
	*x = NestedWithNested{}
	mi := &file_protoconfig_v1_example_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NestedWithNested) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NestedWithNested) ProtoMessage() {}

func (x *NestedWithNested) ProtoReflect() protoreflect.Message {
	mi := &file_protoconfig_v1_example_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NestedWithNested.ProtoReflect.Descriptor instead.
func (*NestedWithNested) Descriptor() ([]byte, []int) {
	return file_protoconfig_v1_example_proto_rawDescGZIP(), []int{2}
}

func (x *NestedWithNested) GetNestedNested() *NestedWithNested_NestedNested {
	if x != nil {
		return x.NestedNested
	}
	return nil
}

type Test struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BoolField     bool             `protobuf:"varint,1,opt,name=bool_field,json=boolField,proto3" json:"bool_field,omitempty"`
	EnumField     Test_ExampleEnum `protobuf:"varint,2,opt,name=enum_field,json=enumField,proto3,enum=protoconfig.v1.Test_ExampleEnum" json:"enum_field,omitempty"`
	Int32Field    int32            `protobuf:"varint,3,opt,name=int32_field,json=int32Field,proto3" json:"int32_field,omitempty"`
	Sint32Field   int32            `protobuf:"zigzag32,4,opt,name=sint32_field,json=sint32Field,proto3" json:"sint32_field,omitempty"`
	Uint32Field   uint32           `protobuf:"varint,5,opt,name=uint32_field,json=uint32Field,proto3" json:"uint32_field,omitempty"`
	Int64Field    int64            `protobuf:"varint,6,opt,name=int64_field,json=int64Field,proto3" json:"int64_field,omitempty"`
	Sint64Field   int64            `protobuf:"zigzag64,7,opt,name=sint64_field,json=sint64Field,proto3" json:"sint64_field,omitempty"`
	Uint64Field   uint64           `protobuf:"varint,8,opt,name=uint64_field,json=uint64Field,proto3" json:"uint64_field,omitempty"`
	Sfixed32Field int32            `protobuf:"fixed32,9,opt,name=sfixed32_field,json=sfixed32Field,proto3" json:"sfixed32_field,omitempty"`
	Fixed32Field  uint32           `protobuf:"fixed32,10,opt,name=fixed32_field,json=fixed32Field,proto3" json:"fixed32_field,omitempty"`
	FloatField    float32          `protobuf:"fixed32,11,opt,name=float_field,json=floatField,proto3" json:"float_field,omitempty"`
	Sfixed64Field int64            `protobuf:"fixed64,12,opt,name=sfixed64_field,json=sfixed64Field,proto3" json:"sfixed64_field,omitempty"`
	Fixed64Field  uint64           `protobuf:"fixed64,13,opt,name=fixed64_field,json=fixed64Field,proto3" json:"fixed64_field,omitempty"`
	DoubleField   float64          `protobuf:"fixed64,14,opt,name=double_field,json=doubleField,proto3" json:"double_field,omitempty"`
	StringField   string           `protobuf:"bytes,15,opt,name=string_field,json=stringField,proto3" json:"string_field,omitempty"`
	BytesField    []byte           `protobuf:"bytes,16,opt,name=bytes_field,json=bytesField,proto3" json:"bytes_field,omitempty"`
	// Interesting to test for infinite recursion.
	MessageField          *Test                    `protobuf:"bytes,17,opt,name=message_field,json=messageField,proto3" json:"message_field,omitempty"`
	NestedMessageField    *Nested                  `protobuf:"bytes,18,opt,name=nested_message_field,json=nestedMessageField,proto3" json:"nested_message_field,omitempty"`
	OptionalInt32Field    *int32                   `protobuf:"varint,19,opt,name=optional_int32_field,json=optionalInt32Field,proto3,oneof" json:"optional_int32_field,omitempty"`
	RepeatedNestedMessage []*Nested2               `protobuf:"bytes,20,rep,name=repeated_nested_message,json=repeatedNestedMessage,proto3" json:"repeated_nested_message,omitempty"`
	StringToMap           map[string]*Nested2      `protobuf:"bytes,21,rep,name=string_to_map,json=stringToMap,proto3" json:"string_to_map,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	ListOfInts            []int32                  `protobuf:"varint,22,rep,packed,name=list_of_ints,json=listOfInts,proto3" json:"list_of_ints,omitempty"`
	ListOfStrings         []string                 `protobuf:"bytes,23,rep,name=list_of_strings,json=listOfStrings,proto3" json:"list_of_strings,omitempty"`
	ListOfEnums           []Test_ExampleEnum       `protobuf:"varint,24,rep,packed,name=list_of_enums,json=listOfEnums,proto3,enum=protoconfig.v1.Test_ExampleEnum" json:"list_of_enums,omitempty"`
	Timestamp             *timestamppb.Timestamp   `protobuf:"bytes,25,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Timestamps            []*timestamppb.Timestamp `protobuf:"bytes,26,rep,name=timestamps,proto3" json:"timestamps,omitempty"`
	OverriddenByEnv       *Nested2                 `protobuf:"bytes,27,opt,name=overridden_by_env,json=overriddenByEnv,proto3" json:"overridden_by_env,omitempty"`
	NestedWithNested      *NestedWithNested        `protobuf:"bytes,28,opt,name=nested_with_nested,json=nestedWithNested,proto3" json:"nested_with_nested,omitempty"`
	PrimitiveMap          map[string]string        `protobuf:"bytes,29,rep,name=primitive_map,json=primitiveMap,proto3" json:"primitive_map,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Test) Reset() {
	*x = Test{}
	mi := &file_protoconfig_v1_example_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Test) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Test) ProtoMessage() {}

func (x *Test) ProtoReflect() protoreflect.Message {
	mi := &file_protoconfig_v1_example_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Test.ProtoReflect.Descriptor instead.
func (*Test) Descriptor() ([]byte, []int) {
	return file_protoconfig_v1_example_proto_rawDescGZIP(), []int{3}
}

func (x *Test) GetBoolField() bool {
	if x != nil {
		return x.BoolField
	}
	return false
}

func (x *Test) GetEnumField() Test_ExampleEnum {
	if x != nil {
		return x.EnumField
	}
	return Test_EXAMPLE_ENUM_UNSPECIFIED
}

func (x *Test) GetInt32Field() int32 {
	if x != nil {
		return x.Int32Field
	}
	return 0
}

func (x *Test) GetSint32Field() int32 {
	if x != nil {
		return x.Sint32Field
	}
	return 0
}

func (x *Test) GetUint32Field() uint32 {
	if x != nil {
		return x.Uint32Field
	}
	return 0
}

func (x *Test) GetInt64Field() int64 {
	if x != nil {
		return x.Int64Field
	}
	return 0
}

func (x *Test) GetSint64Field() int64 {
	if x != nil {
		return x.Sint64Field
	}
	return 0
}

func (x *Test) GetUint64Field() uint64 {
	if x != nil {
		return x.Uint64Field
	}
	return 0
}

func (x *Test) GetSfixed32Field() int32 {
	if x != nil {
		return x.Sfixed32Field
	}
	return 0
}

func (x *Test) GetFixed32Field() uint32 {
	if x != nil {
		return x.Fixed32Field
	}
	return 0
}

func (x *Test) GetFloatField() float32 {
	if x != nil {
		return x.FloatField
	}
	return 0
}

func (x *Test) GetSfixed64Field() int64 {
	if x != nil {
		return x.Sfixed64Field
	}
	return 0
}

func (x *Test) GetFixed64Field() uint64 {
	if x != nil {
		return x.Fixed64Field
	}
	return 0
}

func (x *Test) GetDoubleField() float64 {
	if x != nil {
		return x.DoubleField
	}
	return 0
}

func (x *Test) GetStringField() string {
	if x != nil {
		return x.StringField
	}
	return ""
}

func (x *Test) GetBytesField() []byte {
	if x != nil {
		return x.BytesField
	}
	return nil
}

func (x *Test) GetMessageField() *Test {
	if x != nil {
		return x.MessageField
	}
	return nil
}

func (x *Test) GetNestedMessageField() *Nested {
	if x != nil {
		return x.NestedMessageField
	}
	return nil
}

func (x *Test) GetOptionalInt32Field() int32 {
	if x != nil && x.OptionalInt32Field != nil {
		return *x.OptionalInt32Field
	}
	return 0
}

func (x *Test) GetRepeatedNestedMessage() []*Nested2 {
	if x != nil {
		return x.RepeatedNestedMessage
	}
	return nil
}

func (x *Test) GetStringToMap() map[string]*Nested2 {
	if x != nil {
		return x.StringToMap
	}
	return nil
}

func (x *Test) GetListOfInts() []int32 {
	if x != nil {
		return x.ListOfInts
	}
	return nil
}

func (x *Test) GetListOfStrings() []string {
	if x != nil {
		return x.ListOfStrings
	}
	return nil
}

func (x *Test) GetListOfEnums() []Test_ExampleEnum {
	if x != nil {
		return x.ListOfEnums
	}
	return nil
}

func (x *Test) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *Test) GetTimestamps() []*timestamppb.Timestamp {
	if x != nil {
		return x.Timestamps
	}
	return nil
}

func (x *Test) GetOverriddenByEnv() *Nested2 {
	if x != nil {
		return x.OverriddenByEnv
	}
	return nil
}

func (x *Test) GetNestedWithNested() *NestedWithNested {
	if x != nil {
		return x.NestedWithNested
	}
	return nil
}

func (x *Test) GetPrimitiveMap() map[string]string {
	if x != nil {
		return x.PrimitiveMap
	}
	return nil
}

type NestedWithNested_NestedNested struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DeeplyNestedString string `protobuf:"bytes,1,opt,name=deeply_nested_string,json=deeplyNestedString,proto3" json:"deeply_nested_string,omitempty"`
}

func (x *NestedWithNested_NestedNested) Reset() {
	*x = NestedWithNested_NestedNested{}
	mi := &file_protoconfig_v1_example_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NestedWithNested_NestedNested) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NestedWithNested_NestedNested) ProtoMessage() {}

func (x *NestedWithNested_NestedNested) ProtoReflect() protoreflect.Message {
	mi := &file_protoconfig_v1_example_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NestedWithNested_NestedNested.ProtoReflect.Descriptor instead.
func (*NestedWithNested_NestedNested) Descriptor() ([]byte, []int) {
	return file_protoconfig_v1_example_proto_rawDescGZIP(), []int{2, 0}
}

func (x *NestedWithNested_NestedNested) GetDeeplyNestedString() string {
	if x != nil {
		return x.DeeplyNestedString
	}
	return ""
}

var File_protoconfig_v1_example_proto protoreflect.FileDescriptor

var file_protoconfig_v1_example_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x76, 0x31,
	0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x1a, 0x1c,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x76, 0x31, 0x2f, 0x6f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x79, 0x0a,
	0x06, 0x4e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x12, 0x40, 0x0a, 0x0c, 0x73, 0x74, 0x72, 0x69, 0x6e,
	0x67, 0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x1d, 0xa2,
	0xbc, 0x18, 0x19, 0x0a, 0x17, 0x4e, 0x45, 0x53, 0x54, 0x45, 0x44, 0x5f, 0x53, 0x54, 0x52, 0x49,
	0x4e, 0x47, 0x5f, 0x46, 0x49, 0x45, 0x4c, 0x44, 0x5f, 0x56, 0x41, 0x4c, 0x52, 0x0b, 0x73, 0x74,
	0x72, 0x69, 0x6e, 0x67, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x2d, 0x0a, 0x13, 0x6e, 0x6f, 0x74,
	0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x76, 0x69, 0x61, 0x5f, 0x65, 0x6e, 0x76,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x6e, 0x6f, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x64, 0x56, 0x69, 0x61, 0x45, 0x6e, 0x76, 0x22, 0x5b, 0x0a, 0x07, 0x4e, 0x65, 0x73, 0x74,
	0x65, 0x64, 0x32, 0x12, 0x21, 0x0a, 0x0c, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x5f, 0x66, 0x69,
	0x65, 0x6c, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x74, 0x72, 0x69, 0x6e,
	0x67, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x2d, 0x0a, 0x13, 0x6e, 0x6f, 0x74, 0x5f, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x76, 0x69, 0x61, 0x5f, 0x65, 0x6e, 0x76, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x10, 0x6e, 0x6f, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x56,
	0x69, 0x61, 0x45, 0x6e, 0x76, 0x22, 0xa8, 0x01, 0x0a, 0x10, 0x4e, 0x65, 0x73, 0x74, 0x65, 0x64,
	0x57, 0x69, 0x74, 0x68, 0x4e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x12, 0x52, 0x0a, 0x0d, 0x6e, 0x65,
	0x73, 0x74, 0x65, 0x64, 0x5f, 0x6e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x2d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e,
	0x76, 0x31, 0x2e, 0x4e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x57, 0x69, 0x74, 0x68, 0x4e, 0x65, 0x73,
	0x74, 0x65, 0x64, 0x2e, 0x4e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x4e, 0x65, 0x73, 0x74, 0x65, 0x64,
	0x52, 0x0c, 0x6e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x4e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x1a, 0x40,
	0x0a, 0x0c, 0x4e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x4e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x12, 0x30,
	0x0a, 0x14, 0x64, 0x65, 0x65, 0x70, 0x6c, 0x79, 0x5f, 0x6e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x5f,
	0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x64, 0x65,
	0x65, 0x70, 0x6c, 0x79, 0x4e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67,
	0x22, 0xef, 0x10, 0x0a, 0x04, 0x54, 0x65, 0x73, 0x74, 0x12, 0x33, 0x0a, 0x0a, 0x62, 0x6f, 0x6f,
	0x6c, 0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x42, 0x14, 0xa2,
	0xbc, 0x18, 0x10, 0x0a, 0x0e, 0x42, 0x4f, 0x4f, 0x4c, 0x5f, 0x46, 0x49, 0x45, 0x4c, 0x44, 0x5f,
	0x56, 0x41, 0x4c, 0x52, 0x09, 0x62, 0x6f, 0x6f, 0x6c, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x55,
	0x0a, 0x0a, 0x65, 0x6e, 0x75, 0x6d, 0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x20, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x2e, 0x76, 0x31, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x2e, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65,
	0x45, 0x6e, 0x75, 0x6d, 0x42, 0x14, 0xa2, 0xbc, 0x18, 0x10, 0x0a, 0x0e, 0x45, 0x4e, 0x55, 0x4d,
	0x5f, 0x46, 0x49, 0x45, 0x4c, 0x44, 0x5f, 0x56, 0x41, 0x4c, 0x52, 0x09, 0x65, 0x6e, 0x75, 0x6d,
	0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x36, 0x0a, 0x0b, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x5f, 0x66,
	0x69, 0x65, 0x6c, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x42, 0x15, 0xa2, 0xbc, 0x18, 0x11,
	0x0a, 0x0f, 0x49, 0x4e, 0x54, 0x33, 0x32, 0x5f, 0x46, 0x49, 0x45, 0x4c, 0x44, 0x5f, 0x56, 0x41,
	0x4c, 0x52, 0x0a, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x39, 0x0a,
	0x0c, 0x73, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x11, 0x42, 0x16, 0xa2, 0xbc, 0x18, 0x12, 0x0a, 0x10, 0x53, 0x49, 0x4e, 0x54, 0x33,
	0x32, 0x5f, 0x46, 0x49, 0x45, 0x4c, 0x44, 0x5f, 0x56, 0x41, 0x4c, 0x52, 0x0b, 0x73, 0x69, 0x6e,
	0x74, 0x33, 0x32, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x39, 0x0a, 0x0c, 0x75, 0x69, 0x6e, 0x74,
	0x33, 0x32, 0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x42, 0x16,
	0xa2, 0xbc, 0x18, 0x12, 0x0a, 0x10, 0x55, 0x49, 0x4e, 0x54, 0x33, 0x32, 0x5f, 0x46, 0x49, 0x45,
	0x4c, 0x44, 0x5f, 0x56, 0x41, 0x4c, 0x52, 0x0b, 0x75, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x46, 0x69,
	0x65, 0x6c, 0x64, 0x12, 0x36, 0x0a, 0x0b, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x5f, 0x66, 0x69, 0x65,
	0x6c, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x42, 0x15, 0xa2, 0xbc, 0x18, 0x11, 0x0a, 0x0f,
	0x49, 0x4e, 0x54, 0x36, 0x34, 0x5f, 0x46, 0x49, 0x45, 0x4c, 0x44, 0x5f, 0x56, 0x41, 0x4c, 0x52,
	0x0a, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x39, 0x0a, 0x0c, 0x73,
	0x69, 0x6e, 0x74, 0x36, 0x34, 0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x12, 0x42, 0x16, 0xa2, 0xbc, 0x18, 0x12, 0x0a, 0x10, 0x53, 0x49, 0x4e, 0x54, 0x36, 0x34, 0x5f,
	0x46, 0x49, 0x45, 0x4c, 0x44, 0x5f, 0x56, 0x41, 0x4c, 0x52, 0x0b, 0x73, 0x69, 0x6e, 0x74, 0x36,
	0x34, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x39, 0x0a, 0x0c, 0x75, 0x69, 0x6e, 0x74, 0x36, 0x34,
	0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x04, 0x42, 0x16, 0xa2, 0xbc,
	0x18, 0x12, 0x0a, 0x10, 0x55, 0x49, 0x4e, 0x54, 0x36, 0x34, 0x5f, 0x46, 0x49, 0x45, 0x4c, 0x44,
	0x5f, 0x56, 0x41, 0x4c, 0x52, 0x0b, 0x75, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x46, 0x69, 0x65, 0x6c,
	0x64, 0x12, 0x3f, 0x0a, 0x0e, 0x73, 0x66, 0x69, 0x78, 0x65, 0x64, 0x33, 0x32, 0x5f, 0x66, 0x69,
	0x65, 0x6c, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0f, 0x42, 0x18, 0xa2, 0xbc, 0x18, 0x14, 0x0a,
	0x12, 0x53, 0x46, 0x49, 0x58, 0x45, 0x44, 0x33, 0x32, 0x5f, 0x46, 0x49, 0x45, 0x4c, 0x44, 0x5f,
	0x56, 0x41, 0x4c, 0x52, 0x0d, 0x73, 0x66, 0x69, 0x78, 0x65, 0x64, 0x33, 0x32, 0x46, 0x69, 0x65,
	0x6c, 0x64, 0x12, 0x3c, 0x0a, 0x0d, 0x66, 0x69, 0x78, 0x65, 0x64, 0x33, 0x32, 0x5f, 0x66, 0x69,
	0x65, 0x6c, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x07, 0x42, 0x17, 0xa2, 0xbc, 0x18, 0x13, 0x0a,
	0x11, 0x46, 0x49, 0x58, 0x45, 0x44, 0x33, 0x32, 0x5f, 0x46, 0x49, 0x45, 0x4c, 0x44, 0x5f, 0x56,
	0x41, 0x4c, 0x52, 0x0c, 0x66, 0x69, 0x78, 0x65, 0x64, 0x33, 0x32, 0x46, 0x69, 0x65, 0x6c, 0x64,
	0x12, 0x36, 0x0a, 0x0b, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x18,
	0x0b, 0x20, 0x01, 0x28, 0x02, 0x42, 0x15, 0xa2, 0xbc, 0x18, 0x11, 0x0a, 0x0f, 0x46, 0x4c, 0x4f,
	0x41, 0x54, 0x5f, 0x46, 0x49, 0x45, 0x4c, 0x44, 0x5f, 0x56, 0x41, 0x4c, 0x52, 0x0a, 0x66, 0x6c,
	0x6f, 0x61, 0x74, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x3f, 0x0a, 0x0e, 0x73, 0x66, 0x69, 0x78,
	0x65, 0x64, 0x36, 0x34, 0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x10,
	0x42, 0x18, 0xa2, 0xbc, 0x18, 0x14, 0x0a, 0x12, 0x53, 0x46, 0x49, 0x58, 0x45, 0x44, 0x36, 0x34,
	0x5f, 0x46, 0x49, 0x45, 0x4c, 0x44, 0x5f, 0x56, 0x41, 0x4c, 0x52, 0x0d, 0x73, 0x66, 0x69, 0x78,
	0x65, 0x64, 0x36, 0x34, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x3c, 0x0a, 0x0d, 0x66, 0x69, 0x78,
	0x65, 0x64, 0x36, 0x34, 0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x06,
	0x42, 0x17, 0xa2, 0xbc, 0x18, 0x13, 0x0a, 0x11, 0x46, 0x49, 0x58, 0x45, 0x44, 0x36, 0x34, 0x5f,
	0x46, 0x49, 0x45, 0x4c, 0x44, 0x5f, 0x56, 0x41, 0x4c, 0x52, 0x0c, 0x66, 0x69, 0x78, 0x65, 0x64,
	0x36, 0x34, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x39, 0x0a, 0x0c, 0x64, 0x6f, 0x75, 0x62, 0x6c,
	0x65, 0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x01, 0x42, 0x16, 0xa2,
	0xbc, 0x18, 0x12, 0x0a, 0x10, 0x44, 0x4f, 0x55, 0x42, 0x4c, 0x45, 0x5f, 0x46, 0x49, 0x45, 0x4c,
	0x44, 0x5f, 0x56, 0x41, 0x4c, 0x52, 0x0b, 0x64, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x46, 0x69, 0x65,
	0x6c, 0x64, 0x12, 0x39, 0x0a, 0x0c, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x5f, 0x66, 0x69, 0x65,
	0x6c, 0x64, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x42, 0x16, 0xa2, 0xbc, 0x18, 0x12, 0x0a, 0x10,
	0x53, 0x54, 0x52, 0x49, 0x4e, 0x47, 0x5f, 0x46, 0x49, 0x45, 0x4c, 0x44, 0x5f, 0x56, 0x41, 0x4c,
	0x52, 0x0b, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x36, 0x0a,
	0x0b, 0x62, 0x79, 0x74, 0x65, 0x73, 0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x10, 0x20, 0x01,
	0x28, 0x0c, 0x42, 0x15, 0xa2, 0xbc, 0x18, 0x11, 0x0a, 0x0f, 0x42, 0x59, 0x54, 0x45, 0x53, 0x5f,
	0x46, 0x49, 0x45, 0x4c, 0x44, 0x5f, 0x56, 0x41, 0x4c, 0x52, 0x0a, 0x62, 0x79, 0x74, 0x65, 0x73,
	0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x39, 0x0a, 0x0d, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x11, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x65,
	0x73, 0x74, 0x52, 0x0c, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x46, 0x69, 0x65, 0x6c, 0x64,
	0x12, 0x48, 0x0a, 0x14, 0x6e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x12, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e,
	0x4e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x52, 0x12, 0x6e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x55, 0x0a, 0x14, 0x6f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x5f, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x5f, 0x66, 0x69, 0x65,
	0x6c, 0x64, 0x18, 0x13, 0x20, 0x01, 0x28, 0x05, 0x42, 0x1e, 0xa2, 0xbc, 0x18, 0x1a, 0x0a, 0x18,
	0x4f, 0x50, 0x54, 0x49, 0x4f, 0x4e, 0x41, 0x4c, 0x5f, 0x49, 0x4e, 0x54, 0x33, 0x32, 0x5f, 0x46,
	0x49, 0x45, 0x4c, 0x44, 0x5f, 0x56, 0x41, 0x4c, 0x48, 0x00, 0x52, 0x12, 0x6f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x61, 0x6c, 0x49, 0x6e, 0x74, 0x33, 0x32, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x88, 0x01,
	0x01, 0x12, 0x4f, 0x0a, 0x17, 0x72, 0x65, 0x70, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x6e, 0x65,
	0x73, 0x74, 0x65, 0x64, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x14, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x32, 0x52, 0x15, 0x72, 0x65, 0x70,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x4e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x49, 0x0a, 0x0d, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x5f, 0x74, 0x6f, 0x5f,
	0x6d, 0x61, 0x70, 0x18, 0x15, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x2e,
	0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x54, 0x6f, 0x4d, 0x61, 0x70, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x52, 0x0b, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x54, 0x6f, 0x4d, 0x61, 0x70, 0x12, 0x20, 0x0a,
	0x0c, 0x6c, 0x69, 0x73, 0x74, 0x5f, 0x6f, 0x66, 0x5f, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x16, 0x20,
	0x03, 0x28, 0x05, 0x52, 0x0a, 0x6c, 0x69, 0x73, 0x74, 0x4f, 0x66, 0x49, 0x6e, 0x74, 0x73, 0x12,
	0x26, 0x0a, 0x0f, 0x6c, 0x69, 0x73, 0x74, 0x5f, 0x6f, 0x66, 0x5f, 0x73, 0x74, 0x72, 0x69, 0x6e,
	0x67, 0x73, 0x18, 0x17, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0d, 0x6c, 0x69, 0x73, 0x74, 0x4f, 0x66,
	0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x73, 0x12, 0x44, 0x0a, 0x0d, 0x6c, 0x69, 0x73, 0x74, 0x5f,
	0x6f, 0x66, 0x5f, 0x65, 0x6e, 0x75, 0x6d, 0x73, 0x18, 0x18, 0x20, 0x03, 0x28, 0x0e, 0x32, 0x20,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e,
	0x54, 0x65, 0x73, 0x74, 0x2e, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x45, 0x6e, 0x75, 0x6d,
	0x52, 0x0b, 0x6c, 0x69, 0x73, 0x74, 0x4f, 0x66, 0x45, 0x6e, 0x75, 0x6d, 0x73, 0x12, 0x38, 0x0a,
	0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x19, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x3a, 0x0a, 0x0a, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x73, 0x18, 0x1a, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x73, 0x12, 0x43, 0x0a, 0x11, 0x6f, 0x76, 0x65, 0x72, 0x72, 0x69, 0x64, 0x64, 0x65,
	0x6e, 0x5f, 0x62, 0x79, 0x5f, 0x65, 0x6e, 0x76, 0x18, 0x1b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e,
	0x4e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x32, 0x52, 0x0f, 0x6f, 0x76, 0x65, 0x72, 0x72, 0x69, 0x64,
	0x64, 0x65, 0x6e, 0x42, 0x79, 0x45, 0x6e, 0x76, 0x12, 0x4e, 0x0a, 0x12, 0x6e, 0x65, 0x73, 0x74,
	0x65, 0x64, 0x5f, 0x77, 0x69, 0x74, 0x68, 0x5f, 0x6e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x18, 0x1c,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x57, 0x69, 0x74, 0x68,
	0x4e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x52, 0x10, 0x6e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x57, 0x69,
	0x74, 0x68, 0x4e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x12, 0x4b, 0x0a, 0x0d, 0x70, 0x72, 0x69, 0x6d,
	0x69, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x6d, 0x61, 0x70, 0x18, 0x1d, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x26, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31,
	0x2e, 0x54, 0x65, 0x73, 0x74, 0x2e, 0x50, 0x72, 0x69, 0x6d, 0x69, 0x74, 0x69, 0x76, 0x65, 0x4d,
	0x61, 0x70, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0c, 0x70, 0x72, 0x69, 0x6d, 0x69, 0x74, 0x69,
	0x76, 0x65, 0x4d, 0x61, 0x70, 0x1a, 0x57, 0x0a, 0x10, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x54,
	0x6f, 0x4d, 0x61, 0x70, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x2d, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x65, 0x73, 0x74,
	0x65, 0x64, 0x32, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x3f,
	0x0a, 0x11, 0x50, 0x72, 0x69, 0x6d, 0x69, 0x74, 0x69, 0x76, 0x65, 0x4d, 0x61, 0x70, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22,
	0x3d, 0x0a, 0x08, 0x54, 0x65, 0x73, 0x74, 0x45, 0x6e, 0x75, 0x6d, 0x12, 0x19, 0x0a, 0x15, 0x54,
	0x45, 0x53, 0x54, 0x5f, 0x45, 0x4e, 0x55, 0x4d, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49,
	0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x16, 0x0a, 0x12, 0x54, 0x45, 0x53, 0x54, 0x5f, 0x45,
	0x4e, 0x55, 0x4d, 0x5f, 0x53, 0x4f, 0x4d, 0x45, 0x5f, 0x56, 0x41, 0x4c, 0x10, 0x01, 0x22, 0x49,
	0x0a, 0x0b, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x45, 0x6e, 0x75, 0x6d, 0x12, 0x1c, 0x0a,
	0x18, 0x45, 0x58, 0x41, 0x4d, 0x50, 0x4c, 0x45, 0x5f, 0x45, 0x4e, 0x55, 0x4d, 0x5f, 0x55, 0x4e,
	0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x1c, 0x0a, 0x18, 0x45,
	0x58, 0x41, 0x4d, 0x50, 0x4c, 0x45, 0x5f, 0x45, 0x4e, 0x55, 0x4d, 0x5f, 0x45, 0x58, 0x41, 0x4d,
	0x50, 0x4c, 0x45, 0x5f, 0x56, 0x41, 0x4c, 0x10, 0x01, 0x42, 0x17, 0x0a, 0x15, 0x5f, 0x6f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x5f, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x5f, 0x66, 0x69, 0x65,
	0x6c, 0x64, 0x42, 0xd8, 0x01, 0x0a, 0x12, 0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x42, 0x0c, 0x45, 0x78, 0x61, 0x6d, 0x70,
	0x6c, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x5b, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x69, 0x72, 0x64, 0x61, 0x79, 0x7a, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2d, 0x65, 0x63, 0x6f, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x76, 0x31, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x50, 0x58, 0x58, 0xaa, 0x02, 0x0e, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0e,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5c, 0x56, 0x31, 0xe2, 0x02,
	0x1a, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5c, 0x56, 0x31, 0x5c,
	0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0f, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protoconfig_v1_example_proto_rawDescOnce sync.Once
	file_protoconfig_v1_example_proto_rawDescData = file_protoconfig_v1_example_proto_rawDesc
)

func file_protoconfig_v1_example_proto_rawDescGZIP() []byte {
	file_protoconfig_v1_example_proto_rawDescOnce.Do(func() {
		file_protoconfig_v1_example_proto_rawDescData = protoimpl.X.CompressGZIP(file_protoconfig_v1_example_proto_rawDescData)
	})
	return file_protoconfig_v1_example_proto_rawDescData
}

var file_protoconfig_v1_example_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_protoconfig_v1_example_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_protoconfig_v1_example_proto_goTypes = []any{
	(Test_TestEnum)(0),                    // 0: protoconfig.v1.Test.TestEnum
	(Test_ExampleEnum)(0),                 // 1: protoconfig.v1.Test.ExampleEnum
	(*Nested)(nil),                        // 2: protoconfig.v1.Nested
	(*Nested2)(nil),                       // 3: protoconfig.v1.Nested2
	(*NestedWithNested)(nil),              // 4: protoconfig.v1.NestedWithNested
	(*Test)(nil),                          // 5: protoconfig.v1.Test
	(*NestedWithNested_NestedNested)(nil), // 6: protoconfig.v1.NestedWithNested.NestedNested
	nil,                                   // 7: protoconfig.v1.Test.StringToMapEntry
	nil,                                   // 8: protoconfig.v1.Test.PrimitiveMapEntry
	(*timestamppb.Timestamp)(nil),         // 9: google.protobuf.Timestamp
}
var file_protoconfig_v1_example_proto_depIdxs = []int32{
	6,  // 0: protoconfig.v1.NestedWithNested.nested_nested:type_name -> protoconfig.v1.NestedWithNested.NestedNested
	1,  // 1: protoconfig.v1.Test.enum_field:type_name -> protoconfig.v1.Test.ExampleEnum
	5,  // 2: protoconfig.v1.Test.message_field:type_name -> protoconfig.v1.Test
	2,  // 3: protoconfig.v1.Test.nested_message_field:type_name -> protoconfig.v1.Nested
	3,  // 4: protoconfig.v1.Test.repeated_nested_message:type_name -> protoconfig.v1.Nested2
	7,  // 5: protoconfig.v1.Test.string_to_map:type_name -> protoconfig.v1.Test.StringToMapEntry
	1,  // 6: protoconfig.v1.Test.list_of_enums:type_name -> protoconfig.v1.Test.ExampleEnum
	9,  // 7: protoconfig.v1.Test.timestamp:type_name -> google.protobuf.Timestamp
	9,  // 8: protoconfig.v1.Test.timestamps:type_name -> google.protobuf.Timestamp
	3,  // 9: protoconfig.v1.Test.overridden_by_env:type_name -> protoconfig.v1.Nested2
	4,  // 10: protoconfig.v1.Test.nested_with_nested:type_name -> protoconfig.v1.NestedWithNested
	8,  // 11: protoconfig.v1.Test.primitive_map:type_name -> protoconfig.v1.Test.PrimitiveMapEntry
	3,  // 12: protoconfig.v1.Test.StringToMapEntry.value:type_name -> protoconfig.v1.Nested2
	13, // [13:13] is the sub-list for method output_type
	13, // [13:13] is the sub-list for method input_type
	13, // [13:13] is the sub-list for extension type_name
	13, // [13:13] is the sub-list for extension extendee
	0,  // [0:13] is the sub-list for field type_name
}

func init() { file_protoconfig_v1_example_proto_init() }
func file_protoconfig_v1_example_proto_init() {
	if File_protoconfig_v1_example_proto != nil {
		return
	}
	file_protoconfig_v1_options_proto_init()
	file_protoconfig_v1_example_proto_msgTypes[3].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_protoconfig_v1_example_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_protoconfig_v1_example_proto_goTypes,
		DependencyIndexes: file_protoconfig_v1_example_proto_depIdxs,
		EnumInfos:         file_protoconfig_v1_example_proto_enumTypes,
		MessageInfos:      file_protoconfig_v1_example_proto_msgTypes,
	}.Build()
	File_protoconfig_v1_example_proto = out.File
	file_protoconfig_v1_example_proto_rawDesc = nil
	file_protoconfig_v1_example_proto_goTypes = nil
	file_protoconfig_v1_example_proto_depIdxs = nil
}
