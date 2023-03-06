package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

var ErrInvalidData = errors.New("invalid task id or source index or destination index")

type Task struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id,omitempty"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type TaskService struct {
	DB *sql.DB
}

func (ts TaskService) GetAll(userID int64) ([]Task, error) {
	query := `
        select x.id, content, created_at 
        from taskorder, unnest(value)
        with ordinality as x(id)
        join tasks on tasks.id = x.id where tasks.user_id = $1;
    `
	rows, err := ts.DB.Query(query, userID)
	if err != nil {
		return []Task{}, err
	}
	tasks := []Task{}
	for rows.Next() {
		task := Task{}
		rows.Scan(&task.ID, &task.Content, &task.CreatedAt)
		tasks = append(tasks, task)
	}
	if err := rows.Close(); err != nil {
		return []Task{}, err
	}
	return tasks, nil
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

func (ts TaskService) SortTask(userID, taskID, sourceIndex, destinationIndex int64) error {
	query := `
        select value from taskorder
        where user_id = $1
    `
	row := ts.DB.QueryRow(query, userID)
	ids := []int64{}
	if err := row.Scan(pq.Array(&ids)); err != nil {
		return err
	}
	// TODO: Refactor
	idx, ok := taskIdInArray(taskID, ids)
	if !ok || idx != sourceIndex || destinationIndex > int64(len(ids)-1) {
		return ErrInvalidData
	}
	move(taskID, sourceIndex, destinationIndex, ids)
	queryInsert := `
        update taskorder
        set value = $1
        where user_id = $2
    `
	args := []any{pq.Array(ids), userID}
	_, err := ts.DB.Exec(queryInsert, args...)
	return err
}

func move(taskID, sourceIndex, destinationIndex int64, ids []int64) {
	if sourceIndex < destinationIndex {
		for i := sourceIndex; i < destinationIndex; i++ {
			ids[i] = ids[i+1]
		}
		ids[destinationIndex] = taskID
		return
	} else {
		for i := sourceIndex; i > destinationIndex; i-- {
			ids[i] = ids[i-1]
		}
		ids[destinationIndex] = taskID
		return
	}
}

func taskIdInArray(taskID int64, ids []int64) (int64, bool) {
	for idx, val := range ids {
		if val == taskID {
			return int64(idx), true
		}
	}
	return 0, false
}
