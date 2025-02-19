// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v3.21.12
// source: health.proto

package types

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

type State int32

const (
	State_Invalid_State State = 0
	State_Up            State = 1
	State_Down          State = 2
	State_Unreachable   State = 3
)

// Enum value maps for State.
var (
	State_name = map[int32]string{
		0: "Invalid_State",
		1: "Up",
		2: "Down",
		3: "Unreachable",
	}
	State_value = map[string]int32{
		"Invalid_State": 0,
		"Up":            1,
		"Down":          2,
		"Unreachable":   3,
	}
)

func (x State) Enum() *State {
	p := new(State)
	*p = x
	return p
}

func (x State) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (State) Descriptor() protoreflect.EnumDescriptor {
	return file_health_proto_enumTypes[0].Descriptor()
}

func (State) Type() protoreflect.EnumType {
	return &file_health_proto_enumTypes[0]
}

func (x State) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use State.Descriptor instead.
func (State) EnumDescriptor() ([]byte, []int) {
	return file_health_proto_rawDescGZIP(), []int{0}
}

type HealthPoint struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AUuid    string          `protobuf:"bytes,1,opt,name=a_uuid,json=aUuid,proto3" json:"a_uuid,omitempty"`
	ZUuid    string          `protobuf:"bytes,2,opt,name=z_uuid,json=zUuid,proto3" json:"z_uuid,omitempty"`
	Alias    string          `protobuf:"bytes,3,opt,name=alias,proto3" json:"alias,omitempty"`
	Services map[string]bool `protobuf:"bytes,4,rep,name=services,proto3" json:"services,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	Status   State           `protobuf:"varint,5,opt,name=status,proto3,enum=types.State" json:"status,omitempty"`
}

func (x *HealthPoint) Reset() {
	*x = HealthPoint{}
	if protoimpl.UnsafeEnabled {
		mi := &file_health_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HealthPoint) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthPoint) ProtoMessage() {}

func (x *HealthPoint) ProtoReflect() protoreflect.Message {
	mi := &file_health_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HealthPoint.ProtoReflect.Descriptor instead.
func (*HealthPoint) Descriptor() ([]byte, []int) {
	return file_health_proto_rawDescGZIP(), []int{0}
}

func (x *HealthPoint) GetAUuid() string {
	if x != nil {
		return x.AUuid
	}
	return ""
}

func (x *HealthPoint) GetZUuid() string {
	if x != nil {
		return x.ZUuid
	}
	return ""
}

func (x *HealthPoint) GetAlias() string {
	if x != nil {
		return x.Alias
	}
	return ""
}

func (x *HealthPoint) GetServices() map[string]bool {
	if x != nil {
		return x.Services
	}
	return nil
}

func (x *HealthPoint) GetStatus() State {
	if x != nil {
		return x.Status
	}
	return State_Invalid_State
}

var File_health_proto protoreflect.FileDescriptor

var file_health_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05,
	0x74, 0x79, 0x70, 0x65, 0x73, 0x22, 0xf2, 0x01, 0x0a, 0x0b, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68,
	0x50, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x15, 0x0a, 0x06, 0x61, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x61, 0x55, 0x75, 0x69, 0x64, 0x12, 0x15, 0x0a, 0x06,
	0x7a, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x7a, 0x55,
	0x75, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x12, 0x3c, 0x0a, 0x08, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x74, 0x79,
	0x70, 0x65, 0x73, 0x2e, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x2e,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x12, 0x24, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e,
	0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x1a, 0x3b, 0x0a,
	0x0d, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x2a, 0x3d, 0x0a, 0x05, 0x53, 0x74,
	0x61, 0x74, 0x65, 0x12, 0x11, 0x0a, 0x0d, 0x49, 0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x5f, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x10, 0x00, 0x12, 0x06, 0x0a, 0x02, 0x55, 0x70, 0x10, 0x01, 0x12, 0x08,
	0x0a, 0x04, 0x44, 0x6f, 0x77, 0x6e, 0x10, 0x02, 0x12, 0x0f, 0x0a, 0x0b, 0x55, 0x6e, 0x72, 0x65,
	0x61, 0x63, 0x68, 0x61, 0x62, 0x6c, 0x65, 0x10, 0x03, 0x42, 0x24, 0x0a, 0x10, 0x63, 0x6f, 0x6d,
	0x2e, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x42, 0x05, 0x54,
	0x79, 0x70, 0x65, 0x73, 0x50, 0x01, 0x5a, 0x07, 0x2e, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_health_proto_rawDescOnce sync.Once
	file_health_proto_rawDescData = file_health_proto_rawDesc
)

func file_health_proto_rawDescGZIP() []byte {
	file_health_proto_rawDescOnce.Do(func() {
		file_health_proto_rawDescData = protoimpl.X.CompressGZIP(file_health_proto_rawDescData)
	})
	return file_health_proto_rawDescData
}

var file_health_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_health_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_health_proto_goTypes = []interface{}{
	(State)(0),          // 0: types.State
	(*HealthPoint)(nil), // 1: types.HealthPoint
	nil,                 // 2: types.HealthPoint.ServicesEntry
}
var file_health_proto_depIdxs = []int32{
	2, // 0: types.HealthPoint.services:type_name -> types.HealthPoint.ServicesEntry
	0, // 1: types.HealthPoint.status:type_name -> types.State
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_health_proto_init() }
func file_health_proto_init() {
	if File_health_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_health_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HealthPoint); i {
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
			RawDescriptor: file_health_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_health_proto_goTypes,
		DependencyIndexes: file_health_proto_depIdxs,
		EnumInfos:         file_health_proto_enumTypes,
		MessageInfos:      file_health_proto_msgTypes,
	}.Build()
	File_health_proto = out.File
	file_health_proto_rawDesc = nil
	file_health_proto_goTypes = nil
	file_health_proto_depIdxs = nil
}
