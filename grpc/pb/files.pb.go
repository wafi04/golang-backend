// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.3
// 	protoc        v3.12.4
// source: grpc/pb/files.proto

package pb

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

type FileUploadRequest_FileType int32

const (
	FileUploadRequest_IMAGE    FileUploadRequest_FileType = 0
	FileUploadRequest_VIDEO    FileUploadRequest_FileType = 1
	FileUploadRequest_DOCUMENT FileUploadRequest_FileType = 2
)

// Enum value maps for FileUploadRequest_FileType.
var (
	FileUploadRequest_FileType_name = map[int32]string{
		0: "IMAGE",
		1: "VIDEO",
		2: "DOCUMENT",
	}
	FileUploadRequest_FileType_value = map[string]int32{
		"IMAGE":    0,
		"VIDEO":    1,
		"DOCUMENT": 2,
	}
)

func (x FileUploadRequest_FileType) Enum() *FileUploadRequest_FileType {
	p := new(FileUploadRequest_FileType)
	*p = x
	return p
}

func (x FileUploadRequest_FileType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (FileUploadRequest_FileType) Descriptor() protoreflect.EnumDescriptor {
	return file_grpc_pb_files_proto_enumTypes[0].Descriptor()
}

func (FileUploadRequest_FileType) Type() protoreflect.EnumType {
	return &file_grpc_pb_files_proto_enumTypes[0]
}

func (x FileUploadRequest_FileType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use FileUploadRequest_FileType.Descriptor instead.
func (FileUploadRequest_FileType) EnumDescriptor() ([]byte, []int) {
	return file_grpc_pb_files_proto_rawDescGZIP(), []int{1, 0}
}

type FileUploadResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Url           string                 `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	PublicId      string                 `protobuf:"bytes,2,opt,name=public_id,json=publicId,proto3" json:"public_id,omitempty"`
	Error         string                 `protobuf:"bytes,3,opt,name=error,proto3" json:"error,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FileUploadResponse) Reset() {
	*x = FileUploadResponse{}
	mi := &file_grpc_pb_files_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FileUploadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileUploadResponse) ProtoMessage() {}

func (x *FileUploadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_pb_files_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileUploadResponse.ProtoReflect.Descriptor instead.
func (*FileUploadResponse) Descriptor() ([]byte, []int) {
	return file_grpc_pb_files_proto_rawDescGZIP(), []int{0}
}

func (x *FileUploadResponse) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *FileUploadResponse) GetPublicId() string {
	if x != nil {
		return x.PublicId
	}
	return ""
}

func (x *FileUploadResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type FileUploadRequest struct {
	state         protoimpl.MessageState     `protogen:"open.v1"`
	FileData      []byte                     `protobuf:"bytes,1,opt,name=file_data,json=fileData,proto3" json:"file_data,omitempty"`
	FileName      string                     `protobuf:"bytes,2,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	Folder        string                     `protobuf:"bytes,3,opt,name=folder,proto3" json:"folder,omitempty"`
	PublicId      string                     `protobuf:"bytes,4,opt,name=public_id,json=publicId,proto3" json:"public_id,omitempty"`
	FileType      FileUploadRequest_FileType `protobuf:"varint,5,opt,name=file_type,json=fileType,proto3,enum=pb.FileUploadRequest_FileType" json:"file_type,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FileUploadRequest) Reset() {
	*x = FileUploadRequest{}
	mi := &file_grpc_pb_files_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FileUploadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileUploadRequest) ProtoMessage() {}

func (x *FileUploadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_pb_files_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileUploadRequest.ProtoReflect.Descriptor instead.
func (*FileUploadRequest) Descriptor() ([]byte, []int) {
	return file_grpc_pb_files_proto_rawDescGZIP(), []int{1}
}

func (x *FileUploadRequest) GetFileData() []byte {
	if x != nil {
		return x.FileData
	}
	return nil
}

func (x *FileUploadRequest) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *FileUploadRequest) GetFolder() string {
	if x != nil {
		return x.Folder
	}
	return ""
}

func (x *FileUploadRequest) GetPublicId() string {
	if x != nil {
		return x.PublicId
	}
	return ""
}

func (x *FileUploadRequest) GetFileType() FileUploadRequest_FileType {
	if x != nil {
		return x.FileType
	}
	return FileUploadRequest_IMAGE
}

var File_grpc_pb_files_proto protoreflect.FileDescriptor

var file_grpc_pb_files_proto_rawDesc = []byte{
	0x0a, 0x13, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x62, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x22, 0x59, 0x0a, 0x12, 0x46, 0x69, 0x6c,
	0x65, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72,
	0x6c, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x49, 0x64, 0x12, 0x14,
	0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x22, 0xef, 0x01, 0x0a, 0x11, 0x46, 0x69, 0x6c, 0x65, 0x55, 0x70, 0x6c,
	0x6f, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69,
	0x6c, 0x65, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x66,
	0x69, 0x6c, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x12, 0x1b, 0x0a, 0x09,
	0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x49, 0x64, 0x12, 0x3b, 0x0a, 0x09, 0x66, 0x69, 0x6c,
	0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1e, 0x2e, 0x70,
	0x62, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x08, 0x66, 0x69,
	0x6c, 0x65, 0x54, 0x79, 0x70, 0x65, 0x22, 0x2e, 0x0a, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x09, 0x0a, 0x05, 0x49, 0x4d, 0x41, 0x47, 0x45, 0x10, 0x00, 0x12, 0x09, 0x0a,
	0x05, 0x56, 0x49, 0x44, 0x45, 0x4f, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x44, 0x4f, 0x43, 0x55,
	0x4d, 0x45, 0x4e, 0x54, 0x10, 0x02, 0x32, 0x4c, 0x0a, 0x0b, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3d, 0x0a, 0x0a, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x46,
	0x69, 0x6c, 0x65, 0x12, 0x15, 0x2e, 0x70, 0x62, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x55, 0x70, 0x6c,
	0x6f, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x70, 0x62, 0x2e,
	0x46, 0x69, 0x6c, 0x65, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x42, 0x24, 0x5a, 0x22, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x77, 0x61, 0x66, 0x69, 0x30, 0x34, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2f,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_grpc_pb_files_proto_rawDescOnce sync.Once
	file_grpc_pb_files_proto_rawDescData = file_grpc_pb_files_proto_rawDesc
)

func file_grpc_pb_files_proto_rawDescGZIP() []byte {
	file_grpc_pb_files_proto_rawDescOnce.Do(func() {
		file_grpc_pb_files_proto_rawDescData = protoimpl.X.CompressGZIP(file_grpc_pb_files_proto_rawDescData)
	})
	return file_grpc_pb_files_proto_rawDescData
}

var file_grpc_pb_files_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_grpc_pb_files_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_grpc_pb_files_proto_goTypes = []any{
	(FileUploadRequest_FileType)(0), // 0: pb.FileUploadRequest.FileType
	(*FileUploadResponse)(nil),      // 1: pb.FileUploadResponse
	(*FileUploadRequest)(nil),       // 2: pb.FileUploadRequest
}
var file_grpc_pb_files_proto_depIdxs = []int32{
	0, // 0: pb.FileUploadRequest.file_type:type_name -> pb.FileUploadRequest.FileType
	2, // 1: pb.FileService.UploadFile:input_type -> pb.FileUploadRequest
	1, // 2: pb.FileService.UploadFile:output_type -> pb.FileUploadResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_grpc_pb_files_proto_init() }
func file_grpc_pb_files_proto_init() {
	if File_grpc_pb_files_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_grpc_pb_files_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_grpc_pb_files_proto_goTypes,
		DependencyIndexes: file_grpc_pb_files_proto_depIdxs,
		EnumInfos:         file_grpc_pb_files_proto_enumTypes,
		MessageInfos:      file_grpc_pb_files_proto_msgTypes,
	}.Build()
	File_grpc_pb_files_proto = out.File
	file_grpc_pb_files_proto_rawDesc = nil
	file_grpc_pb_files_proto_goTypes = nil
	file_grpc_pb_files_proto_depIdxs = nil
}
