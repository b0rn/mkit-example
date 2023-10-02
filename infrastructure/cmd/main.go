package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/b0rn/mkit-example/api"
	"github.com/b0rn/mkit-example/domain/aggregates"
	"github.com/b0rn/mkit-example/domain/usecases"
	"github.com/b0rn/mkit-example/infrastructure/config"
	"github.com/b0rn/mkit-example/infrastructure/dataservicefactory/dbdataservice"
	"github.com/b0rn/mkit-example/infrastructure/dataservicefactory/mqdataservice"
	"github.com/b0rn/mkit-example/infrastructure/datastorefactory"
	"github.com/b0rn/mkit-example/infrastructure/servicehelper"
	"github.com/b0rn/mkit-example/infrastructure/usecasefactory"
	"github.com/b0rn/mkit/pkg/mlog"
	"github.com/b0rn/mkit/pkg/service"
	"golang.org/x/sync/errgroup"
)

func main() {
	svc := service.NewService()
	errGroup, cancel := handleSignals(svc)
	svc.LoadEnvVars("env", "./.env")
	var cfg config.Config
	svc.BuildConfig("yaml", "./config.yml", &cfg)
	svc.EnableLogger(cfg.LogConfig, nil)
	setFactories(svc)

	ctx := context.Background()
	usecasesS, err := buildUsecases(ctx, svc, &cfg)
	handleFatalError(err, cancel)

	handleFatalError(api.SetupApi(ctx, svc, &cfg, usecasesS), cancel)
	handleFatalError(svc.ApiManager.ServeAll(ctx), cancel)

	go standaloneFunction(context.Background(), usecasesS)
	handleFatalError(errGroup.Wait(), cancel)

}

func handleFatalError(err error, cancel context.CancelFunc) {
	if err == nil {
		return
	}
	mlog.Logger.Fatal().Err(err).Msg("exiting")
}

func handleSignals(svc *service.Service) (*errgroup.Group, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
		<-c
		cancel()
	}()

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		<-gCtx.Done()
		return svc.GracefulShutdown(context.Background())
	})
	return g, cancel
}

func setFactories(svc *service.Service) {
	dataStoreMgr := svc.DataStoreManager
	dataStoreMgr.SetFactory(config.DATASTORE.SQL, datastorefactory.SqlFactory)
	dataStoreMgr.SetFactory(config.DATASTORE.NATS, datastorefactory.NatsFactory)

	dataServiceMgr := svc.DataServiceManager
	dataServiceMgr.SetFactory(config.DATASERVICE.DB_DATA, dbdataservice.DbDataServiceSQL(svc))
	dataServiceMgr.SetFactory(config.DATASERVICE.MQ_DATA, mqdataservice.MqDataServiceNats(svc))

	usecaseMgr := svc.UsecaseManager
	usecaseMgr.SetFactory(config.USECASE.MANAGE_USERS, usecasefactory.ManageUsersUsecaseFactory(svc))
}

func buildUsecases(ctx context.Context, svc *service.Service, cfg *config.Config) (*usecases.Usecases, error) {
	manageUsersUc, err := servicehelper.GetManageUsersUsecase(ctx, svc.UsecaseManager, cfg)
	if err != nil {
		return nil, err
	}
	return &usecases.Usecases{
		ManageUsers: manageUsersUc,
	}, nil
}

func standaloneFunction(ctx context.Context, usecasesS *usecases.Usecases) {
	manageUsersUc := usecasesS.ManageUsers
	if err := manageUsersUc.DeleteAllUsers(ctx); err != nil {
		mlog.Logger.Error().Err(err).Msg("")
	}
	user := &aggregates.User{
		Username: "b0rn",
	}
	if err := manageUsersUc.CreateUser(ctx, user); err != nil {
		mlog.Logger.Error().Err(err).Msg("")
	}
	user, err := manageUsersUc.GetUser(ctx, "b0rn")
	if err != nil {
		mlog.Logger.Error().Err(err).Msg("")
		return
	} else if user != nil {
		mlog.Logger.Info().Msgf("got result for user b0rn : %v", user)
	} else {
		mlog.Logger.Info().Msg("no result for user b0rn")
		return
	}
	if err := manageUsersUc.DeleteUser(ctx, "b0rn"); err != nil {
		mlog.Logger.Error().Err(err).Msg("")
	}
}
