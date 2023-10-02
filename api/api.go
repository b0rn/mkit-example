package api

import (
	"context"
	"errors"
	"net/http"

	httpapi "github.com/b0rn/mkit-example/api/http"
	"github.com/b0rn/mkit-example/api/nats"
	"github.com/b0rn/mkit-example/api/rest"
	"github.com/b0rn/mkit-example/domain/usecases"
	"github.com/b0rn/mkit-example/infrastructure/config"
	"github.com/b0rn/mkit-example/infrastructure/datastore"
	"github.com/b0rn/mkit-example/infrastructure/datastorefactory"
	"github.com/b0rn/mkit/pkg/api"
	"github.com/b0rn/mkit/pkg/service"
)

func SetupApi(ctx context.Context, svc *service.Service, cfg *config.Config, usecasesS *usecases.Usecases) error {
	apiMgr := svc.ApiManager
	apiMgr.SetFactory(config.API.REST, rest.RESTApiFactory)
	apiMgr.SetFactory(config.API.HTTP, httpapi.HttpApiFactory)
	apiMgr.SetFactory(config.API.NATS, nats.NatsApiFactory)

	// REST
	restHandler, err := setupREST(ctx, apiMgr, &cfg.ApisConfig, usecasesS)
	if err != nil {
		return err
	}
	// HTTP
	_, err = apiMgr.Build(ctx, config.API.HTTP, httpapi.HttpApiFactoryParams{
		Cfg:     cfg.ApisConfig,
		Handler: restHandler,
	})
	if err != nil {
		return err
	}
	// NATS
	if err := setupNats(ctx, svc, &cfg.ApisConfig.MQConfig, usecasesS); err != nil {
		return err
	}
	return nil
}

func setupREST(ctx context.Context, apiMgr api.ApiManager, cfg *config.ApisConfig, usecasesS *usecases.Usecases) (http.Handler, error) {
	api, err := apiMgr.Build(ctx, config.API.REST, rest.RESTApiFactoryParams{
		Config:   cfg,
		Usecases: usecasesS,
	})
	if err != nil {
		return nil, err
	}
	restApi, ok := api.(*rest.RestApi)
	if !ok {
		return nil, errors.New("failed to convert api to RestApi")
	}
	return restApi.Init(ctx)
}

func setupNats(ctx context.Context, svc *service.Service, cfg *config.MQApiConfig, usecasesS *usecases.Usecases) error {
	ds, err := svc.DataStoreManager.Build(ctx, config.DATASTORE.NATS, datastorefactory.NatsFactoryParams{
		Cfg:      cfg.Config,
		Stream:   cfg.Stream,
		Subjects: []string{cfg.CreateUserSubject},
	})
	if err != nil {
		return err
	}
	natsDs, ok := ds.(*datastore.NatsDataStore)
	if !ok {
		return errors.New("failed to convert DataStore to NatsDataStore")
	}
	_, err = svc.ApiManager.Build(ctx, config.API.NATS, nats.NatsApiFactoryParams{
		Cfg:      cfg,
		Nats:     natsDs,
		Usecases: usecasesS,
	})
	return err
}
