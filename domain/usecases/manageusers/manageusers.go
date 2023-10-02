package manageusers

import (
	"context"
	"errors"

	"github.com/b0rn/mkit-example/domain/aggregates"
	"github.com/b0rn/mkit-example/domain/usecases"
	"github.com/b0rn/mkit-example/infrastructure/dataservice"
)

type ManageUsersUseCase struct {
	DbDataService dataservice.DbDataServiceInterface
	MqDataService dataservice.MqDataServiceInterface
}

func (mu *ManageUsersUseCase) CreateUser(ctx context.Context, u *aggregates.User) error {
	if u.Username == "" {
		return &usecases.ErrUseCase{
			StatusCode: 400,
			Err:        errors.New("empty username"),
		}
	}
	err := mu.DbDataService.CreateUser(ctx, u)
	if err != nil {
		return err
	}
	return mu.MqDataService.NotifyUserCreation(ctx, u)
}
func (mu *ManageUsersUseCase) GetUser(ctx context.Context, username string) (*aggregates.User, error) {
	return mu.DbDataService.ReadUser(ctx, username)
}

func (mu *ManageUsersUseCase) DeleteUser(ctx context.Context, username string) error {
	return mu.DbDataService.DeleteUser(ctx, username)
}

func (mu *ManageUsersUseCase) DeleteAllUsers(ctx context.Context) error {
	return mu.DbDataService.DeleteAllUsers(ctx)
}

func (mu *ManageUsersUseCase) GracefulShutdown() error {
	errDb := mu.DbDataService.GracefulShutdown()
	errMq := mu.MqDataService.GracefulShutdown()
	return errors.Join(errDb, errMq)
}
