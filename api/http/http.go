package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/b0rn/mkit-example/infrastructure/config"
	"github.com/b0rn/mkit/pkg/api"
	"github.com/b0rn/mkit/pkg/mlog"
)

type HttpApi struct {
	Address      string
	Port         uint16
	Handler      http.Handler
	errorChannel chan error
	server       *http.Server
	wg           *sync.WaitGroup
}

type HttpApiFactoryParams struct {
	Cfg          config.ApisConfig
	Handler      http.Handler
	ErrorChannel chan error
}

func HttpApiFactory(ctx context.Context, cfg interface{}) (api.Api, error) {
	config, ok := cfg.(HttpApiFactoryParams)
	if !ok {
		return nil, errors.New("could not convert configuration interface to config.ApisConfig")
	}
	return &HttpApi{
		Address:      config.Cfg.Address,
		Port:         config.Cfg.Port,
		Handler:      config.Handler,
		errorChannel: make(chan error, 1),
		wg:           &sync.WaitGroup{},
	}, nil
}

func (h *HttpApi) Serve(ctx context.Context) error {
	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		addr := fmt.Sprintf("%s:%d", h.Address, h.Port)
		h.server = &http.Server{Addr: addr, Handler: h.Handler}
		mlog.Logger.Info().Msgf("Serving HTTP on %s", addr)
		h.errorChannel <- h.server.ListenAndServe()
	}()
	select {
	case err := <-h.errorChannel:
		return err
	case <-time.After(1 * time.Second):
		return nil
	}
}

func (h *HttpApi) GracefulShutdown(ctx context.Context) error {
	if h.server == nil {
		return nil
	}
	mlog.Logger.Info().Msg("shutting down HTTP server")
	err := h.server.Shutdown(ctx)
	if err != nil {
		return err
	}
	h.wg.Wait()
	return nil
}
