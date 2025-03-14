// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.15.8
// source: proto/prysm/storage/version.proto

package ethereum_eth_storage

import (
	reflect "reflect"
	sync "sync"

	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Version int32

const (
	Version_UNKNOWN   Version = 0
	Version_PHASE0    Version = 1
	Version_ALTAIR    Version = 2
	Version_BELLATRIX Version = 3
)

// Enum value maps for Version.
var (
	Version_name = map[int32]string{
		0: "UNKNOWN",
		1: "PHASE0",
		2: "ALTAIR",
		3: "BELLATRIX",
	}
	Version_value = map[string]int32{
		"UNKNOWN":   0,
		"PHASE0":    1,
		"ALTAIR":    2,
		"BELLATRIX": 3,
	}
)

func (x Version) Enum() *Version {
	p := new(Version)
	*p = x
	return p
}

func (x Version) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Version) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_prysm_storage_version_proto_enumTypes[0].Descriptor()
}

func (Version) Type() protoreflect.EnumType {
	return &file_proto_prysm_storage_version_proto_enumTypes[0]
}

func (x Version) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Version.Descriptor instead.
func (Version) EnumDescriptor() ([]byte, []int) {
	return file_proto_prysm_storage_version_proto_rawDescGZIP(), []int{0}
}

var File_proto_prysm_storage_version_proto protoreflect.FileDescriptor

var file_proto_prysm_storage_version_proto_rawDesc = []byte{
	0x0a, 0x21, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x72, 0x79, 0x73, 0x6d, 0x2f, 0x73, 0x74,
	0x6f, 0x72, 0x61, 0x67, 0x65, 0x2f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x14, 0x65, 0x74, 0x68, 0x65, 0x72, 0x65, 0x75, 0x6d, 0x2e, 0x65, 0x74,
	0x68, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2a, 0x3d, 0x0a, 0x07, 0x56, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10,
	0x00, 0x12, 0x0a, 0x0a, 0x06, 0x50, 0x48, 0x41, 0x53, 0x45, 0x30, 0x10, 0x01, 0x12, 0x0a, 0x0a,
	0x06, 0x41, 0x4c, 0x54, 0x41, 0x49, 0x52, 0x10, 0x02, 0x12, 0x0d, 0x0a, 0x09, 0x42, 0x45, 0x4c,
	0x4c, 0x41, 0x54, 0x52, 0x49, 0x58, 0x10, 0x03, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_prysm_storage_version_proto_rawDescOnce sync.Once
	file_proto_prysm_storage_version_proto_rawDescData = file_proto_prysm_storage_version_proto_rawDesc
)

func file_proto_prysm_storage_version_proto_rawDescGZIP() []byte {
	file_proto_prysm_storage_version_proto_rawDescOnce.Do(func() {
		file_proto_prysm_storage_version_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_prysm_storage_version_proto_rawDescData)
	})
	return file_proto_prysm_storage_version_proto_rawDescData
}

var file_proto_prysm_storage_version_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_prysm_storage_version_proto_goTypes = []interface{}{
	(Version)(0), // 0: ethereum.eth.storage.Version
}
var file_proto_prysm_storage_version_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_prysm_storage_version_proto_init() }
func file_proto_prysm_storage_version_proto_init() {
	if File_proto_prysm_storage_version_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_prysm_storage_version_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_prysm_storage_version_proto_goTypes,
		DependencyIndexes: file_proto_prysm_storage_version_proto_depIdxs,
		EnumInfos:         file_proto_prysm_storage_version_proto_enumTypes,
	}.Build()
	File_proto_prysm_storage_version_proto = out.File
	file_proto_prysm_storage_version_proto_rawDesc = nil
	file_proto_prysm_storage_version_proto_goTypes = nil
	file_proto_prysm_storage_version_proto_depIdxs = nil
}
