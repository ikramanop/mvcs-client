package client

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/ikramanop/mvcs-client/app/config"
	"github.com/ikramanop/mvcs-client/app/model"
	"github.com/ikramanop/mvcs-client/app/protobuf/messages"
	"github.com/ikramanop/mvcs-client/app/protobuf/project"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	ProjectClientCreationError = "ошибка при создании ветки"
	BranchCreationError        = "ошибка при создании ветки"
	GetAllBranchesError        = "ошибка при получении веток"
	FileUploadError            = "ошибка при загрузке файла"
	FilesUploadError           = "ошибка при загрузке файлов"
	GetFileVersionError        = "ошибка при получении версии файла"
	UploadFileError            = "ошибка при загрузки файла на remote"
)

type Project interface {
	CreateBranch(ctx context.Context, name string, parentBranchId *int32) (*messages.Branch, error)
	GetAllBranches(ctx context.Context) (*messages.GetBranchesResponse, error)
	//todo
	UploadFile(ctx context.Context, branchId int32, file *messages.FileRequest) (*messages.UploadFileResponse, error)
	//UploadFiles(ctx context.Context, branchId int32, file *model.File) (*messages.UploadFilesResponse, error)
	//GetFileVersions(ctx context.Context) (*messages.FileResponse, error)
}

type ProjectClient struct {
	conn *grpc.ClientConn
	c    project.ProjectClient
}

type ProjectClientStub struct{}

func NewProjectClient(cfg *config.ClientConfig) (*ProjectClient, error) {
	conn, err := grpc.Dial(cfg.ServerHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, model.WrapError(ProjectClientCreationError, err)
	}

	c := project.NewProjectClient(conn)

	return &ProjectClient{
		conn,
		c,
	}, nil
}

func NewProjectStubClient(cfg *config.ClientConfig) (*ProjectClientStub, error) {
	return &ProjectClientStub{}, nil
}

func (cl *ProjectClient) CreateBranch(ctx context.Context, name string, parentBranchId *int32) (*messages.Branch, error) {
	branch, err := cl.c.CreateBranch(ctx, &messages.CreateBranchRequest{
		Name:           name,
		ParentBranchId: parentBranchId,
	})
	if err != nil {
		return nil, model.WrapError(BranchCreationError, err)
	}
	return branch, nil
}

func (cl *ProjectClientStub) CreateBranch(ctx context.Context, name string, parentBranchId *int32) (*messages.Branch, error) {
	return nil, nil
}

func (cl *ProjectClient) GetAllBranches(ctx context.Context) (*messages.GetBranchesResponse, error) {
	branches, err := cl.c.GetAllBranches(ctx, &empty.Empty{})
	if err != nil {
		return nil, model.WrapError(GetAllBranchesError, err)
	}

	return branches, nil
}

func (cl *ProjectClientStub) GetAllBranches(ctx context.Context) (*messages.GetBranchesResponse, error) {
	return nil, nil
}

func (cl *ProjectClient) UploadFile(ctx context.Context, branchId int32, file *messages.FileRequest) (*messages.UploadFileResponse, error) {
	fileResponse, err := cl.c.UploadFile(ctx, &messages.UploadFileRequest{
		BranchId: branchId,
		File:     file,
	})

	if err != nil {
		return nil, model.WrapError(UploadFileError, err)
	}

	return fileResponse, nil
}

func (cl *ProjectClientStub) UploadFile(ctx context.Context, branchId int32, file *messages.FileRequest) (*messages.UploadFileResponse, error) {
	return nil, nil
}
