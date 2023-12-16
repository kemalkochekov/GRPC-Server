// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.25.0
// source: classInfo.proto

package classInfoProto

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

// ClassInfo
type ClassInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	StudentID int64  `protobuf:"varint,2,opt,name=studentID,proto3" json:"studentID,omitempty"`
	ClassName string `protobuf:"bytes,3,opt,name=className,proto3" json:"className,omitempty"`
}

func (x *ClassInfo) Reset() {
	*x = ClassInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_classInfo_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClassInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClassInfo) ProtoMessage() {}

func (x *ClassInfo) ProtoReflect() protoreflect.Message {
	mi := &file_classInfo_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClassInfo.ProtoReflect.Descriptor instead.
func (*ClassInfo) Descriptor() ([]byte, []int) {
	return file_classInfo_proto_rawDescGZIP(), []int{0}
}

func (x *ClassInfo) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *ClassInfo) GetStudentID() int64 {
	if x != nil {
		return x.StudentID
	}
	return 0
}

func (x *ClassInfo) GetClassName() string {
	if x != nil {
		return x.ClassName
	}
	return ""
}

type ClassAddRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StudentID int64  `protobuf:"varint,1,opt,name=studentID,proto3" json:"studentID,omitempty"`
	ClassName string `protobuf:"bytes,2,opt,name=className,proto3" json:"className,omitempty"`
}

func (x *ClassAddRequest) Reset() {
	*x = ClassAddRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_classInfo_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClassAddRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClassAddRequest) ProtoMessage() {}

func (x *ClassAddRequest) ProtoReflect() protoreflect.Message {
	mi := &file_classInfo_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClassAddRequest.ProtoReflect.Descriptor instead.
func (*ClassAddRequest) Descriptor() ([]byte, []int) {
	return file_classInfo_proto_rawDescGZIP(), []int{1}
}

func (x *ClassAddRequest) GetStudentID() int64 {
	if x != nil {
		return x.StudentID
	}
	return 0
}

func (x *ClassAddRequest) GetClassName() string {
	if x != nil {
		return x.ClassName
	}
	return ""
}

type ClassAddResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClassInfo *ClassInfo `protobuf:"bytes,1,opt,name=classInfo,proto3" json:"classInfo,omitempty"`
	Message   string     `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *ClassAddResponse) Reset() {
	*x = ClassAddResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_classInfo_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClassAddResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClassAddResponse) ProtoMessage() {}

func (x *ClassAddResponse) ProtoReflect() protoreflect.Message {
	mi := &file_classInfo_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClassAddResponse.ProtoReflect.Descriptor instead.
func (*ClassAddResponse) Descriptor() ([]byte, []int) {
	return file_classInfo_proto_rawDescGZIP(), []int{2}
}

func (x *ClassAddResponse) GetClassInfo() *ClassInfo {
	if x != nil {
		return x.ClassInfo
	}
	return nil
}

func (x *ClassAddResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type ClassUpdateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StudentID int64  `protobuf:"varint,1,opt,name=studentID,proto3" json:"studentID,omitempty"`
	ClassName string `protobuf:"bytes,2,opt,name=className,proto3" json:"className,omitempty"`
}

func (x *ClassUpdateRequest) Reset() {
	*x = ClassUpdateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_classInfo_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClassUpdateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClassUpdateRequest) ProtoMessage() {}

func (x *ClassUpdateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_classInfo_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClassUpdateRequest.ProtoReflect.Descriptor instead.
func (*ClassUpdateRequest) Descriptor() ([]byte, []int) {
	return file_classInfo_proto_rawDescGZIP(), []int{3}
}

func (x *ClassUpdateRequest) GetStudentID() int64 {
	if x != nil {
		return x.StudentID
	}
	return 0
}

func (x *ClassUpdateRequest) GetClassName() string {
	if x != nil {
		return x.ClassName
	}
	return ""
}

type ClassUpdateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *ClassUpdateResponse) Reset() {
	*x = ClassUpdateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_classInfo_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClassUpdateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClassUpdateResponse) ProtoMessage() {}

func (x *ClassUpdateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_classInfo_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClassUpdateResponse.ProtoReflect.Descriptor instead.
func (*ClassUpdateResponse) Descriptor() ([]byte, []int) {
	return file_classInfo_proto_rawDescGZIP(), []int{4}
}

func (x *ClassUpdateResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type ClassDeleteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StudentID int64 `protobuf:"varint,1,opt,name=studentID,proto3" json:"studentID,omitempty"`
}

func (x *ClassDeleteRequest) Reset() {
	*x = ClassDeleteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_classInfo_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClassDeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClassDeleteRequest) ProtoMessage() {}

func (x *ClassDeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_classInfo_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClassDeleteRequest.ProtoReflect.Descriptor instead.
func (*ClassDeleteRequest) Descriptor() ([]byte, []int) {
	return file_classInfo_proto_rawDescGZIP(), []int{5}
}

func (x *ClassDeleteRequest) GetStudentID() int64 {
	if x != nil {
		return x.StudentID
	}
	return 0
}

type ClassDeleteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *ClassDeleteResponse) Reset() {
	*x = ClassDeleteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_classInfo_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClassDeleteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClassDeleteResponse) ProtoMessage() {}

func (x *ClassDeleteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_classInfo_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClassDeleteResponse.ProtoReflect.Descriptor instead.
func (*ClassDeleteResponse) Descriptor() ([]byte, []int) {
	return file_classInfo_proto_rawDescGZIP(), []int{6}
}

func (x *ClassDeleteResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type ClassesGetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StudentID int64 `protobuf:"varint,1,opt,name=studentID,proto3" json:"studentID,omitempty"`
}

func (x *ClassesGetRequest) Reset() {
	*x = ClassesGetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_classInfo_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClassesGetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClassesGetRequest) ProtoMessage() {}

func (x *ClassesGetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_classInfo_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClassesGetRequest.ProtoReflect.Descriptor instead.
func (*ClassesGetRequest) Descriptor() ([]byte, []int) {
	return file_classInfo_proto_rawDescGZIP(), []int{7}
}

func (x *ClassesGetRequest) GetStudentID() int64 {
	if x != nil {
		return x.StudentID
	}
	return 0
}

type ClassesGetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClassInfo []*ClassInfo `protobuf:"bytes,1,rep,name=classInfo,proto3" json:"classInfo,omitempty"`
	Message   string       `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *ClassesGetResponse) Reset() {
	*x = ClassesGetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_classInfo_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClassesGetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClassesGetResponse) ProtoMessage() {}

