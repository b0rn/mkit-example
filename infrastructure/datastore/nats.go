package datastore

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NatsDataStore struct {
	Nc *nats.Conn
	Js jetstream.JetStream
}
