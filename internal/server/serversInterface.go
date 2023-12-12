package server

import (
	"GRPC_Server/proto/gRPCService/classInfoProto"
	"GRPC_Server/proto/gRPCService/studentProto"
	"context"
)

// StudentHandlerInterface defines the methods required for handling student-related requests.
type StudentHandlerInterface interface {
	Create(ctx context.Context, req *studentProto.StudentCreateRequest) *studentProto.StudentCreateResponse
	Update(ctx context.Context, req *studentProto.StudentUpdateRequest) *studentProto.StudentUpdateResponse
	Get(ctx context.Context, req *studentProto.StudentGetRequest) *studentProto.StudentGetResponse
	Delete(ctx context.Context, req *studentProto.StudentDeleteRequest) *studentProto.StudentDeleteResponse
}

// ClassInfoHandlerInterface defines the methods required for handling class information-related requests.
type ClassInfoHandlerInterface interface {
	AddClass(ctx context.Context, req *classInfoProto.ClassAddRequest) *classInfoProto.ClassAddResponse
	UpdateClass(ctx context.Context, req *classInfoProto.ClassUpdateRequest) *classInfoProto.ClassUpdateResponse
	DeleteClassByStudent(ctx context.Context, req *classInfoProto.ClassDeleteRequest) *classInfoProto.ClassDeleteResponse
	GetAllClassesByStudent(ctx context.Context, req *classInfoProto.ClassesGetRequest) *classInfoProto.ClassesGetResponse
}
