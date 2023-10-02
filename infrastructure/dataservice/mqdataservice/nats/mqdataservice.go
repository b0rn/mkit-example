package nats

import (
	"context"
	"encoding/json"

	"github.com/b0rn/mkit-example/domain/aggregates"
	"github.com/b0rn/mkit-example/infrastructure/datastore"
	"github.com/b0rn/mkit/pkg/mlog"
)

type MqDataServiceNats struct {
	DataStore          *datastore.NatsDataStore
	UserCreatedSubject string
}

func (mqds *MqDataServiceNats) NotifyUserCreation(ctx context.Context, u *aggregates.User) error {
	str, err := json.Marshal(u)
	if err != nil {
		return err
	}
	mlog.Logger.Debug().Msg("publishing message on topic " + mqds.UserCreatedSubject)
	_, err = mqds.DataStore.Js.Publish(ctx, mqds.UserCreatedSubject, []byte(str))
	return err
}

func (mqds *MqDataServiceNats) GracefulShutdown() error {
	mlog.Logger.Info().Msg("shutting down message queue connection")
	mqds.DataStore.Nc.Close()
	return nil
}
