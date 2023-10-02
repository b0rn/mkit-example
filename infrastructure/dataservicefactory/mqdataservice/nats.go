package mqdataservice

import (
	"context"
	"errors"

	"github.com/b0rn/mkit-example/infrastructure/config"
	"github.com/b0rn/mkit-example/infrastructure/dataservice/mqdataservice/nats"
	"github.com/b0rn/mkit-example/infrastructure/datastore"
	"github.com/b0rn/mkit-example/infrastructure/datastorefactory"
	"github.com/b0rn/mkit/pkg/dataservice"
	"github.com/b0rn/mkit/pkg/service"
)

func MqDataServiceNats(svc *service.Service) func(ctx context.Context, cfg interface{}) (dataservice.DataService, error) {
	return func(ctx context.Context, cfg interface{}) (dataservice.DataService, error) {
		dataServiceCfg, ok := cfg.(config.MQDataServiceConfig)
		if !ok {
			return nil, errors.New("failed to convert configuration interface to config.DataServiceConfig")
		}
		dataStoreCfg := dataServiceCfg.DataStoreConfig
		var subjects []string
		subjects = append(subjects, dataServiceCfg.UserCreatedSubject)
		params := datastorefactory.NatsFactoryParams{
			Cfg:      dataStoreCfg,
			Stream:   dataServiceCfg.StreamName,
			Subjects: subjects,
		}
		ds, err := svc.DataStoreManager.Build(ctx, dataStoreCfg.Code, params)
		if err != nil {
			return nil, err
		}
		natsDs, ok := ds.(*datastore.NatsDataStore)
		if !ok {
			return nil, errors.New("failed to convert DataStore to datastore.NatsDataStore")
		}

		return &nats.MqDataServiceNats{
			DataStore:          natsDs,
			UserCreatedSubject: dataServiceCfg.UserCreatedSubject,
		}, nil
	}
}
