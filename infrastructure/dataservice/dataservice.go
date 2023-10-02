package dataservice

import (
	"context"

	"github.com/b0rn/mkit-example/domain/aggregates"
)

type DbDataServiceInterface interface {
	CreateUser(ctx context.Context, u *aggregates.User) error
	ReadUser(ctx context.Context, username string) (*aggregates.User, error)
	DeleteUser(ctx context.Context, username string) error
	DeleteAllUsers(ctx context.Context) error
	GracefulShutdown() error
}

type MqDataServiceInterface interface {
	NotifyUserCreation(ctx context.Context, u *aggregates.User) error
	GracefulShutdown() error
}
