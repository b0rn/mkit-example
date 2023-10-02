package dbdataservice

import (
	"context"
	"database/sql"
	"errors"

	"github.com/b0rn/mkit-example/infrastructure/config"
	dbdataservicesql "github.com/b0rn/mkit-example/infrastructure/dataservice/dbdataservice/sql"
	"github.com/b0rn/mkit/pkg/dataservice"
	"github.com/b0rn/mkit/pkg/service"
)

func DbDataServiceSQL(svc *service.Service) func(ctx context.Context, cfg interface{}) (dataservice.DataService, error) {
	return func(ctx context.Context, cfg interface{}) (dataservice.DataService, error) {
		dataServiceCfg, ok := cfg.(config.DbDataServiceConfig)
		if !ok {
			return nil, errors.New("failed to convert configuration interface to config.DataServiceConfig")
		}
		dataStoreCfg := dataServiceCfg.DataStoreConfig
		ds, err := svc.DataStoreManager.Build(ctx, dataStoreCfg.Code, dataStoreCfg)
		if err != nil {
			return nil, err
		}
		sqlDb, ok := ds.(*sql.DB)
		if !ok {
			return nil, errors.New("failed to convert DataStore to sql.DB")
		}
		return &dbdataservicesql.DbDataServiceSql{
			DataStore: sqlDb,
		}, nil
	}
}
