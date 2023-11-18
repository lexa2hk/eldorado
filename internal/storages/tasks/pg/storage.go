package pg

import (
	"context"
	"database/sql"

	"github.com/romankravchuk/eldorado/internal/data"
	"github.com/romankravchuk/eldorado/internal/storages"
	"github.com/romankravchuk/eldorado/internal/storages/tasks"
)

// TasksStorage is a postgres implementation of tasks.Storage.
type TasksStorage struct {
	db *sql.DB
}

// New returns new TasksStorage instance with postgres db pool.
//
// If db is nil returns storages.ErrNilDBPool.
func New(db *sql.DB) (*TasksStorage, error) {
	if db == nil {
		return nil, storages.ErrNilDBPool
	}

	return &TasksStorage{db: db}, nil
}

// UncompletedStatistic returns 5 or low uncompleted tasks for each user.
func (s *TasksStorage) UncompletedStatistic(ctx context.Context) ([]data.StatisticTask, error) {
	const query = "WITH ranked_tasks AS (SELECT id, user_id, title, created_on, ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY created_on DESC) AS rank FROM tasks WHERE is_completed = false) SELECT u.email, rt.title, rt.created_on FROM users u JOIN ranked_tasks rt ON rt.user_id = u.id WHERE rt.rank <= 5 ORDER BY u.email, rt.rank"

	prepareCtx, cancel := context.WithTimeout(ctx, storages.PrepareTimeout)
	defer cancel()

	stmt, err := s.db.PrepareContext(prepareCtx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	var tasks []data.StatisticTask
	for rows.Next() {
		var st data.StatisticTask
		if err = rows.Scan(&st.Email, &st.Title, &st.CreatedOn); err != nil {
			break
		}
		tasks = append(tasks, st)
	}

	if closeErr := rows.Close(); closeErr != nil {
		return nil, closeErr
	}

	if err != nil {
		return nil, err
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// FindByUserID returns a list of tasks for a given user.
func (s *TasksStorage) FindByUserID(ctx context.Context, userID string) ([]data.Task, error) {
	const query = "SELECT id, user_id, title, description, is_completed, created_on FROM tasks WHERE user_id = $1 AND is_deleted = false"

	prepareCtx, cancel := context.WithTimeout(ctx, storages.PrepareTimeout)
	defer cancel()

	stmt, err := s.db.PrepareContext(prepareCtx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		return nil, err
	}

	var tasks []data.Task
	for rows.Next() {
		var task data.Task
		if err = rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.IsCompleted, &task.CreatedOn); err != nil {
			break
		}
		tasks = append(tasks, task)
	}

	if closeErr := rows.Close(); closeErr != nil {
		return nil, closeErr
	}

	if err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// Save saves a tasks to the database.
//
// If save succeeds ID, IsCompleted and CreatedOn fields are filled.
func (s *TasksStorage) Save(ctx context.Context, t *data.Task) error {
	const query = "INSERT INTO tasks (user_id, title, description) VALUES ($1, $2, $3) RETURNING id, is_completed, created_on"

	prepareCtx, cancel := context.WithTimeout(ctx, storages.PrepareTimeout)
	defer cancel()

	stmt, err := s.db.PrepareContext(prepareCtx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, t.UserID, t.Title, t.Description).
		Scan(&t.ID, &t.IsCompleted, &t.CreatedOn)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes a task from the database.
//
// Actually set is_delete = true.
// If count of affected rows is not 1 returns tasks.ErrNotFound.
func (s *TasksStorage) Delete(ctx context.Context, id string) error {
	const query = "UPDATE tasks SET is_deleted = true WHERE id = $1"

	prepareCtx, cancel := context.WithTimeout(ctx, storages.PrepareTimeout)
	defer cancel()

	stmt, err := s.db.PrepareContext(prepareCtx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count != 1 {
		return tasks.ErrNotFound
	}

	return nil
}

// Update updates a task in the database.
//
// If count of affected rows is not 1 returns tasks.ErrNotFound.
func (s *TasksStorage) Update(ctx context.Context, t *data.Task) error {
	query := "UPDATE tasks SET title = $1, description = $2, is_completed = $3 WHERE id = $4"

	prepareCtx, cancel := context.WithTimeout(ctx, storages.PrepareTimeout)
	defer cancel()

	stmt, err := s.db.PrepareContext(prepareCtx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, t.Title, t.Description, t.IsCompleted, t.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count != 1 {
		return tasks.ErrNotFound
	}

	return nil
}
