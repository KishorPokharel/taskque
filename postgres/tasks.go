package postgres

import (
	"context"
	"database/sql"
	"time"
)

type Task struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type TaskService struct {
	DB *sql.DB
}

func (ts TaskService) Insert(task *Task) error {
	queryInsertTask := `
        insert into tasks (user_id, content)
        values ($1, $2) returning id, created_at
    `
	args := []any{task.UserID, task.Content}
	tx, err := ts.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	taskRow := tx.QueryRowContext(context.Background(), queryInsertTask, args...)
	err = taskRow.Scan(&task.ID, &task.CreatedAt)
	if err != nil {
		return err
	}
	queryInsertOrder := `
        update taskorder set value = array_append(value, $1)
        where user_id = $2;
    `
	_, err = tx.ExecContext(context.Background(), queryInsertOrder, task.ID, task.UserID)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil

}
