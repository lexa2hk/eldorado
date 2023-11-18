package data

import (
	"time"
)

type key string

const (
	ContextKeyUser key = "user"
)

type User struct {
	ID                string    `db:"id"`
	Email             string    `db:"email"`
	Username          string    `db:"username"`
	Name              string    `db:"name"`
	EncryptedPassword string    `db:"encrypted_password"`
	CreatedOn         time.Time `db:"created_on"`
	DeletedOn         time.Time `db:"deleted_on"`
}
