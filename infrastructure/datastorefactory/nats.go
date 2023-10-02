package datastorefactory

import (
	"context"
	"errors"

	"github.com/b0rn/mkit-example/infrastructure/config"
	"github.com/b0rn/mkit-example/infrastructure/datastore"
	"github.com/b0rn/mkit-example/infrastructure/tool/natstools"
	ds "github.com/b0rn/mkit/pkg/datastore"
	"github.com/b0rn/mkit/pkg/mlog"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NatsFactoryParams struct {
	Cfg      config.MQDataStoreConfig
	Stream   string
	Subjects []string
}

func NatsFactory(ctx context.Context, cfg interface{}) (ds.DataStore, error) {
	params, ok := cfg.(NatsFactoryParams)
	if !ok {
		return nil, errors.New("failed to convert configuration interface to NatsFactoryParams")
	}
	natsCfg := params.Cfg
	nc, err := nats.Connect(natsCfg.URL)
	if err != nil {
		return nil, err
	}
	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}
	mlog.Logger.Info().Msg("connected to JetStream")
	_, err = natstools.GetOrCreateStream(ctx, js, params.Stream, params.Subjects)
	if err != nil {
		return nil, err
	}
	return &datastore.NatsDataStore{
		Nc: nc,
		Js: js,
	}, nil
}
