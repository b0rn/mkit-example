package nats

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/b0rn/mkit-example/domain/aggregates"
	"github.com/b0rn/mkit-example/domain/usecases"
	"github.com/b0rn/mkit-example/infrastructure/config"
	"github.com/b0rn/mkit-example/infrastructure/datastore"
	"github.com/b0rn/mkit-example/infrastructure/tool/natstools"
	"github.com/b0rn/mkit/pkg/api"
	"github.com/nats-io/nats.go/jetstream"
)

type natsApi struct {
	dataStore *datastore.NatsDataStore
	consumer  jetstream.Consumer
	usecases  *usecases.Usecases
}

type NatsApiFactoryParams struct {
	Cfg      *config.MQApiConfig
	Nats     *datastore.NatsDataStore
	Usecases *usecases.Usecases
}

func NatsApiFactory(ctx context.Context, cfg interface{}) (api.Api, error) {
	params, ok := cfg.(NatsApiFactoryParams)
	if !ok {
		return nil, errors.New("failed to convert configuration interface to NatsApiFactoryParams")
	}
	_, err := natstools.GetOrCreateStream(ctx, params.Nats.Js, params.Cfg.Stream, []string{params.Cfg.CreateUserSubject})
	if err != nil {
		return nil, err
	}
	cons, err := params.Nats.Js.CreateOrUpdateConsumer(ctx, params.Cfg.Stream, jetstream.ConsumerConfig{
		AckPolicy:     jetstream.AckExplicitPolicy,
		FilterSubject: params.Cfg.CreateUserSubject,
	})
	if err != nil {
		return nil, err
	}
	return &natsApi{
		dataStore: params.Nats,
		consumer:  cons,
		usecases:  params.Usecases,
	}, nil
}

func (n *natsApi) Serve(ctx context.Context) error {
	n.consumer.Consume(func(msg jetstream.Msg) {
		var u aggregates.User
		if json.Unmarshal(msg.Data(), &u) != nil {
			n.usecases.ManageUsers.CreateUser(ctx, &u)
		}
	})
	return nil
}

func (n *natsApi) GracefulShutdown(ctx context.Context) error {
	n.dataStore.Nc.Close()
	return nil
}
