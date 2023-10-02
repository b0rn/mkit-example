package rest

import (
	"context"
	"errors"

	"github.com/b0rn/mkit-example/domain/usecases"
	"github.com/b0rn/mkit-example/infrastructure/config"
	"github.com/b0rn/mkit/pkg/api"
)

type RESTApiFactoryParams struct {
	Config   *config.ApisConfig
	Usecases *usecases.Usecases
}

func RESTApiFactory(ctx context.Context, cfg interface{}) (api.Api, error) {
	params, ok := cfg.(RESTApiFactoryParams)
	if !ok {
		return nil, errors.New("failed to convert configuration interface to RESTApiFactoryParams")
	}
	if params.Config == nil {
		return nil, errors.New("configuration is nil")
	}
	if params.Usecases == nil {
		return nil, errors.New("usescases is nil")
	}
	restCfg := params.Config.RESTConfig
	if !restCfg.Enabled {
		return nil, nil
	}
	r := &RestApi{
		ApisRoot: params.Config.Root,
		Config:   &params.Config.RESTConfig,
		UseCases: params.Usecases,
	}
	return r, nil
}
