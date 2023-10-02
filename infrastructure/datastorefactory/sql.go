package datastorefactory

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/b0rn/mkit-example/infrastructure/config"
	"github.com/b0rn/mkit/pkg/datastore"
	"github.com/b0rn/mkit/pkg/mlog"
	_ "github.com/lib/pq"
)

func SqlFactory(ctx context.Context, cfg interface{}) (datastore.DataStore, error) {
	sqlCfg, ok := cfg.(config.DbDataStoreConfig)
	if !ok {
		return nil, errors.New("failed to convert configuration interface to config.DbDataStoreConfig")
	}
	str := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", sqlCfg.Host, sqlCfg.Port, sqlCfg.User, sqlCfg.Password, sqlCfg.DbName)
	conn, err := sql.Open(sqlCfg.DriverName, str)
	if err != nil {
		return nil, err
	}
	mlog.Logger.Info().Msg("connected to sql database")
	return conn, nil
}
