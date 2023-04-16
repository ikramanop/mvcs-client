package app

import (
	"context"
	"github.com/ikramanop/mvcs-client/app/client"
	"github.com/ikramanop/mvcs-client/app/config"
	"github.com/ikramanop/mvcs-client/app/model"
	"github.com/ikramanop/mvcs-client/app/repository"
	"google.golang.org/grpc/metadata"
	"time"
)

const (
	GetConfigError     = "ошибка в получении конфигурации"
	GetAuthClientError = "ошибка при создании клиента авторизации"
	GetRepositoryError = "ошибка при подключении к БД"
)

type Container struct {
	Cfg  *config.ClientConfig
	Auth client.Auth
	DB   repository.Repository
}

func NewContainer(init bool, skipDB bool) (*Container, error) {
	cfg, err := config.NewConfigFromOS()
	if err != nil {
		return nil, model.WrapError(GetConfigError, err)
	}

	authClient, err := client.NewAuthClient(cfg)
	if err != nil {
		return nil, model.WrapError(GetAuthClientError, err)
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

	return &Container{
		Cfg:  cfg,
		Auth: authClient,
		DB:   DB,
	}, nil
}

func (c *Container) GetContext(auth bool) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*c.Cfg.RequestTimeout)

	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", c.Cfg.Auth.Key)

	return ctx, cancel
}
