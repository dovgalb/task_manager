package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	tmv1 "github.com/dovgalb/taskmanager_proto/gen/go/task_manager"
)

type gRPCServerApi struct {
	tmv1.UnimplementedTaskCategoryServer
}

func Register(gRPC *grpc.Server) {
	tmv1.RegisterTaskCategoryServer(gRPC, &gRPCServerApi{})
}

func (tm *gRPCServerApi) CreateTaskCategory(ctx context.Context, request *tmv1.CreateTaskCategoryRequest) (*tmv1.CreateTaskCategoryResponse, error) {
	//TODO implement me
	return &tmv1.CreateTaskCategoryResponse{TaskCategoryId: 1}, nil

}

func (tm *gRPCServerApi) ReadTaskCategory(ctx context.Context, request *tmv1.ReadTaskCategoryRequest) (*tmv1.TaskCategoryResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (tm *gRPCServerApi) UpdateTaskCategory(ctx context.Context, request *tmv1.UpdateTaskCategoryRequest) (*tmv1.TaskCategoryResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (tm *gRPCServerApi) DeleteTaskCategory(ctx context.Context, request *tmv1.DeleteTaskCategoryRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}
