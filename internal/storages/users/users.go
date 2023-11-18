package users

import (
	"context"
	"errors"

	"github.com/romankravchuk/eldorado/internal/data"
)

var (
	ErrNotFound      = errors.New("the user was not found")
	ErrAlreadyExists = errors.New("the user already exists")
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name Storage
type Storage interface {
	FindByUsername(ctx context.Context, username string) (data.User, error)
	FindByEmail(ctx context.Context, email string) (data.User, error)
	Save(ctx context.Context, u *data.User) error
}