func (x *ClassesGetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_classInfo_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClassesGetResponse.ProtoReflect.Descriptor instead.
func (*ClassesGetResponse) Descriptor() ([]byte, []int) {
	return file_classInfo_proto_rawDescGZIP(), []int{8}
}

func (x *ClassesGetResponse) GetClassInfo() []*ClassInfo {
	if x != nil {
		return x.ClassInfo
	}
	return nil
}

func (x *ClassesGetResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_classInfo_proto protoreflect.FileDescriptor

var file_classInfo_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e,
	0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x57, 0x0a, 0x09, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1c, 0x0a, 0x09,
	0x73, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x09, 0x73, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x12, 0x1c, 0x0a, 0x09, 0x63, 0x6c,
	0x61, 0x73, 0x73, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63,
	0x6c, 0x61, 0x73, 0x73, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x4d, 0x0a, 0x0f, 0x43, 0x6c, 0x61, 0x73,
	0x73, 0x41, 0x64, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x73,
	0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09,
	0x73, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x12, 0x1c, 0x0a, 0x09, 0x63, 0x6c, 0x61,
	0x73, 0x73, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x6c,
	0x61, 0x73, 0x73, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x56, 0x0a, 0x10, 0x43, 0x6c, 0x61, 0x73, 0x73,
	0x41, 0x64, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x28, 0x0a, 0x09, 0x63,
	0x6c, 0x61, 0x73, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a,
	0x2e, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x09, 0x63, 0x6c, 0x61, 0x73,
	0x73, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22,
	0x50, 0x0a, 0x12, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74,
	0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x73, 0x74, 0x75, 0x64, 0x65, 0x6e,
	0x74, 0x49, 0x44, 0x12, 0x1c, 0x0a, 0x09, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x4e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x4e, 0x61, 0x6d,
	0x65, 0x22, 0x2f, 0x0a, 0x13, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x22, 0x32, 0x0a, 0x12, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x74, 0x75, 0x64,
	0x65, 0x6e, 0x74, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x73, 0x74, 0x75,
	0x64, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x22, 0x2f, 0x0a, 0x13, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x31, 0x0a, 0x11, 0x43, 0x6c, 0x61, 0x73, 0x73,
	0x65, 0x73, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09,
	0x73, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x09, 0x73, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x22, 0x58, 0x0a, 0x12, 0x43, 0x6c,
	0x61, 0x73, 0x73, 0x65, 0x73, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x28, 0x0a, 0x09, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x52,
	0x09, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x32, 0x81, 0x03, 0x0a, 0x10, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x49, 0x6e,
	0x66, 0x6f, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4a, 0x0a, 0x08, 0x41, 0x64, 0x64,
	0x43, 0x6c, 0x61, 0x73, 0x73, 0x12, 0x10, 0x2e, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x41, 0x64, 0x64,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x41,
	0x64, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x19, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x13, 0x3a, 0x01, 0x2a, 0x22, 0x0e, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6c, 0x61, 0x73, 0x73,
	0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x12, 0x53, 0x0a, 0x0b, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43,
	0x6c, 0x61, 0x73, 0x73, 0x12, 0x13, 0x2e, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x43, 0x6c, 0x61, 0x73,
	0x73, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x19, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x13, 0x3a, 0x01, 0x2a, 0x1a, 0x0e, 0x2f, 0x76, 0x31, 0x2f,
	0x63, 0x6c, 0x61, 0x73, 0x73, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x12, 0x65, 0x0a, 0x14, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x42, 0x79, 0x53, 0x74, 0x75, 0x64, 0x65,
	0x6e, 0x74, 0x12, 0x13, 0x2e, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x22, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x1c, 0x2a, 0x1a, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6c, 0x61, 0x73, 0x73,
	0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x2f, 0x7b, 0x73, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x49, 0x44,
	0x7d, 0x12, 0x65, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x43, 0x6c, 0x61, 0x73, 0x73,
	0x65, 0x73, 0x42, 0x79, 0x53, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x12, 0x12, 0x2e, 0x43, 0x6c,
	0x61, 0x73, 0x73, 0x65, 0x73, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x13, 0x2e, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x65, 0x73, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x22, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1c, 0x12, 0x1a, 0x2f, 0x76,
	0x31, 0x2f, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x2f, 0x7b, 0x73, 0x74,
	0x75, 0x64, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x7d, 0x42, 0x1c, 0x5a, 0x1a, 0x67, 0x52, 0x50, 0x43,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x49, 0x6e, 0x66,
	0x6f, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_classInfo_proto_rawDescOnce sync.Once
	file_classInfo_proto_rawDescData = file_classInfo_proto_rawDesc
)

func file_classInfo_proto_rawDescGZIP() []byte {
	file_classInfo_proto_rawDescOnce.Do(func() {
		file_classInfo_proto_rawDescData = protoimpl.X.CompressGZIP(file_classInfo_proto_rawDescData)
	})
	return file_classInfo_proto_rawDescData
}

var file_classInfo_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_classInfo_proto_goTypes = []interface{}{
	(*ClassInfo)(nil),           // 0: ClassInfo
	(*ClassAddRequest)(nil),     // 1: ClassAddRequest
	(*ClassAddResponse)(nil),    // 2: ClassAddResponse
	(*ClassUpdateRequest)(nil),  // 3: ClassUpdateRequest
	(*ClassUpdateResponse)(nil), // 4: ClassUpdateResponse
	(*ClassDeleteRequest)(nil),  // 5: ClassDeleteRequest
	(*ClassDeleteResponse)(nil), // 6: ClassDeleteResponse
	(*ClassesGetRequest)(nil),   // 7: ClassesGetRequest
	(*ClassesGetResponse)(nil),  // 8: ClassesGetResponse
}
var file_classInfo_proto_depIdxs = []int32{
	0, // 0: ClassAddResponse.classInfo:type_name -> ClassInfo
	0, // 1: ClassesGetResponse.classInfo:type_name -> ClassInfo
	1, // 2: ClassInfoService.AddClass:input_type -> ClassAddRequest
	3, // 3: ClassInfoService.UpdateClass:input_type -> ClassUpdateRequest
	5, // 4: ClassInfoService.DeleteClassByStudent:input_type -> ClassDeleteRequest
	7, // 5: ClassInfoService.GetAllClassesByStudent:input_type -> ClassesGetRequest
	2, // 6: ClassInfoService.AddClass:output_type -> ClassAddResponse
	4, // 7: ClassInfoService.UpdateClass:output_type -> ClassUpdateResponse
	6, // 8: ClassInfoService.DeleteClassByStudent:output_type -> ClassDeleteResponse
	8, // 9: ClassInfoService.GetAllClassesByStudent:output_type -> ClassesGetResponse
	6, // [6:10] is the sub-list for method output_type
	2, // [2:6] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_classInfo_proto_init() }
func file_classInfo_proto_init() {
	if File_classInfo_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_classInfo_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClassInfo); i {
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
		file_classInfo_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClassAddRequest); i {
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
		file_classInfo_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClassAddResponse); i {
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
		file_classInfo_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClassUpdateRequest); i {
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
		file_classInfo_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClassUpdateResponse); i {
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
		file_classInfo_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClassDeleteRequest); i {
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
		file_classInfo_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClassDeleteResponse); i {
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
		file_classInfo_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClassesGetRequest); i {
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
		file_classInfo_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClassesGetResponse); i {
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
			RawDescriptor: file_classInfo_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_classInfo_proto_goTypes,
		DependencyIndexes: file_classInfo_proto_depIdxs,
		MessageInfos:      file_classInfo_proto_msgTypes,
	}.Build()
	File_classInfo_proto = out.File
	file_classInfo_proto_rawDesc = nil
	file_classInfo_proto_goTypes = nil
	file_classInfo_proto_depIdxs = nil
}