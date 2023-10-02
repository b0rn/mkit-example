package natstools

import (
	"context"
	"errors"

	"github.com/b0rn/mkit/pkg/mlog"
	"github.com/nats-io/nats.go/jetstream"
)

func GetOrCreateStream(ctx context.Context, js jetstream.JetStream, name string, subjects []string) (jetstream.Stream, error) {
	var s jetstream.Stream
	_, err := js.Stream(ctx, name)
	if errors.Is(err, jetstream.ErrStreamNotFound) {
		mlog.Logger.Debug().Msgf("creating stream %s with subjects %v", name, subjects)
		s, err = js.CreateStream(ctx, jetstream.StreamConfig{
			Name:     name,
			Subjects: subjects,
		})
	} else {
		mlog.Logger.Debug().Msgf("updating stream %s with subjects %v", name, subjects)
		s, err = js.UpdateStream(ctx, jetstream.StreamConfig{
			Name:     name,
			Subjects: subjects,
		})
	}
	return s, err
}
