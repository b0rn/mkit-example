package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/b0rn/mkit-example/domain/aggregates"
	"github.com/b0rn/mkit-example/domain/usecases"
	"github.com/b0rn/mkit-example/infrastructure/config"
	"github.com/b0rn/mkit/pkg/mlog"
	"github.com/emicklei/go-restful/v3"
)

type RestApi struct {
	ApisRoot string
	Config   *config.RESTConfig
	UseCases *usecases.Usecases
}

type restResponse[T any] struct {
	Error    string `json:"error,omitempty"`
	Warning  string `json:"warning,omitempty"`
	Response *T     `json:"response,omitempty"`
}

type restEmptyResponse restResponse[struct{}]

func (r *RestApi) Init(ctx context.Context) (http.Handler, error) {
	path := fmt.Sprintf("%s%s", r.ApisRoot, r.Config.Path)
	mlog.Logger.Info().Msgf("configured REST API on path %s", path)
	wsContainer := restful.NewContainer()
	createLogFilters(wsContainer)
	wsContainer.EnableContentEncoding(true)

	ws := new(restful.WebService)
	ws.Path(path).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/health").To(handlerWrapper(r.health)).Writes(restEmptyResponse{}))
	ws.Route(ws.POST("/user").To(handlerWrapper(r.createUser)).Reads(
		aggregates.User{}).Writes(restEmptyResponse{}))
	ws.Route(ws.GET("/user/{username}").To(handlerWrapper(r.readUser)).
		Param(ws.PathParameter("username", "identifier of the user").DataType("string")).
		Writes(aggregates.User{}))
	ws.Route(ws.DELETE("/user/{username}").To(handlerWrapper(r.deleteUser)).Writes(restEmptyResponse{}))

	wsContainer.Add(ws)
	return wsContainer, nil
}

func (r *RestApi) Serve(ctx context.Context) error {
	return nil
}

func (r *RestApi) GracefulShutdown(ctx context.Context) error {
	return nil
}

func (r *RestApi) health(request *restful.Request, response *restful.Response) error {
	return response.WriteHeaderAndEntity(http.StatusOK, restEmptyResponse{})
}

func (r *RestApi) createUser(request *restful.Request, response *restful.Response) error {
	var createUserRequest aggregates.User
	err := request.ReadEntity(&createUserRequest)
	if err != nil {
		return &usecases.ErrUseCase{
			StatusCode: 400,
			Err:        err,
		}
	}
	err = r.UseCases.ManageUsers.CreateUser(request.Request.Context(), &createUserRequest)
	if err != nil {
		return err
	}
	return response.WriteHeaderAndEntity(http.StatusOK, restEmptyResponse{})
}

func (r *RestApi) readUser(request *restful.Request, response *restful.Response) error {
	username := request.PathParameter("username")
	mlog.Logger.Debug().Msg("username is " + username)
	u, err := r.UseCases.ManageUsers.GetUser(request.Request.Context(), username)
	if err != nil {
		return err
	}
	return response.WriteHeaderAndEntity(http.StatusOK, u)
}

func (r *RestApi) deleteUser(request *restful.Request, response *restful.Response) error {
	username := request.PathParameter("username")
	err := r.UseCases.ManageUsers.DeleteUser(request.Request.Context(), username)
	if err != nil {
		return err
	}
	return response.WriteHeaderAndEntity(http.StatusOK, restEmptyResponse{})
}
