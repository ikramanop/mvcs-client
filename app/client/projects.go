package client

import (
	"context"
	"github.com/ikramanop/mvcs-client/app/config"
	"github.com/ikramanop/mvcs-client/app/model"
	"github.com/ikramanop/mvcs-client/app/protobuf/projects"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	ProjectsClientCreationError = "ошибка при создании клиента проекта"
	ProjectCreationError        = "ошибка при создании проекта"
	ProjectFindError            = "ошибка при поиске проекта"
)

type Projects interface {
	Create(ctx context.Context, name string) (*projects.Project, error)
	//todo
	//UploadFile(ctx context.Context, branchId int32, file *FileRequest) (*projects.UploadFileResponse, error)
	//UploadFiles(ctx context.Context, branchId int32, files *FileRequest[]) (*projects.UploadFilesResponse, error)
	//GetFileVersions(ctx context.Context, branchId int32, ...) (*projects.FileResponse, error)
	Find(ctx context.Context, name string) (*projects.Project, error)
}

type ProjectsClient struct {
	conn *grpc.ClientConn
	c    projects.ProjectsClient
}

type ProjectsClientStub struct{}

func NewProjectsClient(cfg *config.ClientConfig) (*ProjectsClient, error) {
	conn, err := grpc.Dial(cfg.ServerHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, model.WrapError(ProjectsClientCreationError, err)
	}

	c := projects.NewProjectsClient(conn)

	return &ProjectsClient{
		conn,
		c,
	}, nil
}

func NewProjectsStubClient(cfg *config.ClientConfig) (*ProjectsClientStub, error) {
	return &ProjectsClientStub{}, nil
}

func (cl *ProjectsClient) Create(ctx context.Context, name string) (*projects.Project, error) {
	project, err := cl.c.Create(ctx, &projects.CreateRequest{
		Name: name,
	})
	if err != nil {
		return nil, model.WrapError(ProjectCreationError, err)
	}

	return project, nil
}

func (cl *ProjectsClientStub) Create(ctx context.Context, name string) (*projects.Project, error) {
	return nil, nil
}

func (cl *ProjectsClient) Find(ctx context.Context, name string) (*projects.Project, error) {
	project, err := cl.c.Find(ctx, &projects.FindRequest{
		Name: name,
	})
	if err != nil {
		return nil, model.WrapError(ProjectFindError, err)
	}

	return project, nil
}

func (cl *ProjectsClientStub) Find(ctx context.Context, name string) (*projects.Project, error) {
	return nil, nil
}
