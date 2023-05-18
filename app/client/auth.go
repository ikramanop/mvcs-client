package client

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/ikramanop/mvcs-client/app/config"
	"github.com/ikramanop/mvcs-client/app/model"
	"github.com/ikramanop/mvcs-client/app/protobuf/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	AuthClientCreationError = "ошибка при создании клиента авторизации"
	AuthClientRegisterError = "ошибка при регистрации"
	AuthClientMeError       = "ошибка при проверке пользователя"
)

type Auth interface {
	Register(ctx context.Context, login string, key string) error
	Me(ctx context.Context) (*auth.MeMessage, error)
}

type AuthClient struct {
	conn *grpc.ClientConn
	c    auth.AuthClient
}

type AuthClientStub struct{}

func NewAuthClient(cfg *config.ClientConfig) (*AuthClient, error) {
	conn, err := grpc.Dial(cfg.ServerHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, model.WrapError(AuthClientCreationError, err)
	}

	c := auth.NewAuthClient(conn)

	return &AuthClient{
		conn,
		c,
	}, nil
}

func NewAuthStubClient(cfg *config.ClientConfig) (*AuthClientStub, error) {
	return &AuthClientStub{}, nil
}

func (cl *AuthClient) Register(ctx context.Context, login string, key string) error {
	_, err := cl.c.RegistrationWithKey(ctx, &auth.RegistrationWithKeyRequest{
		Username: login,
		Key:      key,
	})
	if err != nil {
		return model.WrapError(AuthClientRegisterError, err)
	}

	return nil
}

func (cl *AuthClientStub) Register(ctx context.Context, login string, key string) error {
	return nil
}

func (cl *AuthClient) Me(ctx context.Context) (*auth.MeMessage, error) {
	message, err := cl.c.Me(ctx, &empty.Empty{})
	if err != nil {
		return nil, model.WrapError(AuthClientMeError, err)
	}

	return message, nil
}

func (cl *AuthClientStub) Me(ctx context.Context) (*auth.MeMessage, error) {
	return &auth.MeMessage{Username: ""}, nil
}
