package usecasefactory

import (
	"context"
	"errors"

	"github.com/b0rn/mkit-example/domain/usecases/manageusers"
	"github.com/b0rn/mkit-example/infrastructure/config"
	"github.com/b0rn/mkit-example/infrastructure/dataservice"
	"github.com/b0rn/mkit/pkg/service"
	"github.com/b0rn/mkit/pkg/usecase"
)

func ManageUsersUsecaseFactory(svc *service.Service) func(ctx context.Context, cfg interface{}) (usecase.UseCase, error) {
	return func(ctx context.Context, cfg interface{}) (usecase.UseCase, error) {
		appCfg, ok := cfg.(*config.Config)
		if !ok {
			return nil, errors.New("failed to convert configuration interface to config.Config")
		}
		ucCfg := appCfg.UsecasesConfig.ManageUsersUsecaseConfig
		if !ok {
			return nil, errors.New("failed to convert configuration interface to config.ManageUsersUsecaseConfig")
		}

		// Db DataService
		dbDsCfg := ucCfg.DbDataServiceConfig
		dbDs, err := svc.DataServiceManager.Build(ctx, dbDsCfg.Code, dbDsCfg)
		if err != nil {
			return nil, err
		}
		db, ok := dbDs.(dataservice.DbDataServiceInterface)
		if !ok {
			return nil, errors.New("failed to convert DataService to dataservice.DbDataServiceInterface")
		}

		// Mq DataService
		mqDsCfg := ucCfg.MqDataServiceConfig
		mqDs, err := svc.DataServiceManager.Build(ctx, mqDsCfg.Code, mqDsCfg)
		if err != nil {
			return nil, err
		}
		mq, ok := mqDs.(dataservice.MqDataServiceInterface)
		if !ok {
			return nil, errors.New("failed to convert DataService to dataservice.MqDataServiceInterface")
		}

		return &manageusers.ManageUsersUseCase{
			DbDataService: db,
			MqDataService: mq,
		}, nil
	}
}
