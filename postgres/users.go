package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

const round = 12

var (
	ErrDuplicateEmail = errors.New("email already exists")
	ErrUserNotFound   = errors.New("user not found")
)

type UserService struct {
	DB *sql.DB
}

type User struct {
	ID        int64
	Username  string
	Email     string
	Password  password
	CreatedAt time.Time
}

type password struct {
	plainText *string
	hash      []byte
}

func (p *password) Set(pwd string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), round)
	if err != nil {
		return err
	}
	p.plainText = &pwd
	p.hash = hash
	return nil
}

func (p *password) Matches(pwd string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(pwd))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func (us UserService) Create(user *User) error {
	queryInsertUser := `
        insert into users (username, email, password)
        values ($1, $2, $3)
        returning id, created_at
    `
	argsInsertUser := []any{user.Username, user.Email, user.Password.hash}

	tx, err := us.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	userRow := tx.QueryRowContext(context.Background(), queryInsertUser, argsInsertUser...)
	err = userRow.Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		switch e := err.(type) {
		case *pq.Error:
			if e.Code == "23505" {
				return ErrDuplicateEmail
			}
		default:
			return err
		}
		return err
	}

	queryInsertList := `
        insert into taskorder (user_id, value)
        values ($1, array[]::bigint[])
    `
	_, err = tx.ExecContext(context.Background(), queryInsertList, user.ID)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
