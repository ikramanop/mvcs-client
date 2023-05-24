package client

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/ikramanop/mvcs-client/app/config"
	"github.com/ikramanop/mvcs-client/app/model"
	"github.com/ikramanop/mvcs-client/app/protobuf/messages"
	"github.com/ikramanop/mvcs-client/app/protobuf/projects"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	ProjectsClientCreationError = "ошибка при создании клиента проекта"
	ProjectCreationError        = "ошибка при создании проекта"
	ProjectFindError            = "ошибка при поиске проекта"
	ProjectGetAllError          = "ошибка при получении списка проектов"
)

type Projects interface {
	Create(ctx context.Context, name string) (*messages.Project, error)
	Find(ctx context.Context, name string) (*messages.Project, error)
	GetAll(ctx context.Context) (*messages.ProjectList, error)
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

func (cl *ProjectsClient) Create(ctx context.Context, name string) (*messages.Project, error) {
	project, err := cl.c.Create(ctx, &messages.CreateRequest{
		Name: name,
	})

	if err != nil {
		return nil, model.WrapError(ProjectCreationError, err)
	}

	return project, nil
}

func (cl *ProjectsClientStub) Create(ctx context.Context, name string) (*messages.Project, error) {
	return nil, nil
}

func (cl *ProjectsClient) Find(ctx context.Context, name string) (*messages.Project, error) {
	project, err := cl.c.Find(ctx, &messages.FindRequest{
		Name: name,
	})
	if project == nil {
		return nil, nil
	}
	if err != nil {
		return nil, model.WrapError(ProjectFindError, err)
	}

	return project, nil
}

func (cl *ProjectsClientStub) Find(ctx context.Context, name string) (*messages.Project, error) {
	return nil, nil
}

func (cl *ProjectsClient) GetAll(ctx context.Context) (*messages.ProjectList, error) {
	projectsList, err := cl.c.GetAll(ctx, &empty.Empty{})
	if err != nil {
		return nil, model.WrapError(ProjectGetAllError, err)
	}

	return projectsList, nil
}

func (cl *ProjectsClientStub) GetAll(ctx context.Context) (*messages.ProjectList, error) {
	return nil, nil
}
