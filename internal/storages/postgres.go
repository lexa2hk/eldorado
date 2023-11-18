package storages

import (
	"database/sql"

	_ "github.com/lib/pq"
)

const UniqueViolationCode = "23505"

// NewDBPool returns a connection pool for databsae.
func NewDBPool(driver, url string) (*sql.DB, error) {
	db, err := sql.Open(driver, url)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, err
}
