package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/b0rn/mkit-example/domain/usecases"
	"github.com/b0rn/mkit-example/infrastructure/config"
	"github.com/b0rn/mkit/pkg/mlog"
	restful "github.com/emicklei/go-restful/v3"
	"github.com/rs/zerolog/hlog"
)

type customRouteFunction func(*restful.Request, *restful.Response) error
type ctxIdStartTime struct{}

func createLogFilters(ws *restful.Container) {
	logger := mlog.Logger.With().Str("api-type", config.API.REST)
	ws.Filter(httpHandlerChainToFilter(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			r = r.WithContext(context.WithValue(ctx, ctxIdStartTime{}, time.Now()))
			next.ServeHTTP(w, r)
		})
	}))
	ws.Filter(httpHandlerChainToFilter(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			l := logger.Logger()
			r = r.WithContext(l.WithContext(ctx))
			next.ServeHTTP(w, r)
		})
	}))
	ws.Filter(httpHandlerChainToFilter(hlog.RemoteAddrHandler("ip")))
	ws.Filter(httpHandlerChainToFilter(hlog.UserAgentHandler("user_agent")))
	ws.Filter(httpHandlerChainToFilter(hlog.RefererHandler("referer")))
}

func handlerWrapper(f customRouteFunction) restful.RouteFunction {
	return func(request *restful.Request, response *restful.Response) {
		err := f(request, response)
		if e := handleError(response, err); e != nil {
			err = errors.Join(err, e)
		}
		ctx := request.Request.Context()
		start := ctx.Value(ctxIdStartTime{}).(time.Time)
		logCall(request, response, err, time.Since(start))
	}
}

func handleRecover(ctx context.Context, request *restful.Request, response *restful.Response, rec any) {
	start := ctx.Value(ctxIdStartTime{}).(time.Time)
	err := response.WriteHeaderAndEntity(http.StatusInternalServerError, restEmptyResponse{
		Error: "Internal Server Error",
	})
	if err != nil {
		err = errors.Join(err, fmt.Errorf("panicked with : %v", rec))
	} else {
		err = fmt.Errorf("panicked with : %v", rec)
	}
	logCall(request, response, err, time.Since(start))
}

func handleError(response *restful.Response, err error) error {
	if err == nil {
		return nil
	}
	var useCaseError *usecases.ErrUseCase
	var status int
	if errors.As(err, &useCaseError) {
		status = useCaseError.StatusCode
		res := restEmptyResponse{
			Error:    useCaseError.Error(),
			Response: nil,
		}
		return response.WriteHeaderAndEntity(status, res)
	} else {
		status = http.StatusInternalServerError
		res := restEmptyResponse{
			Error:    err.Error(),
			Response: nil,
		}
		return response.WriteHeaderAndEntity(status, res)
	}
}

func logCall(request *restful.Request, response *restful.Response, err error, duration time.Duration) {
	l := hlog.FromRequest(request.Request)
	e := l.Info()
	if response.StatusCode() >= 500 {
		e = l.Error()
	}
	e.Err(err).
		Str("method", request.Request.Method).
		Stringer("url", request.Request.URL).
		Int("status", response.StatusCode()).
		Int("size", response.ContentLength()).
		Dur("duration", duration).
		Msg("")
}

func httpHandlerChainToFilter(f func(http.Handler) http.Handler) restful.FilterFunction {
	return func(req *restful.Request, res *restful.Response, chain *restful.FilterChain) {
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					handleRecover(r.Context(), req, res, rec)
				}
			}()
			req.Request = r
			chain.ProcessFilter(req, res)
		})
		f(h).ServeHTTP(res.ResponseWriter, req.Request)
	}
}
