package servicehelper

import (
	"context"
	"errors"

	"github.com/b0rn/mkit-example/domain/usecases"
	"github.com/b0rn/mkit-example/infrastructure/config"
	"github.com/b0rn/mkit/pkg/usecase"
)

func GetManageUsersUsecase(ctx context.Context, usecaseManager usecase.UseCaseManager, cfg *config.Config) (usecases.ManageUsersUseCaseInterface, error) {
	uc, err := usecaseManager.Build(ctx, config.USECASE.MANAGE_USERS, cfg)
	if err != nil {
		return nil, err
	}
	mu, ok := uc.(usecases.ManageUsersUseCaseInterface)
	if !ok {
		return nil, errors.New("failed to convert Usecase to ManageUsersUseCaseInterface")
	}
	return mu, nil
}
