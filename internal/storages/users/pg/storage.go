package pg

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/romankravchuk/eldorado/internal/data"
	"github.com/romankravchuk/eldorado/internal/storages"
	"github.com/romankravchuk/eldorado/internal/storages/users"
)

// UsersStorage is a postgres implementation of users.Storage.
type UsersStorage struct {
	db *sql.DB
}

// New retuns new UsersStorage instance with postgres db pool.
//
// If db is nil returns storages.ErrNilDBPool.
func New(db *sql.DB) (*UsersStorage, error) {
	if db == nil {
		return nil, storages.ErrNilDBPool
	}

	return &UsersStorage{
		db: db,
	}, nil
}

// FindByUsername returns user by given username.
//
// If user is not found returns users.ErrNotFound.
func (s *UsersStorage) FindByUsername(ctx context.Context, username string) (data.User, error) {
	const query = "SELECT id, email, username, encrypted_password, name, created_on FROM users WHERE username = $1 AND deleted_on IS NULL"

	return s.findUser(ctx, query, username)
}

// FindByUsername returns user by given username.
//
// If user is not found returns users.ErrNotFound.
func (s *UsersStorage) FindByEmail(ctx context.Context, email string) (data.User, error) {
	const query = "SELECT id, email, username, encrypted_password, name, created_on FROM users WHERE email = $1 AND deleted_on IS NULL"

	return s.findUser(ctx, query, email)
}

// Save saves a given user in database.
//
// If user with given email or username already exists returns users.ErrAlreadyExists.
func (s *UsersStorage) Save(ctx context.Context, u *data.User) error {
	const query = "INSERT INTO users (email, username, name, encrypted_password) VALUES ($1, $2, $3, $4) RETURNING id"

	prepareCtx, cancel := context.WithTimeout(ctx, storages.PrepareTimeout)
	defer cancel()

	stmt, err := s.db.PrepareContext(prepareCtx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, u.Email, u.Username, u.Name, u.EncryptedPassword).Scan(&u.ID)
	if err != nil {
		if psqlErr, ok := err.(*pq.Error); ok && psqlErr.Code == storages.UniqueViolationCode {
			return users.ErrAlreadyExists
		}

		return err
	}

	return nil
}

func (s *UsersStorage) findUser(ctx context.Context, query, param string) (data.User, error) {
	prepareCtx, cancel := context.WithTimeout(ctx, storages.PrepareTimeout)
	defer cancel()

	stmt, err := s.db.PrepareContext(prepareCtx, query)
	if err != nil {
		return data.User{}, err
	}
	defer stmt.Close()

	var u data.User
	err = stmt.QueryRowContext(ctx, param).
		Scan(&u.ID, &u.Email, &u.Username, &u.EncryptedPassword, &u.Name, &u.CreatedOn)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return data.User{}, users.ErrNotFound
		}

		return data.User{}, err
	}

	return u, nil
}
