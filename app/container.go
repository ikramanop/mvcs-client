package app

import (
	"context"
	"github.com/ikramanop/mvcs-client/app/client"
	"github.com/ikramanop/mvcs-client/app/config"
	"github.com/ikramanop/mvcs-client/app/difftool"
	"github.com/ikramanop/mvcs-client/app/model"
	"github.com/ikramanop/mvcs-client/app/repository"
	"google.golang.org/grpc/metadata"
	"time"
)

const (
	GetConfigError         = "ошибка в получении конфигурации"
	GetAuthClientError     = "ошибка при создании клиента авторизации"
	GetRepositoryError     = "ошибка при подключении к БД"
	GetDiffToolError       = "ошибка при инициализации файловой системы"
	GetProjectsClientError = "ошибка при создании клиента проектов"
	GetProjectClientError  = "ошибка при создании клиента проекта"
)

type Container struct {
	Cfg      *config.ClientConfig
	Auth     client.Auth
	DB       repository.Repository
	Diff     difftool.Tool
	Projects client.Projects
	Project  client.Project
}

func NewContainer(init bool, skipDB bool, skipClients bool) (*Container, error) {
	cfg, err := config.NewConfigFromOS()
	if err != nil {
		return nil, model.WrapError(GetConfigError, err)
	}

	var authClient client.Auth
	if skipClients {
		authClient, _ = client.NewAuthStubClient(cfg)
	} else {
		authClient, err = client.NewAuthClient(cfg)
		if err != nil {
			return nil, model.WrapError(GetAuthClientError, err)
		}
	}

	var DB repository.Repository
	if skipDB {
		DB = nil
	} else {
		DB, err = repository.NewRepository(init)
		if err != nil {
			return nil, model.WrapError(GetRepositoryError, err)
		}
	}

	diffTool, err := difftool.NewDiffTool()
	if err != nil {
		return nil, model.WrapError(GetDiffToolError, err)
	}

	var projectsClient client.Projects
	if skipClients {
		projectsClient, _ = client.NewProjectsStubClient(cfg)
	} else {
		projectsClient, err = client.NewProjectsClient(cfg)
		if err != nil {
			return nil, model.WrapError(GetProjectsClientError, err)
		}
	}

	var projectClient client.Project
	if skipClients {
		projectClient, _ = client.NewProjectStubClient(cfg)
	} else {
		projectClient, err = client.NewProjectClient(cfg)
		if err != nil {
			return nil, model.WrapError(GetProjectClientError, err)
		}
	}

	return &Container{
		Cfg:      cfg,
		Auth:     authClient,
		DB:       DB,
		Diff:     diffTool,
		Projects: projectsClient,
		Project:  projectClient,
	}, nil
}

func (c *Container) GetContext(auth bool, projectId ...string) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*c.Cfg.RequestTimeout)

	if auth {
		ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", c.Cfg.Auth.Key)
	}
	if projectId != nil {
		ctx = metadata.AppendToOutgoingContext(ctx, "ProjectId", projectId[0])
	}
	ctx = metadata.AppendToOutgoingContext(ctx, "Content-Type", "application/grpc")

	return ctx, cancel
}
