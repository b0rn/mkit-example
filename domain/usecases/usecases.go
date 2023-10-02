package usecases

import (
	"context"

	"github.com/b0rn/mkit-example/domain/aggregates"
)

type Usecases struct {
	ManageUsers ManageUsersUseCaseInterface
}

type ManageUsersUseCaseInterface interface {
	// CreateUser creates a user in the database and dispatches
	// an async message to notify that a user has been created
	CreateUser(ctx context.Context, u *aggregates.User) error
	// GetUser retrieves a user from the database
	GetUser(ctx context.Context, username string) (*aggregates.User, error)
	// Deletes a user
	DeleteUser(ctx context.Context, username string) error
	// Deletes all users
	DeleteAllUsers(ctx context.Context) error

	GracefulShutdown() error
}
