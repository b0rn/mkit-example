package sql

import (
	"context"
	"database/sql"

	"github.com/b0rn/mkit-example/domain/aggregates"
	"github.com/b0rn/mkit/pkg/mlog"
)

type DbDataServiceSql struct {
	DataStore *sql.DB
}

func (ds *DbDataServiceSql) CreateUser(ctx context.Context, u *aggregates.User) error {
	mlog.Logger.Debug().Msg("creating user " + u.Username)
	_, err := ds.DataStore.ExecContext(ctx, "INSERT INTO users(username) VALUES($1)", u.Username)
	return err
}
func (ds *DbDataServiceSql) ReadUser(ctx context.Context, username string) (*aggregates.User, error) {
	mlog.Logger.Debug().Msg("reading user " + username)
	var u aggregates.User
	err := ds.DataStore.QueryRowContext(ctx, "SELECT * FROM users WHERE username = $1", username).Scan(&u.Username)
	return &u, err
}

func (ds *DbDataServiceSql) DeleteUser(ctx context.Context, username string) error {
	mlog.Logger.Debug().Msg("deleting user " + username)
	_, err := ds.DataStore.ExecContext(ctx, "DELETE FROM users WHERE username=$1", username)
	return err
}

func (ds *DbDataServiceSql) DeleteAllUsers(ctx context.Context) error {
	mlog.Logger.Debug().Msg("deleting all users")
	_, err := ds.DataStore.ExecContext(ctx, "DELETE FROM users")
	return err
}

func (ds *DbDataServiceSql) GracefulShutdown() error {
	mlog.Logger.Info().Msg("shutting down database connection")
	return ds.DataStore.Close()
}
