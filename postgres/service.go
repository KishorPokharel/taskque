package postgres

import "database/sql"

type Service struct {
	User UserService
}

func NewService(db *sql.DB) Service {
	s := Service{
		User: UserService{DB: db},
	}
	return s
}
