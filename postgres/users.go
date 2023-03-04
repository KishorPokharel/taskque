package postgres

import (
	"context"
	"database/sql"
	"time"
)

type UserService struct {
	DB *sql.DB
}

type User struct {
	ID        int64
	Username  string
	Email     string
	Password  Password
	CreatedAt time.Time
}

type Password struct {
	PlainText string
	Hash      []byte
}

func (us UserService) CreateUser(user *User) (*User, error) {
	queryInsertUser := `
        insert into users (username, email, password)
        values ($1, $2, $3)
        returning id, username, email, created_at
    `
	argsInsertUser := []any{user.Username, user.Email, user.Password.Hash}

	tx, err := us.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	userRow := tx.QueryRowContext(context.Background(), queryInsertUser, argsInsertUser...)
	err = userRow.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	queryInsertList := `
        insert into taskorder (user_id, value)
        values ($1, array[]::bigint[])
    `
	_, err = tx.ExecContext(context.Background(), queryInsertList, user.ID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}
