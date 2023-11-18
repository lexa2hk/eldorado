package data

import "time"

type Task struct {
	ID          string    `db:"id"`
	UserID      string    `db:"user_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	IsCompleted bool      `db:"is_completed"`
	IsDeleted   bool      `db:"is_deleted"`
	CreatedOn   time.Time `db:"created_on"`
}

type StatisticTask struct {
	Email     string    `db:"email"`
	Title     string    `db:"title"`
	CreatedOn time.Time `db:"created_on"`
}
